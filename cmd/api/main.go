package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/dim2k2006/correlateapp-be/cmd/api/middleware"
	"github.com/dim2k2006/correlateapp-be/cmd/api/schemas"
	"github.com/dim2k2006/correlateapp-be/pkg/domain/measurement"
	"github.com/dim2k2006/correlateapp-be/pkg/domain/parameter"
	"github.com/dim2k2006/correlateapp-be/pkg/domain/user"
	"github.com/getsentry/sentry-go"
	sentryfiber "github.com/getsentry/sentry-go/fiber"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

//nolint:gocyclo // It is fine to keep all routes in one place for now
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

	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		log.Fatal("APP_ENV is empty")
	}

	sentryDsn := os.Getenv("SENTRY_DSN")
	if sentryDsn == "" {
		log.Fatal("SENTRY_DSN is empty")
	}

	cosmosDBConnectionString := os.Getenv("COSMOS_DB_CONNECTION_STRING")
	if cosmosDBConnectionString == "" {
		log.Fatal("COSMOS_DB_CONNECTION_STRING is empty")
	}

	isProduction := appEnv == "production"

	userRepository, userRepositoryErr := user.NewCosmosUserRepository(cosmosDBConnectionString)
	if userRepositoryErr != nil {
		log.Fatalf("failed to create user repository: %v", userRepositoryErr)
	}
	userService := user.NewService(userRepository)

	parameterRepository, parameterRepositoryErr := parameter.NewCosmosParameterRepository(cosmosDBConnectionString)
	if parameterRepositoryErr != nil {
		log.Fatalf("failed to create parameter repository: %v", parameterRepositoryErr)
	}
	parameterService := parameter.NewService(parameterRepository)

	measurementRepository, measurementRepositoryErr := measurement.NewCosmosMeasurementRepository(cosmosDBConnectionString)
	if measurementRepositoryErr != nil {
		log.Fatalf("failed to create measurement repository: %v", measurementRepositoryErr)
	}
	measurementService := measurement.NewService(measurementRepository, parameterService)

	if isProduction {
		if err := sentry.Init(sentry.ClientOptions{
			Dsn:              sentryDsn,
			TracesSampleRate: 1.0,
		}); err != nil {
			log.Printf("Sentry initialization failed: %v\n", err)
		}
	}

	sentryHandler := sentryfiber.New(sentryfiber.Options{
		Repanic:         true,
		WaitForDelivery: true,
	})

	app := fiber.New()

	app.Use(sentryHandler)

	app.Use(cors.New())

	app.All("/error", func(_ *fiber.Ctx) error {
		panic("y tho")
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		now := time.Now()

		return c.SendString("It is alive 🔥🔥🔥. Now: " + now.Format(time.RFC3339))
	})

	api := app.Group("/api")

	if isProduction {
		api.Use(middleware.VerifySignatureMiddleware(secretKeys))
	}

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
		updatedUser, updateUserErr := userService.UpdateUser(ctx, input)
		if updateUserErr != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": updateUserErr.Error(),
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
		var req schemas.CreateParameterRequest

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

	parameters.Get("/user/:userId", func(c *fiber.Ctx) error {
		userIDStr := c.Params("userId")
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid user ID",
			})
		}

		ctx := context.Background()
		parametersData, err := parameterService.ListParametersByUser(ctx, userID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// initialize response with empty array
		response := []schemas.ParameterResponse{}
		for _, p := range parametersData {
			response = append(response, schemas.NewParameterResponse(p))
		}

		return c.JSON(response)
	})

	parameters.Put("/:id", func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, uuidParseErr := uuid.Parse(idStr)
		if uuidParseErr != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid parameter ID",
			})
		}

		var req schemas.UpdateParameterRequest
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

		input := parameter.UpdateParameterInput{
			ID:          id,
			Name:        req.Name,
			Description: req.Description,
			Unit:        req.Unit,
		}

		ctx := context.Background()
		updatedParameter, updateParameterErr := parameterService.UpdateParameter(ctx, input)
		if updateParameterErr != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": updateParameterErr.Error(),
			})
		}

		return c.JSON(schemas.NewParameterResponse(updatedParameter))
	})

	parameters.Delete("/:id", func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, uuidParseErr := uuid.Parse(idStr)
		if uuidParseErr != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid parameter ID",
			})
		}

		ctx := context.Background()
		if err := parameterService.DeleteParameter(ctx, id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Return no content status upon successful deletion.
		return c.SendStatus(fiber.StatusNoContent)
	})

	measurements := api.Group("/measurements")

	measurements.Post("/", func(c *fiber.Ctx) error {
		var req schemas.CreateMeasurementRequest

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

		input := measurement.CreateMeasurementInput{
			ParameterID: req.ParameterID,
			Notes:       req.Notes,
			Value:       req.Value,
			Timestamp:   req.Timestamp,
		}

		ctx := context.Background()
		createdMeasurement, err := measurementService.CreateMeasurement(ctx, input)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(schemas.NewMeasurementResponse(createdMeasurement))
	})

	measurements.Get("/user/:userId", func(c *fiber.Ctx) error {
		userIDStr := c.Params("userId")
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid user ID",
			})
		}

		ctx := context.Background()
		measurementsData, err := measurementService.ListMeasurementsByUser(ctx, userID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		response := []schemas.MeasurementResponse{}
		for _, measurementItem := range measurementsData {
			response = append(response, schemas.NewMeasurementResponse(measurementItem))
		}

		return c.JSON(response)
	})

	measurements.Get("/parameter/:parameterId", func(c *fiber.Ctx) error {
		parameterIDStr := c.Params("parameterId")
		parameterID, err := uuid.Parse(parameterIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid parameter ID",
			})
		}

		ctx := context.Background()
		measurementsData, err := measurementService.ListMeasurementsByParameter(ctx, parameterID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		response := []schemas.MeasurementResponse{}
		for _, measurementItem := range measurementsData {
			response = append(response, schemas.NewMeasurementResponse(measurementItem))
		}

		return c.JSON(response)
	})

	measurements.Delete("/:id", func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, uuidParseErr := uuid.Parse(idStr)
		if uuidParseErr != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid measurement ID",
			})
		}

		ctx := context.Background()
		if err := measurementService.DeleteMeasurement(ctx, id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Return no content status upon successful deletion.
		return c.SendStatus(fiber.StatusNoContent)
	})

	// -------------------------
	// Start the server in a goroutine
	// -------------------------
	log.Println("Starting server on port", port)
	go func() {
		if err := app.Listen(":" + port); err != nil {
			// If Listen fails, we log and send a signal, or exit
			log.Fatalf("could not start server: %v", err)
		}
	}()

	// -------------------------
	// Listen for kill signals (graceful shutdown)
	// -------------------------
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit // Block until we get a signal

	log.Println("Gracefully shutting down server...")
	if err := app.Shutdown(); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exiting")
}
