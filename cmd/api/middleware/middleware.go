package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func VerifySignatureMiddleware(secretKeys []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get headers
		timestampHeader := c.Get("X-Timestamp")
		signatureHeader := c.Get("X-Signature")

		if timestampHeader == "" || signatureHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authentication headers",
			})
		}

		// Convert timestamp to int64
		timestamp, err := strconv.ParseInt(timestampHeader, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid timestamp",
			})
		}

		// Prevent replay attacks (allow max 5 minutes)
		if time.Since(time.Unix(timestamp, 0)) > 5*time.Minute {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Request expired",
			})
		}

		// Get the request body
		payload := string(c.Body())
		if c.Method() == fiber.MethodGet || c.Method() == fiber.MethodDelete {
			if len(payload) == 0 { // Only override if there's no body
				payload = "" // Ensure an empty payload for GET & DELETE without a body
			}
		}

		// Check if the signature matches any of the provided API keys
		for _, key := range secretKeys {
			expectedSignature := computeHMAC(key, payload, timestamp)
			if signatureHeader == expectedSignature {
				// Valid request, continue processing
				return c.Next()
			}
		}

		// If none of the keys match, reject the request
		log.Printf("Invalid signature: %s", signatureHeader)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid API key or signature",
		})
	}
}

func computeHMAC(secret, payload string, timestamp int64) string {
	message := fmt.Sprintf("%s|%d", payload, timestamp)
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
