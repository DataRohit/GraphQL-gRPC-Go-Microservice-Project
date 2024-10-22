package product

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	elastic "gopkg.in/olivere/elastic.v5"
)

type ProductRepository interface {
	Close()
	CreateProduct(ctx context.Context, name, description string, price float64) (*Product, error)
	GetProductByID(ctx context.Context, id string) (*Product, error)
	ListProducts(ctx context.Context, offset, limit int) ([]Product, error)
	ListProductsWithIDs(ctx context.Context, ids []string) ([]Product, error)
	SearchProducts(ctx context.Context, query string, offset, limit int) ([]Product, error)
}

type elasticRepository struct {
	client *elastic.Client
}

type productDocument struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

func NewElasticRepository(url string) (ProductRepository, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetSniff(false),
	)
	if err != nil {
		return nil, err
	}
	return &elasticRepository{client}, nil
}

func (r *elasticRepository) Close() {}

func (r *elasticRepository) CreateProduct(ctx context.Context, name, description string, price float64) (*Product, error) {
	productID := uuid.NewString()
	_, err := r.client.Index().
		Index("catalog").
		Type("product").
		Id(productID).
		BodyJson(productDocument{
			Name:        name,
			Description: description,
			Price:       price,
		}).
		Do(ctx)

	if err != nil {
		return nil, err
	}

	return &Product{
		ID:          productID,
		Name:        name,
		Description: description,
		Price:       price,
	}, nil
}

func (r *elasticRepository) GetProductByID(ctx context.Context, id string) (*Product, error) {
	res, err := r.client.Get().
		Index("catalog").
		Type("product").
		Id(id).
		Do(ctx)
	if err != nil {
		return nil, err
	}
	if !res.Found {
		return nil, fmt.Errorf("failed to get product by id: %w", err)
	}
	var p productDocument
	if err = json.Unmarshal(*res.Source, &p); err != nil {
		return nil, err
	}
	return &Product{
		ID:          id,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
	}, nil
}

func (r *elasticRepository) ListProducts(ctx context.Context, offset, limit int) ([]Product, error) {
	res, err := r.client.Search().
		Index("catalog").
		Type("product").
		Query(elastic.NewMatchAllQuery()).
		From(int(offset)).
		Size(int(limit)).
		Do(ctx)
	if err != nil {
		return nil, err
	}
	products := []Product{}
	for _, hit := range res.Hits.Hits {
		var p productDocument
		if err = json.Unmarshal(*hit.Source, &p); err == nil {
			products = append(products, Product{
				ID:          hit.Id,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
			})
		}
	}
	return products, nil
}

func (r *elasticRepository) ListProductsWithIDs(ctx context.Context, ids []string) ([]Product, error) {
	items := []*elastic.MultiGetItem{}
	for _, id := range ids {
		items = append(
			items,
			elastic.NewMultiGetItem().
				Index("catalog").
				Type("product").
				Id(id),
		)
	}
	res, err := r.client.MultiGet().
		Add(items...).
		Do(ctx)
	if err != nil {
		return nil, err
	}
	products := []Product{}
	for _, doc := range res.Docs {
		var p productDocument
		if err = json.Unmarshal(*doc.Source, &p); err == nil {
			products = append(products, Product{
				ID:          doc.Id,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
			})
		}
	}
	return products, nil
}

func (r *elasticRepository) SearchProducts(ctx context.Context, query string, offset, limit int) ([]Product, error) {
	res, err := r.client.Search().
		Index("catalog").
		Type("product").
		Query(elastic.NewMultiMatchQuery(query, "name", "description")).
		From(int(offset)).
		Size(int(limit)).
		Do(ctx)
	if err != nil {
		return nil, err
	}
	products := []Product{}
	for _, hit := range res.Hits.Hits {
		var p productDocument
		if err = json.Unmarshal(*hit.Source, &p); err == nil {
			products = append(products, Product{
				ID:          hit.Id,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
			})
		}
	}
	return products, nil
}
