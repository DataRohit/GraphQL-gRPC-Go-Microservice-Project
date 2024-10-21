.PHONY: docker-up docker-build help

docker-up:
	@docker compose up -d
	@echo "Docker containers started in detached mode."

docker-build:
	@docker compose up -d --build
	@echo "Docker containers built and started in detached mode."

help:
	@echo "Available Commands:"
	@echo "  make docker-up    - Start Docker containers"
	@echo "  make docker-build - Build and start Docker containers"
	@echo "  make help         - Show this help message"
