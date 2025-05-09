package grpcservers

import (
	"context"
	"go-keeper/internal/server/dto"
	pb "go-keeper/proto"
)

var _ pb.AuthServer = (*AuthServer)(nil)

type AuthService interface {
	Register(ctx context.Context, creds dto.Creds) (token string, err error)
	Login(ctx context.Context, creds dto.Creds) (token string, err error)
}

type AuthServer struct {
	pb.UnimplementedAuthServer
	service AuthService
}

func NewAuthServer(service AuthService) *AuthServer {
	return &AuthServer{
		service: service,
	}
}

func (s *AuthServer) Register(ctx context.Context, request *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	creds := dto.Creds{
		Username: request.GetLogin(),
		Password: request.GetPassword(),
	}
	token, err := s.service.Register(ctx, creds)
	if err != nil {
		errMsg := err.Error()
		return pb.RegisterResponse_builder{
			Error: &errMsg,
		}.Build(), nil
	}
	return pb.RegisterResponse_builder{
		Token: &token,
	}.Build(), nil
}

func (s *AuthServer) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {
	creds := dto.Creds{
		Username: request.GetLogin(),
		Password: request.GetPassword(),
	}
	token, err := s.service.Login(ctx, creds)
	if err != nil {
		errMsg := err.Error()
		return pb.LoginResponse_builder{
			Error: &errMsg,
		}.Build(), nil
	}
	return pb.LoginResponse_builder{
		Token: &token,
	}.Build(), nil
}
