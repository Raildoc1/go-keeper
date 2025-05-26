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
	ErrInvalidCreds = errors.New("invalid creds")
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

	err = s.statusCodeToError(resp.StatusCode())
	if err != nil {
		return "", fmt.Errorf("register request failed: %w", err)
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

	err = s.statusCodeToError(resp.StatusCode())
	if err != nil {
		return "", fmt.Errorf("login request failed: %w", err)
	}

	tknHeader := resp.Header().Get("Authorization")
	tkn, _ = strings.CutPrefix(tknHeader, "Bearer ")
	return tkn, nil
}

func (s *AuthService) statusCodeToError(statusCode int) error {
	switch statusCode {
	case http.StatusOK:
		return nil
	case http.StatusBadRequest:
		return ErrInvalidCreds
	case http.StatusUnauthorized:
		return ErrInvalidCreds
	default:
		return fmt.Errorf("unexpected status code: %v", statusCode)
	}
}
