package usecases

import (
	"first-api/internal/model"
)

type CustomerRepository interface {
	GetCustomers() ([]model.Customer, error)
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
