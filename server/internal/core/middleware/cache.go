package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
)

func NewCacheMiddleware() fiber.Handler {
	return cache.New(cache.Config{
		Expiration:   30 * time.Minute,
		CacheControl: true,
	})
}
