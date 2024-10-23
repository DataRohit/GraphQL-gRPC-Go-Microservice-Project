package account

import (
	"context"
)

type AccountService interface {
	CreateAccount(ctx context.Context, email, name string) (*Account, error)
	GetAccountByID(ctx context.Context, id string) (*Account, error)
	GetAccountByEmail(ctx context.Context, email string) (*Account, error)
	ListAccounts(ctx context.Context, limit, offset uint32) ([]Account, error)
}

type accountService struct {
	repository AccountRepository
}

func NewAccountService(repository AccountRepository) (AccountService, error) {
	return &accountService{repository: repository}, nil
}

func (service *accountService) CreateAccount(ctx context.Context, email, name string) (*Account, error) {
	if err := service.repository.CreateAccount(ctx, email, name); err != nil {
		return nil, err
	}

	var account Account
	account, err := service.repository.GetAccountByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (service *accountService) GetAccountByID(ctx context.Context, id string) (*Account, error) {
	account, err := service.repository.GetAccountByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (service *accountService) GetAccountByEmail(ctx context.Context, email string) (*Account, error) {
	account, err := service.repository.GetAccountByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (service *accountService) ListAccounts(ctx context.Context, limit, offset uint32) ([]Account, error) {
	accounts, err := service.repository.ListAccounts(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}
