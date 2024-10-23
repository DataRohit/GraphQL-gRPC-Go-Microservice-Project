package account

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"graphql-grpc-go-microservice-project/account/protobuf"

	"graphql-grpc-go-microservice-project/common"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type AccountClient struct {
	conn    *grpc.ClientConn
	service protobuf.AccountServiceClient
	logger  *zap.Logger
}

func NewAccountClient(url string, secure bool) (*AccountClient, error) {
	logger := common.GetLogger()

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
		logger.Error("Failed to connect to gRPC server", zap.String("url", url), zap.String("error", err.Error()))
		return nil, err
	}

	logger.Info("Connected to gRPC server", zap.String("url", url))

	client := protobuf.NewAccountServiceClient(conn)

	return &AccountClient{conn: conn, service: client, logger: logger}, nil
}

func (c *AccountClient) Close() error {
	c.logger.Info("Closing gRPC connection")
	return c.conn.Close()
}

func (c *AccountClient) CreateAccount(ctx context.Context, email, name string) (*Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	c.logger.Info("CreateAccount request received", zap.String("email", email), zap.String("name", name))

	r, err := c.service.CreateAccount(ctx, &protobuf.CreateAccountRequest{
		Email: email,
		Name:  name,
	})
	if err != nil {
		c.logger.Error("Failed to create account", zap.String("email", email), zap.String("error", err.Error()))
		return nil, err
	}

	c.logger.Info("Account created successfully", zap.String("account_id", r.GetAccount().GetId()), zap.String("email", email), zap.String("name", name))

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

	c.logger.Info("GetAccountByID request received", zap.String("account_id", id))

	r, err := c.service.GetAccountByID(ctx, &protobuf.GetAccountByIDRequest{Id: id})
	if err != nil {
		c.logger.Error("Failed to fetch account", zap.String("account_id", id), zap.String("error", err.Error()))
		return nil, err
	}

	c.logger.Info("Account fetched successfully", zap.String("account_id", r.GetAccount().GetId()), zap.String("email", r.GetAccount().GetEmail()), zap.String("name", r.GetAccount().GetName()))

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

	c.logger.Info("GetAccountByEmail request received", zap.String("email", email))

	r, err := c.service.GetAccountByEmail(ctx, &protobuf.GetAccountByEmailRequest{Email: email})
	if err != nil {
		c.logger.Error("Failed to fetch account", zap.String("email", email), zap.String("error", err.Error()))
		return nil, err
	}

	c.logger.Info("Account fetched successfully", zap.String("account_id", r.GetAccount().GetId()), zap.String("email", email), zap.String("name", r.GetAccount().GetName()))

	return &Account{
		ID:        uuid.MustParse(r.GetAccount().GetId()),
		Name:      r.GetAccount().GetName(),
		Email:     r.GetAccount().GetEmail(),
		CreatedAt: r.GetAccount().GetCreatedAt().AsTime(),
		UpdatedAt: r.GetAccount().GetUpdatedAt().AsTime(),
	}, nil
}

func (c *AccountClient) ListAccounts(ctx context.Context, limit, offset uint32) ([]Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	c.logger.Info("ListAccounts request received", zap.Uint32("limit", limit), zap.Uint32("offset", offset))

	r, err := c.service.ListAccounts(ctx, &protobuf.ListAccountsRequest{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		c.logger.Error("Failed to list accounts", zap.String("error", err.Error()))
		return nil, err
	}

	c.logger.Info("Accounts listed successfully", zap.Int("account_count", len(r.GetAccounts())))

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
