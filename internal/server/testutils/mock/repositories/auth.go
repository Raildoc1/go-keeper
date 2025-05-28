package repositories

import (
	"context"
	"go-keeper/internal/server/data"
	"go-keeper/internal/server/dto"
	"go-keeper/internal/server/services"
)

var _ services.AuthRepository = (*AuthRepositoryMock)(nil)

type AuthRepositoryMock struct {
	users []dto.Creds
}

func NewAuthRepositoryMock() *AuthRepositoryMock {
	return &AuthRepositoryMock{
		users: make([]dto.Creds, 0),
	}
}

func (t *AuthRepositoryMock) InsertUser(ctx context.Context, creds dto.Creds) (userID int, err error) {
	for _, user := range t.users {
		if user.Username == creds.Username {
			return -1, data.ErrUniqueConstraintViolation
		}
	}

	t.users = append(t.users, creds)
	return len(t.users) - 1, nil
}

func (t *AuthRepositoryMock) ValidateUser(ctx context.Context, creds dto.Creds) (userID int, err error) {
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
