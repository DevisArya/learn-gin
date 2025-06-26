package helper

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func GetUserId(c *gin.Context) (uint, error) {

	val, exist := c.Get("claims")
	if !exist {
		return 0, NewCustomError(http.StatusUnauthorized, []string{"unauthorized: claims not found"})
	}

	claims, ok := val.(jwt.MapClaims)
	if !ok {
		return 0, NewCustomError(http.StatusUnauthorized, []string{"invalid token claims"})
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, NewCustomError(http.StatusUnauthorized, []string{"unauthorized: user_id not found in token"})
	}

	return uint(userIDFloat), nil
}
