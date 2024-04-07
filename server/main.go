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
		Format:     "${time} | ${latency} | ${locals:requestid} | ${status} | ${method} ${path}​\n",
	}))

	key := config.Paseto.Key
	pv, err := core.NewPasetoVerifier(key)
	if err != nil {
		slog.Error("Unable to create new paseto verifier", "err", err)
	}

	payloadmw := middleware.NewExtractPayloadMiddleware(pv)
	genRandLim := middleware.NewGenerateRandomUrlLimiter()
	genLim := middleware.NewGenerateUrlLimiter()

	pmw := middleware.NewAuthRequiredMiddleware(pv)
	gl := middleware.NewGuestPlanLimiter()
	fl := middleware.NewFreePlanLimiter()
	nbl := middleware.NewNoBrainerPlanLimiter()
	pl := middleware.NewProPlanLimiter()
	rmw := middleware.NewSubdomainRouterMiddleware()

	app.Use(payloadmw)
	// app.Use(fl)
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

	authmw := middleware.NewAuthRequiredMiddleware(pv)
	adc := admin.NewAdminController()
	endpointService := url.NewUrlService(queries, config)
	urlHandler := url.NewUrlController(endpointService)

	adc.RegisterRoutes(app, &pmw)
	ac.RegisterRoutes(app)
	urlHandler.RegisterRoutes(app, authmw, gl, fl, nbl, pl, genLim, genRandLim)

	app.Listen(":3000")
}
