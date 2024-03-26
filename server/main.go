package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/humanbeeng/checkpost/server/internal/admin"
	"github.com/humanbeeng/checkpost/server/internal/auth"
	"github.com/humanbeeng/checkpost/server/internal/url"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Unable to load .env file. %v", err)
	}

	app := fiber.New()
	app.Use(cors.New())
	pmw := auth.NewPasetoMiddleware()

	ac, err := auth.NewGithubAuthHandler()
	if err != nil {
		log.Fatalf("Unable to init auth controller. %v", err)
	}
	ac.RegisterRoutes(app)

	adc := admin.NewAdminController()
	adc.RegisterRoutes(app, &pmw)

	url := url.NewURLHandler()
	url.RegisterRoutes(app)

	app.Listen(":3000")
}
