package main

import (
	"context"
	"errors"
	"graphql-grpc-go-microservice-project/gateway/models"
)

type accountResolver struct {
	server *GatewayServer
}

func (r *accountResolver) ID(ctx context.Context, obj *models.Account) (string, error) {
	if obj == nil {
		return "", errors.New("account object is nil")
	}
	return obj.ID.String(), nil
}
