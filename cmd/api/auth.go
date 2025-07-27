package main

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/square/square-go-sdk"
	"github.com/square/square-go-sdk/client"
	"github.com/square/square-go-sdk/option"

	"github.com/sasirura/restaurant-api/internal/models"
	"gorm.io/gorm"
)

func Authenticate(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if token == "" {
			log.Error("Authorization token missing")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authorization token required"})
		}

		client := client.NewClient(
			option.WithToken(token),
			option.WithBaseURL(
				square.Environments.Sandbox,
			),
		)

		tokenStatus, err := client.OAuth.RetrieveTokenStatus(context.TODO())
		if err != nil {
			log.Error("Failed to retrieve token status", "error", err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve token status"})
		}

		// Get location ID
		locations, err := client.Locations.List(context.TODO())
		if err != nil || len(locations.Locations) == 0 {
			log.Error("Failed to fetch locations", "error", err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch Square locations"})
		}
		locationID := locations.Locations[0].ID

		// Check if restaurant exists in database
		var restaurant models.Restaurant
		err = db.Where(models.Restaurant{
			SqaureToken: token,
		}).First(&restaurant).Error
		if err == gorm.ErrRecordNotFound {
			// Create new restaurant record
			restaurant = models.Restaurant{
				Name:        fmt.Sprintf("Restaurant-%s", *tokenStatus.MerchantID),
				SqaureToken: token,
				LocationID:  *locationID,
			}
			if err := db.Create(&restaurant).Error; err != nil {
				log.Error("Failed to create restaurant", "error", err.Error(), "merchant_id", *tokenStatus.MerchantID)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create restaurant"})
			}
			log.Info("Created new restaurant", "restaurant_id", fmt.Sprintf("%d", restaurant.ID), "merchant_id", *tokenStatus.MerchantID)
		} else if err != nil {
			log.Error("Database error", "error", err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
		}

		log.Info("Restaurant authenticated", "restaurant_id", fmt.Sprintf("%d", restaurant.ID), "merchant_id", *tokenStatus.MerchantID)
		c.Locals("restaurant", restaurant)
		c.Locals("client", client)
		return c.Next()
	}
}
