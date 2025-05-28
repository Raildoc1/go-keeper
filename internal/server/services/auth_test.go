package services

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go-keeper/internal/server/dto"
	"go-keeper/internal/server/testutils/mock/repositories"
	"testing"
)

func TestAuthService(t *testing.T) {
	authRepository := repositories.NewAuthRepositoryMock()
	tokenFactory := repositories.NewTokenFactoryMock()

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

	t.Run("Username Success", func(t *testing.T) {
		_, err := authService.Login(
			context.Background(),
			dto.Creds{
				Username: "test_user",
				Password: "test_pwd",
			},
		)
		assert.NoError(t, err)
	})

	t.Run("Username Non-Existent User", func(t *testing.T) {
		_, err := authService.Login(
			context.Background(),
			dto.Creds{
				Username: "wrong_user",
				Password: "test_pwd",
			},
		)
		assert.ErrorIs(t, err, ErrInvalidCredentials)
	})

	t.Run("Username Wrong Password", func(t *testing.T) {
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
