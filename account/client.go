package account

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/google/uuid"

	"graphql-grpc-go-microservice-project/account/protobuf"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type AccountClient struct {
	conn    *grpc.ClientConn
	service protobuf.AccountServiceClient
}

func NewAccountClient(url string, secure bool) (*AccountClient, error) {
	var opts []grpc.DialOption
	if secure {
		creds := credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: false,
		})
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.NewClient(url, opts...)
	if err != nil {
		return nil, err
	}

	client := protobuf.NewAccountServiceClient(conn)
	return &AccountClient{conn: conn, service: client}, nil
}

func (c *AccountClient) Close() error {
	return c.conn.Close()
}

func (c *AccountClient) CreateAccount(ctx context.Context, email, name string) (*Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	r, err := c.service.CreateAccount(ctx, &protobuf.CreateAccountRequest{
		Email: email,
		Name:  name,
	})
	if err != nil {
		return nil, err
	}

	return &Account{
		ID:        uuid.MustParse(r.GetAccount().GetId()),
		Name:      r.GetAccount().GetName(),
		Email:     r.GetAccount().GetEmail(),
		CreatedAt: r.GetAccount().GetCreatedAt().AsTime(),
		UpdatedAt: r.GetAccount().GetUpdatedAt().AsTime(),
	}, nil
}

func (c *AccountClient) GetAccountByID(ctx context.Context, id string) (*Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	r, err := c.service.GetAccountByID(ctx, &protobuf.GetAccountByIDRequest{Id: id})
	if err != nil {
		return nil, err
	}

	return &Account{
		ID:        uuid.MustParse(r.GetAccount().GetId()),
		Name:      r.GetAccount().GetName(),
		Email:     r.GetAccount().GetEmail(),
		CreatedAt: r.GetAccount().GetCreatedAt().AsTime(),
		UpdatedAt: r.GetAccount().GetUpdatedAt().AsTime(),
	}, nil
}

func (c *AccountClient) GetAccountByEmail(ctx context.Context, email string) (*Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	r, err := c.service.GetAccountByEmail(ctx, &protobuf.GetAccountByEmailRequest{Email: email})
	if err != nil {
		return nil, err
	}

	return &Account{
		ID:        uuid.MustParse(r.GetAccount().GetId()),
		Name:      r.GetAccount().GetName(),
		Email:     r.GetAccount().GetEmail(),
		CreatedAt: r.GetAccount().GetCreatedAt().AsTime(),
		UpdatedAt: r.GetAccount().GetUpdatedAt().AsTime(),
	}, nil
}

func (c *AccountClient) ListAccounts(ctx context.Context, limit, offset int32) ([]Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	r, err := c.service.ListAccounts(ctx, &protobuf.ListAccountsRequest{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	var accounts []Account
	for _, acc := range r.GetAccounts() {
		accounts = append(accounts, Account{
			ID:        uuid.MustParse(acc.GetId()),
			Name:      acc.GetName(),
			Email:     acc.GetEmail(),
			CreatedAt: acc.GetCreatedAt().AsTime(),
			UpdatedAt: acc.GetUpdatedAt().AsTime(),
		})
	}

	return accounts, nil
}
