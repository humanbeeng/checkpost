package main

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	db "github.com/humanbeeng/checkpost/server/db/sqlc"
	"github.com/humanbeeng/checkpost/server/internal/admin"
	"github.com/humanbeeng/checkpost/server/internal/auth"
	"github.com/humanbeeng/checkpost/server/internal/url"
	"github.com/jackc/pgx/v5"

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
	ctx := context.Background()

	connUrl := os.Getenv("POSTGRES_URL")

	conn, err := pgx.Connect(ctx, connUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(ctx)

	queries := db.New(conn)

	ac, err := auth.NewGithubAuthHandler(queries)
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
