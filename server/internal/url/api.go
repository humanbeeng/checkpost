package url

import (
	"encoding/json"
	"strings"

	"github.com/gofiber/fiber/v2"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type URLController struct{}

func NewURLHandler() *URLController {
	return &URLController{}
}

func (uc *URLController) RegisterRoutes(app *fiber.App) {
	urlGroup := app.Group("/url")
	urlGroup.Post("/generate", uc.GenerateURLHandler)
	urlGroup.All("/hook/:path", uc.HookHandler)
	urlGroup.Get("/:path/:requestId", uc.RequestDetailsHandler)
	urlGroup.Get("/stats", uc.StatsHandler)
}

// Returns status of a given endpoint and request-id
func (uc *URLController) StatsHandler(c *fiber.Ctx) error {
	url := c.Query("url", "")
	if url == "" {
		return fiber.ErrBadRequest
	}

	res := map[string]string{
		"req": url,
	}
	return c.JSON(res)
}

func (uc *URLController) RequestDetailsHandler(c *fiber.Ctx) error {
	path := c.Params("path")
	reqId := c.Params("requestId")
	res := map[string]string{
		"path": path,
		"req":  reqId,
	}
	return c.JSON(res)
}

type GenerateURLRequest struct {
	Endpoint string `json:"endpoint"`
}

func (uc *URLController) GenerateURLHandler(c *fiber.Ctx) error {
	var req GenerateURLRequest

	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	id, _ := gonanoid.New()
	return c.SendString(id)
}

func (uc *URLController) HookHandler(c *fiber.Ctx) error {
	var req any

	_ = c.BodyParser(&req)

	strBytes, _ := json.Marshal(req)

	body := string(strBytes)
	ip := c.Query("ip", "Unknown")
	path := c.Path()
	path, _ = strings.CutPrefix(path, "/url/hook")
	method := c.Method()

	res := map[string]string{
		"body":   body,
		"ip":     ip,
		"method": method,
		"path":   path,
	}
	return c.JSON(res)
}
