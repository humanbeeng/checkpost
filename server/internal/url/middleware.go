package url

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
			return fiber.ErrBadGateway
		}
		subdomain := subdomains[0]
		switch subdomain {
		// TODO: Add reserved subdomain routing
		case "api":
			{
				return c.Next()
			}
		default:
			{
				hook := fmt.Sprintf("/url/hook/%v%v", subdomain, c.Path())
				c.Path(hook)
			}
		}
		return c.Next()
	}
}
