package middleware

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"onlineShop/config"
	"onlineShop/response"
	"onlineShop/util"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func Trace() fiber.Handler {
	// Create a global logger with configuration
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)            // Set appropriate log level
	logger.SetFormatter(&logrus.JSONFormatter{}) // Use JSON formatting for structured data
	file, err := os.OpenFile("my_app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o660)
	if err != nil {
		// Handle file opening error gracefully
		fmt.Println("Error opening log file:", err)
		os.Exit(1)
	}
	logger.SetOutput(file)

	return func(c *fiber.Ctx) error {
		// Generate and store trace ID
		traceID := uuid.New()
		c.Set("X-Trace-ID", traceID.String())

		// Start request tracking and log with fields
		startTime := time.Now()
		logger.WithFields(logrus.Fields{
			"method":    c.Route().Method,
			"path":      string(c.Context().URI().Path()),
			"trace_id":  traceID,
			"timestamp": startTime,
		}).Info("Start request")

		// Execute the next middleware and capture errors
		err := c.Next()

		// Log request completion and duration
		duration := time.Since(startTime)

		// Handle potential error reading response body
		bodyBytes := c.Body() // Use c.Body() for efficient reading

		// Modify response body (removing double quotes)
		modifiedBody := strings.ReplaceAll(string(bodyBytes), "\"", "")

		// Log request details with modified response body
		logger.WithFields(logrus.Fields{
			"trace_id":      traceID,
			"duration":      duration,
			"error":         err,
			"response_body": modifiedBody,
			"timestamp":     time.Now(),
		}).Info("Processed request")

		return err // Return original error (if any)
	}
}

func CheckAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authorization := c.Get("Authorization")
		if authorization == "" {
			return response.NewResponse(
				response.WithError(response.ErrorUnauthorized),
			).Send(c)
		}

		bearer := strings.Split(authorization, "Bearer ")
		if len(bearer) != 2 {
			log.Println("token invalid")
			return response.NewResponse(
				response.WithError(response.ErrorUnauthorized),
			).Send(c)
		}

		token := bearer[1]

		publicId, role, err := util.ValidateToken(token, config.Cfg.App.Encryption.JWTSecret)
		if err != nil {
			log.Println(err.Error())
			return response.NewResponse(
				response.WithError(response.ErrorUnauthorized),
			).Send(c)
		}

		c.Locals("ROLE", role)
		c.Locals("PUBLIC_ID", publicId)

		return c.Next()
	}
}

func CheckRoles(authorizedRoles []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := fmt.Sprintf("%v", c.Locals("ROLE"))

		isExists := false
		for _, authorizedRole := range authorizedRoles {
			if role == authorizedRole {
				isExists = true
				break
			}
		}

		if !isExists {
			return response.NewResponse(
				response.WithError(response.ErrorForbiddenAccess),
			).Send(c)
		}

		return c.Next()
	}
}
