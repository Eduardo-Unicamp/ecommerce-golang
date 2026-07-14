package model

import (
	"time"

	"github.com/google/uuid"
)

type LoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterDTO struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type TokenResponseDTO struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    int       `json:"expires_in"`
	CustomerID   uuid.UUID `json:"customer_id"`
}

type RefreshTokenDTO struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshToken struct {
	Token      string    `json:"token"`
	CustomerID uuid.UUID `json:"customer_id"`
	ExpiresAt  time.Time `json:"expires_at"`
	Revoked    bool      `json:"revoked"`
}
