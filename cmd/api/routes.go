package main

import (
	"github.com/sasirura/restaurant-api/internal/handlers"
)

func (a *App) Routes() {

	// API routes
	v1 := a.fiber.Group("/v1")
	{
		// Authenticated routes
		auth := v1.Group("/", Authenticate(a.db))
		auth.Post("/orders", handlers.CreateOrder(a.squareService))
		auth.Get("/orders/:id", handlers.GetOrderByID(a.squareService))
		auth.Get("/orders/table/:tableNumber", handlers.GetOrdersByTable(a.squareService))
		auth.Post("/orders/:orderId/pay", handlers.ProcessPayment(a.squareService))
	}
}
