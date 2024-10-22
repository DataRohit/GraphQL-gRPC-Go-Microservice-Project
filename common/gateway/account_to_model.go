package utils

import (
	"graphql-grpc-go-microservice-project/account"
	"graphql-grpc-go-microservice-project/gateway/models"
	"time"
)

func ConvertAccountToModel(account *account.Account) *models.Account {
	return &models.Account{
		ID:        account.ID,
		Email:     account.Email,
		Name:      account.Name,
		CreatedAt: account.CreatedAt.Format(time.RFC3339),
		UpdatedAt: account.UpdatedAt.Format(time.RFC3339),
	}
}
