package main

import (
	"log"

	handler "github.com/SubhamMurarka/Schotky/Handlers"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// Initialize the Fiber app
	app := fiber.New()

	// Route to handle dashboard import
	app.Get("/dashboard/:url", handler.ImportDashboardHandler)

	// Start the Fiber app on port 3000
	log.Fatal(app.Listen(":7000"))
}
