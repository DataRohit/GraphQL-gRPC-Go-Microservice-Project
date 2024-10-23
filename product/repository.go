package product

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/google/uuid"
)

type productDocument struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type ProductRepository interface {
	Close()
	CreateProduct(ctx context.Context, name, description string, price float64) (*Product, error)
	GetProductByID(ctx context.Context, id string) (*Product, error)
	ListProducts(ctx context.Context, offset, limit int) ([]Product, error)
	ListProductsWithIDs(ctx context.Context, ids []string) ([]Product, error)
	SearchProducts(ctx context.Context, query string, offset, limit int) ([]Product, error)
}

type elasticRepository struct {
	client *elasticsearch.Client
}

func NewElasticRepository(url string) (ProductRepository, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{url},
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &elasticRepository{client: client}, nil
}

func (r *elasticRepository) Close() {}

func (r *elasticRepository) CreateProduct(ctx context.Context, name, description string, price float64) (*Product, error) {
	productID := uuid.NewString()
	product := productDocument{
		Name:        name,
		Description: description,
		Price:       price,
	}

	body, err := json.Marshal(product)
	if err != nil {
		return nil, err
	}

	res, err := r.client.Index(
		"catalog",
		bytes.NewReader(body),
		r.client.Index.WithDocumentID(productID),
		r.client.Index.WithContext(ctx),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("failed to index product: %s", res.String())
	}

	return &Product{
		ID:          productID,
		Name:        name,
		Description: description,
		Price:       price,
	}, nil
}

func (r *elasticRepository) GetProductByID(ctx context.Context, id string) (*Product, error) {
	res, err := r.client.Get(
		"catalog",
		id,
		r.client.Get.WithContext(ctx),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("failed to get product by id: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	source, ok := result["_source"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to retrieve _source for document %s", id)
	}

	p := productDocument{
		Name:        source["name"].(string),
		Description: source["description"].(string),
		Price:       source["price"].(float64),
	}

	return &Product{
		ID:          id,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
	}, nil
}

func (r *elasticRepository) ListProducts(ctx context.Context, offset, limit int) ([]Product, error) {
	query := `{
		"from": %d,
		"size": %d,
		"query": {
			"match_all": {}
		}
	}`
	searchBody := fmt.Sprintf(query, offset, limit)
	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex("catalog"),
		r.client.Search.WithBody(strings.NewReader(searchBody)),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("failed to list products: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	products := []Product{}
	hits := result["hits"].(map[string]interface{})["hits"].([]interface{})
	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"]
		product := Product{
			ID:          hit.(map[string]interface{})["_id"].(string),
			Name:        source.(map[string]interface{})["name"].(string),
			Description: source.(map[string]interface{})["description"].(string),
			Price:       source.(map[string]interface{})["price"].(float64),
		}
		products = append(products, product)
	}

	return products, nil
}

func (r *elasticRepository) ListProductsWithIDs(ctx context.Context, ids []string) ([]Product, error) {
	var buf bytes.Buffer

	query := `{
        "query": {
            "ids": {
                "values": %s
            }
        }
    }`

	idsStr, _ := json.Marshal(ids)
	searchBody := fmt.Sprintf(query, string(idsStr))

	buf.WriteString(searchBody)

	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex("catalog"),
		r.client.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("failed to retrieve products by IDs: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	products := []Product{}
	hits := result["hits"].(map[string]interface{})["hits"].([]interface{})
	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"]
		product := Product{
			ID:          hit.(map[string]interface{})["_id"].(string),
			Name:        source.(map[string]interface{})["name"].(string),
			Description: source.(map[string]interface{})["description"].(string),
			Price:       source.(map[string]interface{})["price"].(float64),
		}
		products = append(products, product)
	}

	return products, nil
}

func (r *elasticRepository) SearchProducts(ctx context.Context, query string, offset, limit int) ([]Product, error) {
	var buf bytes.Buffer

	searchQuery := `{
		"from": %d,
		"size": %d,
		"query": {
			"multi_match": {
				"query": "%s",
				"fields": ["name", "description"]
			}
		}
	}`
	searchBody := fmt.Sprintf(searchQuery, offset, limit, query)
	buf.WriteString(searchBody)

	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex("catalog"),
		r.client.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("failed to search products: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	products := []Product{}
	hits := result["hits"].(map[string]interface{})["hits"].([]interface{})
	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"]
		product := Product{
			ID:          hit.(map[string]interface{})["_id"].(string),
			Name:        source.(map[string]interface{})["name"].(string),
			Description: source.(map[string]interface{})["description"].(string),
			Price:       source.(map[string]interface{})["price"].(float64),
		}
		products = append(products, product)
	}

	return products, nil
}
