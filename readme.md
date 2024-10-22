# GraphQL-gRPC-Go-Microservice-Project

![Project Architecture Design](./assets/project-architecture-diagram.png)

[Project Architecture Design - Whimsical](https://whimsical.com/graphql-grpc-go-microservice-JGUJXyUsLacNEHpCCxCpcC)

## Project Overview

This is a microservice-based project that implements an order management system. The project consists of the following main components:

1. **Account Service**: Handles account-related operations such as account creation, retrieval by ID or email, and listing accounts.
2. **Gateway Service**: Provides a GraphQL API as the entry point for client applications to interact with the microservices.

Each service has its own database:

- **Account Service**: Uses PostgreSQL v16.

## Project Structure

The project is organized into the following directories:

- `account/`: Contains the code for the Account Service, including the GraphQL API implementation.
- `gateway/`: Contains the code for the Gateway Service, which provides the GraphQL API.

## GraphQL API Usage

The GraphQL API for the Account Service is documented in the [account/README.md](./account/README.md) file.

## Technologies Used

- Go programming language
- Docker for containerization
- PostgreSQL v16 for the Account Service database
- GraphQL for the API
- gRPC for inter-service communication
- Protocol Buffers (protobuf) for defining the service interfaces

## Getting Started

To run the project locally, you'll need to have the following installed:

- Go
- Docker

1. Clone the repository: `git clone https://github.com/datarohit/GraphQL-gRPC-Go-Microservice-Project.git`
2. Navigate to the project directory: `cd GraphQL-gRPC-Go-Microservice-Project`
3. Start the services using Docker Compose: `docker-compose up`

- The GraphQL API will be available at `http://localhost:8080/graphql`.
- The GraphQL API playground will be available at `http://localhost:8080/playground`.

## Contributing

If you'd like to contribute to this project, please follow the standard GitHub workflow:

1. Fork the repository
2. Create a new branch: `git checkout -b feature/your-feature-name`
3. Make your changes and commit them: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin feature/your-feature-name`
5. Create a new Pull Request
