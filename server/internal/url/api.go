package url

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
)

// TODO: Add better error messages
type URLController struct {
	service *UrlService
}

func NewEndpointHandler(service *UrlService) *URLController {
	return &URLController{service: service}
}

func (uc *URLController) RegisterRoutes(app *fiber.App, tmw *fiber.Handler) {
	// TODO: Add rate limiter
	urlGroup := app.Group("/url")
	urlGroup.Post("/generate", *tmw, uc.GenerateURLHandler)
	urlGroup.All("/hook/:endpoint/:path?", uc.HookHandler)
	urlGroup.Get("/:path/:request-id", uc.RequestDetailsHandler)
	urlGroup.Get("/stats", uc.StatsHandler)
}

// Returns status of a given endpoint and request-id
func (uc *URLController) StatsHandler(c *fiber.Ctx) error {
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
func (uc *URLController) RequestDetailsHandler(c *fiber.Ctx) error {
	return fiber.ErrBadGateway
	path := c.Params("path")
	reqId := c.Params("requestId")
	res := map[string]string{
		"path": path,
		"req":  reqId,
	}
	return c.JSON(res)
}

type GenerateUrlRequest struct {
	Endpoint string `json:"endpoint"`
}

type GenerateUrlResponse struct {
	Url       string    `json:"url"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (uc *URLController) GenerateURLHandler(c *fiber.Ctx) error {
	var req GenerateUrlRequest

	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	var username string
	usernameLocal, ok := c.Locals("username").(string)
	if !ok {
		return fiber.ErrInternalServerError
	}

	if usernameLocal == "" {
		slog.Info("Generate url request received from guest user")
	} else {
		username = usernameLocal
		slog.Info("Generate url request received", "username", username)
	}

	url, err := uc.service.GenerateUrl(c.Context(), username, req.Endpoint)
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

func (uc *URLController) HookHandler(c *fiber.Ctx) error {
	err := uc.service.StoreRequestDetails(c)
	if err != nil {
		return &fiber.Error{
			Code:    err.Code,
			Message: err.Message,
		}
	}
	return c.SendStatus(fiber.StatusOK)
}
