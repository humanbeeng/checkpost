package middleware

import (
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

func NewSubdomainRouterMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		subdomains := c.Subdomains()
		if len(subdomains) == 0 {
			slog.Error("No subdomains found")
			return c.Next()
		}
		subdomain := subdomains[0]
		switch subdomain {
		// TODO: Add reserved subdomain routing
		case "api", "localhost:3000":
			{
				return c.Next()
			}
		default:
			{
				hook := fmt.Sprintf("/endpoint/hook/%s%s", subdomain, c.Path())
				c.Path(hook)
			}
		}
		return c.Next()
	}
}
