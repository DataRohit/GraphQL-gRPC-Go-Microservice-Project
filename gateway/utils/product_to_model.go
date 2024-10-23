package utils

import (
	"graphql-grpc-go-microservice-project/gateway/models"
	"graphql-grpc-go-microservice-project/product"
)

func ConvertProductToModel(product *product.Product) *models.Product {
	return &models.Product{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
	}
}
