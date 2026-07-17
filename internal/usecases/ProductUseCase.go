package usecases

import (
	"context"
	"first-api/internal/model"
)

type ProductRepository interface {
	GetProducts(context.Context) (*[]model.Product, error)
	GetProductByID(context.Context, string) (*model.Product, error)
	CreateProduct(context.Context, *model.Product) error
	UpdateProduct(context.Context, string, *model.Product) error
	DeleteProduct(context.Context, string) error
}

type ProductUseCase struct {
	repository ProductRepository
}

func NewProductUseCase(repository ProductRepository) *ProductUseCase {
	return &ProductUseCase{
		repository: repository,
	}
}

func (pu *ProductUseCase) GetProducts(ctx context.Context) (*[]model.Product, error) {
	products, err := pu.repository.GetProducts(ctx)
	if err != nil {
		return &[]model.Product{}, err
	}

	return products, err
}

func (pu *ProductUseCase) GetProductByID(ctx context.Context, productID string) (*model.Product, error) {
	product, err := pu.repository.GetProductByID(ctx, productID)

	return product, err

}

func (pu *ProductUseCase) CreateProduct(ctx context.Context, request *model.CreateProductRequest) (*model.Product, error) {
	product, err := model.NewProduct(request.Name, request.Price, request.Stock)
	if err != nil {
		return nil, err
	}

	if err := pu.repository.CreateProduct(ctx, product); err != nil {
		return nil, err
	}

	return product, nil

}

func (pu *ProductUseCase) UpdateProduct(ctx context.Context, productID string, request *model.UpdateProductRequest) (*model.Product, error) {
	product, err := model.NewProduct(request.Name, request.Price, request.Stock)
	if err != nil {
		return product, err
	}

	err = pu.repository.UpdateProduct(ctx, productID, product)

	return product, err

}

func (pu *ProductUseCase) DeleteProduct(ctx context.Context, productID string) error {
	err := pu.repository.DeleteProduct(ctx, productID)
	if err != nil {
		return err
	}
	return nil

}
