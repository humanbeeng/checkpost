package main

import (
	"encoding/json"
	"fmt"

	fiber "github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("From: %v", c.Subdomains())
		return c.SendString(msg)
	})

	app.Get("/v1/auth/google", func(c *fiber.Ctx) error {
		return c.SendString("Google login")
	})

	app.Get("/dashboard", func(c *fiber.Ctx) error {
		return c.SendString("Dashboard")
	})

	app.Get("/profile", func(c *fiber.Ctx) error {
		return c.SendString("Profile")
	})

	app.All("/v1/endpoint/:path", func(c *fiber.Ctx) error {
		var req any

		_ = c.BodyParser(&req)

		str, err := json.Marshal(req)
		if err != nil {
			fmt.Println(err)
		}

		msg := fmt.Sprintf("Path: %v \nBody: %v", c.Path(), string(str))
		fmt.Println(msg)

		return c.SendString(msg)
	})

	app.Listen(":8080")
}
