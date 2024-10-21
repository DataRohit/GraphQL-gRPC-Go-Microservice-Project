PROTOC_GEN_GO_OPTS=--go_out=. --go_opt=paths=source_relative
PROTOC_GEN_GO_GRPC_OPTS=--go-grpc_out=. --go-grpc_opt=paths=source_relative

.PHONY: proto-gen help

proto-gen:
	@protoc $(PROTOC_GEN_GO_OPTS) $(PROTOC_GEN_GO_GRPC_OPTS) protobuf/account.proto
	@echo "Protobuf files compiled successfully."

help:
	@echo "Available Commands:"
	@echo "  make proto-gen    - Compile protobuf files"
	@echo "  make help         - Show this help message"
