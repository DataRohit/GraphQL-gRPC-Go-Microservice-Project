package main

import (
	"log"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	ACCOUNT_SERVICE_URL string `envconfig:"ACCOUNT_SERVICE_URL" required:"true"`
	PRODUCT_SERVICE_URL string `envconfig:"PRODUCT_SERVICE_URL" required:"true"`
	PORT                string `envconfig:"PORT" default:"8080"`
}

func main() {
	var cfg AppConfig
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}

	server, err := NewGraphQLServer(cfg.ACCOUNT_SERVICE_URL, cfg.PRODUCT_SERVICE_URL, false)
	if err != nil {
		log.Fatalf("Failed to create GraphQL server: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/graphql", handler.NewDefaultServer(server.ToExecutableSchema()))
	mux.Handle("/playground", playground.Handler("GraphQL Playground", "/graphql"))

	srv := &http.Server{
		Addr:         ":" + cfg.PORT,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Starting server on port %s", cfg.PORT)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}
}
