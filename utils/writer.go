package utils

import (
	"github.com/gofiber/fiber/v2"
)

// APIResponse is the standard response wrapper for the API.
type APIResponse struct {
	// Indicates response status
	Status bool `json:"status"`

	// Message describing the response
	Message string `json:"message"`

	// Optional data returned by the endpoint
	Data interface{} `json:"data"`
}

// WriteResponse is a helper for sending consistent JSON responses
func WriteResponse(c *fiber.Ctx, statusCode int, status bool, message string, data interface{}) error { //nolint:typecheck
	return c.Status(statusCode).JSON(APIResponse{
		Status:  status,
		Message: message,
		Data:    data,
	})
}
