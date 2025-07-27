package main

import "os"

// Run starts the application
func (a *App) Serve() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	a.logger.Info("Starting server", "port", port)
	return a.fiber.Listen(":" + port)
}
