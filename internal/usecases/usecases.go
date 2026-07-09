package usecases

import (
	"context"
	"encoding/json"
	"first-api/internal/model"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type CustomerRepository interface {
	GetCustomers(context.Context) ([]model.Customer, error)
	CreateCustomer(context.Context, *model.Customer) error
	UpdateCustomer(context.Context, string, *model.Customer) error
	DeleteCustomer(context.Context, string) error
}

type CustomerUseCase struct {
	repository CustomerRepository
}

func NewCustomerUseCase(repository CustomerRepository) *CustomerUseCase {
	return &CustomerUseCase{
		repository: repository,
	}
}

func (pu *CustomerUseCase) GetCustomers(ctx context.Context) ([]model.Customer, error) {
	customers, err := pu.repository.GetCustomers(ctx)
	if err != nil {
		return []model.Customer{}, err
	}

	return customers, err
}

func (pu *CustomerUseCase) CreateCustomer(ctx context.Context, r *http.Request) (*model.Customer, error) {
	var request model.CreateCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	customer, err := model.NewCustomer(request.Name, request.Email, request.Phone)
	if err != nil {
		return nil, err
	}

	if err := pu.repository.CreateCustomer(ctx, customer); err != nil {
		return nil, err
	}

	return customer, nil

}

func (cu *CustomerUseCase) UpdateCustomer(ctx context.Context, r *http.Request) (*model.Customer, error) {
	customerId := chi.URLParam(r, "customerId")
	var request model.UpdateCustomerRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return &model.Customer{}, err
	}

	customer, err := model.NewCustomer(request.Name, request.Email, request.Phone)
	if err != nil {
		return customer, err
	}

	err = cu.repository.UpdateCustomer(ctx, customerId, customer)

	return customer, err

}

func (pu *CustomerUseCase) DeleteCustomer(ctx context.Context, r *http.Request) error {
	customerId := chi.URLParam(r, "customerId")
	err := pu.repository.DeleteCustomer(ctx, customerId)
	if err != nil {
		return err
	}
	return nil

}
