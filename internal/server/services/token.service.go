package services

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthTokenClaims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uuid.UUID, jwtSecret string) (string, error) {
	claims := AuthTokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		}}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenStr, nil
}

func GetBearToken(r *http.Request) (string, error) {
	tokenHeader := r.Header.Get("Authorization")
	if tokenHeader == "" {
		return "", fmt.Errorf("authentication token missing")
	}

	tokenParts := strings.Split(tokenHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return "", fmt.Errorf("invalid authentication token")
	}

	return tokenParts[1], nil
}

func ParseTokenClaims(tokenStr, jwtSecret string) (AuthTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &AuthTokenClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(jwtSecret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil {
		return AuthTokenClaims{}, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*AuthTokenClaims)
	if !ok {
		return AuthTokenClaims{}, fmt.Errorf("unknown claims type")
	}

	return *claims, nil
}
