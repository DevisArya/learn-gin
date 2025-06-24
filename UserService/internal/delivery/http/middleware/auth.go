package middleware

// import (
// 	"context"
// 	"net/http"
// 	"os"
// 	"strings"

// 	"github.com/DevisArya/BE-challenge-syn/UserService/internal/dto"
// 	"github.com/DevisArya/BE-challenge-syn/UserService/internal/helper"
// 	"github.com/DevisArya/BE-challenge-syn/UserService/internal/usecase"
// 	"github.com/gofiber/fiber/v2"
// 	"github.com/golang-jwt/jwt/v5"
// )

// func NewAuth(userUseCase usecase.UserUseCase) fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		authHeader := c.Get("Authorization")
// 		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
// 			return helper.SendErrorResponse(c, http.StatusUnauthorized, []string{"missing or invalid authorization header"})
// 		}

// 		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

// 		secretKey := os.Getenv("SECRET_KEY")
// 		if secretKey == "" {
// 			return helper.SendErrorResponse(c, http.StatusInternalServerError, []string{"internal server error"})
// 		}

// 		claims := jwt.MapClaims{}
// 		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
// 			return []byte(secretKey), nil
// 		})
// 		if err != nil || !token.Valid {
// 			return helper.SendErrorResponse(c, http.StatusUnauthorized, []string{"invalid token"})
// 		}

// 		//user id from claims
// 		userIDFloat, ok := claims["user_id"].(float64)
// 		if !ok {
// 			return helper.SendErrorResponse(c, http.StatusUnauthorized, []string{"invalid token payload"})
// 		}
// 		userID := uint(userIDFloat)

// 		VerTokenReq := &dto.VerifyTokenRequest{
// 			Token: tokenString,
// 		}

// 		// check token in DB
// 		if err := userUseCase.VerifyToken(context.Background(), VerTokenReq, userID); err != nil {
// 			return helper.SendCustomErrorResponse(c, err)
// 		}

// 		c.Locals("claims", claims)
// 		return c.Next()
// 	}
// }
