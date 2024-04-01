package auth

import (
	"log/slog"
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
			slog.Error("Unable to verify token", "err", err)
			return fiber.ErrUnauthorized
		}
		slog.Info("Payload", "subject", payload.Subject)
		c.Locals("email", payload.Subject)
		return c.Next()
	}
}
