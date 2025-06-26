package security

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/DevisArya/BE-challenge-syn/UserService/internal/helper"
	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(userId uint, name string) (string, error) {

	claims := jwt.MapClaims{
		"exp":     jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		"iss":     "BE-challenge-syn",
		"sub":     fmt.Sprintf("%d", userId),
		"iat":     jwt.NewNumericDate(time.Now()),
		"user_id": userId,
		"name":    name,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		return "", helper.NewCustomError(http.StatusUnauthorized, []string{"missing secret key"})
	}

	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
