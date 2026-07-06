package handler

import (
	"encoding/json"
	"net/http"

	"first-api/internal/model"
)

type CustomerUseCase interface {
	GetCustomers() ([]model.Customer, error)
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
