package auth

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
)

func NewPasetoMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		key := os.Getenv("PASETO_KEY")
		pv, err := NewPasetoVerifier(key)
		if err != nil {
			return err
		}

		token := c.Cookies("token", "")
		payload, err := pv.VerifyToken(token)
		if err != nil {
			fmt.Println("Error:", err)
			return fiber.NewError(401)
		}

		c.Locals("username", payload.Subject)
		return c.Next()

	}
}
