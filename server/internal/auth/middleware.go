package auth

import (
	"log/slog"
	"os"
	"strconv"

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
		userId, err := strconv.ParseInt(payload.Subject, 10, 64)
		if err != nil {
			return fiber.ErrUnauthorized
		}

		c.Locals("userId", userId)
		c.Locals("username", payload.Get("username"))
		return c.Next()
	}
}
