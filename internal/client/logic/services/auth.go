package services

import (
	"errors"
	"fmt"
	"go-keeper/internal/client/logic/requester"
	"go-keeper/internal/common/protocol"
)

var (
	ErrInvalidInput = errors.New("invalid input")
)

type AuthService struct {
	req *requester.Requester
}

func NewAuthService(req *requester.Requester) *AuthService {
	return &AuthService{
		req: req,
	}
}

func (s *AuthService) Register(username, password string) (tkn string, err error) {
	resp, err := s.req.Post("/api/user/register", protocol.Creds{
		Login:    username,
		Password: password,
	})

	if err != nil {
		if errors.Is(err, requester.ErrBadRequest) {
			return "", ErrInvalidInput
		}
		return "", fmt.Errorf("post request failed: %w", err)
	}

	tkn = resp.Header().Get("Authorization")
	return tkn, nil
}

func (s *AuthService) Login(username, password string) (tkn string, err error) {
	resp, err := s.req.Post("/api/user/login", protocol.Creds{
		Login:    username,
		Password: password,
	})

	if err != nil {
		if errors.Is(err, requester.ErrBadRequest) {
			return "", ErrInvalidInput
		}
		return "", fmt.Errorf("post request failed: %w", err)
	}

	tkn = resp.Header().Get("Authorization")
	return tkn, nil
}
