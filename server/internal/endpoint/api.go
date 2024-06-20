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
	"github.com/google/uuid"
	"github.com/humanbeeng/checkpost/server/internal/core"
)

type EndpointController struct {
	clients   *sync.Map
	wsManager *WSManager
	tokens    *sync.Map
	pv        *core.PasetoVerifier
	service   *EndpointService
}

func (ec *EndpointController) AddListener(endpoint string, conn *websocket.Conn) {
	// TODO: Add limit to number of active connections.
	sessionId := conn.Locals("session_id").(string)

	listenersAny, loaded := ec.clients.LoadOrStore(endpoint, map[string]*websocket.Conn{sessionId: conn})
	if loaded {
		listeners := listenersAny.(map[string]*websocket.Conn)
		listeners[sessionId] = conn
	}

	closeHandler := func(code int, text string) error {
		slog.Info("Closing connection", "endpoint", endpoint, "session_id", sessionId)
		_ = conn.WriteControl(websocket.CloseMessage, nil, time.Now().Add(time.Second))
		ec.removeListener(endpoint, sessionId)
		return nil
	}

	ec.clients.Store(endpoint, listenersAny)
	conn.SetCloseHandler(closeHandler)

	slog.Info("Listener added", "endpoint", endpoint, "session_id", conn.Locals("session_id"), "num_listeners", len(listenersAny.(map[string]*websocket.Conn)))
}

func NewEndpointController(service *EndpointService, wsManager *WSManager, pv *core.PasetoVerifier) *EndpointController {
	return &EndpointController{clients: &sync.Map{}, service: service, tokens: &sync.Map{}, pv: pv, wsManager: wsManager}
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

	// endpointGroup.Get("/inspect/:endpoint", websocket.New(ec.InspectRequestsHandler))
	endpointGroup.Get("/inspect/:endpoint", websocket.New(ec.Inspect))
}

func (ec *EndpointController) Inspect(c *websocket.Conn) {
	endpoint := c.Params("endpoint")
	// ec.wsManager.AddConn(endpoint, c)

	// client := WSClient{
	// 	conn:      c,
	// 	endpoint:  endpoint,
	// 	sessionId: c.Locals("requestid").(string),
	// 	manager:   ec.wsManager,
	// 	egress:    make(chan WSMessage),
	// }
	// cl := NewWSClient(c.Locals("requestid").(string), endpoint, c, ec.wsManager)

	// cl.conn.Conn.SetPongHandler(cl.pongHandler)

	go func(c *websocket.Conn) error {

		_, _, err := c.ReadMessage()
		if err != nil {
			slog.Info("Unable to read message from ws conn", "endpoint", endpoint, "session_id", c.Locals("requestid"), "err", err)
			return err
		}
		return nil
	}(c)

	go func(c *websocket.Conn) error {

		t := time.NewTicker(pingInterval)
		defer func(t *time.Ticker) {
			t.Stop()
			// c.manager.RemoveConn(endpoint, c.Locals("sessionid").(string))
		}(t)

		for {
			select {
			case <-t.C:
				{
					fmt.Println("Sending ping", "session_id", c.Locals("requestid"))
					if err := c.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
						slog.Error("unable to send ping", "session_id", c.Locals("requestid"), "err", err)
						// Trigger cleanup func
						return err
					}
				}
			}
		}
	}(c)
	// go cl.readMessages()
	// go cl.writeMessage()

	c.SetPongHandler(func(appData string) error {
		fmt.Println("pong")
		c.SetReadDeadline(time.Now().Add(time.Second))
		return nil
	})
	go func() {
		t := time.NewTicker(time.Second)
		for {
			select {
			case <-t.C:
				{
					slog.Info("Ticker")
				}
			}
		}
	}()
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

	payload, err := ec.pv.VerifyToken(token)
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
			Code: fiber.StatusInternalServerError,
		})
		c.Close()
	}

	sessionId, err := uuid.NewV7()
	if err != nil {
		slog.Error("unable to generate uuid for ws connection", err)
		return
	}

	c.Locals("session_id", sessionId.String())
	c.Locals("username", payload.Get("username"))
	c.Locals("plan", payload.Get("plan"))
	c.Locals("role", payload.Get("role"))

	ec.AddListener(endpoint, c)

	for {
		// Listen for ping message = ""
		_, msg, err := c.ReadMessage()
		if err != nil {
			slog.Info("Unable to read message from ws conn", "endpoint", endpoint, "session_id", sessionId.String(), "err", err)
			break
		}

		if string(msg) == "" {
			err = c.WriteMessage(websocket.TextMessage, []byte(""))
			if err != nil {
				slog.Error("unable to send pong", "endpoint", endpoint, "err", err)
			}
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
	listenersAny, ok := ec.clients.Load(endpoint)
	if !ok {
		slog.Info("No active listeners found", "endpoint", endpoint)
		return
	}

	listeners := listenersAny.(map[string]*websocket.Conn)
	for _, c := range listeners {
		slog.Info("Broadcasting json", "endpoint", endpoint, "session_id", c.Locals("session_id"))
		err := c.WriteJSON(data)
		if err != nil {
			slog.Error("unable to broadcast json msg", "dest", c.RemoteAddr(), "err", err)
			sessionId := c.Locals("session_id").(string)
			ec.removeListener(endpoint, sessionId)
		}
	}
}

func (ec *EndpointController) removeListener(endpoint string, sessionId string) error {
	listenersAny, loaded := ec.clients.Load(endpoint)
	if loaded {
		listeners := listenersAny.(map[string]*websocket.Conn)
		delete(listeners, sessionId)
		ec.clients.Store(endpoint, listeners)
		slog.Info("Removed listener", "endpoint", endpoint, "session_id", sessionId, "num_listeners", len(listeners))
	}
	return nil
}
