package usecases

import (
	"context"
	"encoding/json"
	"first-api/internal/model"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type OrderRepository interface {
	CreateOrder(context.Context, *model.Order) (*model.Order, error)
	GetOrders(context.Context, int, int, string) (*[]model.Order, error)
	GetOrderByID(context.Context, string) (*model.Order, error)
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

func (ou *OrderUseCase) CreateOrder(ctx context.Context, r *http.Request) (*model.Order, error) {

	var request model.NewOrderDTO
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return &model.Order{}, err
	}

	productInfo, err := ou.pr.GetProductInfo(ctx, request.Items)

	newOrder, err := model.NewOrder(request.CustomerID, request.Items, productInfo)

	if err != nil {
		return newOrder, err
	}

	newOrder, err = ou.orderRepository.CreateOrder(ctx, newOrder)

	return newOrder, err

}

func (ou *OrderUseCase) GetOrders(ctx context.Context, r *http.Request) (*[]model.Order, error) {
	limit, offset := extractLimitAndOffset(r)
	customerID := chi.URLParam(r, "customer_id")
	return ou.orderRepository.GetOrders(ctx, limit, offset, customerID)

}

func (ou *OrderUseCase) GetOrderByID(ctx context.Context, r *http.Request) (*model.Order, error) {
	orderId := chi.URLParam(r, "order_id")

	if err := uuid.Validate(orderId); err != nil {
		return &model.Order{}, err
	}

	order, _ := ou.orderRepository.GetOrderByID(ctx, orderId)
	if order.Items == nil {
		return order, model.ErrOrderNotFound
	}
	return order, nil
}

func extractLimitAndOffset(r *http.Request) (int, int) {
	urlParams := r.URL.Query()
	offsetStr, limitStr := urlParams.Get("offset"), urlParams.Get("limit")
	offset, limit := 0, 10 //valores default

	if offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil {
			offset = parsedOffset
		}
	}
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			limit = parsedLimit
		}
	}

	return limit, offset
}

func (ou *OrderUseCase) PayOrder(ctx context.Context, r *http.Request) error {
	orderID := chi.URLParam(r, "order_id")
	order, err := ou.orderRepository.GetOrderByID(ctx, orderID)
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

func (ou *OrderUseCase) CancelOrder(ctx context.Context, r *http.Request) error {
	orderID := chi.URLParam(r, "order_id")
	order, err := ou.orderRepository.GetOrderByID(ctx, orderID)
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
