package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/DevisArya/BE-challenge-syn/NotificationService/internal/helper"
	"github.com/DevisArya/BE-challenge-syn/NotificationService/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func NewAuth(userUseCase usecase.NotificationUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			helper.SendErrorResponse(c, http.StatusUnauthorized, []string{"missing or invalid authorization header"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		secretKey := os.Getenv("SECRET_KEY")
		if secretKey == "" {
			helper.SendErrorResponse(c, http.StatusInternalServerError, []string{"internal server error"})
			c.Abort()
			return
		}

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})
		if err != nil || !token.Valid {
			helper.SendErrorResponse(c, http.StatusUnauthorized, []string{"invalid token"})
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}
