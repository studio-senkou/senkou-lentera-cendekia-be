package cmd

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/studio-senkou/lentera-cendekia-be/app/routes"
	"github.com/studio-senkou/lentera-cendekia-be/config"
	"github.com/studio-senkou/lentera-cendekia-be/database"
	"github.com/studio-senkou/lentera-cendekia-be/utils/app"
	"github.com/studio-senkou/lentera-cendekia-be/utils/cache"
)

type Application interface {
	Run() error
}

type application struct {
}

func NewApplication() Application {
	return &application{}
}

func (a *application) Run() error {
	if err := database.InitializeDatabase(); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer database.CloseDatabase()

	if err := cache.InitRedis(); err != nil {
		return fmt.Errorf("failed to initialize Redis: %w", err)
	}
	defer cache.CloseRedis()

	fiberApp := fiber.New(*config.NewFiberConfig())

	fiberApp.Use(config.NewLoggerConfig())
	fiberApp.Use(config.NewCORSConfig())

	router := fiberApp.Group("/api/v1")
	routes.SetupUserRoutes(router)
	routes.SetupAuthRoutes(router)
	routes.SetupMeetingSessionRoutes(router)
	routes.SetupTestimonyRoutes(router)
	routes.SetupStaticAssetRoutes(router)
	routes.SetupBlogRoutes(router)

	fiberApp.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to Lentera Cendekia API")
	})

	if err := fiberApp.Listen(fmt.Sprintf(":%s", app.GetEnv("APP_PORT", "9000"))); err != nil {
		return errors.New("failed to start server: " + err.Error())
	}

	return nil
}
