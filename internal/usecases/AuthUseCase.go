package usecases

import (
	"context"
	"encoding/json"
	"errors"
	"first-api/internal/auth"
	"first-api/internal/model" // Importe seu repositório de clientes
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type CustomerAuthRepository interface {
	GetCustomers(context.Context) ([]model.Customer, error)
	CreateCustomer(context.Context, *model.Customer) error
	GetCustomerByField(context.Context, string, string) (*model.Customer, error)

	SaveRefreshToken(context.Context, *model.RefreshToken) error
	GetRefreshToken(context.Context, string) (*model.RefreshToken, error)
	RevokeRefreshToken(context.Context, string) error
}

type AuthUseCase struct {
	CustomerRepo CustomerAuthRepository
	jwtConfig    *auth.JWTConfig
}

func NewAuthUseCase(cr CustomerAuthRepository, config *auth.JWTConfig) *AuthUseCase {
	return &AuthUseCase{
		CustomerRepo: cr,
		jwtConfig:    config,
	}
}

func (au *AuthUseCase) Register(ctx context.Context, r *http.Request) (*model.TokenResponseDTO, error) {
	var request model.CreateCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, fmt.Errorf("Error while parsing the json:%w", err)
	}

	_, err := au.CustomerRepo.GetCustomerByField(ctx, "email", request.Email)
	if err == nil {
		return nil, model.ErrEmailTaken
	} //achou, email repetido
	if !errors.Is(err, model.CustomerNotFound) {
		return nil, err
	} //se deu algum outro erro

	customer, err := model.NewCustomer(request.Name, request.Email, request.Phone, request.Password)
	if err != nil {
		return nil, err
	}

	if err := au.CustomerRepo.CreateCustomer(ctx, customer); err != nil {
		return nil, err
	}

	return au.GenerateTokenResponse(ctx, customer.ID)

}

func (au *AuthUseCase) Login(ctx context.Context, r *http.Request) (*model.TokenResponseDTO, error) {
	var request model.LoginDTO
	json.NewDecoder(r.Body).Decode(&request)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, model.ErrReadingJSON
	}

	customer, err := au.CustomerRepo.GetCustomerByField(ctx, "email", request.Email)
	if err != nil {
		return nil, model.ErrInvalidPassword //sim, é email, mas tem aquela regrinha de nao dizer se é email ou senha que errou
	}

	err = bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(request.Password))
	if err != nil {
		return nil, model.ErrInvalidPassword
	}

	return au.GenerateTokenResponse(ctx, customer.ID)
}

func (au *AuthUseCase) GenerateTokenResponse(ctx context.Context, customerID uuid.UUID) (*model.TokenResponseDTO, error) {
	accessTokenStr, err := auth.GenerateToken(customerID, au.jwtConfig)
	if err != nil {
		return nil, err
	}

	refreshTokenStr, err := auth.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	//pra ser salvo no banco
	refreshToken := model.RefreshToken{
		Token:      refreshTokenStr,
		CustomerID: customerID,
		ExpiresAt:  time.Now().Add(time.Hour * 24 * time.Duration(au.jwtConfig.RefreshExpirationDays)),
	}
	if err := au.CustomerRepo.SaveRefreshToken(ctx, &refreshToken); err != nil {
		return nil, err
	}

	//devolvido pro user
	return &model.TokenResponseDTO{
		AccessToken:  accessTokenStr,
		RefreshToken: refreshTokenStr,
		ExpiresIn:    au.jwtConfig.ExpirationMinutes * 60,
		CustomerID:   customerID,
	}, nil
}

func (au *AuthUseCase) RefreshAccessToken(ctx context.Context, r *http.Request) (*model.TokenResponseDTO, error) {
	var request model.RefreshTokenDTO
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	if request.RefreshToken == "" {
		return nil, model.ErrRefreshTokenRequired
	}

	refreshToken, err := au.CustomerRepo.GetRefreshToken(ctx, request.RefreshToken)
	if err != nil {
		return nil, err
	}

	//checa se o token está revogado ou expirado
	if refreshToken.Revoked == true || time.Now().After(refreshToken.ExpiresAt) {
		return nil, model.ErrInvalidRefreshToken
	}

	//revoga o token antigo e gera um outro a cada vez que o user gera um novo access token(refresh token rotation) por segurança
	if err = au.CustomerRepo.RevokeRefreshToken(ctx, refreshToken.Token); err != nil {
		return nil, err
	}
	return au.GenerateTokenResponse(ctx, refreshToken.CustomerID)

}
