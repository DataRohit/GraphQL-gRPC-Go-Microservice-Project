package main

import (
	"log"
	"time"

	"graphql-grpc-go-microservice-project/account"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type Config struct {
	ACCOUNT_GRPC_SERVER_PORT int    `envconfig:"ACCOUNT_GRPC_SERVER_PORT" default:"8080"`
	ACCOUNT_DATABASE_URL     string `envconfig:"ACCOUNT_DATABASE_URL"`
}

func main() {
	err := godotenv.Load(".envs/.account.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("Error processing environment variables: %v", err)
	}

	var repo account.AccountRepository
	retry.ForeverSleep(2*time.Second, func(_ int) error {
		var err error
		repo, err = account.NewAccountRepository(cfg.ACCOUNT_DATABASE_URL)
		if err != nil {
			log.Printf("Database connection failed: %v", err)
		}
		return err
	})

	defer func() {
		if err := repo.Close(); err != nil {
			log.Printf("Error closing repository: %v", err)
		}
	}()

	log.Println("Initializing account service...")
	service, err := account.NewAccountService(repo)
	if err != nil {
		log.Fatalf("Failed to create account service: %v", err)
	}

	log.Printf("Starting gRPC server on port %d...", cfg.ACCOUNT_GRPC_SERVER_PORT)
	if err := account.ListenGRPC(service, cfg.ACCOUNT_GRPC_SERVER_PORT, false); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
}
