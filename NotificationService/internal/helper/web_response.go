package helper

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// body response success without data
func BodySuccessResponse(message string) map[string]interface{} {
	return map[string]interface{}{
		"message": message,
	}
}

// body response success with data
func BodySuccessResponseWithData(message string, data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"message": message,
		"data":    data,
	}
}

// Body response error
func BodyErrorResponse(errors []string) map[string]interface{} {
	return map[string]interface{}{
		"errors": errors,
	}
}

// send success response
func SendSuccessResponse(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).
		JSON(BodySuccessResponse(message))
}

// send success response with data
func SendSuccessResponseWithData(c *fiber.Ctx, statusCode int, message string, data interface{}) error {
	return c.Status(statusCode).
		JSON(BodySuccessResponseWithData(message, data))
}

// send error response
func SendErrorResponse(c *fiber.Ctx, statusCode int, errors []string) error {
	return c.Status(statusCode).
		JSON(BodyErrorResponse(errors))
}

// send custom error response
func SendCustomErrorResponse(c *fiber.Ctx, err error) error {

	// Cek error validasi
	if customErr, ok := err.(*CustomErrors); ok {
		return SendErrorResponse(c, customErr.StatusCode, customErr.Messages)
	}

	// Cek error timeout
	if err == context.DeadlineExceeded {
		return SendErrorResponse(c, http.StatusRequestTimeout, []string{"request timed out"})
	}

	//cek error not found
	if err == gorm.ErrRecordNotFound {
		return SendErrorResponse(c, http.StatusNotFound, []string{"record not found"})
	}

	// Untuk error lainnya
	return SendErrorResponse(c, http.StatusInternalServerError, []string{"internal server error"})
}
