package product

import (
	"context"
	"crypto/tls"
	"fmt"
	"graphql-grpc-go-microservice-project/common"
	"graphql-grpc-go-microservice-project/product/protobuf"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	grpcResponseCodes "google.golang.org/grpc/codes"
	grpcResponseStatus "google.golang.org/grpc/status"
)

type productGrpcServer struct {
	protobuf.UnimplementedProductServiceServer
	service ProductService
	logger  *zap.Logger
}

func ListenGRPC(s ProductService, port int, secure bool) error {
	logger := common.GetLogger()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen on port %d: %v", port, err)
	}

	var opts []grpc.ServerOption
	keepAliveParams := keepalive.ServerParameters{
		Time:    5 * time.Minute,
		Timeout: 20 * time.Second,
	}

	opts = append(opts, grpc.KeepaliveParams(keepAliveParams))

	if secure {
		creds := credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: false,
		})
		opts = append(opts, grpc.Creds(creds))
	} else {
		opts = append(opts, grpc.Creds(insecure.NewCredentials()))
	}

	serv := grpc.NewServer(opts...)
	productServer := &productGrpcServer{
		UnimplementedProductServiceServer: protobuf.UnimplementedProductServiceServer{},
		service:                           s,
		logger:                            logger,
	}
	protobuf.RegisterProductServiceServer(serv, productServer)
	reflection.Register(serv)

	errChan := make(chan error)
	go func() {
		if err := serv.Serve(lis); err != nil {
			errChan <- fmt.Errorf("failed to serve gRPC server: %v", err)
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-signalChan:
		logger.Info("Received signal, shutting down gRPC server", zap.String("signal", sig.String()))
		serv.GracefulStop()
	case err := <-errChan:
		return err
	}

	return nil
}

func (s *productGrpcServer) CreateProduct(ctx context.Context, r *protobuf.CreateProductRequest) (*protobuf.CreateProductResponse, error) {
	s.logger.Info("CreateProduct request received", zap.String("name", r.Name), zap.String("description", r.Description), zap.Float64("price", r.Price))

	p, err := s.service.CreateProduct(ctx, r.Name, r.Description, r.Price)
	if err != nil {
		s.logger.Error("Failed to create product", zap.String("name", r.Name), zap.String("error", err.Error()))
		return &protobuf.CreateProductResponse{
			Result: &protobuf.CreateProductResponse_Error{Error: err.Error()},
		}, grpcResponseStatus.Errorf(grpcResponseCodes.Internal, err.Error())
	}

	s.logger.Info("Product created successfully", zap.String("name", r.Name), zap.String("description", r.Description), zap.Float64("price", r.Price))

	return &protobuf.CreateProductResponse{
		Result: &protobuf.CreateProductResponse_Product{
			Product: &protobuf.Product{
				Id:          p.ID,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
			},
		},
	}, nil
}

func (s *productGrpcServer) GetProductByID(ctx context.Context, r *protobuf.GetProductByIDRequest) (*protobuf.GetProductByIDResponse, error) {
	s.logger.Info("GetProductByID request received", zap.String("product_id", r.Id))

	p, err := s.service.GetProductByID(ctx, r.Id)
	if err != nil {
		s.logger.Error("Failed to fetch product by ID", zap.String("product_id", r.Id), zap.String("error", err.Error()))
		return &protobuf.GetProductByIDResponse{
			Result: &protobuf.GetProductByIDResponse_Error{Error: err.Error()},
		}, grpcResponseStatus.Errorf(grpcResponseCodes.Internal, err.Error())
	}

	s.logger.Info("Product fetched successfully", zap.String("name", p.ID), zap.String("name", p.Name), zap.String("description", p.Description), zap.Float64("price", p.Price))

	return &protobuf.GetProductByIDResponse{
		Result: &protobuf.GetProductByIDResponse_Product{
			Product: &protobuf.Product{
				Id:          p.ID,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
			},
		},
	}, nil
}

func (s *productGrpcServer) ListProducts(ctx context.Context, r *protobuf.ListProductsRequest) (*protobuf.ListProductsResponse, error) {
	s.logger.Info("ListProducts request received", zap.Uint32("offset", r.Offset), zap.Uint32("limit", r.Limit))

	p, err := s.service.ListProducts(ctx, r.Offset, r.Limit)
	if err != nil {
		s.logger.Error("Failed to list products", zap.String("error", err.Error()))
		return &protobuf.ListProductsResponse{
			Error: err.Error(),
		}, grpcResponseStatus.Errorf(grpcResponseCodes.Internal, err.Error())
	}

	var products []*protobuf.Product
	for _, p := range p {
		products = append(products, &protobuf.Product{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		})
	}

	s.logger.Info("Products listed successfully", zap.Int("count", len(p)))
	return &protobuf.ListProductsResponse{Products: products}, nil
}

func (s *productGrpcServer) ListProductsWithIDs(ctx context.Context, r *protobuf.ListProductsWithIDsRequest) (*protobuf.ListProductsWithIDsResponse, error) {
	s.logger.Info("ListProductsWithIDs request received", zap.Strings("ids", r.Ids))

	p, err := s.service.ListProductsWithIDs(ctx, r.Ids)
	if err != nil {
		s.logger.Error("Failed to list products with IDs", zap.Strings("ids", r.Ids), zap.String("error", err.Error()))
		return &protobuf.ListProductsWithIDsResponse{
			Error: err.Error(),
		}, grpcResponseStatus.Errorf(grpcResponseCodes.Internal, err.Error())
	}

	var products []*protobuf.Product
	for _, p := range p {
		products = append(products, &protobuf.Product{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		})
	}

	s.logger.Info("Products listed successfully", zap.Int("count", len(p)))
	return &protobuf.ListProductsWithIDsResponse{Products: products}, nil
}

func (s *productGrpcServer) SearchProducts(ctx context.Context, r *protobuf.SearchProductsRequest) (*protobuf.SearchProductsResponse, error) {
	s.logger.Info("SearchProducts request received", zap.String("query", r.Query), zap.Uint32("offset", r.Offset), zap.Uint32("limit", r.Limit))

	p, err := s.service.SearchProducts(ctx, r.Query, r.Offset, r.Limit)
	if err != nil {
		s.logger.Error("Failed to search products", zap.String("query", r.Query), zap.String("error", err.Error()))
		return &protobuf.SearchProductsResponse{
			Error: err.Error(),
		}, grpcResponseStatus.Errorf(grpcResponseCodes.Internal, err.Error())
	}

	var products []*protobuf.Product
	for _, p := range p {
		products = append(products, &protobuf.Product{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		})
	}

	s.logger.Info("Products searched successfully", zap.Int("count", len(p)))
	return &protobuf.SearchProductsResponse{Products: products}, nil
}
