package admin

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

type AdminController struct{}

func NewAdminController() *AdminController {
	return &AdminController{}
}

func (ac *AdminController) RegisterRoutes(app *fiber.App, pasetoMiddleware *fiber.Handler) {
	apiGroup := app.Group("admin", *pasetoMiddleware)
	apiGroup.Get("/dashboard", ac.AdminHandler)
}

func (ac *AdminController) AdminHandler(c *fiber.Ctx) error {
	var req any

	_ = c.BodyParser(&req)

	// strBytes, _ := json.Marshal(req)
	// reqBody := string(strBytes)

	ip := c.Query("ip", "unknown")
	path := c.Path()
	path, _ = strings.CutPrefix(path, "/url")
	method := c.Method()
	username := c.Locals("username").(string)

	// headers, _ := json.MarshalIndent(c.GetReqHeaders(), "", "  ")

	res := struct {
		IP       string `json:"ip"`
		Path     string `json:"path"`
		Method   string `json:"method"`
		Username string `json:"username"`
	}{
		IP:       ip,
		Path:     path,
		Method:   method,
		Username: username,
	}

	return c.JSON(res)
}
