package admin

import (
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
	return c.SendString("admin")
}
