package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		now := time.Now()

		return c.SendString("It is alive 🔥🔥🔥. Now: " + now.Format(time.RFC3339))
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	log.Println("Starting server on " + port)

	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
