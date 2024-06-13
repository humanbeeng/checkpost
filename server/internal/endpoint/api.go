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

	"github.com/gofiber/contrib/socketio"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/humanbeeng/checkpost/server/internal/core"
)

type EndpointController struct {
	conns   *sync.Map
	tokens  *sync.Map
	pv      *core.PasetoVerifier
	service *EndpointService
}

func (ec *EndpointController) AddRequestListener(endpoint string, conn *websocket.Conn) {

	// TODO: Add limit to number of active connections.
	_, loaded := ec.conns.LoadOrStore(endpoint, []*websocket.Conn{conn})
	if loaded {
		// Replace existing connection.
		// TODO: Add support for multiple connections

		// slog.Info("Closing existing ws connections", "num_connections", len(conns.([]*websocket.Conn)))
		// for _, conn := range conns.([]*websocket.Conn) {
		// 	conn.Close()
		// }

		c := append([]*websocket.Conn{}, conn)
		ec.conns.Store(endpoint, c)
	}

	slog.Info("Listener connection added", "endpoint", endpoint)
}

func NewEndpointController(service *EndpointService, pv *core.PasetoVerifier) *EndpointController {
	return &EndpointController{conns: &sync.Map{}, service: service, tokens: &sync.Map{}, pv: pv}
}

func (ec *EndpointController) RegisterRoutes(app *fiber.App, authmw, cache fiber.Handler) {
	endpointGroup := app.Group("/endpoint")

	endpointGroup.Get("/", authmw, ec.GetUserEndpointsHandler)

	endpointGroup.Get("/exists/:endpoint", cache, ec.CheckSubdomainExistsHandler)

	endpointGroup.Post("/generate", authmw, ec.GenerateEndpointHandler)

	endpointGroup.All("/hook/:endpoint/*", ec.HookHandler)

	endpointGroup.Get("/history/:endpoint", authmw, ec.GetEndpointHistoryHandler)
	endpointGroup.Get("/request/:uuid", authmw, ec.RequestDetailsUUIDHandler)

	endpointGroup.Get("/stats/:endpoint", authmw, ec.StatsHandler)

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

	endpointGroup.Get("/inspect/:endpoint", socketio.New(func(kws *socketio.Websocket) {

		endpoint := kws.Params("endpoint")

		// Add the connection to the list of the connected clients
		// The UUID is generated randomly and is the key that allow
		// socketio to manage Emit/EmitTo/Broadcast

		// Every websocket connection has an optional session key => value storage
		kws.SetAttribute("endpoint", endpoint)

		slog.Info("Inspecting endpoint", "endpoint", endpoint)

		kws.Emit([]byte("Hello"), socketio.TextMessage)

	}))

	socketio.On(socketio.EventConnect, func(ep *socketio.EventPayload) {
		fmt.Println("Connected")
	})
	socketio.On(socketio.EventDisconnect, func(ep *socketio.EventPayload) {
		fmt.Println("Disconnected")
	})
	socketio.On(socketio.EventClose, func(ep *socketio.EventPayload) {
		fmt.Println("Closed")
	})
}

func (ec *EndpointController) InspectRequestsHandler(c *websocket.Conn) {
	endpoint := c.Params("endpoint", "")
	if endpoint == "" {
		c.WriteJSON(fiber.Error{
			Code:    fiber.StatusNotFound,
			Message: "Endpoint has either expired or not yet created.",
		})
		c.Close()
		return
	}

	endpoint = strings.ToLower(endpoint)

	// Authorize websocket connection by checking token
	token := c.Query("token", "")
	if token == "" {
		slog.Warn("No token passed. Unauthorized access", "endpoint", endpoint)
		c.WriteJSON(fiber.Error{
			Code:    fiber.StatusUnauthorized,
			Message: fiber.ErrUnauthorized.Message,
		})
		c.Close()
		return
	}

	_, err := ec.pv.VerifyToken(token)

	if err != nil {
		c.WriteJSON(fiber.Error{
			Code:    fiber.StatusUnauthorized,
			Message: fiber.ErrUnauthorized.Message,
		})
		c.Close()
		return
	}

	// Check if endpoint exists
	exists, err := ec.service.endpointq.CheckEndpointExists(context.Background(), endpoint)
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
		c.WriteJSON(WebsocketPayload{
			Code:    fiber.StatusInternalServerError,
			Message: fiber.ErrInternalServerError.Message,
		})
		c.Close()
	}

	ec.AddRequestListener(endpoint, c)

	for {
		if _, _, err := c.ReadMessage(); err != nil {
			break
		}
	}
}

