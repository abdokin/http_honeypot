package main

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func CustomBasicAuth(p_username string, p_password string) fiber.Handler {

	// Return the middleware function
	return func(c *fiber.Ctx) error {
		// Get the Authorization header from the request
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			// No Authorization header found, return unauthorized
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}

		// Check if the Authorization header starts with "Basic "
		if !strings.HasPrefix(authHeader, "Basic ") {
			// Invalid Authorization header, return unauthorized
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}

		// Extract the base64-encoded username:password string
		encodedCreds := strings.TrimPrefix(authHeader, "Basic ")
		decodedCreds, err := base64.StdEncoding.DecodeString(encodedCreds)
		if err != nil {
			// Error decoding credentials, return unauthorized
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}

		// Split the decoded credentials into username and password
		creds := strings.SplitN(string(decodedCreds), ":", 2)
		if len(creds) != 2 {
			// Invalid credentials format, return unauthorized
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}

		username := creds[0]
		password := creds[1]
		status := false
		if username == p_username && p_password == password {
			status = true
		}
		logInfo(c, username, password, status)
		if !status {
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}
		return c.Next()
	}
}
func logInfo(c *fiber.Ctx, username string, password string, status bool) {
	fmt.Printf("Username: %s, Password: %s, Status: %t, Remote IP: %s, Time: %s, User Agent: %s\n",
		username, password, status, c.Context().RemoteIP(), c.Context().Time(), c.Context().UserAgent())
}

func main() {
	// Create a new Fiber instance
	app := fiber.New()
	app.Use(CustomBasicAuth("admin", "admin"))

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
