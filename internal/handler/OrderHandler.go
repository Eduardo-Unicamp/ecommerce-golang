package handler

import (
	"context"
	"strconv"

	"encoding/json"

	"first-api/internal/middleware"
	"first-api/internal/model"

	"net/http"

	"github.com/go-chi/chi/v5"
)

type OrderUseCase interface {
	CreateOrder(ctx context.Context, request model.NewOrderDTO) (*model.Order, error)
	GetOrders(ctx context.Context, limit int, offset int) (*[]model.Order, error)
	GetOrderByID(ctx context.Context, orderID string, customerID string) (*model.Order, error)
	PayOrder(ctx context.Context, orderID string) error
	CancelOrder(ctx context.Context, orderID string) error
}

type OrderHandler struct {
	UseCase OrderUseCase
}

func NewOrderHandler(orderUseCase OrderUseCase) *OrderHandler {
	return &OrderHandler{UseCase: orderUseCase}
}

func (oh *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	//parte request to a dto
	var request model.NewOrderDTO
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		WriteOrderError(w, err)
	}
	//
	order, err := oh.UseCase.CreateOrder(ctx, request)
	if err != nil {
		WriteOrderError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(*order)

}

func (oh *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	limit, offset := extractLimitAndOffset(r)
	orders, err := oh.UseCase.GetOrders(ctx, limit, offset)
	if err != nil {
		WriteOrderError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(*orders)

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

func (oh *OrderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	orderID := chi.URLParam(r, "order_id")
	customerID := middleware.GetUserIDFromToken(ctx) //user autenticado
	order, err := oh.UseCase.GetOrderByID(ctx, orderID, customerID)
	if err != nil {
		WriteOrderError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(*order)

}

func (oh *OrderHandler) PayOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	orderID := chi.URLParam(r, "order_id")

	if err := oh.UseCase.PayOrder(ctx, orderID); err != nil {
		WriteOrderError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

}

func (oh *OrderHandler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	orderID := chi.URLParam(r, "order_id")

	if err := oh.UseCase.CancelOrder(ctx, orderID); err != nil {
		WriteOrderError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

}
