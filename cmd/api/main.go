package main

import (
	"context"
	"errors"
	"log"
	"os"
	"strings"
	"time"

	"github.com/dim2k2006/correlateapp-be/cmd/api/middleware"
	"github.com/dim2k2006/correlateapp-be/cmd/api/schemas"
	"github.com/dim2k2006/correlateapp-be/pkg/domain/parameter"
	"github.com/dim2k2006/correlateapp-be/pkg/domain/user"
	"github.com/go-playground/validator/v10"
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

	parameterRepository := parameter.NewInMemoryRepository()
	parameterService := parameter.NewService(parameterRepository)

	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		now := time.Now()

		return c.SendString("It is alive 🔥🔥🔥. Now: " + now.Format(time.RFC3339))
	})

	api := app.Group("/api", middleware.VerifySignatureMiddleware(secretKeys))

	users := api.Group("/users")

	users.Post("/", func(c *fiber.Ctx) error {
		var req schemas.CreateUserRequest

		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid input: " + err.Error(),
			})
		}

		if err := req.Validate(); err != nil {
			var validationErrors validator.ValidationErrors
			errors.As(err, &validationErrors)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Validation failed",
				"details": validationErrors.Error(),
			})
		}

		input := user.CreateUserInput{
			ExternalID: req.ExternalID,
			FirstName:  req.FirstName,
			LastName:   req.LastName,
		}

		ctx := context.Background()
		createdUser, err := userService.CreateUser(ctx, input)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(schemas.NewUserResponse(createdUser))
	})

	users.Get("/:id", func(c *fiber.Ctx) error {
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

		return c.JSON(schemas.NewUserResponse(userData))
	})

	users.Get("/external/:externalId", func(c *fiber.Ctx) error {
		externalID := c.Params("externalId")
		ctx := context.Background()
		userData, err := userService.GetUserByExternalID(ctx, externalID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(schemas.NewUserResponse(userData))
	})

	users.Put("/:id", func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, uuidParseErr := uuid.Parse(idStr)
		if uuidParseErr != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid user ID",
			})
		}

		var req schemas.UpdateUserRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid input: " + err.Error(),
			})
		}

		if err := req.Validate(); err != nil {
			var validationErrors validator.ValidationErrors
			errors.As(err, &validationErrors)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Validation failed",
				"details": validationErrors.Error(),
			})
		}

		input := user.UpdateUserInput{
			ID:        id,
			FirstName: req.FirstName,
			LastName:  req.LastName,
		}

		ctx := context.Background()
		updatedUser, uuidParseErr := userService.UpdateUser(ctx, input)
		if uuidParseErr != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": uuidParseErr.Error(),
			})
		}

		return c.JSON(schemas.NewUserResponse(updatedUser))
	})

	users.Delete("/:id", func(c *fiber.Ctx) error {
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

	parameters := api.Group("/parameters")

	parameters.Post("/", func(c *fiber.Ctx) error {
		var req schemas.CreateParameterInput

		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid input: " + err.Error(),
			})
		}

		if err := req.Validate(); err != nil {
			var validationErrors validator.ValidationErrors
			errors.As(err, &validationErrors)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Validation failed",
				"details": validationErrors.Error(),
			})
		}

		input := parameter.CreateParameterInput{
			UserID:      req.UserID,
			Name:        req.Name,
			Description: req.Description,
			DataType:    req.DataType,
			Unit:        req.Unit,
		}

		ctx := context.Background()
		createdParameter, err := parameterService.CreateParameter(ctx, input)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(schemas.NewParameterResponse(createdParameter))
	})

	parameters.Get("/:id", func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid parameter ID",
			})
		}

		ctx := context.Background()
		parameterData, err := parameterService.GetParameterByID(ctx, id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(schemas.NewParameterResponse(parameterData))
	})

	// TODO implement routes for measurements

	log.Println("Starting server on " + port)

	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
