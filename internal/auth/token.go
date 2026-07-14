package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrMissingSecret = errors.New("JWT_SECRET environment variable is not set")
	ErrInvalidToken  = errors.New("Invalid or expired token")
)

type JWTConfig struct {
	Secret                []byte
	ExpirationMinutes     int
	RefreshExpirationDays int
}

func LoadJWTConfig() (*JWTConfig, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, ErrMissingSecret
	}

	expirationMinutesStr := os.Getenv("JWT_EXPIRATION_MINUTES")
	if expirationMinutesStr == "" { //fallback
		expirationMinutesStr = "60"
	}
	expirationMinutes, err := strconv.Atoi(expirationMinutesStr)
	if err != nil {
		return nil, errors.New("Invalid format for token expiration environment variable")
	}
	refreshExpiration, err := strconv.Atoi(os.Getenv("REFRESH_EXPIRATION_DAYS"))
	if err != nil {
		return nil, errors.New("Invalid format for refresher token expiration value")
	}

	return &JWTConfig{
		Secret:                []byte(secret),
		ExpirationMinutes:     expirationMinutes,
		RefreshExpirationDays: refreshExpiration,
	}, nil
}

func GenerateToken(userID uuid.UUID, config *JWTConfig) (string, error) {
	now := time.Now()
	expiration := now.Add(time.Minute * time.Duration(config.ExpirationMinutes))

	claims := jwt.RegisteredClaims{
		Subject:   userID.String(),
		ExpiresAt: jwt.NewNumericDate(expiration),
		IssuedAt:  jwt.NewNumericDate(now),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.Secret)
}

func ValidateToken(tokenStr string, config *JWTConfig) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return config.Secret, nil
	})
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
