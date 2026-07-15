package handler

import (
	"context"
	"encoding/json"
	"first-api/internal/model"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
)

type AuthUseCase interface {
	Register(context.Context, *http.Request) (*model.TokenResponseDTO, error)
	Login(context.Context, *http.Request) (*model.TokenResponseDTO, error)
	RefreshAccessToken(context.Context, *http.Request) (*model.TokenResponseDTO, error)
	SocialLogin(context.Context, goth.User) (*model.TokenResponseDTO, error)
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

//OAUTH2

func (au *AuthHandler) BeginOAuth(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(r.Context(), "provider", provider))

	//redireciona para o github
	gothic.BeginAuthHandler(w, r)

}

func (ah *AuthHandler) CallbackOAuth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(ctx, "provider", provider))

	gothUser, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		WriteOrderError(w, model.ErrAuthorizationFailed)
		return
	}

	tokenResponse, err := ah.AuthUseCase.SocialLogin(ctx, gothUser)
	if err != nil {
		WriteOrderError(w, err)
		return
	}
	baseURL := os.Getenv("OAUTH_SUCCESS_REDIRECT_TO") //url de onde é pra redirecionar pós login
	redirectURL := fmt.Sprintf("%s?access_token=%s&refresh_token=%s", baseURL, tokenResponse.AccessToken, tokenResponse.RefreshToken)

	http.Redirect(w, r, redirectURL, http.StatusFound)
}
