package cmd

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/studio-senkou/lentera-quiz-be/config"
	"github.com/studio-senkou/lentera-quiz-be/utils"
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
	app := fiber.New(*config.NewFiberConfig())

	app.Use(config.NewLoggerConfig())
	app.Use(config.NewCORSConfig())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to Lentera Quiz API")
	})

	if err := app.Listen(fmt.Sprintf(":%s", utils.GetEnv("APP_PORT", "9000"))); err != nil {
		return errors.New("failed to start server: " + err.Error())
	}

	return nil
}
