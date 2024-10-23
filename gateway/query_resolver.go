package main

import (
	"context"
	"errors"
	gatewayUtils "graphql-grpc-go-microservice-project/common/gateway"
	"graphql-grpc-go-microservice-project/gateway/models"

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

	return gatewayUtils.ConvertAccountToModel(account), nil
}

func (r *queryResolver) GetAccountByEmail(ctx context.Context, email string) (*models.Account, error) {
	account, err := r.server.AccountClient.GetAccountByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return gatewayUtils.ConvertAccountToModel(account), nil
}

func (r *queryResolver) ListAccounts(ctx context.Context, pagination *models.PaginationInput) ([]*models.Account, error) {
	limit, offset := 0, 0
	if pagination != nil {
		limit, offset = pagination.Limit, pagination.Offset
	}

	accounts, err := r.server.AccountClient.ListAccounts(ctx, int32(limit), int32(offset))
	if err != nil {
		return nil, err
	}

	var accountList []*models.Account
	for _, acc := range accounts {
		accountList = append(accountList, gatewayUtils.ConvertAccountToModel(&acc))
	}
	return accountList, nil
}

func (r *queryResolver) GetProductByID(ctx context.Context, id string) (*models.Product, error) {
	product, err := r.server.ProductClient.GetProductByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return gatewayUtils.ConvertProductToModel(product), nil
}

func (r *queryResolver) ListProducts(ctx context.Context, pagination *models.PaginationInput) ([]*models.Product, error) {
	limit, offset := 0, 0
	if pagination != nil {
		limit, offset = pagination.Limit, pagination.Offset
	}

	products, err := r.server.ProductClient.ListProducts(ctx, int32(limit), int32(offset))
	if err != nil {
		return nil, err
	}

	var productList []*models.Product
	for _, prod := range products {
		productList = append(productList, gatewayUtils.ConvertProductToModel(prod))
	}
	return productList, nil
}

func (r *queryResolver) ListProductsWithIDs(ctx context.Context, ids []string, pagination *models.PaginationInput) ([]*models.Product, error) {
	limit, offset := 0, 0
	if pagination != nil {
		limit, offset = pagination.Limit, pagination.Offset
	}

	products, err := r.server.ProductClient.ListProductsWithIDs(ctx, ids, int32(limit), int32(offset))
	if err != nil {
		return nil, err
	}

	var productList []*models.Product
	for _, prod := range products {
		productList = append(productList, gatewayUtils.ConvertProductToModel(prod))
	}
	return productList, nil
}

func (r *queryResolver) SearchProducts(ctx context.Context, search string, pagination *models.PaginationInput) ([]*models.Product, error) {
	limit, offset := 0, 0
	if pagination != nil {
		limit, offset = pagination.Limit, pagination.Offset
	}

	products, err := r.server.ProductClient.SearchProducts(ctx, search, int32(limit), int32(offset))
	if err != nil {
		return nil, err
	}

	var productList []*models.Product
	for _, prod := range products {
		productList = append(productList, gatewayUtils.ConvertProductToModel(prod))
	}
	return productList, nil
}
