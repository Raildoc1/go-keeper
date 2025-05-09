package server

import (
	"fmt"
	"go-keeper/internal/server/grpcservers"
	pb "go-keeper/proto"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"net"
)

type GRPCConfig struct {
	Port uint16
}

type Server struct {
	cfg     GRPCConfig
	auth    *grpc.Server
	storage *grpc.Server
}

func NewServer(cfg GRPCConfig) *Server {
	authServer := grpc.NewServer()
	pb.RegisterAuthServer(authServer, grpcservers.NewAuthServer())

	storageServer := grpc.NewServer()
	pb.RegisterStorageServer(storageServer, grpcservers.NewStorageServer())

	return &Server{
		cfg:     cfg,
		auth:    authServer,
		storage: storageServer,
	}
}

func (s *Server) Run() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", s.cfg.Port))
	if err != nil {
		return fmt.Errorf("failed to start listen: %w", err)
	}

	g := &errgroup.Group{}

	g.Go(func() error {
		if err := s.auth.Serve(listener); err != nil {
			return fmt.Errorf("auth server error: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		if err := s.storage.Serve(listener); err != nil {
			return fmt.Errorf("storage server error: %w", err)
		}
		return nil
	})

	return g.Wait()
}

func (s *Server) Stop() error {
	g := &errgroup.Group{}

	g.Go(func() error {
		s.auth.GracefulStop()
		return nil
	})

	g.Go(func() error {
		s.storage.GracefulStop()
		return nil
	})

	return g.Wait()
}
