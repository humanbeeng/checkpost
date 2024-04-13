package url

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
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

	slog.Info("Ws connection added", "endpoint", endpoint)
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

func NewUrlController(service *UrlService) *UrlController {
	return &UrlController{conns: &sync.Map{}, service: service}
}

func (uc *UrlController) RegisterRoutes(app *fiber.App, authmw, gl, fl, bl, pl, genLim, genRandLim, endpointCheckLim, cache fiber.Handler) {
	urlGroup := app.Group("/url")

	urlGroup.Get("/exists/:endpoint", endpointCheckLim, cache, uc.CheckEndpointExistsHandler)

	urlGroup.Post("/generate", authmw, genLim, uc.GenerateUrlHandler)
	urlGroup.Get("/generate/random", genRandLim, uc.GenerateGuestUrlHandler)

	urlGroup.All("/hook/:endpoint/:path?", gl, fl, bl, pl, uc.HookHandler)

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
	endpoint := c.Params("endpoint")
	// Check if endpoint exists
	// TODO: Revisit this context

	exists, err := uc.service.urlq.CheckEndpointExists(context.TODO(), endpoint)
	if !exists {
		slog.Info("No endpoint found", "endpoint", endpoint)
		// TODO: Format this into response
		c.WriteMessage(websocket.TextMessage, []byte("Endpoint does not exist or has expired"))
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

func (uc *UrlController) GenerateGuestUrlHandler(c *fiber.Ctx) error {
	endpoint, err := uc.service.CreateGuestUrl(c.Context())
	if err != nil {
		return fiber.NewError(err.Code, err.Message)
	}

	return c.JSON(GenerateUrlResponse{
		Url:  fmt.Sprintf("https://%v.checkpost.io", endpoint.Endpoint),
		Plan: string(endpoint.Plan),
	})
}

func (uc *UrlController) GenerateUrlHandler(c *fiber.Ctx) error {
	var req GenerateUrlRequest

	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	username, ok := c.Locals("username").(string)
	if !ok {
		return fiber.ErrInternalServerError
	}

	slog.Info("Generate url request received", "username", username)

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
	req, err := uc.service.StoreRequestDetails(c)
	if err != nil {
		return &fiber.Error{
			Code:    err.Code,
			Message: err.Message,
		}
	}

	endpoint := c.Params("endpoint", "")
	uc.BroadcastJSON(endpoint, req)

	return c.SendStatus(fiber.StatusOK)
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
	if len(endpoint) < 4 {
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Must be atleast 4 characters",
		}
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
			Message:  "Its available 🐱. Sign up and make it yours!",
		})
	} else {
		return c.JSON(CheckEndpointExistsResponse{
			Endpoint: endpoint,
			Exists:   exists,
			Message:  "That endpoint is already taken 😿. Try something else. Maybe your workplace name?",
		})
	}

}
