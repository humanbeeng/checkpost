package url

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type URLController struct{}

func NewURLController() *URLController {
	return &URLController{}
}

func (uc *URLController) RegisterRoutes(router fiber.Router) {
	router.Get("/url", uc.AdminHandler)
}

func (uc *URLController) AdminHandler(c *fiber.Ctx) error {
	fmt.Println("Received")
	return c.JSON("received")
}
