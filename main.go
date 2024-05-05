package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

const HOST = "http://127.0.0.1:8000/api/http"

var (
	HOST_PASSWORD = "admin"
	HOST_USERNAME = "admin"
)

func CustomBasicAuth() fiber.Handler {

	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}

		if !strings.HasPrefix(authHeader, "Basic ") {
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}

		encodedCreds := strings.TrimPrefix(authHeader, "Basic ")
		decodedCreds, err := base64.StdEncoding.DecodeString(encodedCreds)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}

		creds := strings.SplitN(string(decodedCreds), ":", 2)
		if len(creds) != 2 {
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}

		username := creds[0]
		password := creds[1]

		logInfo(c, username, password)

		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")

	}
}
func logInfo(c *fiber.Ctx, username string, password string) {
	pwned := HOST_PASSWORD == password && HOST_USERNAME == username
	data := map[string]interface{}{
		"remoteAddr":     c.Context().RemoteIP().String(),
		"username":       username,
		"password":       password,
		"client_version": string(c.Context().UserAgent()),
		"pwned":          pwned,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	req, err := http.NewRequest("POST", HOST, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Username: %s, Password: %s, Status: %t, Remote IP: %s, Time: %s, User Agent: %s\n",
		username, password, pwned, c.Context().RemoteIP(), c.Context().Time(), c.Context().UserAgent())
}

func main() {
	// Create a new Fiber instance
	app := fiber.New()
	app.Use(CustomBasicAuth())

	// Define a handler function for the root route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome !!")
	})

	// Start the Fiber server on port 8080
	fmt.Println("Server is listening on http://localhost:9000")
	if err := app.Listen(":9000"); err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
