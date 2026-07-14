package handler

import (
	"context"
	"encoding/json"
	"first-api/internal/model"
	"net/http"
)

type AuthUseCase interface {
	Register(context.Context, *http.Request) (*model.TokenResponseDTO, error)
	Login(context.Context, *http.Request) (*model.TokenResponseDTO, error)
	RefreshAccessToken(context.Context, *http.Request) (*model.TokenResponseDTO, error)
}

type AuthHandler struct {
	AuthUseCase AuthUseCase
}

func NewAuthHandler(authUseCase AuthUseCase) *AuthHandler {
	return &AuthHandler{AuthUseCase: authUseCase}
}

func (ah *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tokenResponse, err := ah.AuthUseCase.Register(ctx, r)
	if err != nil {
		WriteOrderError(w, err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tokenResponse)
}

func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tokenResponse, err := ah.AuthUseCase.Login(ctx, r)
	if err != nil {
		WriteOrderError(w, err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tokenResponse)
}

func (ah *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	refreshResponse, err := ah.AuthUseCase.RefreshAccessToken(ctx, r)
	if err != nil {
		WriteOrderError(w, err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(refreshResponse)
}
