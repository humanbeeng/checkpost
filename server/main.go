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
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robfig/cron/v3"

	"github.com/humanbeeng/checkpost/server/config"
	db "github.com/humanbeeng/checkpost/server/db/sqlc"
	"github.com/humanbeeng/checkpost/server/internal/auth"
	"github.com/humanbeeng/checkpost/server/internal/core"
	"github.com/humanbeeng/checkpost/server/internal/core/jobs"
	"github.com/humanbeeng/checkpost/server/internal/core/middleware"
	"github.com/humanbeeng/checkpost/server/internal/endpoint"
	"github.com/humanbeeng/checkpost/server/internal/user"
)

// TODO: Implement graceful shutdown
func main() {
	config, err := config.GetAppConfig()
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	app.Use(requestid.New())
	// TODO: Revisit this configuration and slog configuration
	app.Use(logger.New(logger.Config{
		// For more options, see the Config section
		TimeFormat: "2006/03/01 15:04:05",
		Format:     "${time} | ${latency} | ${locals:requestid} | ${status} | ${method} ${path}\n",
	}))

	app.Use(cors.New())
	key := config.Paseto.Key

	pasetoVerifier, err := core.NewPasetoVerifier(key)
	if err != nil {
		slog.Error("unable to create new paseto verifier", "err", err)
	}

	authmw := middleware.NewAuthRequiredMiddleware(pasetoVerifier)
	routermw := middleware.NewSubdomainRouterMiddleware()

	app.Use(routermw)

	ctx := context.Background()

	connectionString := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", config.Postgres.User, config.Postgres.Password, config.Postgres.Host, config.Postgres.Port, config.Postgres.Database)

	conn, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	slog.Info("Connection established", "host", config.Postgres.Host, "database", config.Postgres.Database, "user", config.Postgres.User)

	queries := db.New(conn)

	ac, err := auth.NewAuthHandler(config, queries)
	if err != nil {
		log.Fatalf("unable to init auth controller. %v", err)
	}

	endpointStore := endpoint.NewEndpointStore(queries)
	userStore := user.NewUserStore(queries)
	endpointService := endpoint.NewEndpointService(endpointStore, userStore)
	wsManager := endpoint.NewWSManager()
	endpointHandler := endpoint.NewEndpointController(endpointService, wsManager, pasetoVerifier)

	cachemw := middleware.NewCacheMiddleware()

	userc := user.NewUserController(userStore)
	userc.RegisterRoutes(app, authmw)

	ac.RegisterRoutes(app)
	endpointHandler.RegisterRoutes(app, authmw, cachemw)

	re := jobs.NewExpiredRequestsRemover(cron.New(), *endpointStore)
	re.Start()

	err = app.Listen(":3000")
	if err != nil {
		slog.Error("unable to start fiber server", "err", err)
		return
	}

}
