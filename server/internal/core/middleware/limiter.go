package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// TODO: Implement endpoint plan based rate limiting
func NewFreePlanLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		// Next: func(c *fiber.Ctx) bool {
		// 	return c.IsFromLocal()
		// },
		Max:               3,
		Expiration:        time.Second,
		LimiterMiddleware: limiter.FixedWindow{},
	})
}

func NewBasicPlanLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.IsFromLocal()
		},
		Max:               10,
		Expiration:        time.Second,
		LimiterMiddleware: limiter.FixedWindow{},
	})
}

func NewProPlanLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.IsFromLocal()
		},
		Max:               20,
		Expiration:        time.Second,
		LimiterMiddleware: limiter.FixedWindow{},
	})
}

func NewGlobalLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		// Next: func(c *fiber.Ctx) bool {
		// 	return c.IsFromLocal()
		// },
		Max:               20,
		Expiration:        10 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
	})
}

func NewGenerateUrlLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.IsFromLocal()
		},
		Max:               100,
		Expiration:        time.Minute,
		LimiterMiddleware: limiter.SlidingWindow{},
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Locals("username").(string)
		},
	})
}

func NewEndpointCheckLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:               10,
		Expiration:        time.Minute,
		LimiterMiddleware: limiter.SlidingWindow{},
	})
}
