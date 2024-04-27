package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	db "github.com/humanbeeng/checkpost/server/db/sqlc"
)

func NewGuestPlanLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Locals("plan").(string) != string(db.PlanGuest)
		},
		Max:               2,
		Expiration:        time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
	})
}

func NewFreePlanLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Locals("plan").(string) != string(db.PlanFree)
		},
		Max:               5,
		Expiration:        time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Locals("username").(string)
		},
	})
}

func NewHobbyPlanLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Locals("plan").(string) != string(db.PlanFree)
		},
		Max:               10,
		Expiration:        time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Locals("username").(string)
		},
	})
}

func NewProPlanLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Locals("plan").(string) != string(db.PlanFree)
		},
		Max:               20,
		Expiration:        time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Locals("username").(string)
		},
	})
}

func NewDefaultLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.BaseURL() == "http://api.checkpost.local:3000"
		},
		Max:               60,
		Expiration:        time.Minute,
		LimiterMiddleware: limiter.SlidingWindow{},
	})
}

func NewGenerateUrlLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.BaseURL() == "http://api.checkpost.local:3000"
		},
		Max:               100,
		Expiration:        time.Minute,
		LimiterMiddleware: limiter.SlidingWindow{},
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Locals("username").(string)
		},
	})
}

func NewGenerateRandomUrlLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.BaseURL() == "http://api.checkpost.local:3000"
		},
		Max:               1,
		Expiration:        time.Minute,
		LimiterMiddleware: limiter.SlidingWindow{},
	})
}

func NewEndpointCheckLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:               10,
		Expiration:        time.Minute,
		LimiterMiddleware: limiter.SlidingWindow{},
	})
}
