package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"

	db "github.com/humanbeeng/checkpost/server/db/sqlc"
	"github.com/humanbeeng/checkpost/server/internal/admin"
	"github.com/humanbeeng/checkpost/server/internal/auth"
	"github.com/humanbeeng/checkpost/server/internal/url"
)

func main() {
	env := os.Getenv("ENVIRONMENT")
	err := godotenv.Load()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) && env == "" {
			log.Fatal(err)
		}
	}

	app := fiber.New()
	app.Use(cors.New())
	pmw := auth.NewPasetoMiddleware()
	rmw := url.NewSubdomainRouterMiddleware()
	app.Use(rmw)
	ctx := context.Background()

	connUrl := os.Getenv("POSTGRES_URL")

	conn, err := pgx.Connect(ctx, connUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(ctx)
	slog.Info("Connection established with db", "conn", connUrl)

	queries := db.New(conn)

	ac, err := auth.NewGithubAuthHandler(queries)
	if err != nil {
		log.Fatalf("Unable to init auth controller. %v", err)
	}
	ac.RegisterRoutes(app)

	adc := admin.NewAdminController()
	adc.RegisterRoutes(app, &pmw)

	endpointService := url.NewUrlService(queries)
	urlHandler := url.NewEndpointHandler(endpointService)
	urlHandler.RegisterRoutes(app, &pmw)

	// TODO: Fetch port from config
	app.Listen(":3000")
}
