package repository

import (
	"database/sql"
	"first-api/internal/model"
	"fmt"
)

type CustomerRepository struct {
	connection *sql.DB
}

func NewCustomerRepository(connection *sql.DB) *CustomerRepository {
	return &CustomerRepository{
		connection: connection,
	}
}

func (pr *CustomerRepository) GetCustomers() ([]model.Customer, error) {
	query := "SELECT id, name,email,phone FROM customer"

	rows, err := pr.connection.Query(query)
	if err != nil {
		fmt.Println(err)
		return []model.Customer{}, err
	}

	var customerList []model.Customer
	var customerObj model.Customer
	for rows.Next() {
		err = rows.Scan(
			&customerObj.ID,
			&customerObj.Name,
			&customerObj.Email,
			&customerObj.Phone)
		if err != nil {
			fmt.Println(err)
			return []model.Customer{}, err
		}

		customerList = append(customerList, customerObj)
	}
	rows.Close()

	return customerList, nil
}
