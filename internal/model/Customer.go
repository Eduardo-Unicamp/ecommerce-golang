package model

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewCustomer(name string, email string, phone string) (*Customer, error) {

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

	customerId, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}
	var customer Customer = Customer{ID: customerId, Name: name, Email: email, Phone: phone}

	return &customer, nil
}

type CreateCustomerRequest struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Phone string    `json:"phone"`
}

type UpdateCustomerRequest struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Phone string    `json:"phone"`
}
