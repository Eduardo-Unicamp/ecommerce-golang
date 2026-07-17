package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"first-api/internal/model"

	"github.com/go-chi/chi/v5"
)

type ProductUseCase interface {
	GetProducts(context.Context) (*[]model.Product, error)
	GetProductByID(context.Context, string) (*model.Product, error)
	CreateProduct(context.Context, *model.CreateProductRequest) (*model.Product, error)
	UpdateProduct(context.Context, string, *model.UpdateProductRequest) (*model.Product, error)
	DeleteProduct(context.Context, string) error
}

type ProductHandler struct {
	UseCase ProductUseCase
}

func NewProductHandler(useCase ProductUseCase) *ProductHandler {
	return &ProductHandler{UseCase: useCase}
}

func (p *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	products, err := p.UseCase.GetProducts(ctx)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

func (p *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	productID := chi.URLParam(r, "product_id")
	product, err := p.UseCase.GetProductByID(ctx, productID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

func (p *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var request model.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		WriteOrderError(w, err)
		return
	}

	product, err := p.UseCase.CreateProduct(ctx, &request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(*product)

}

func (p *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	productID := chi.URLParam(r, "product_id")
	var request model.UpdateProductRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		WriteOrderError(w, err)
		return
	}
	product, err := p.UseCase.UpdateProduct(ctx, productID, &request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(*product)
}

func (p *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	productID := chi.URLParam(r, "product_id")
	err := p.UseCase.DeleteProduct(ctx, productID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)

}
