package main

import (
	fiber "github.com/gofiber/fiber/v2"
)

func main() {

	// Load environment variables
	LoadEnv()

	_aws := NewAWS()

	// Create AWS session
	sess := _aws.ConnectAWS()

	// Initialize appication
	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("sess", sess)
		return c.Next()
	})

	app.Post("/upload", _aws.HandlerFileUpload)

	app.Listen(":9090")
}
