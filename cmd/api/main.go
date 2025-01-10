package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	log.Println("Starting server on :8080")
	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
