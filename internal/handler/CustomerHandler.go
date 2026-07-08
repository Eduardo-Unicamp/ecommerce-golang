package handler

import (
	"encoding/json"
	"net/http"

	"first-api/internal/model"
)

type CustomerUseCase interface {
	GetCustomers() ([]model.Customer, error)
	CreateCustomer(*http.Request) (*model.Customer, error)
	UpdateCustomer(*http.Request) (*model.Customer, error)
	DeleteCustomer(*http.Request) error
}

type CustomerHandler struct {
	useCase CustomerUseCase
}

func NewCustomerHandler(useCase CustomerUseCase) *CustomerHandler {
	return &CustomerHandler{useCase: useCase}
}

func (c *CustomerHandler) GetCustomers(w http.ResponseWriter, r *http.Request) {
	customers, err := c.useCase.GetCustomers()
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
	customer, err := c.useCase.CreateCustomer(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(*customer)

}

func (c *CustomerHandler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	customer, err := c.useCase.UpdateCustomer(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(*customer)
}

func (c *CustomerHandler) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	err := c.useCase.DeleteCustomer(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)

}
