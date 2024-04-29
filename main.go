package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
)

func main() {
	// Create a new Fiber instance
	app := fiber.New()

	// Define basic authentication credentials
	credentials := basicauth.Config{
		Users: map[string]string{
			"admin":  "admin",
		},
	}

	// Apply basic authentication middleware to the entire app
	app.Use(basicauth.New(credentials))

	// Define a handler function for the root route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome !!")
	})

	// Start the Fiber server on port 8080
	fmt.Println("Server is listening on http://localhost:8080")
	if err := app.Listen(":8080"); err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
