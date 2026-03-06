package utils

import (
	"errors"
	"go-mongo-project/config"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(email string, userId string) (string, error) {
	secretKey := os.Getenv(config.ACCESS_TOKEN_SECRET)

	if secretKey == "" {
		return "", errors.New("ACCESS_TOKEN_SECRET environment variable is not set")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		config.TOKEN_EMAIL:      email,
		config.USER_NAME:        userId,
		config.TOKEN_EXPIRATION: time.Now().Add(time.Hour * 2).Unix(),
	})

	return token.SignedString([]byte(secretKey))
}

func VerifyToken(tokenString string) (jwt.MapClaims, error) {
	secretKey := os.Getenv(config.ACCESS_TOKEN_SECRET)

	if secretKey == "" {
		return nil, errors.New("ACCESS_TOKEN_SECRET environment variable is not set")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenMalformed
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrTokenInvalidClaims
}
