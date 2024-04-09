package url

import (
	"log/slog"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

// TODO: Add better error messages
type UrlController struct {
	service *UrlService
}

func NewUrlController(service *UrlService) *UrlController {
	return &UrlController{service: service}
}

func (uc *UrlController) RegisterRoutes(app *fiber.App, authmw, gl, fl, nbl, pl, genLim, genRandLim fiber.Handler) {
	urlGroup := app.Group("/url")

	urlGroup.Post("/generate", authmw, genLim, uc.GenerateUrlHandler)
	urlGroup.Get("/generate/random", genRandLim, uc.GenerateRandomUrlHandler)

	urlGroup.All("/hook/:endpoint/:path?", gl, fl, nbl, pl, uc.HookHandler)

	urlGroup.Get("/history/:endpoint", uc.GetEndpointHistoryHandler)
	urlGroup.Get("/request/:requestid", uc.RequestDetailsHandler)

	urlGroup.Get("/stats", uc.StatsHandler)
}

// Returns status of a given endpoint and request-id
func (uc *UrlController) StatsHandler(c *fiber.Ctx) error {
	// url := c.Query("url", "")
	// if url == "" {
	// 	return fiber.ErrBadRequest
	// }
	//
	// res := map[string]string{
	// 	"req": url,
	// }
	// return c.JSON(res)
	return c.SendString("Ok")
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
}

func (uc *UrlController) GenerateRandomUrlHandler(c *fiber.Ctx) error {
	url, err := uc.service.CreateRandomUrl(c.Context(), nil)
	if err != nil {
		return fiber.NewError(err.Code, err.Message)
	}

	return c.JSON(GenerateUrlResponse{
		Url: url,
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

	url, err := uc.service.CreateUrl(c.Context(), username, req.Endpoint)
	if err != nil {
		return &fiber.Error{
			Code:    err.Code,
			Message: err.Message,
		}
	}

	res := GenerateUrlResponse{
		Url: url,
		// TODO: Add plan based expiry
		ExpiresAt: time.Now().Add(time.Hour * 24),
	}
	return c.JSON(res)
}

func (uc *UrlController) HookHandler(c *fiber.Ctx) error {
	err := uc.service.StoreRequestDetails(c)
	if err != nil {
		return &fiber.Error{
			Code:    err.Code,
			Message: err.Message,
		}
	}
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
