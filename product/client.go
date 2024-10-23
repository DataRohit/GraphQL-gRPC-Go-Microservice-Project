package product

import (
	"context"
	"crypto/tls"
	"graphql-grpc-go-microservice-project/common"
	"graphql-grpc-go-microservice-project/product/protobuf"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type ProductClient struct {
	conn    *grpc.ClientConn
	service protobuf.ProductServiceClient
	logger  *zap.Logger
}

func NewProductClient(url string, secure bool) (*ProductClient, error) {
	logger := common.GetLogger()

	var opts []grpc.DialOption
	if secure {
		creds := credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: false,
		})
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.NewClient(url, opts...)
	if err != nil {
		logger.Error("Failed to connect to gRPC server", zap.String("url", url), zap.String("error", err.Error()))
		return nil, err
	}

	logger.Info("Connected to gRPC server", zap.String("url", url))

	client := protobuf.NewProductServiceClient(conn)

	return &ProductClient{conn: conn, service: client, logger: logger}, nil
}

func (c *ProductClient) Close() error {
	c.logger.Info("Closing gRPC connection")
	return c.conn.Close()
}

func (c *ProductClient) CreateProduct(ctx context.Context, name, description string, price float64) (*Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	c.logger.Info("CreateProduct request received", zap.String("name", name), zap.String("description", description), zap.Float64("price", price))

	r, err := c.service.CreateProduct(ctx, &protobuf.CreateProductRequest{
		Name:        name,
		Description: description,
		Price:       price,
	})
	if err != nil {
		c.logger.Error("Failed to create product", zap.String("name", name), zap.String("error", err.Error()))
		return nil, err
	}

	c.logger.Info("Product created successfully", zap.String("product_id", r.GetProduct().GetId()), zap.String("name", name), zap.String("description", description), zap.Float64("price", price))

	return &Product{
		ID:          r.GetProduct().GetId(),
		Name:        r.GetProduct().GetName(),
		Description: r.GetProduct().GetDescription(),
		Price:       r.GetProduct().GetPrice(),
	}, nil
}

func (c *ProductClient) GetProductByID(ctx context.Context, id string) (*Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	c.logger.Info("GetProduct request received", zap.String("product_id", id))

	r, err := c.service.GetProductByID(ctx, &protobuf.GetProductByIDRequest{
		Id: id,
	})
	if err != nil {
		c.logger.Error("Failed to fetch product", zap.String("product_id", id), zap.String("error", err.Error()))
		return nil, err
	}

	c.logger.Info("Product fetched successfully", zap.String("product_id", r.GetProduct().GetId()), zap.String("name", r.GetProduct().GetName()), zap.String("description", r.GetProduct().GetDescription()), zap.Float64("price", r.GetProduct().GetPrice()))

	return &Product{
		ID:          r.GetProduct().GetId(),
		Name:        r.GetProduct().GetName(),
		Description: r.GetProduct().GetDescription(),
		Price:       r.GetProduct().GetPrice(),
	}, nil
}

func (c *ProductClient) ListProducts(ctx context.Context, limit, offset int32) ([]*Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	c.logger.Info("ListProducts request received", zap.Int32("limit", limit), zap.Int32("offset", offset))

	r, err := c.service.ListProducts(ctx, &protobuf.ListProductsRequest{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		c.logger.Error("Failed to fetch products", zap.String("error", err.Error()))
		return nil, err
	}

	products := make([]*Product, 0)
	for _, p := range r.GetProducts() {
		products = append(products, &Product{
			ID:          p.GetId(),
			Name:        p.GetName(),
			Description: p.GetDescription(),
			Price:       p.GetPrice(),
		})
	}

	c.logger.Info("Products fetched successfully", zap.Int("product_count", len(products)))

	return products, nil
}

func (c *ProductClient) ListProductsWithIDs(ctx context.Context, ids []string, limit, offset int32) ([]*Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	c.logger.Info("ListProductsWithIDs request received", zap.Strings("product_ids", ids), zap.Int32("limit", limit), zap.Int32("offset", offset))

	r, err := c.service.ListProductsWithIDs(ctx, &protobuf.ListProductsWithIDsRequest{
		Ids:    ids,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		c.logger.Error("Failed to fetch products with IDs", zap.Strings("product_ids", ids), zap.String("error", err.Error()))
		return nil, err
	}

	products := make([]*Product, 0)
	for _, p := range r.GetProducts() {
		products = append(products, &Product{
			ID:          p.GetId(),
			Name:        p.GetName(),
			Description: p.GetDescription(),
			Price:       p.GetPrice(),
		})
	}

	c.logger.Info("Products fetched successfully", zap.Int("product_count", len(products)))

	return products, nil
}

func (c *ProductClient) SearchProducts(ctx context.Context, query string, limit, offset int32) ([]*Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	c.logger.Info("SearchProducts request received", zap.String("query", query), zap.Int32("limit", limit), zap.Int32("offset", offset))

	r, err := c.service.SearchProducts(ctx, &protobuf.SearchProductsRequest{
		Query:  query,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		c.logger.Error("Failed to search products", zap.String("query", query), zap.String("error", err.Error()))
		return nil, err
	}

	products := make([]*Product, 0)
	for _, p := range r.GetProducts() {
		products = append(products, &Product{
			ID:          p.GetId(),
			Name:        p.GetName(),
			Description: p.GetDescription(),
			Price:       p.GetPrice(),
		})
	}

	c.logger.Info("Products searched successfully", zap.Int("product_count", len(products)))

	return products, nil
}
