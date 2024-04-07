package middleware

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/gofiber/fiber/v2"
	db "github.com/humanbeeng/checkpost/server/db/sqlc"
	"github.com/humanbeeng/checkpost/server/internal/core"
)

func NewExtractPayloadMiddleware(pv *core.PasetoVerifier) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Cookies("token", "")
		if token == "" {
			c.Locals("plan", string(db.PlanGuest))
			return c.Next()
		}

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
		c.Locals("plan", payload.Get("plan"))
		c.Locals("role", payload.Get("role"))
		fmt.Println(c.Locals("userId"), c.Locals("username"))
		return c.Next()
	}
}
