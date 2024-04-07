package url

import (
	"log/slog"
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
	urlGroup.Get("/generate/random", genRandLim, uc.GenerateRandomUrlHandler)
	urlGroup.Post("/generate", authmw, genLim, uc.GenerateUrlHandler)
	urlGroup.All("/hook/:endpoint/:path?", gl, fl, nbl, pl, uc.HookHandler)
	urlGroup.Get("/:path/:request-id", uc.RequestDetailsHandler)
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
	return fiber.ErrBadGateway
	path := c.Params("path")
	reqId := c.Params("request-id")
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
