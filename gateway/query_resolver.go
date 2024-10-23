package main

import (
	"context"
	"errors"
	"fmt"
	"graphql-grpc-go-microservice-project/gateway/models"
	"graphql-grpc-go-microservice-project/gateway/utils"

	"github.com/google/uuid"
)

type queryResolver struct {
	server *GatewayServer
}

func (r *queryResolver) GetAccountByID(ctx context.Context, id string) (*models.Account, error) {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid ID format")
	}

	account, err := r.server.AccountClient.GetAccountByID(ctx, uuidID.String())
	if err != nil {
		return nil, err
	}

	return utils.ConvertAccountToModel(account), nil
}

func (r *queryResolver) GetAccountByEmail(ctx context.Context, email string) (*models.Account, error) {
	account, err := r.server.AccountClient.GetAccountByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return utils.ConvertAccountToModel(account), nil
}

func (r *queryResolver) ListAccounts(ctx context.Context, pagination *models.PaginationInput) ([]*models.Account, error) {
	var limit, offset uint32
	limit, offset = 0, 0

	if pagination != nil {
		limit = uint32(pagination.Limit)
		offset = uint32(pagination.Offset)

		if limit > 20 {
			return nil, fmt.Errorf("failed to fetch accounts: limit %d is greater than 20", limit)
		}
	}

	accounts, err := r.server.AccountClient.ListAccounts(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	var accountList []*models.Account
	for _, acc := range accounts {
		accountList = append(accountList, utils.ConvertAccountToModel(&acc))
	}
	return accountList, nil
}

func (r *queryResolver) GetProductByID(ctx context.Context, id string) (*models.Product, error) {
	product, err := r.server.ProductClient.GetProductByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return utils.ConvertProductToModel(product), nil
}

func (r *queryResolver) ListProducts(ctx context.Context, pagination *models.PaginationInput) ([]*models.Product, error) {
	var limit, offset uint32
	limit, offset = 0, 0

	if pagination != nil {
		limit = uint32(pagination.Limit)
		offset = uint32(pagination.Offset)

		if limit > 20 {
			return nil, fmt.Errorf("failed to fetch accounts: limit %d is greater than 20", limit)
		}
	}

	products, err := r.server.ProductClient.ListProducts(ctx, uint32(limit), uint32(offset))
	if err != nil {
		return nil, err
	}

	var productList []*models.Product
	for _, prod := range products {
		productList = append(productList, utils.ConvertProductToModel(prod))
	}
	return productList, nil
}

func (r *queryResolver) ListProductsWithIDs(ctx context.Context, ids []string, pagination *models.PaginationInput) ([]*models.Product, error) {
	var limit, offset uint32
	limit, offset = 0, 0

	if pagination != nil {
		limit = uint32(pagination.Limit)
		offset = uint32(pagination.Offset)

		if limit > 20 {
			return nil, fmt.Errorf("failed to fetch accounts: limit %d is greater than 20", limit)
		}
	}

	products, err := r.server.ProductClient.ListProductsWithIDs(ctx, ids, uint32(limit), uint32(offset))
	if err != nil {
		return nil, err
	}

	var productList []*models.Product
	for _, prod := range products {
		productList = append(productList, utils.ConvertProductToModel(prod))
	}
	return productList, nil
}

func (r *queryResolver) SearchProducts(ctx context.Context, search string, pagination *models.PaginationInput) ([]*models.Product, error) {
	var limit, offset uint32
	limit, offset = 0, 0

	if pagination != nil {
		limit = uint32(pagination.Limit)
		offset = uint32(pagination.Offset)

		if limit > 20 {
			return nil, fmt.Errorf("failed to fetch accounts: limit %d is greater than 20", limit)
		}
	}

	products, err := r.server.ProductClient.SearchProducts(ctx, search, uint32(limit), uint32(offset))
	if err != nil {
		return nil, err
	}

	var productList []*models.Product
	for _, prod := range products {
		productList = append(productList, utils.ConvertProductToModel(prod))
	}
	return productList, nil
}
