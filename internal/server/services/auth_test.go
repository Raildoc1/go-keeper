package services

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go-keeper/internal/server/data"
	"go-keeper/internal/server/dto"
	"testing"
)

var _ AuthRepository = (*TestAuthRepository)(nil)

type TestAuthRepository struct {
	users []dto.Creds
}

func NewTestAuthRepository() *TestAuthRepository {
	return &TestAuthRepository{
		users: make([]dto.Creds, 0),
	}
}

func (t *TestAuthRepository) InsertUser(ctx context.Context, creds dto.Creds) (userID int, err error) {
	for _, user := range t.users {
		if user.Username == creds.Username {
			return -1, data.ErrUniqueConstraintViolation
		}
	}

	t.users = append(t.users, creds)
	return len(t.users) - 1, nil
}

func (t *TestAuthRepository) ValidateUser(ctx context.Context, creds dto.Creds) (userID int, err error) {
	for userId, user := range t.users {
		if user.Username == creds.Username {
			if user.Password == creds.Password {
				return userId, nil
			} else {
				return -1, data.ErrInvalidPassword
			}
		}
	}
	return -1, data.ErrInvalidLogin
}

var _ TokenFactory = (*TestTokenFactory)(nil)

type TestTokenFactory struct{}

func NewTestTokenFactory() *TestTokenFactory {
	return &TestTokenFactory{}
}

func (t *TestTokenFactory) Generate(extraClaims map[string]string) (string, error) {
	return "test_token", nil
}

func TestAuthService(t *testing.T) {
	authRepository := NewTestAuthRepository()
	tokenFactory := NewTestTokenFactory()

	authService := NewAuthService(authRepository, tokenFactory)

	t.Run("Register Success", func(t *testing.T) {
		_, err := authService.Register(
			context.Background(),
			dto.Creds{
				Username: "test_user",
				Password: "test_pwd",
			},
		)
		assert.NoError(t, err)
	})

	t.Run("Register Same User", func(t *testing.T) {
		_, err := authService.Register(
			context.Background(),
			dto.Creds{
				Username: "test_user",
				Password: "test_pwd",
			},
		)
		assert.ErrorIs(t, err, ErrLoginTaken)
	})

	t.Run("Login Success", func(t *testing.T) {
		_, err := authService.Login(
			context.Background(),
			dto.Creds{
				Username: "test_user",
				Password: "test_pwd",
			},
		)
		assert.NoError(t, err)
	})

	t.Run("Login Non-Existent User", func(t *testing.T) {
		_, err := authService.Login(
			context.Background(),
			dto.Creds{
				Username: "wrong_user",
				Password: "test_pwd",
			},
		)
		assert.ErrorIs(t, err, ErrInvalidCredentials)
	})

	t.Run("Login Wrong Password", func(t *testing.T) {
		_, err := authService.Login(
			context.Background(),
			dto.Creds{
				Username: "test_user",
				Password: "wrong_pwd",
			},
		)
		assert.ErrorIs(t, err, ErrInvalidCredentials)
	})
}
