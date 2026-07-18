package repository

import (
	"context"
	"database/sql"
	"first-api/internal/model"
	"fmt"

	"github.com/shopspring/decimal"
)

type ProductRepository struct {
	connection ConnectionPool
}

func NewProductRepository(connection ConnectionPool) *ProductRepository {
	return &ProductRepository{
		connection: connection,
	}
}

func (pr *ProductRepository) GetProducts(ctx context.Context) (*[]model.Product, error) {
	query := "SELECT id, name,price,stock FROM products"

	rows, err := pr.connection.Query(ctx, query)
	if err != nil {
		fmt.Println(err)
		return &[]model.Product{}, err
	}

	var productList []model.Product
	var productObj model.Product
	for rows.Next() {
		err = rows.Scan(
			&productObj.ID,
			&productObj.Name,
			&productObj.Price,
			&productObj.Stock,
		)
		if err != nil {
			fmt.Println(err)
			return &[]model.Product{}, err
		}

		productList = append(productList, productObj)
	}
	rows.Close()

	return &productList, nil
}

func (pr *ProductRepository) GetProductByID(ctx context.Context, productId string) (*model.Product, error) {
	query := `SELECT * from products WHERE id=$1`
	var product model.Product
	row := pr.connection.QueryRow(ctx, query, productId)
	err := row.Scan(&product.ID, &product.Name, &product.Price, &product.Stock)
	if err != nil {
		if err == sql.ErrNoRows {
			return &product, model.ErrProductNotFound
		}
		//se for outro erro
		return &product, err
	}
	return &product, nil
}

func (pr *ProductRepository) CreateProduct(ctx context.Context, product *model.Product) error {
	query := `INSERT INTO products (id,name,price,stock)
	VALUES ($1, $2, $3, $4)`
	_, err := pr.connection.Exec(ctx,
		query,
		product.ID,
		product.Name,
		product.Price,
		product.Stock,
	)

	if err != nil {
		return err
	}

	return nil

}

func (pr *ProductRepository) UpdateProduct(ctx context.Context, productId string, product *model.Product) error {
	_, err := pr.GetProductByID(ctx, productId)
	if err != nil {
		return err
	}

	query := `UPDATE products
	SET name=$1,
		price=$2,
		stock=$3
	WHERE id=$4;
	`

	if _, err := pr.connection.Exec(ctx, query, product.Name, product.Price, product.Stock, productId); err != nil {
		return err
	}
	return nil

}

func (pr *ProductRepository) DeleteProduct(ctx context.Context, productId string) error {
	query := `DELETE FROM products WHERE products.id = $1`

	if _, err := pr.connection.Exec(ctx, query, productId); err != nil {
		return err
	}
	return nil

}

func (pr *ProductRepository) GetProductInfo(ctx context.Context, itemDTOs []model.NewOrderItemDTO) (*model.ProductInfo, error) {
	var ids []string
	results := model.NewProductInfo()

	for _, itemDTO := range itemDTOs {
		ids = append(ids, itemDTO.ProductID.String())
	}

	query := `SELECT id,price,stock FROM products WHERE id=ANY($1)`
	rows, err := pr.connection.Query(ctx, query, ids)
	if err != nil {
		return &model.ProductInfo{}, err
	}
	var id string
	var price decimal.Decimal
	var stock int
	for rows.Next() {
		if err := rows.Scan(&id, &price, &stock); err != nil {
			return &model.ProductInfo{}, err
		}

		results.CurrentPrices[id] = price
		results.CurrentStock[id] = stock

	}
	return results, nil
}
