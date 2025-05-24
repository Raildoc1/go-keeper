package services

import (
	"errors"
	"fmt"
	"go-keeper/internal/client/logic/requester"
	"go-keeper/internal/common/protocol"
	"net/http"
	"strings"
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
		return "", fmt.Errorf("post request failed: %w", err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		break
	case http.StatusBadRequest:
		return "", ErrInvalidInput
	default:
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	tknHeader := resp.Header().Get("Authorization")
	tkn, _ = strings.CutPrefix(tknHeader, "Bearer ")
	return tkn, nil
}

func (s *AuthService) Login(username, password string) (tkn string, err error) {
	resp, err := s.req.Post("/api/user/login", protocol.Creds{
		Login:    username,
		Password: password,
	})

	if err != nil {
		return "", fmt.Errorf("post request failed: %w", err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		break
	case http.StatusBadRequest:
		return "", ErrInvalidInput
	default:
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	tknHeader := resp.Header().Get("Authorization")
	tkn, _ = strings.CutPrefix(tknHeader, "Bearer ")
	return tkn, nil
}
