package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/dim2k2006/correlateapp-be/cmd/api/middleware"
	"github.com/dim2k2006/correlateapp-be/pkg/domain/user"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	dotEnvErr := godotenv.Load()
	if dotEnvErr != nil {
		log.Println("Warning: No .env file found, using system environment variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	secretKeysString := os.Getenv("API_SECRET_KEYS")
	if secretKeysString == "" {
		log.Fatal("API_SECRET_KEYS is empty string")
	}

	secretKeys := strings.Split(secretKeysString, ",")
	if len(secretKeys) == 0 {
		log.Fatal("API_SECRET_KEYS is empty")
	}

	userRepository := user.NewInMemoryRepository()
	userService := user.NewService(userRepository)

	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		now := time.Now()

		return c.SendString("It is alive 🔥🔥🔥. Now: " + now.Format(time.RFC3339))
	})

	api := app.Group("/api", middleware.VerifySignatureMiddleware(secretKeys))

	api.Post("/users", func(c *fiber.Ctx) error {
		var input user.CreateUserInput
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid input: " + err.Error(),
			})
		}

		ctx := context.Background()
		createdUser, err := userService.CreateUser(ctx, input)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(createdUser)
	})

	api.Get("/users/:id", func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid user ID",
			})
		}

		ctx := context.Background()
		userData, err := userService.GetUserByID(ctx, id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(userData)
	})

	api.Get("/users/external/:externalId", func(c *fiber.Ctx) error {
		externalID := c.Params("externalId")
		ctx := context.Background()
		userData, err := userService.GetUserByExternalID(ctx, externalID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(userData)
	})

	api.Put("/users/:id", func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, uuidParseErr := uuid.Parse(idStr)
		if uuidParseErr != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid user ID",
			})
		}

		var input user.UpdateUserInput
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid input: " + err.Error(),
			})
		}

		// Optionally, ensure that the URL id matches the body id.
		// In this example we override the input.ID with the id from the URL.
		input.ID = id

		ctx := context.Background()
		updatedUser, uuidParseErr := userService.UpdateUser(ctx, input)
		if uuidParseErr != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": uuidParseErr.Error(),
			})
		}

		return c.JSON(updatedUser)
	})

	api.Delete("/users/:id", func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, uuidParseErr := uuid.Parse(idStr)
		if uuidParseErr != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid user ID",
			})
		}

		ctx := context.Background()
		if err := userService.DeleteUser(ctx, id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Return no content status upon successful deletion.
		return c.SendStatus(fiber.StatusNoContent)
	})

	log.Println("Starting server on " + port)

	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
