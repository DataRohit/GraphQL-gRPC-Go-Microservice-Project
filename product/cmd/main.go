package main

import (
	"graphql-grpc-go-microservice-project/product"
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type Config struct {
	PRODUCT_GRPC_SERVER_PORT int    `envconfig:"PRODUCT_GRPC_SERVER_PORT" default:"8080"`
	PRODUCT_DATABASE_URL     string `envconfig:"PRODUCT_DATABASE_URL"`
}

func main() {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("Error processing environment variables: %v", err)
	}

	var repo product.ProductRepository
	retry.ForeverSleep(2*time.Second, func(_ int) error {
		var err error
		repo, err = product.NewElasticRepository(cfg.PRODUCT_DATABASE_URL)
		if err != nil {
			log.Printf("Database connection failed: %v", err)
		}
		return err
	})
	defer repo.Close()

	log.Println("Initializing product service...")
	service, err := product.NewProductService(repo)
	if err != nil {
		log.Fatalf("Failed to create product service: %v", err)
	}

	log.Printf("Starting gRPC server on port %d...", cfg.PRODUCT_GRPC_SERVER_PORT)
	if err := product.ListenGRPC(service, cfg.PRODUCT_GRPC_SERVER_PORT, false); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
}
