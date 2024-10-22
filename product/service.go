package product

import "context"

type ProductService interface {
	CreateProduct(ctx context.Context, name, description string, price float64) (*Product, error)
	GetProductByID(ctx context.Context, id string) (*Product, error)
	ListProducts(ctx context.Context, limit, offset int) ([]Product, error)
	ListProductsWithIDs(ctx context.Context, ids []string) ([]Product, error)
	SearchProducts(ctx context.Context, query string, limit, offset int) ([]Product, error)
}

type productService struct {
	repository ProductRepository
}

func NewProductService(repository ProductRepository) (ProductService, error) {
	return &productService{repository: repository}, nil
}

func (service *productService) CreateProduct(ctx context.Context, name, description string, price float64) (*Product, error) {
	product, err := service.repository.CreateProduct(ctx, name, description, price)
	if err == nil {
		return product, nil
	}

	return nil, err
}

func (service *productService) GetProductByID(ctx context.Context, id string) (*Product, error) {
	product, err := service.repository.GetProductByID(ctx, id)
	if err == nil {
		return product, nil
	}

	return nil, err
}

func (service *productService) ListProducts(ctx context.Context, limit, offset int) ([]Product, error) {
	products, err := service.repository.ListProducts(ctx, limit, offset)
	if err == nil {
		return products, nil
	}

	return nil, err
}

func (service *productService) ListProductsWithIDs(ctx context.Context, ids []string) ([]Product, error) {
	products, err := service.repository.ListProductsWithIDs(ctx, ids)
	if err == nil {
		return products, nil
	}

	return nil, err
}

func (service *productService) SearchProducts(ctx context.Context, query string, limit, offset int) ([]Product, error) {
	products, err := service.repository.SearchProducts(ctx, query, limit, offset)
	if err == nil {
		return products, nil
	}

	return nil, err
}
