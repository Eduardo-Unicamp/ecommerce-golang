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
	query := "SELECT id, name,email,phone,created_at,updated_at FROM customer"

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
			&customerObj.Phone,
			&customerObj.CreatedAt,
			&customerObj.UpdatedAt,
		)
		if err != nil {
			fmt.Println(err)
			return []model.Customer{}, err
		}

		customerList = append(customerList, customerObj)
	}
	rows.Close()

	return customerList, nil
}

func (cr *CustomerRepository) GetCustomerById(customerId string) (model.Customer, error) {
	query := `SELECT * from customer WHERE id=$1`
	var customer model.Customer
	row := cr.connection.QueryRow(query, customerId)
	err := row.Scan(&customer.ID, &customer.Name, &customer.Email, &customer.Phone, &customer.CreatedAt, &customer.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return customer, model.CustomerNotFound
		}
		//se for outro erro
		return customer, err
	}
	return customer, nil
}

func (pr *CustomerRepository) CreateCustomer(customer *model.Customer) error {
	query := `INSERT INTO customer (id,name,email,phone)
	VALUES ($1, $2, $3, $4)`
	_, err := pr.connection.Exec(
		query,
		customer.ID,
		customer.Name,
		customer.Email,
		customer.Phone,
	)

	if err != nil {
		return err
	}

	return nil

}

func (cr *CustomerRepository) UpdateCustomer(customerId string, customer *model.Customer) error {
	_, err := cr.GetCustomerById(customerId)
	if err != nil {
		return err
	}

	query := `UPDATE customer
	SET name=$1,
		email=$2,
		phone=$3
	WHERE id=$4;
	`

	if _, err := cr.connection.Exec(query, customer.Name, customer.Email, customer.Phone, customerId); err != nil {
		return err
	}
	return nil

}

func (cr *CustomerRepository) DeleteCustomer(customerId string) error {
	query := `DELETE FROM customer WHERE customer.id = $1`

	if _, err := cr.connection.Exec(query, customerId); err != nil {
		return err
	}
	return nil

}
