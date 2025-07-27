package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/sasirura/restaurant-api/internal/logger"
	"github.com/sasirura/restaurant-api/internal/services"
	"gorm.io/gorm"
)

// App holds application dependencies
type App struct {
	fiber         *fiber.App
	db            *gorm.DB
	squareService *services.SquareService
	logger        *logger.Logger
}

func main() {
	app, err := Initialize()
	if err != nil {
		log.Fatal("Failed to initialize app", "error", err)
		return
	}

	app.Routes()

	if err := app.Serve(); err != nil {
		app.logger.Fatal("Failed to run app", "error", err)
	}
}
