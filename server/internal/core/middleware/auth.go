package middleware

import (
	"log/slog"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/humanbeeng/checkpost/server/internal/core"
)

// Allows only signed in users to access a given API
func NewAuthRequiredMiddleware(pv *core.PasetoVerifier) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Cookies("token", "")
		if token == "" {
			slog.Info("Received empty token")
			return fiber.ErrUnauthorized
		}

		payload, err := pv.VerifyToken(token)
		if err != nil {
			slog.Error("unable to verify token", "err", err)
			return fiber.ErrUnauthorized
		}
		userId, err := strconv.ParseInt(payload.Subject, 10, 64)
		if err != nil {
			return fiber.ErrUnauthorized
		}

		c.Locals("userId", userId)
		c.Locals("username", payload.Get("username"))
		c.Locals("plan", payload.Get("plan"))
		c.Locals("role", payload.Get("role"))
		return c.Next()
	}
}
