package model

import (
	"errors"
)

var ErrNameRequired = errors.New("O nome não pode ser vazio")

var ErrEmailRequired = errors.New("O email não pode ser vazio")

var ErrPhoneRequired = errors.New("O telefone não pode ser vazio")

var CustomerNotFound = errors.New("Cliente não encontrado")
