package utils

import (
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func VerifyToken(token string) (string, string, error) {
	secretKey := os.Getenv("SECRET_KEY")
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method.")
		}

		return []byte(secretKey), nil
	})

	if err != nil {
		if err.Error() == "token has invalid claims: token is expired" {
			return "", "", errors.New("Token is expired.")
		}
		return "", "", errors.New("Could not parse token.")
	}

	tokenIsValid := parsedToken.Valid
	if !tokenIsValid {
		return "", "", errors.New("Invalid token.")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", errors.New("Invalid token claims.")
	}

	email := claims["email"].(string)
	userID := claims["userID"].(string)
	return userID, email, nil
}
