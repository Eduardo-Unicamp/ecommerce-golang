package usecases

import (
	"encoding/json"
	"first-api/internal/model"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type CustomerRepository interface {
	GetCustomers() ([]model.Customer, error)
	CreateCustomer(*model.Customer) error
	UpdateCustomer(string, *model.Customer) error
	DeleteCustomer(string) error
}

type CustomerUseCase struct {
	repository CustomerRepository
}

func NewCustomerUseCase(repository CustomerRepository) *CustomerUseCase {
	return &CustomerUseCase{
		repository: repository,
	}
}

func (pu *CustomerUseCase) GetCustomers() ([]model.Customer, error) {
	customers, err := pu.repository.GetCustomers()
	if err != nil {
		return []model.Customer{}, err
	}

	return customers, err
}

func (pu *CustomerUseCase) CreateCustomer(r *http.Request) (*model.Customer, error) {
	var request model.CreateCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	customer, err := model.NewCustomer(request.Name, request.Email, request.Phone)
	if err != nil {
		return nil, err
	}

	if err := pu.repository.CreateCustomer(customer); err != nil {
		return nil, err
	}

	return customer, nil

}

func (cu *CustomerUseCase) UpdateCustomer(r *http.Request) (*model.Customer, error) {
	customerId := chi.URLParam(r, "customerId")
	var request model.UpdateCustomerRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return &model.Customer{}, err
	}

	customer, err := model.NewCustomer(request.Name, request.Email, request.Phone)
	if err != nil {
		return customer, err
	}

	err = cu.repository.UpdateCustomer(customerId, customer)

	return customer, err

}

func (pu *CustomerUseCase) DeleteCustomer(r *http.Request) error {
	customerId := chi.URLParam(r, "customerId")
	err := pu.repository.DeleteCustomer(customerId)
	if err != nil {
		return err
	}
	return nil

}
