package account

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"graphql-grpc-go-microservice-project/account/protobuf"
	"graphql-grpc-go-microservice-project/common"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"

	"go.uber.org/zap"
	grpcResponseCodes "google.golang.org/grpc/codes"
	grpcResponseStatus "google.golang.org/grpc/status"
)

type accountGrpcServer struct {
	protobuf.UnimplementedAccountServiceServer
	service AccountService
	logger  *zap.Logger
}

func ListenGRPC(s AccountService, port int, secure bool) error {
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
	accountServer := &accountGrpcServer{
		UnimplementedAccountServiceServer: protobuf.UnimplementedAccountServiceServer{},
		service:                           s,
		logger:                            logger,
	}
	protobuf.RegisterAccountServiceServer(serv, accountServer)
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

func (s *accountGrpcServer) CreateAccount(ctx context.Context, r *protobuf.CreateAccountRequest) (*protobuf.CreateAccountResponse, error) {
	s.logger.Info("CreateAccount request received", zap.String("email", r.Email), zap.String("name", r.Name))

	a, err := s.service.CreateAccount(ctx, r.Email, r.Name)
	if err != nil {
		s.logger.Error("Failed to create account", zap.String("email", r.Email), zap.String("error", err.Error()))
		return &protobuf.CreateAccountResponse{
			Result: &protobuf.CreateAccountResponse_Error{Error: err.Error()},
		}, grpcResponseStatus.Errorf(grpcResponseCodes.Internal, err.Error())
	}

	s.logger.Info("Account created successfully", zap.String("account_id", a.ID.String()), zap.String("email", a.Email), zap.String("name", a.Name))

	return &protobuf.CreateAccountResponse{
		Result: &protobuf.CreateAccountResponse_Account{
			Account: &protobuf.Account{
				Id:        a.ID.String(),
				Email:     a.Email,
				Name:      a.Name,
				CreatedAt: timestamppb.New(a.CreatedAt),
				UpdatedAt: timestamppb.New(a.UpdatedAt),
			},
		},
	}, nil
}

func (s *accountGrpcServer) GetAccountByID(ctx context.Context, r *protobuf.GetAccountByIDRequest) (*protobuf.GetAccountByIDResponse, error) {
	s.logger.Info("GetAccountByID request received", zap.String("account_id", r.Id))

	a, err := s.service.GetAccountByID(ctx, r.Id)
	if err != nil {
		s.logger.Error("Failed to fetch account by ID", zap.String("account_id", r.Id), zap.String("error", err.Error()))
		return &protobuf.GetAccountByIDResponse{
			Result: &protobuf.GetAccountByIDResponse_Error{Error: err.Error()},
		}, grpcResponseStatus.Errorf(grpcResponseCodes.NotFound, err.Error())
	}

	s.logger.Info("Account fetched successfully", zap.String("account_id", a.ID.String()), zap.String("email", a.Email), zap.String("name", a.Name))

	return &protobuf.GetAccountByIDResponse{
		Result: &protobuf.GetAccountByIDResponse_Account{
			Account: &protobuf.Account{
				Id:        a.ID.String(),
				Email:     a.Email,
				Name:      a.Name,
				CreatedAt: timestamppb.New(a.CreatedAt),
				UpdatedAt: timestamppb.New(a.UpdatedAt),
			},
		},
	}, nil
}

func (s *accountGrpcServer) GetAccountByEmail(ctx context.Context, r *protobuf.GetAccountByEmailRequest) (*protobuf.GetAccountByEmailResponse, error) {
	s.logger.Info("GetAccountByEmail request received", zap.String("email", r.Email))

	a, err := s.service.GetAccountByEmail(ctx, r.Email)
	if err != nil {
		s.logger.Error("Failed to fetch account by email", zap.String("email", r.Email), zap.String("error", err.Error()))
		return &protobuf.GetAccountByEmailResponse{
			Result: &protobuf.GetAccountByEmailResponse_Error{Error: err.Error()},
		}, grpcResponseStatus.Errorf(grpcResponseCodes.NotFound, err.Error())
	}

	s.logger.Info("Account fetched successfully", zap.String("account_id", a.ID.String()), zap.String("email", r.Email), zap.String("name", a.Name))

	return &protobuf.GetAccountByEmailResponse{
		Result: &protobuf.GetAccountByEmailResponse_Account{
			Account: &protobuf.Account{
				Id:        a.ID.String(),
				Email:     a.Email,
				Name:      a.Name,
				CreatedAt: timestamppb.New(a.CreatedAt),
				UpdatedAt: timestamppb.New(a.UpdatedAt),
			},
		},
	}, nil
}

func (s *accountGrpcServer) ListAccounts(ctx context.Context, r *protobuf.ListAccountsRequest) (*protobuf.ListAccountsResponse, error) {
	s.logger.Info("ListAccounts request received", zap.Uint32("limit", r.Limit), zap.Uint32("offset", r.Offset))

	accountsList, err := s.service.ListAccounts(ctx, r.Limit, r.Offset)
	if err != nil {
		s.logger.Error("Failed to list accounts", zap.String("error", err.Error()))
		return &protobuf.ListAccountsResponse{
			Error: err.Error(),
		}, grpcResponseStatus.Errorf(grpcResponseCodes.Internal, err.Error())
	}

	var accounts []*protobuf.Account
	for _, a := range accountsList {
		accounts = append(accounts, &protobuf.Account{
			Id:        a.ID.String(),
			Email:     a.Email,
			Name:      a.Name,
			CreatedAt: timestamppb.New(a.CreatedAt),
			UpdatedAt: timestamppb.New(a.UpdatedAt),
		})
	}

	s.logger.Info("Accounts listed successfully", zap.Int("account_count", len(accounts)))
	return &protobuf.ListAccountsResponse{Accounts: accounts}, nil
}
