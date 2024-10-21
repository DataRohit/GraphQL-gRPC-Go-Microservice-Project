package account

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"graphql-grpc-go-microservice-project/account/protobuf"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type accountGrpcServer struct {
	protobuf.UnimplementedAccountServiceServer
	service AccountService
}

func ListenGRPC(s AccountService, port int, secure bool) error {
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
	protobuf.RegisterAccountServiceServer(serv, &accountGrpcServer{
		UnimplementedAccountServiceServer: protobuf.UnimplementedAccountServiceServer{},
		service:                           s,
	})
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
		fmt.Printf("Received signal: %s. Shutting down gRPC server...\n", sig)
		serv.GracefulStop()
	case err := <-errChan:
		return err
	}

	return nil
}

func (s *accountGrpcServer) CreateAccount(ctx context.Context, r *protobuf.CreateAccountRequest) (*protobuf.CreateAccountResponse, error) {
	a, err := s.service.CreateAccount(ctx, r.Email, r.Name)
	if err != nil {
		log.Printf("failed to create account: %v", err)

		return &protobuf.CreateAccountResponse{
			Result: &protobuf.CreateAccountResponse_Error{Error: err.Error()},
		}, nil
	}

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
	a, err := s.service.GetAccountByID(ctx, r.Id)
	if err != nil {
		return &protobuf.GetAccountByIDResponse{
			Result: &protobuf.GetAccountByIDResponse_Error{Error: err.Error()},
		}, nil
	}

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
	a, err := s.service.GetAccountByEmail(ctx, r.Email)
	if err != nil {
		return &protobuf.GetAccountByEmailResponse{
			Result: &protobuf.GetAccountByEmailResponse_Error{Error: err.Error()},
		}, nil
	}

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
	accountsList, err := s.service.ListAccounts(ctx, int(r.Limit), int(r.Offset))
	if err != nil {
		return &protobuf.ListAccountsResponse{
			Error: err.Error(),
		}, nil
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

	return &protobuf.ListAccountsResponse{Accounts: accounts}, nil
}
