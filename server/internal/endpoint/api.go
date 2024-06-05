package endpoint

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

// TODO: Add better error messages
type EndpointController struct {
	conns   *sync.Map
	service *EndpointService
}

func (uc *EndpointController) AddRequestListener(endpoint string, conn *websocket.Conn) {
	// TODO: Resolve concurrency issue.
	// TODO: Add limit to number of active connections.
	_, loaded := uc.conns.LoadOrStore(endpoint, []*websocket.Conn{conn})
	if loaded {
		// Replace existing connection.
		// TODO: Add support for multiple connections
		c := append([]*websocket.Conn{}, conn)
		uc.conns.Store(endpoint, c)
	}

	slog.Info("Listener connection added", "endpoint", endpoint)
}

func NewEndpointController(service *EndpointService) *EndpointController {
	return &EndpointController{conns: &sync.Map{}, service: service}
}

func (uc *EndpointController) RegisterRoutes(app *fiber.App, authmw, cache fiber.Handler) {
	endpointGroup := app.Group("/endpoint")

	endpointGroup.Get("/", authmw, uc.GetUserEndpointsHandler)

	endpointGroup.Get("/exists/:endpoint", cache, uc.CheckSubdomainExistsHandler)

	endpointGroup.Post("/generate", authmw, uc.GenerateEndpointHandler)

	endpointGroup.All("/hook/:endpoint/*", uc.HookHandler)

	endpointGroup.Get("/history/:endpoint", authmw, uc.GetEndpointHistoryHandler)
	endpointGroup.Get("/request/:uuid", authmw, uc.RequestDetailsUUIDHandler)

	endpointGroup.Get("/stats/:endpoint", authmw, uc.StatsHandler)

	// TODO: Add rate/conn limiter
	endpointGroup.Use("/inspect", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	endpointGroup.Get("/inspect/:endpoint", authmw, websocket.New(uc.InspectRequestsHandler))
}

func (uc *EndpointController) InspectRequestsHandler(c *websocket.Conn) {
	slog.Info("Received websocket connection", "username", c.Locals("username").(string))

	endpoint := c.Params("endpoint", "")
	endpoint = strings.ToLower(endpoint)
	if endpoint == "" {
		slog.Info("No endpoint found", "endpoint", endpoint)
		c.WriteJSON(fiber.Error{
			Code:    fiber.StatusNotFound,
			Message: "Endpoint has either expired or not yet created.",
		})
		c.Close()
	}

	// Check if endpoint exists
	// TODO: Revisit this context
	exists, err := uc.service.endpointq.CheckEndpointExists(context.TODO(), endpoint)
	if !exists {
		slog.Info("No endpoint found", "endpoint", endpoint)

		c.WriteJSON(fiber.Error{
			Code:    fiber.StatusNotFound,
			Message: "Endpoint has either expired or not yet created.",
		})
		c.Close()
		return
	}

	if err != nil {
		slog.Error("unable to check if endpoint exists", "err", err)
		c.WriteMessage(websocket.TextMessage, []byte("Oops! Something went wrong."))
		c.Close()
	}

	uc.AddRequestListener(endpoint, c)

	for {
		if _, _, err := c.ReadMessage(); err != nil {
			break
		}
	}
}

// Returns status of a given endpoint
func (uc *EndpointController) StatsHandler(c *fiber.Ctx) error {
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

func (uc *EndpointController) RequestDetailsHandler(c *fiber.Ctx) error {
	reqIdStr := c.Params("requestid", "")
	if reqIdStr == "" {
		return fiber.NewError(
			fiber.StatusNotFound,
			"No request id found",
		)
	}

	reqId, parseErr := strconv.ParseInt(reqIdStr, 10, 64)
	if parseErr != nil {
		slog.Error("unable to convert request id from path to int", "err", parseErr)
		return fiber.ErrBadRequest
	}

	req, err := uc.service.GetRequestDetails(c.Context(), reqId)
	if err != nil {
		return &fiber.Error{Code: err.Code, Message: err.Message}
	}

	return c.JSON(req)
}

func (uc *EndpointController) RequestDetailsUUIDHandler(c *fiber.Ctx) error {
	uuid := c.Params("uuid", "")
	if uuid == "" {
		return fiber.NewError(
			fiber.StatusNotFound,
			"No uuid found",
		)
	}

	req, err := uc.service.GetRequestByUUID(c.Context(), uuid)
	if err != nil {
		return &fiber.Error{Code: err.Code, Message: err.Message}
	}

	return c.JSON(req)
}

type GenerateEndpointRequest struct {
	Endpoint string `json:"endpoint"`
}

type GenerateEndpointResponse struct {
	Endpoint  string    `json:"endpoint"`
	ExpiresAt time.Time `json:"expires_at"`
	Plan      string    `json:"plan"`
}

func (uc *EndpointController) GenerateEndpointHandler(c *fiber.Ctx) error {
	var req GenerateEndpointRequest

	if err := c.BodyParser(&req); err != nil {
		slog.Error("Malformed request payload", "err", err)
		return fiber.ErrBadRequest
	}

	username, ok := c.Locals("username").(string)
	if !ok {
		return fiber.ErrInternalServerError
	}

	endpoint, err := uc.service.CreateEndpoint(c.Context(), username, req.Endpoint)
	if err != nil {
		return &fiber.Error{
			Code:    err.Code,
			Message: err.Message,
		}
	}

	res := GenerateEndpointResponse{
		Endpoint:  endpoint.Endpoint,
		ExpiresAt: endpoint.ExpiresAt.Time,
	}
	return c.JSON(res)
}

func (uc *EndpointController) HookHandler(c *fiber.Ctx) error {
	// TODO: return request details
	// Get the Content-Type header from the request
	contentType := c.Get(fiber.HeaderContentType)

	// Print the Content-Type
	fmt.Println("Content-Type:", contentType)

	endpoint := c.Params("endpoint", "")

	if endpoint == "" {
		return &fiber.Error{
			Code:    http.StatusNotFound,
			Message: "Endpoint has either expired or not created",
		}
	}

	body := c.Body()

	// Note: key is string and value is []string
	headers := c.GetReqHeaders()
	var ip string

	// Specific to railway.app deployment
	envoyAddr, ok := headers["X-Envoy-External-Address"]
	if ok {
		ip = envoyAddr[0]
	} else {
		ip = c.IP()
	}
	path := c.Params("*", "/")

	method := c.Method()
	query := c.Queries()

	hookReq := HookRequest{
		Endpoint:     endpoint,
		UUID:         c.Locals("requestid").(string),
		Path:         path,
		Headers:      headers,
		QueryParams:  query,
		SourceIp:     ip,
		Method:       method,
		Content:      string(body),
		ContentSize:  int32(len(body)),
		ResponseCode: http.StatusOK,
	}

	requestRecord, endpointErr := uc.service.StoreRequestDetails(c.Context(), hookReq)
	if endpointErr != nil {
		return &fiber.Error{
			Code:    endpointErr.Code,
			Message: endpointErr.Message,
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

func (uc *EndpointController) GetUserEndpointsHandler(c *fiber.Ctx) error {
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
	Requests []HookRequest `json:"requests"`
}

func (uc *EndpointController) GetEndpointHistoryHandler(c *fiber.Ctx) error {
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

	offsetStr := c.Query("limit", "0")
	offset, err := strconv.ParseInt(offsetStr, 10, 32)
	if err != nil {
		return fiber.ErrBadRequest
	}
	userId := c.Locals("userId").(int64)

	reqs, serviceErr := uc.service.GetEndpointRequestHistory(c.Context(), endpoint, userId, int32(limit), int32(offset))
	if serviceErr != nil {
		return &fiber.Error{
			Code:    serviceErr.Code,
			Message: serviceErr.Message,
		}
	}

	res := GetEndpointsHistoryResponse{
		Requests: reqs,
	}
	fmt.Println("Returning", len(res.Requests), "num of requests")
	return c.JSON(res)
}

type CheckSubdomainExistsResponse struct {
	Endpoint string `json:"endpoint"`
	Exists   bool   `json:"exists"`
	Message  string `json:"message"`
}

func (uc *EndpointController) CheckSubdomainExistsHandler(c *fiber.Ctx) error {
	endpoint := c.Params("endpoint", "")

	if endpoint == "" {
		return fiber.ErrBadRequest
	}

	subdomainExists, err := uc.service.CheckEndpointExists(c.Context(), endpoint)
	if err != nil {
		return &fiber.Error{
			Code:    err.Code,
			Message: err.Message,
		}
	}

	switch subdomainExists {
	case Available:
		{
			return c.JSON(CheckSubdomainExistsResponse{
				Endpoint: endpoint,
				Exists:   false,
				Message:  string(Available),
			})
		}
	case Taken:
		{
			return c.JSON(CheckSubdomainExistsResponse{
				Endpoint: endpoint,
				Exists:   true,
				Message:  string(Taken),
			})
		}
	case ReservedCompany:
		{
			return c.JSON(CheckSubdomainExistsResponse{
				Endpoint: endpoint,
				Exists:   false,
				Message:  string(ReservedCompany),
			})
		}
	default:
		{
			return c.JSON(CheckSubdomainExistsResponse{
				Endpoint: endpoint,
				Exists:   true,
				Message:  string(Error),
			})
		}
	}
}

func (uc *EndpointController) BroadcastJSON(endpoint string, data any) {
	slog.Info("Broadcasting incoming request", "endpoint", endpoint)
	connAny, ok := uc.conns.Load(endpoint)
	if !ok {
		slog.Info("No active listeners found", "endpoint", endpoint)
		return
	}

	conns := connAny.([]*websocket.Conn)
	for _, c := range conns {
		err := c.WriteJSON(data)
		if err != nil {
			slog.Error("unable to broadcast json msg", "dest", c.RemoteAddr(), "err", err)
		}
	}
}
