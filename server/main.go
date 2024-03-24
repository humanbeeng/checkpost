package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/humanbeeng/checkpost/server/internal/auth"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Unable to load .env file. %v", err)
	}

	app := fiber.New()
	ac, err := auth.NewGithubAuthController()
	if err != nil {
		log.Fatalf("Unable to init auth controller. %v", err)
	}
	ac.RegisterRoutes(app)

	app.Listen(":3000")

}
