package services

import (
	"context"
	"go-keeper/internal/server/dto"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) Register(ctx context.Context, creds dto.Creds) (token string, err error) {
	//TODO implement me
	panic("implement me")
}

func (s *AuthService) Login(ctx context.Context, creds dto.Creds) (token string, err error) {
	//TODO implement me
	panic("implement me")
}
