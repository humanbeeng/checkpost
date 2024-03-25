package main

import (
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/humanbeeng/checkpost/server/internal/admin"
	"github.com/humanbeeng/checkpost/server/internal/auth"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Unable to load .env file. %v", err)
	}

	app := fiber.New()

	app.Use(cors.New())
	adc := admin.NewAdminController()

	pmw := auth.NewPasetoMiddleware()

	apiGroup := app.Group("admin", pmw)

	apiGroup.Get("/dashboard", adc.AdminHandler)

	ac, err := auth.NewGithubAuthController()
	if err != nil {
		log.Fatalf("Unable to init auth controller. %v", err)
	}
	ac.RegisterRoutes(app)

	url := app.Group("url")

	// // Apply rate limiting middleware to url endpoint
	// url.Use(limiter.New(limiter.Config{
	// 	Max:        2, // limit each IP to 2 requests per second
	// 	Expiration: 1 * time.Second,
	// }))

	url.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(map[string]string{"message": "received"})
	})

	app.Get("/dashboard", func(c *fiber.Ctx) error {
		var req any

		_ = c.BodyParser(&req)

		// strBytes, _ := json.Marshal(req)
		// reqBody := string(strBytes)

		ip := c.Query("ip", "unknown")
		path := c.Path()
		path, _ = strings.CutPrefix(path, "/url")
		method := c.Method()

		// headers, _ := json.MarshalIndent(c.GetReqHeaders(), "", "  ")

		res := struct {
			IP     string `json:"ip"`
			Path   string `json:"path"`
			Method string `json:"method"`
		}{
			IP:     ip,
			Path:   path,
			Method: method,
		}

		return c.JSON(res)

	})

	app.Listen(":3000")
}