// Returns status of a given endpoint
func (ec *EndpointController) StatsHandler(c *fiber.Ctx) error {
	endpoint := c.Params("endpoint", "")
	endpoint = strings.ToLower(endpoint)
	if endpoint == "" {
		return fiber.ErrBadRequest
	}

	stats, err := ec.service.GetEndpointStats(c.Context(), endpoint)
	if err != nil {
		return &fiber.Error{
			Code:    err.Code,
			Message: err.Message,
		}
	}
	return c.JSON(stats)
}

func (ec *EndpointController) RequestDetailsHandler(c *fiber.Ctx) error {
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

	req, err := ec.service.GetRequestDetails(c.Context(), reqId)
	if err != nil {
		return &fiber.Error{Code: err.Code, Message: err.Message}
	}

	return c.JSON(req)
}

func (ec *EndpointController) RequestDetailsUUIDHandler(c *fiber.Ctx) error {
	uuid := c.Params("uuid", "")
	if uuid == "" {
		return fiber.NewError(
			fiber.StatusNotFound,
			"No uuid found",
		)
	}

	req, err := ec.service.GetRequestByUUID(c.Context(), uuid)
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

func (ec *EndpointController) GenerateEndpointHandler(c *fiber.Ctx) error {
	var req GenerateEndpointRequest

	if err := c.BodyParser(&req); err != nil {
		slog.Error("Malformed request payload", "err", err)
		return fiber.ErrBadRequest
	}

	username, ok := c.Locals("username").(string)
	if !ok {
		return fiber.ErrInternalServerError
	}

	endpoint, err := ec.service.CreateEndpoint(c.Context(), username, req.Endpoint)
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

type WebsocketPayload struct {
	Code        int         `json:"code"`
	HookRequest HookRequest `json:"hook_request"`
	Message     string      `json:"message"`
}

func (ec *EndpointController) HookHandler(c *fiber.Ctx) error {
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

	requestRecord, endpointErr := ec.service.StoreRequestDetails(c.Context(), hookReq)
	if endpointErr != nil {
		return &fiber.Error{
			Code:    endpointErr.Code,
			Message: endpointErr.Message,
		}
	}

	hookReq.ExpiresAt = requestRecord.ExpiresAt.Time
	hookReq.CreatedAt = requestRecord.CreatedAt.Time

	payload := WebsocketPayload{
		HookRequest: hookReq,
		Code:        200,
	}

	ec.BroadcastJSON(endpoint, payload)

	return c.SendStatus(fiber.StatusOK)
}

type GetUserEndpointsResponse struct {
	Endpoints []Endpoint `json:"endpoints"`
}

func (ec *EndpointController) GetUserEndpointsHandler(c *fiber.Ctx) error {
	userId, ok := c.Locals("userId").(int64)
	if !ok {
		return fiber.ErrBadRequest
	}
	slog.Info("Requesting user endpoints", "userId", userId)

	endpoints, err := ec.service.GetUserEndpoints(c.Context(), userId)
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

func (ec *EndpointController) GetEndpointHistoryHandler(c *fiber.Ctx) error {
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

	reqs, serviceErr := ec.service.GetEndpointRequestHistory(c.Context(), endpoint, userId, int32(limit), int32(offset))
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

func (ec *EndpointController) CheckSubdomainExistsHandler(c *fiber.Ctx) error {
	endpoint := c.Params("endpoint", "")

	if endpoint == "" {
		return fiber.ErrBadRequest
	}

	subdomainExists, err := ec.service.CheckEndpointExists(c.Context(), endpoint)
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

func (ec *EndpointController) BroadcastJSON(endpoint string, data any) {
	connsAny, ok := ec.conns.Load(endpoint)
	if !ok {
		slog.Info("No active listeners found", "endpoint", endpoint)
		return
	}

	conns := connsAny.([]*websocket.Conn)

	slog.Info("Broadcasting incoming request", "endpoint", endpoint, "num_listeners", len(conns))

	for _, c := range conns {
		err := c.WriteJSON(data)
		if err != nil {
			slog.Error("unable to broadcast json msg", "dest", c.RemoteAddr(), "err", err)
		}
	}
}
