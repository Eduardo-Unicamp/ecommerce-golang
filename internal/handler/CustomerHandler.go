package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"first-api/internal/model"
)

type CustomerUseCase interface {
	GetCustomers(context.Context) ([]model.Customer, error)
	CreateCustomer(context.Context, *http.Request) (*model.Customer, error)
	UpdateCustomer(context.Context, *http.Request) (*model.Customer, error)
	DeleteCustomer(context.Context, *http.Request) error
}

type CustomerHandler struct {
	useCase CustomerUseCase
}

func NewCustomerHandler(useCase CustomerUseCase) *CustomerHandler {
	return &CustomerHandler{useCase: useCase}
}

func (c *CustomerHandler) GetCustomers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	customers, err := c.useCase.GetCustomers(ctx)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(customers)
}

func (c *CustomerHandler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	customer, err := c.useCase.CreateCustomer(ctx, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(*customer)

}

func (c *CustomerHandler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	customer, err := c.useCase.UpdateCustomer(ctx, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(*customer)
}

func (c *CustomerHandler) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := c.useCase.DeleteCustomer(ctx, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)

}
