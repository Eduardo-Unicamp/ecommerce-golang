package usecases

import (
	"context"
	"database/sql"
	"errors"
	"first-api/internal/middleware"
	"first-api/internal/model"

	"github.com/google/uuid"
)

type OrderRepository interface {
	CreateOrder(context.Context, *model.Order) (*model.Order, error)
	GetOrders(context.Context, int, int, string) (*[]model.Order, error)
	GetOrderByID(context.Context, string, string) (*model.Order, error)
	UpdateOrderStatus(context.Context, *model.Order) error
	CancelOrder(context.Context, *model.Order) error
}

type ProductRepositoryForOrder interface {
	GetProductInfo(context.Context, []model.NewOrderItemDTO) (*model.ProductInfo, error)
}

type OrderUseCase struct {
	orderRepository OrderRepository
	pr              ProductRepositoryForOrder
}

func NewOrderUseCase(orderRepository OrderRepository, pr ProductRepositoryForOrder) *OrderUseCase {
	return &OrderUseCase{orderRepository: orderRepository, pr: pr}
}

func (ou *OrderUseCase) CreateOrder(ctx context.Context, request model.NewOrderDTO) (*model.Order, error) {

	//getting uuid string from token and parsing to UUID
	customerID := middleware.GetUserIDFromToken(ctx)
	parsedCustomerID, err := uuid.Parse(customerID)
	if err != nil {
		return nil, model.ErrInvalidIDFormat
	}

	//get product info
	productInfo, err := ou.pr.GetProductInfo(ctx, request.Items)

	//create new order
	newOrder, err := model.NewOrder(parsedCustomerID, request.Items, productInfo)

	if err != nil {
		return newOrder, err
	}

	newOrder, err = ou.orderRepository.CreateOrder(ctx, newOrder)

	return newOrder, err

}

func (ou *OrderUseCase) GetOrders(ctx context.Context, limit int, offset int) (*[]model.Order, error) {

	customerID := middleware.GetUserIDFromToken(ctx)

	return ou.orderRepository.GetOrders(ctx, limit, offset, customerID)

}

func (ou *OrderUseCase) GetOrderByID(ctx context.Context, orderID string) (*model.Order, error) {

	customerID := middleware.GetUserIDFromToken(ctx) //user autenticado
	if customerID == "" {
		return nil, model.ErrAuthorizationFailed
	}

	if err := uuid.Validate(orderID); err != nil {
		return nil, err
	}

	order, err := ou.orderRepository.GetOrderByID(ctx, orderID, customerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrOrderNotFound
		}
	}
	if order.Items == nil {
		return order, model.ErrEmptyOrder
	}

	return order, nil
}

func (ou *OrderUseCase) PayOrder(ctx context.Context, orderID string) error {

	customerID := middleware.GetUserIDFromToken(ctx)

	order, err := ou.orderRepository.GetOrderByID(ctx, orderID, customerID)

	if err != nil {
		return err
	}
	//regra de negocio: cancelado ou pago nao pode mudar de status
	if order.Status != model.PENDING {
		return model.ErrUnableToPay
	}

	order.Pay()
	err = ou.orderRepository.UpdateOrderStatus(ctx, order)
	return err
}

func (ou *OrderUseCase) CancelOrder(ctx context.Context, orderID string) error {
	customerID := middleware.GetUserIDFromToken(ctx)
	order, err := ou.orderRepository.GetOrderByID(ctx, orderID, customerID)
	if err != nil {
		return err
	}
	//regra de negocio: cancelado ou pago nao pode mudar de status
	if order.Status != model.PENDING {
		return model.ErrUnableToCancel
	}
	order.Cancel()
	if err := ou.orderRepository.CancelOrder(ctx, order); err != nil {
		return err
	}

	return nil
}
