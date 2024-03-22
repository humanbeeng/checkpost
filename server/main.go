package main

import (
	// "encoding/json"
	// "fmt"
	// "strings"

	"fmt"
	"time"

	"github.com/humanbeeng/checkpost/server/internal/auth"
	// _ "github.com/jackc/pgx/v5"
	// fiber "github.com/gofiber/fiber/v2"
)

func main() {
	pv, _ := auth.NewPasetoVerifier("TW9XIRA4EXFJ7ZJTA5G6CV1KZSMMUGYH")

	token, err := pv.CreateToken("nithin", time.Hour)
	if err != nil {
		return
	}

	fmt.Println(token)
	// app := fiber.New()
	//
	// app.Get("/dashboard", func(c *fiber.Ctx) error {
	// 	return c.SendString("Dashboard")
	// })
	//
	// app.Get("/profile", func(c *fiber.Ctx) error {
	// 	return c.SendString("Profile")
	// })
	//
	// app.All("/url/*", func(c *fiber.Ctx) error {
	// 	var req any
	//
	// 	_ = c.BodyParser(&req)
	//
	// 	strBytes, _ := json.Marshal(req)
	//
	// 	str := string(strBytes)
	// 	ip := c.Query("ip", "unknown")
	// 	path := c.Path()
	// 	path, _ = strings.CutPrefix(path, "/url")
	// 	method := c.Method()
	//
	// 	msg := fmt.Sprintf("Path: %v \nBody: %v\nSource IP: %v\nMethod: %v", path, str, ip, method)
	//
	// 	return c.SendString(msg)
	// })
	//
	// app.Listen(":8080")
}
