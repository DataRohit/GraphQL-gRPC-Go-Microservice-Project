package main

import (
	"context"
	"graphql-grpc-go-microservice-project/gateway/models"
	"graphql-grpc-go-microservice-project/gateway/utils"
	"time"
)

type mutationResolver struct {
	server *GatewayServer
}

func (r *mutationResolver) CreateAccount(ctx context.Context, in models.AccountInput) (*models.Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	account, err := r.server.AccountClient.CreateAccount(ctx, in.Email, in.Name)
	if err != nil {
		return nil, err
	}

	return utils.ConvertAccountToModel(account), nil
}

func (r *mutationResolver) CreateProduct(ctx context.Context, in models.ProductInput) (*models.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	product, err := r.server.ProductClient.CreateProduct(ctx, in.Name, in.Description, in.Price)
	if err != nil {
		return nil, err
	}

	return utils.ConvertProductToModel(product), nil
}
