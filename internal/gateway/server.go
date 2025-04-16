package gateway

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/saleh-ghazimoradi/EcoBay/config"
	"github.com/saleh-ghazimoradi/EcoBay/internal/gateway/routes"
	"github.com/saleh-ghazimoradi/EcoBay/slg"
	"github.com/saleh-ghazimoradi/EcoBay/utils"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func Server() error {
	db, err := utils.DBConnection(utils.DBMigrator)
	if err != nil {
		slg.Logger.Error(err.Error())
		return err
	}

	postDB, err := db.DB()
	if err != nil {
		slg.Logger.Error("Unable to get sql.DB", "error", err)
		return err
	}

	defer func() {
		if err := postDB.Close(); err != nil {
			slg.Logger.Error("Error closing database connection", "error", err)
		}
	}()

	app := fiber.New(fiber.Config{
		BodyLimit:    config.AppConfig.Server.BodyLimit,
		WriteTimeout: config.AppConfig.Server.WriteTimeout,
		ReadTimeout:  config.AppConfig.Server.ReadTimeout,
		IdleTimeout:  config.AppConfig.Server.IdleTimeout,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3030",
		AllowHeaders: "Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, PATCH, DELETE, OPTIONS",
	}))

	app.Use(limiter.New(limiter.Config{
		Max:        config.AppConfig.Server.RateLimit,
		Expiration: config.AppConfig.Server.RateLimitExp,
		KeyGenerator: func(ctx *fiber.Ctx) string {
			return ctx.IP()
		},
		LimitReached: func(ctx *fiber.Ctx) error {
			return ctx.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Rate limit exceeded",
			})
		},
	}))

	app.Use(recover.New())
	app.Use(logger.New())

	routes.RegisterRoutes(app, db)

	slg.Logger.Info("Starting server", "port", config.AppConfig.Server.Port)

	var wg sync.WaitGroup
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := app.Listen(config.AppConfig.Server.Port); err != nil {
			return
		}
	}()

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), config.AppConfig.Server.Timeout)
	defer cancel()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := app.ShutdownWithContext(ctx); err != nil {
			return
		}
	}()

	wg.Wait()

	return nil
}
