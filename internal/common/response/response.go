package response

import (
	"github.com/JustDoItBetter/FITS-backend/internal/common/errors"
	"github.com/gofiber/fiber/v2"
)

// SuccessResponse represents a successful API response
type SuccessResponse struct {
	Success bool        `json:"success" example:"true"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty" example:"operation successful"`
}

// ErrorResponse represents an error API response
type ErrorResponse struct {
	Success bool   `json:"success" example:"false"`
	Error   string `json:"error" example:"invalid request"`
	Details string `json:"details,omitempty" example:"field validation failed"`
	Code    int    `json:"code" example:"400"`
}

// Success sends a successful JSON response
func Success(c *fiber.Ctx, data interface{}) error {
	return c.JSON(SuccessResponse{
		Success: true,
		Data:    data,
	})
}

// SuccessWithMessage sends a successful JSON response with a message
func SuccessWithMessage(c *fiber.Ctx, message string, data interface{}) error {
	return c.JSON(SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Created sends a 201 Created response
func Created(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(SuccessResponse{
		Success: true,
		Data:    data,
	})
}

// NoContent sends a 204 No Content response
func NoContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

// Error sends an error JSON response with proper status code
// Converts AppError to ErrorResponse format and sets appropriate HTTP status
func Error(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	errMsg := "Internal Server Error"
	details := err.Error()

	// Check if it's our custom AppError type
	if appErr, ok := err.(*errors.AppError); ok {
		code = appErr.Code
		errMsg = appErr.Message
		details = appErr.Details
	} else if fiberErr, ok := err.(*fiber.Error); ok {
		// Handle Fiber's built-in errors
		code = fiberErr.Code
		errMsg = fiberErr.Message
		details = ""
	}

	return c.Status(code).JSON(ErrorResponse{
		Success: false,
		Error:   errMsg,
		Details: details,
		Code:    code,
	})
}
