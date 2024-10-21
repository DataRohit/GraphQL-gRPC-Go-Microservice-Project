package main

import (
	"context"
	"errors"
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

	return convertToModel(account), nil
}

func (r *queryResolver) GetAccountByEmail(ctx context.Context, email string) (*models.Account, error) {
	account, err := r.server.AccountClient.GetAccountByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return convertToModel(account), nil
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
		accountList = append(accountList, convertToModel(&acc))
	}
	return accountList, nil
}
