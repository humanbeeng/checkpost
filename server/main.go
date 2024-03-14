package main

import (
	"encoding/json"
	"fmt"

	fiber "github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/dashboard", func(c *fiber.Ctx) error {
		return c.SendString("Dashboard")
	})

	app.Get("/profile", func(c *fiber.Ctx) error {
		return c.SendString("Profile")
	})

	app.All("/url/*", func(c *fiber.Ctx) error {
		var req any

		_ = c.BodyParser(&req)

		str, _ := json.Marshal(req)

		msg := fmt.Sprintf("Path: %v \nBody: %v", c.Path(), string(str))
		fmt.Println(msg)

		return c.SendString(msg)
	})

	app.Listen(":8080")
}
