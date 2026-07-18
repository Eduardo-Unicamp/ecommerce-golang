package repository

import (
	"context"
	"first-api/internal/model"
	"strings"
)

type OrderRepository struct {
	connection ConnectionPool
}

func NewOrderRepository(connection ConnectionPool) *OrderRepository {
	return &OrderRepository{connection: connection}
}

func (or *OrderRepository) CreateOrder(ctx context.Context, newOrder *model.Order) (*model.Order, error) {

	transaction, err := or.connection.Begin(ctx)
	if err != nil {
		return &model.Order{}, err
	}
	defer transaction.Rollback(ctx)

	query := `INSERT INTO orders (id,status,customer_id)
VALUES ($1,$2,$3);
`

	if _, err := transaction.Exec(ctx, query, newOrder.ID, newOrder.Status, newOrder.CustomerID); err != nil {
		if strings.Contains(err.Error(), "SQLSTATE 23503") {
			return nil, model.CustomerNotFound
		}
		return nil, err
	}

	for _, item := range newOrder.Items {

		//Contabilização do estoque do produto
		query = `
		UPDATE products
		SET stock = stock-$1
		WHERE id = $2;
		`
		_, err := transaction.Exec(ctx, query, item.UnitsOrdered, item.ProductID)
		if err != nil {
			return &model.Order{}, err
		}

		//Add item ao banco
		query = `INSERT INTO order_items (id,selling_price,units,product_id,order_id)
		VALUES ($1,$2,$3,$4,$5);`

		_, err = transaction.Exec(ctx, query, item.ID.String(), item.SellingPrice, item.UnitsOrdered, item.ProductID, newOrder.ID)
		if err != nil {
			return &model.Order{}, err
		}

	}

	//se chegou até nao deu erro, entao commita as mudanças
	//o que foi commitado o rollback não mexe
	if err := transaction.Commit(ctx); err != nil {
		return &model.Order{}, err
	}

	return newOrder, err

}

func (or *OrderRepository) CancelOrder(ctx context.Context, order *model.Order) error {
	//start transaction
	tx, err := or.connection.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	//repoe o estoque do pedido cancelado
	for _, item := range order.Items {
		units := item.UnitsOrdered
		query := `
		UPDATE products
		SET stock = stock+$1
		WHERE id=$2
		`
		if _, err := tx.Exec(ctx, query, units, item.ProductID); err != nil {
			return err
		}

	}
	//atualiza o order status
	query := `
	UPDATE orders
	SET status = $1
	WHERE id = $2;
	`
	if _, err := tx.Exec(ctx, query, order.Status, order.ID); err != nil {
		return err
	}

	tx.Commit(ctx)

	return nil

}

func (or *OrderRepository) UpdateOrderStatus(ctx context.Context, order *model.Order) error {

	query := `
	UPDATE orders
	SET status = $1
	WHERE id = $2;
	`
	if _, err := or.connection.Exec(ctx, query, order.Status, order.ID); err != nil {
		return err
	}

	return nil
}

func (or *OrderRepository) GetOrders(ctx context.Context, limit int, offset int, customerID string) (*[]model.Order, error) {
	var orders []model.Order

	tx, err := or.connection.Begin(ctx)
	if err != nil {
		return &[]model.Order{}, err
	}
	defer tx.Rollback(ctx)

	//first loops once and gets products cause pgx dont allow new .Scan whithout closing last one so cant nest
	query := `SELECT * FROM orders WHERE customer_id = $1 LIMIT $2 OFFSET $3;`

	rows, err := tx.Query(ctx, query, customerID, limit, offset)

	for rows.Next() {
		var newOrder model.Order
		err = rows.Scan(
			&newOrder.ID,
			&newOrder.Status,
			&newOrder.CustomerID,
		)
		if err != nil {
			return nil, err
		}

		orders = append(orders, newOrder)

	}
	rows.Close()

	//no fill items for each product
	for i := 0; i < len(orders); i++ {
		order := &orders[i]
		query = `SELECT id,selling_price,units,product_id FROM order_items WHERE order_id=$1`
		rowsItem, err := tx.Query(ctx, query, order.ID)

		for rowsItem.Next() {
			newItem := model.OrderItem{}

			err = rowsItem.Scan(
				&newItem.ID,
				&newItem.SellingPrice,
				&newItem.UnitsOrdered,
				&newItem.ProductID,
			)
			if err != nil {
				return &[]model.Order{}, err
			}

			order.Items = append(order.Items, newItem)
		}
		rowsItem.Close()
	}

	err = tx.Commit(ctx)
	return &orders, nil

}

func (or *OrderRepository) GetOrderByID(ctx context.Context, orderID string, customerID string) (*model.Order, error) {
	var newOrder model.Order

	tx, err := or.connection.Begin(ctx)
	if err != nil {
		return &model.Order{}, err
	}
	defer tx.Rollback(ctx)

	//add order
	query := `SELECT * FROM orders WHERE id=$1 AND customer_id=$2;`
	row := tx.QueryRow(ctx, query, orderID, customerID)
	if err := row.Scan(&newOrder.ID, &newOrder.Status, &newOrder.CustomerID); err != nil {
		return nil, err
	}

	//add items
	query = `SELECT id,selling_price,units,product_id FROM order_items WHERE order_id=$1`
	rows, err := tx.Query(ctx, query, orderID)

	for rows.Next() {
		newItem := model.OrderItem{}

		err = rows.Scan(
			&newItem.ID,
			&newItem.SellingPrice,
			&newItem.UnitsOrdered,
			&newItem.ProductID,
		)
		if err != nil {
			return nil, err
		}

		newOrder.Items = append(newOrder.Items, newItem)
	}

	err = tx.Commit(ctx)

	return &newOrder, err

}
