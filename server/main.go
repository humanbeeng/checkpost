package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	db "github.com/humanbeeng/checkpost/server/db/sqlc"
	"github.com/humanbeeng/checkpost/server/internal/admin"
	"github.com/humanbeeng/checkpost/server/internal/auth"
	"github.com/humanbeeng/checkpost/server/internal/core"
	"github.com/humanbeeng/checkpost/server/internal/core/middleware"
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

	key := os.Getenv("PASETO_KEY")
	pv, err := core.NewPasetoVerifier(key)
	if err != nil {
		slog.Error("Unable to create new paseto verifier", "err", err)
	}
	pmw := middleware.NewPasetoMiddleware(pv)
	tmw := middleware.NewGuestMiddleware(pv)
	rmw := middleware.NewSubdomainRouterMiddleware()
	app.Use(rmw)
	ctx := context.Background()

	connectionString := os.Getenv("POSTGRES_URL")

	conn, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	slog.Info("Connection established", "conn", connectionString)

	runDBMigration("file://db/migration", connectionString)

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
	urlHandler.RegisterRoutes(app, &tmw)

	// TODO: Fetch port from config
	app.Listen(":3000")
}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		slog.Error("Unable to create new migrate instance", "err", err)
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		slog.Error("Unable to run migrate up", "err", err)
	}

	slog.Info("DB migrated successfully")
}
