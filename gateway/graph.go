package main

import (
	"graphql-grpc-go-microservice-project/account"
	gatewayGraphQL "graphql-grpc-go-microservice-project/gateway/graphql"
)

type GatewayServer struct {
	AccountClient *account.AccountClient
}

func NewGraphQLServer(accountServiceURL string, secure bool) (*GatewayServer, error) {
	accountClient, err := account.NewAccountClient(accountServiceURL, secure)
	if err != nil {
		return nil, err
	}

	return &GatewayServer{
		AccountClient: accountClient,
	}, nil
}

func (s *GatewayServer) Mutation() gatewayGraphQL.MutationResolver {
	return &mutationResolver{
		server: s,
	}
}
