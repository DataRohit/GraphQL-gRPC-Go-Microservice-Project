package main

import (
	"graphql-grpc-go-microservice-project/account"
	gatewayGraphQL "graphql-grpc-go-microservice-project/gateway/graphql"
	"graphql-grpc-go-microservice-project/product"

	"github.com/99designs/gqlgen/graphql"
)

type GatewayServer struct {
	AccountClient *account.AccountClient
	ProductClient *product.ProductClient
}

func NewGraphQLServer(accountServiceURL string, productServiceURL string, secure bool) (*GatewayServer, error) {
	accountClient, err := account.NewAccountClient(accountServiceURL, secure)
	if err != nil {
		accountClient.Close()
		return nil, err
	}

	productClient, err := product.NewProductClient(productServiceURL, secure)
	if err != nil {
		accountClient.Close()
		return nil, err
	}

	return &GatewayServer{
		AccountClient: accountClient,
		ProductClient: productClient,
	}, nil
}

func (s *GatewayServer) Mutation() gatewayGraphQL.MutationResolver {
	return &mutationResolver{
		server: s,
	}
}

func (s *GatewayServer) Query() gatewayGraphQL.QueryResolver {
	return &queryResolver{
		server: s,
	}
}

func (s *GatewayServer) Account() gatewayGraphQL.AccountResolver {
	return &accountResolver{
		server: s,
	}
}

func (s *GatewayServer) ToExecutableSchema() graphql.ExecutableSchema {
	return gatewayGraphQL.NewExecutableSchema(gatewayGraphQL.Config{
		Resolvers: s,
	})
}
