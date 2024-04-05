package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/humanbeeng/checkpost/server/config"
	db "github.com/humanbeeng/checkpost/server/db/sqlc"
	"github.com/humanbeeng/checkpost/server/internal/admin"
	"github.com/humanbeeng/checkpost/server/internal/auth"
	"github.com/humanbeeng/checkpost/server/internal/core"
	"github.com/humanbeeng/checkpost/server/internal/core/middleware"
	"github.com/humanbeeng/checkpost/server/internal/url"
)

func main() {
	config, err := config.GetAppConfig()
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New()
	app.Use(cors.New())
	app.Use(requestid.New())
	// TODO: Revisit this configuration and slog configuration
	app.Use(logger.New(logger.Config{
		// For more options, see the Config section
		TimeFormat: "2006/03/01 15:04:05",
		Format:     "${time} | ${locals:requestid} | ${status} | ${latency} | ${method} ${path}​\n",
	}))

	key := config.Paseto.Key
	pv, err := core.NewPasetoVerifier(key)
	if err != nil {
		slog.Error("Unable to create new paseto verifier", "err", err)
	}
	pmw := middleware.NewPasetoMiddleware(pv)
	tmw := middleware.NewGuestMiddleware(pv)
	rmw := middleware.NewSubdomainRouterMiddleware()
	app.Use(rmw)
	ctx := context.Background()

	connectionString := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", config.Postgres.User, config.Postgres.Password, config.Postgres.Host, config.Postgres.Port, config.Postgres.Database)

	conn, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	slog.Info("Connection established", "host", config.Postgres.Host, "database", config.Postgres.Database, "user", config.Postgres.User)

	queries := db.New(conn)

	ac, err := auth.NewGithubAuthHandler(config, queries)
	if err != nil {
		log.Fatalf("Unable to init auth controller. %v", err)
	}

	adc := admin.NewAdminController()
	endpointService := url.NewUrlService(queries)
	urlHandler := url.NewEndpointHandler(endpointService)

	adc.RegisterRoutes(app, &pmw)
	ac.RegisterRoutes(app)
	urlHandler.RegisterRoutes(app, &tmw)

	// TODO: Fetch port from config
	app.Listen(":3000")
}
