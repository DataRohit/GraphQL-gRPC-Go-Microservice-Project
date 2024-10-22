package main

import (
	"context"
	"errors"
	gatewayUtils "graphql-grpc-go-microservice-project/common/gateway"
	"graphql-grpc-go-microservice-project/gateway/models"
)

type mutationResolver struct {
	server *GatewayServer
}

func (r *mutationResolver) CreateAccount(ctx context.Context, email, name string) (*models.Account, error) {
	if email == "" || name == "" {
		return nil, errors.New("email and name cannot be empty")
	}

	account, err := r.server.AccountClient.CreateAccount(ctx, email, name)
	if err != nil {
		return nil, err
	}

	return gatewayUtils.ConvertAccountToModel(account), nil
}
