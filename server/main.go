package main

import (
	"encoding/json"
	"fmt"

	fiber "github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/about", func(c *fiber.Ctx) error {
		return c.SendString("about page")
	})

	app.All(":path", func(c *fiber.Ctx) error {
		var req any

		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.ErrBadRequest.Code, fiber.ErrBadRequest.Message)
		}

		str, err := json.Marshal(req)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(string(str))

		return c.SendString(c.Params("path"))
	})

	app.Listen(":8080")
}
