package repository

import (
	"context"
	"database/sql"
	"errors"
	"first-api/internal/model"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CustomerRepository struct {
	connection *pgxpool.Pool
}

func NewCustomerRepository(connection *pgxpool.Pool) *CustomerRepository {
	return &CustomerRepository{
		connection: connection,
	}
}

func (pr *CustomerRepository) GetCustomers(ctx context.Context) ([]model.Customer, error) {
	query := "SELECT id, name,email,phone,created_at,updated_at FROM customers"

	rows, err := pr.connection.Query(ctx, query)
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

func (cr *CustomerRepository) GetCustomerById(ctx context.Context, customerId string) (*model.Customer, error) {
	query := `SELECT id,name,email,phone,created_at,updated_at from customers WHERE id=$1`
	var customer model.Customer
	row := cr.connection.QueryRow(ctx, query, customerId)
	err := row.Scan(&customer.ID, &customer.Name, &customer.Email, &customer.Phone, &customer.CreatedAt, &customer.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.CustomerNotFound
		}
		//se for outro erro
		return nil, err
	}
	return &customer, nil
}

func (pr *CustomerRepository) CreateCustomer(ctx context.Context, customer *model.Customer) error {
	query := `INSERT INTO customers (id,name,email,phone,password)
	VALUES ($1, $2, $3, $4, $5)`
	_, err := pr.connection.Exec(ctx,
		query,
		customer.ID,
		customer.Name,
		customer.Email,
		customer.Phone,
		customer.Password,
	)

	if err != nil {
		return err
	}

	return nil

}

func (cr *CustomerRepository) UpdateCustomer(ctx context.Context, customerId string, customer *model.Customer) error {
	_, err := cr.GetCustomerById(ctx, customerId)
	if err != nil {
		return err
	}

	query := `UPDATE customers
	SET name=$1,
		email=$2,
		phone=$3,
		password=$4
	WHERE id=$5;
	`

	if _, err := cr.connection.Exec(ctx, query, customer.Name, customer.Email, customer.Phone, customer.Password, customerId); err != nil {
		return err
	}
	return nil

}

func (cr *CustomerRepository) DeleteCustomer(ctx context.Context, customerId string) error {
	query := `DELETE FROM customers WHERE customers.id = $1`

	if _, err := cr.connection.Exec(ctx, query, customerId); err != nil {
		return err
	}
	return nil

}

func (cr *CustomerRepository) GetCustomerByField(ctx context.Context, field string, value string) (*model.Customer, error) {
	var customer model.Customer
	//WHITELIST pra proteger contra injection
	if field != "id" && field != "name" && field != "email" && field != "phone" {
		return &customer, model.ErrInvalidField
	}

	query := fmt.Sprintf(`SELECT id,name,email,phone FROM customers WHERE %s = $1`, field)
	err := cr.connection.QueryRow(ctx, query, value).Scan(&customer.ID, &customer.Name, &customer.Email, &customer.Phone)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.CustomerNotFound
		}
		//se for outro erro
		return nil, err
	}

	return &customer, err
}

func (cu *CustomerRepository) SaveRefreshToken(ctx context.Context, token *model.RefreshToken) error {
	query := "INSERT INTO refresh_tokens (token,customer_id,expires_at,revoked) VALUES ($1,$2,$3,$4)"
	_, err := cu.connection.Exec(ctx, query, token.Token, token.CustomerID, token.ExpiresAt, token.Revoked)

	return err
}

func (cu *CustomerRepository) GetRefreshToken(ctx context.Context, tokenStr string) (*model.RefreshToken, error) {
	var refreshToken model.RefreshToken

	query := "SELECT token, customer_id,expires_at,revoked FROM refresh_tokens WHERE token=$1"

	if err := cu.connection.QueryRow(ctx, query, tokenStr).Scan(&refreshToken.Token, &refreshToken.CustomerID, &refreshToken.ExpiresAt, &refreshToken.Revoked); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrInvalidRefreshToken
		}
		return nil, err
	}

	return &refreshToken, nil

}

func (cu *CustomerRepository) RevokeRefreshToken(ctx context.Context, tokenStr string) error {
	query := "UPDATE refresh_tokens SET revoked=TRUE WHERE token=$1"
	if _, err := cu.connection.Exec(ctx, query, tokenStr); err != nil {
		return err
	}

	return nil
}
