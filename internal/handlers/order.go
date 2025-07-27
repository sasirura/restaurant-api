// Package handlers
package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sasirura/restaurant-api/internal/models"
	"github.com/sasirura/restaurant-api/internal/services"
	"github.com/square/square-go-sdk/client"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": err.Error(),
	})
}

// CreateOrder creates a new order
func CreateOrder(squareService *services.SquareService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		restaurant := c.Locals("restaurant").(models.Restaurant)
		client := c.Locals("client").(*client.Client)
		var req struct {
			TableNumber string             `json:"tableNumber"`
			Items       []models.OrderItem `json:"items"`
		}

		if err := c.BodyParser(&req); err != nil {
			squareService.Logger.Error("Invalid request body", "error", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}

		order, err := squareService.CreateOrder(c.Context(), restaurant, client, req.TableNumber, req.Items)
		if err != nil {
			squareService.Logger.Error("Failed to create order", "error", err, "restaurant_id", restaurant.ID)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		squareService.Logger.Info("Order created successfully", "order_id", order.ID, "restaurant_id", restaurant.ID)
		return c.JSON(order)
	}
}

// GetOrdersByTable retrieves orders by table number
func GetOrdersByTable(squareService *services.SquareService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		restaurant := c.Locals("restaurant").(models.Restaurant)
		tableNumber := c.Params("tableNumber")

		orders, err := squareService.GetOrdersByTable(c.Context(), restaurant, tableNumber)
		if err != nil {
			squareService.Logger.Error("Failed to fetch orders by table", "error", err, "table_number", tableNumber)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		squareService.Logger.Info("Fetched orders by table", "table_number", tableNumber, "count", len(orders))
		return c.JSON(orders)
	}
}

// GetOrderByID retrieves an order by ID
func GetOrderByID(squareService *services.SquareService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		restaurant := c.Locals("restaurant").(models.Restaurant)
		orderID := c.Params("id")

		order, err := squareService.GetOrderByID(c.Context(), restaurant, orderID)
		if err != nil {
			squareService.Logger.Error("Failed to fetch order by ID", "error", err, "order_id", orderID)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Order not found"})
		}

		squareService.Logger.Info("Fetched order by ID", "order_id", orderID)
		return c.JSON(order)
	}
}

// ProcessPayment processes a payment for an order
func ProcessPayment(squareService *services.SquareService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		restaurant := c.Locals("restaurant").(models.Restaurant)
		client := c.Locals("client").(*client.Client)
		orderID := c.Params("orderId")
		var req models.PaymentRequest

		if err := c.BodyParser(&req); err != nil {
			squareService.Logger.Error("Invalid payment request body", "error", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}

		if err := squareService.ProcessPayment(c.Context(), restaurant, client, orderID, req); err != nil {
			squareService.Logger.Error("Failed to process payment", "error", err, "order_id", orderID)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		squareService.Logger.Info("Payment processed successfully", "order_id", orderID, "payment_id", req.PaymentID)
		return c.JSON(fiber.Map{"status": "Payment processed successfully"})
	}
}
