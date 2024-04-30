package url

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/humanbeeng/checkpost/server/internal/core"
)

// TODO: Add better error messages
type UrlController struct {
	conns   *sync.Map
	service *UrlService
}

func (uc *UrlController) AddRequestListener(endpoint string, conn *websocket.Conn) {
	// TODO: Resolve concurrency issue.
	// TODO: Add limit to number of active connections.
	conns, loaded := uc.conns.LoadOrStore(endpoint, []*websocket.Conn{conn})
	if loaded {
		c := conns.([]*websocket.Conn)
		c = append(c, conn)
		uc.conns.Store(endpoint, c)
	}

	slog.Info("Listener connection added", "endpoint", endpoint)
}

func NewUrlController(service *UrlService) *UrlController {
	return &UrlController{conns: &sync.Map{}, service: service}
}

func (uc *UrlController) RegisterRoutes(app *fiber.App, authmw, freeLim, basicLim, proLim, generateUrlLim, endpointCheckLim, cache fiber.Handler) {
	urlGroup := app.Group("/url")

	urlGroup.Get("/", authmw, uc.GetUserEndpointsHandler)

	urlGroup.Get("/exists/:endpoint", endpointCheckLim, cache, uc.CheckEndpointExistsHandler)

	urlGroup.Post("/generate", authmw, generateUrlLim, uc.GenerateUrlHandler)

	urlGroup.All("/hook/:endpoint/*", freeLim, basicLim, proLim, uc.HookHandler)

	urlGroup.Get("/history/:endpoint", uc.GetEndpointHistoryHandler)
	urlGroup.Get("/request/:requestid", uc.RequestDetailsHandler)

	urlGroup.Get("/stats/:endpoint", uc.StatsHandler)

	// TODO: Add rate/conn limiter
	urlGroup.Use("/inspect", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	urlGroup.Get("/inspect/:endpoint", websocket.New(uc.InspectRequestsHandler))
}

func (uc *UrlController) InspectRequestsHandler(c *websocket.Conn) {
	endpoint := c.Params("endpoint", "")
	endpoint = strings.ToLower(endpoint)
	if endpoint == "" {
		slog.Info("No endpoint found", "endpoint", endpoint)
		c.WriteJSON(fiber.Error{
			Code:    fiber.StatusNotFound,
			Message: "URL has either expired or not yet created.",
		})
		c.Close()
	}

	// Check if endpoint exists
	// TODO: Revisit this context
	exists, err := uc.service.urlq.CheckEndpointExists(context.TODO(), endpoint)
	if !exists {
		slog.Info("No endpoint found", "endpoint", endpoint)

		c.WriteJSON(fiber.Error{
			Code:    fiber.StatusNotFound,
			Message: "URL has either expired or not yet created.",
		})
		c.Close()
		return
	}

	if err != nil {
		slog.Error("Unable to check if endpoint exists", "err", err)
		c.WriteMessage(websocket.TextMessage, []byte("Oops! Something went wrong."))
		c.Close()
	}

	// TODO: Add authorization
	uc.AddRequestListener(endpoint, c)

	for {
		if _, _, err := c.ReadMessage(); err != nil {
			break
		}
	}
}

// Returns status of a given endpoint
func (uc *UrlController) StatsHandler(c *fiber.Ctx) error {
	endpoint := c.Params("endpoint", "")
	endpoint = strings.ToLower(endpoint)
	if endpoint == "" {
		return fiber.ErrBadRequest
	}

	stats, err := uc.service.GetEndpointStats(c.Context(), endpoint)
	if err != nil {
		return &fiber.Error{
			Code:    err.Code,
			Message: err.Message,
		}
	}
	return c.JSON(stats)
}

// TODO: Implement this
func (uc *UrlController) RequestDetailsHandler(c *fiber.Ctx) error {
	reqIdStr := c.Params("requestid", "")
	if reqIdStr == "" {
		return fiber.NewError(
			fiber.StatusNotFound,
			"No request id found",
		)
	}

	reqId, parseErr := strconv.ParseInt(reqIdStr, 10, 64)
	if parseErr != nil {
		slog.Error("Unable to convert request id from path to int", "err", parseErr)
		return fiber.ErrBadRequest
	}

	req, err := uc.service.GetRequestDetails(c.Context(), reqId)
	if err != nil {
		return &fiber.Error{Code: err.Code, Message: err.Message}
	}

	return c.JSON(req)
}

type GenerateUrlRequest struct {
	Endpoint string `json:"endpoint"`
}

type GenerateUrlResponse struct {
	Url       string    `json:"url"`
	ExpiresAt time.Time `json:"expires_at"`
	Plan      string    `json:"plan"`
}

func (uc *UrlController) GenerateUrlHandler(c *fiber.Ctx) error {
	var req GenerateUrlRequest

	if err := c.BodyParser(&req); err != nil {
		slog.Error("Malformed request payload", "err", err)
		return fiber.ErrBadRequest
	}

	username, ok := c.Locals("username").(string)
	if !ok {
		return fiber.ErrInternalServerError
	}

	endpoint, err := uc.service.CreateUrl(c.Context(), username, req.Endpoint)
	if err != nil {
		return &fiber.Error{
			Code:    err.Code,
			Message: err.Message,
		}
	}

	res := GenerateUrlResponse{
		Url:       endpoint.Endpoint,
		ExpiresAt: endpoint.ExpiresAt.Time,
	}
	return c.JSON(res)
}

func (uc *UrlController) HookHandler(c *fiber.Ctx) error {
	// TODO: return request details
	// Get the Content-Type header from the request
	contentType := c.Get(fiber.HeaderContentType)

	// Print the Content-Type
	fmt.Println("Content-Type:", contentType)

	endpoint := c.Params("endpoint", "")

	if endpoint == "" {
		return &fiber.Error{
			Code:    http.StatusNotFound,
			Message: "URL has either expired or not created",
		}
	}

	var req any
	_ = c.BodyParser(&req)

	strBytes, _ := json.Marshal(req)
	body := string(strBytes)
	// Note: key is string and value is []string
	headers := c.GetReqHeaders()
	ip := c.IP()
	path := c.Params("+", "/")

	method := c.Method()
	query := c.Queries()

	hookReq := HookRequest{
		Endpoint:     endpoint,
		Path:         path,
		Headers:      headers,
		Query:        query,
		SourceIp:     ip,
		Method:       method,
		Content:      body,
		ContentSize:  len(body),
		ResponseCode: http.StatusOK,
	}

	requestRecord, urlErr := uc.service.StoreRequestDetails(c.Context(), hookReq)
	if urlErr != nil {
		return &fiber.Error{
			Code:    urlErr.Code,
			Message: urlErr.Message,
		}
	}

	hookReq.ExpiresAt = requestRecord.ExpiresAt.Time
	hookReq.CreatedAt = requestRecord.CreatedAt.Time

	uc.BroadcastJSON(endpoint, hookReq)

	return c.SendStatus(fiber.StatusOK)
}

type GetUserEndpointsResponse struct {
	Endpoints []Endpoint `json:"endpoints"`
}

func (uc *UrlController) GetUserEndpointsHandler(c *fiber.Ctx) error {
	userId, ok := c.Locals("userId").(int64)
	if !ok {
		return fiber.ErrBadRequest
	}
	slog.Info("Requesting user endpoints", "userId", userId)

	endpoints, err := uc.service.GetUserEndpoints(c.Context(), userId)
	if err != nil {
		return &fiber.Error{
			Code:    err.Code,
			Message: err.Message,
		}
	}

	res := GetUserEndpointsResponse{
		Endpoints: endpoints,
	}

	return c.JSON(res)
}

type GetEndpointsHistoryResponse struct {
	Requests []Request `json:"requests"`
}

func (uc *UrlController) GetEndpointHistoryHandler(c *fiber.Ctx) error {
	endpoint := c.Params("endpoint", "")
	if endpoint == "" {
		slog.Info("No endpoint found in path params")
		return fiber.ErrBadRequest
	}

	limitStr := c.Query("limit", "20")
	limit, err := strconv.ParseInt(limitStr, 10, 32)
	if err != nil {
		return fiber.ErrBadRequest
	}

	offsetStr := c.Query("limit", "1")
	offset, err := strconv.ParseInt(offsetStr, 10, 32)
	if err != nil {
		return fiber.ErrBadRequest
	}

	reqs, serviceErr := uc.service.GetEndpointRequestHistory(c.Context(), endpoint, int32(limit), int32(offset))
	if serviceErr != nil {
		return &fiber.Error{
			Code:    serviceErr.Code,
			Message: serviceErr.Message,
		}
	}

	res := GetEndpointsHistoryResponse{
		Requests: reqs,
	}
	return c.JSON(res)
}

type CheckEndpointExistsResponse struct {
	Endpoint string `json:"endpoint"`
	Exists   bool   `json:"exists"`
	Message  string `json:"message"`
}

func (uc *UrlController) CheckEndpointExistsHandler(c *fiber.Ctx) error {
	endpoint := c.Params("endpoint", "")

	if endpoint == "" {
		return fiber.ErrBadRequest
	}

	if len(endpoint) < 4 || len(endpoint) > 10 {
		return &fiber.Error{
			Code:    http.StatusBadRequest,
			Message: "Subdomain should be 4 to 10 characters.",
		}
	}

	if _, ok := core.ReservedSubdomains[endpoint]; ok {
		return c.JSON(CheckEndpointExistsResponse{
			Endpoint: endpoint,
			Exists:   true,
			Message:  fmt.Sprintf("Subdomain %s is reserved.", endpoint),
		})
	}

	exists, err := uc.service.CheckEndpointExists(c.Context(), endpoint)
	if err != nil {
		return &fiber.Error{
			Code:    err.Code,
			Message: err.Message,
		}
	}

	if exists {
		return c.JSON(CheckEndpointExistsResponse{
			Endpoint: endpoint,
			Exists:   exists,
			Message:  "That subdomain is already taken 😿. Try something else?",
		})
	}

	return c.JSON(CheckEndpointExistsResponse{
		Endpoint: endpoint,
		Exists:   exists,
		Message:  "Its available 🐱. Sign up and make it yours!",
	})
}

func (uc *UrlController) BroadcastJSON(endpoint string, data any) {
	slog.Info("Broadcasting JSON", "endpoint", endpoint)
	connAny, ok := uc.conns.Load(endpoint)
	if !ok {
		slog.Info("No active listeners found", "endpoint", endpoint)
		return
	}

	conns := connAny.([]*websocket.Conn)

	for _, c := range conns {
		err := c.WriteJSON(data)
		if err != nil {
			slog.Error("Unable to broadcast json msg", "dest", c.RemoteAddr(), "err", err)
		}
	}
}
