package helper

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func GetUserId(c *fiber.Ctx) (uint, error) {
	claims, ok := c.Locals("claims").(jwt.MapClaims)
	if !ok {
		return 0, NewCustomError(http.StatusUnauthorized, []string{"unauthorized: claims not found"})
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, NewCustomError(http.StatusUnauthorized, []string{"unauthorized: user_id not found in token"})
	}

	return uint(userIDFloat), nil
}
