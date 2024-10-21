PROTOC_GEN_GO_OPTS=--go_out=. --go_opt=paths=source_relative
PROTOC_GEN_GO_GRPC_OPTS=--go-grpc_out=. --go-grpc_opt=paths=source_relative

.PHONY: proto-gen docker-up docker-build help

proto-gen:
	@protoc $(PROTOC_GEN_GO_OPTS) $(PROTOC_GEN_GO_GRPC_OPTS) protobuf/account.proto
	@echo "Protobuf files generate successfully."

docker-up:
	@docker compose up -d
	@echo "Docker containers started in detached mode."

docker-build:
	@docker compose up -d --build
	@echo "Docker containers built and started in detached mode."

help:
	@echo "Available Commands:"
	@echo "  make proto-gen    - Generate protobuf files"
	@echo "  make docker-up    - Start Docker containers"
	@echo "  make docker-build - Build and start Docker containers"
	@echo "  make help         - Show this help message"
