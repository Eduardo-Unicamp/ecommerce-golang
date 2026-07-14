package model

import (
	"errors"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

var ErrNameRequired = errors.New("O nome não pode ser vazio")

var ErrEmailRequired = errors.New("O email não pode ser vazio")
var ErrEmailTaken = errors.New("O email informado já está sendo utilizado.")

var ErrPhoneRequired = errors.New("O telefone não pode ser vazio")

var CustomerNotFound = errors.New("Cliente não encontrado")

var ErrEmptyString = errors.New("String vazia")

var ErrInvalidPrice = errors.New("O preço informado não é válido")

var ErrInvalidStockQuantity = errors.New("Quantidade informada não é válida")

var ErrEmptyOrder = errors.New("O pedido precisa conter ao menos um item")

var ErrInvalidOrderStatus = errors.New("Status de pedido inválido.")

var ErrProductNotFound = errors.New("Produto não encontrado")

var ErrOrderNotFound = errors.New("Pedido não encontrado")

var ErrInsufficientStock = errors.New("Estoque insuficiente")

var ErrInvalidField = errors.New("Value inserted for field is not accepted")

var ErrUnableToCancel = errors.New("Only PENDING orders can be canceled!")

var ErrUnableToPay = errors.New("Only PENDING orders can be canceled!")

var ErrInvalidPassword = errors.New("Invalid password(must be between 6 and 20 characters)")

var ErrReadingJSON = errors.New("Error while reading the json")

var ErrInvalidRefreshToken = errors.New("Invalid refresh token")
var ErrRefreshTokenRequired = errors.New("Refresh token required")
