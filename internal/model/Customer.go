package model

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Customer struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewCustomer(name string, email string, phone string, password string) (*Customer, error) {

	name, email, phone = strings.TrimSpace(name), strings.TrimSpace(email), strings.TrimSpace(phone)

	if name == "" {
		return &Customer{}, ErrNameRequired
	}
	if email == "" {
		return &Customer{}, ErrEmailRequired
	}
	if phone == "" {
		return &Customer{}, ErrPhoneRequired
	}
	if err := validatePassword(password); err != nil {
		return nil, err
	}

	customerId, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 14) //"first we write code that works"Uncle Bob's son
	if err != nil {
		return nil, err
	}

	var customer Customer = Customer{ID: customerId, Name: name, Email: email, Phone: phone, Password: string(passwordHash)}

	return &customer, nil
}

func NewCustomerThroughSocial(name string, email string) (*Customer, error) {
	name, email = strings.TrimSpace(name), strings.TrimSpace(email)
	if name == "" {
		return &Customer{}, ErrNameRequired
	}
	if email == "" {
		return &Customer{}, ErrEmailRequired
	}

	customerID, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	//gera uma password pro banco que nao será usada pelo user
	randomPassword, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(randomPassword.String()), 14)
	if err != nil {
		return nil, err
	}

	return &Customer{
		ID:       customerID,
		Name:     name,
		Email:    email,
		Phone:    "",
		Password: string(passwordHash),
	}, nil

}

type CreateCustomerRequest struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Phone    string    `json:"phone"`
	Password string    `json:"password"`
}

type UpdateCustomerRequest struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Phone    string    `json:"phone"`
	Password string    `json:"password"`
}

func validatePassword(password string) error {
	if len(password) < 6 || len(password) > 20 {
		return ErrInvalidPassword
	}

	return nil
}
