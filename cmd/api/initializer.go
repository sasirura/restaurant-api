package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	middlewareLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/sasirura/restaurant-api/internal/handlers"
	"github.com/sasirura/restaurant-api/internal/logger"
	"github.com/sasirura/restaurant-api/internal/models"
	"github.com/sasirura/restaurant-api/internal/services"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Init initializes the application
func Initialize() (*App, error) {
	// Initialize logger
	log := logger.New(log.LevelInfo, os.Stdout)

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Error("Failed to load .env file", "error", err)
		return nil, err
	}

	// Initialize database
	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		log.Error("Failed to connect to database", "error", err)
		return nil, err
	}

	// Auto-migrate models
	if err := db.AutoMigrate(&models.Restaurant{}, &models.Order{}, &models.OrderItem{},
		&models.Discount{}, &models.Modifier{}, &models.OrderTotals{}, models.PaymentRequest{}); err != nil {
		log.Error("Failed to migrate database", "error", err)
		return nil, err
	}

	// Initialize Fiber
	app := fiber.New(fiber.Config{
		ErrorHandler: handlers.ErrorHandler,
	})

	// Initialize Middleware

	// Cors
	app.Use(cors.New())
	// Logger
	ml := middlewareLogger.New(middlewareLogger.Config{
		Format: "${time} ${method} ${path} - ${status} ${latency}\n",
	})
	app.Use(ml)
	// Rate limit
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 60,
	}))

	// Initialize services
	squareService := services.New(db, log)

	return &App{
		fiber:         app,
		db:            db,
		squareService: squareService,
		logger:        log,
	}, nil
}
