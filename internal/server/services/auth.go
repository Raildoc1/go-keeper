package services

import (
	"context"
	"errors"
	"fmt"
	"go-keeper/internal/server/data"
	"go-keeper/internal/server/dto"
	"strconv"
)

var (
	ErrLoginTaken         = errors.New("login is already taken")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

var (
	UserIDClaimName = "user_id"
)

type AuthRepository interface {
	InsertUser(ctx context.Context, login, password string) (userID int, err error)
	ValidateUser(ctx context.Context, login, password string) (userID int, err error)
}

type TokenFactory interface {
	Generate(extraClaims map[string]string) (string, error)
}

type AuthService struct {
	repository   AuthRepository
	tokenFactory TokenFactory
}

func NewAuthService(repository AuthRepository, tokenFactory TokenFactory) *AuthService {
	return &AuthService{
		repository:   repository,
		tokenFactory: tokenFactory,
	}
}

func (s *AuthService) Register(ctx context.Context, creds dto.Creds) (string, error) {
	userID, err := s.repository.InsertUser(ctx, creds.Username, creds.Password)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrUniqueConstraintViolation):
			return "", ErrLoginTaken
		default:
			return "", fmt.Errorf("error inserting user: %w", err)
		}
	}

	payload := map[string]string{
		UserIDClaimName: strconv.Itoa(userID),
	}
	tkn, err := s.tokenFactory.Generate(payload)
	if err != nil {
		return "", fmt.Errorf("error generating token: %w", err)
	}

	return tkn, nil
}

func (s *AuthService) Login(ctx context.Context, creds dto.Creds) (string, error) {
	userID, err := s.repository.ValidateUser(ctx, creds.Username, creds.Password)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrInvalidPassword):
			return "", ErrInvalidCredentials
		case errors.Is(err, data.ErrInvalidLogin):
			return "", ErrInvalidCredentials
		default:
			return "", fmt.Errorf("error inserting user: %w", err)
		}
	}

	payload := map[string]string{
		UserIDClaimName: strconv.Itoa(userID),
	}
	tkn, err := s.tokenFactory.Generate(payload)
	if err != nil {
		return "", fmt.Errorf("error generating token: %w", err)
	}

	return tkn, nil
}
