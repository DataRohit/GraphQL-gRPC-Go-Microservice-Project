package account

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountRepository interface {
	Close() error
	CreateAccount(ctx context.Context, email, name string) error
	GetAccountByID(ctx context.Context, id string) (Account, error)
	GetAccountByEmail(ctx context.Context, email string) (Account, error)
	ListAccounts(ctx context.Context, limit, offset uint32) ([]Account, error)
}

type accountRepository struct {
	db *pgxpool.Pool
}

func NewAccountRepository(connString string) (AccountRepository, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database connection string: %w", err)
	}

	config.MaxConns = 25
	config.MaxConnIdleTime = 5 * time.Minute

	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := db.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &accountRepository{db}, nil
}

func (repository *accountRepository) Close() error {
	repository.db.Close()
	return nil
}

func (repository *accountRepository) CreateAccount(ctx context.Context, email, name string) error {
	query := `
        INSERT INTO accounts (email, name)
        VALUES ($1, $2)`
	_, err := repository.db.Exec(ctx, query, email, name)
	if err != nil {
		return fmt.Errorf("failed to create account: %w", err)
	}
	return nil
}

func (repository *accountRepository) GetAccountByID(ctx context.Context, id string) (Account, error) {
	var account Account
	query := "SELECT id, email, name, created_at, updated_at FROM accounts WHERE id = $1"
	err := repository.db.QueryRow(ctx, query, id).Scan(&account.ID, &account.Email, &account.Name, &account.CreatedAt, &account.UpdatedAt)
	if err != nil {
		return Account{}, fmt.Errorf("failed to get account by id: %w", err)
	}
	return account, nil
}

func (repository *accountRepository) GetAccountByEmail(ctx context.Context, email string) (Account, error) {
	var account Account
	query := "SELECT id, email, name, created_at, updated_at FROM accounts WHERE email = $1"
	err := repository.db.QueryRow(ctx, query, email).Scan(&account.ID, &account.Email, &account.Name, &account.CreatedAt, &account.UpdatedAt)
	if err != nil {
		return Account{}, fmt.Errorf("failed to get account by email: %w", err)
	}
	return account, nil
}

func (repository *accountRepository) ListAccounts(ctx context.Context, limit, offset uint32) ([]Account, error) {
	query := `
        SELECT id, email, name, created_at, updated_at
        FROM accounts
        ORDER BY created_at ASC
        LIMIT $1 OFFSET $2`
	rows, err := repository.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list accounts: %w", err)
	}
	defer rows.Close()

	var accounts []Account
	for rows.Next() {
		var account Account
		if err := rows.Scan(&account.ID, &account.Email, &account.Name, &account.CreatedAt, &account.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan account row: %w", err)
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}
