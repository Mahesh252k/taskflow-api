package utils

import (
	"fmt"
	"time"

	"taskflow-api/config"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID uint) (string, error) {
	jwtSecret := []byte(config.GetEnv("JWT_SECRET"))

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}

func ValidateToken(tokenString string) (uint, error) {
	jwtSecret := []byte(config.GetEnv("JWT_SECRET"))

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method")
		}

		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("invalid token claims")
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid user id")
	}

	return uint(userIDFloat), nil
}
