package repository

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"go-keeper/internal/server/data"
	"go-keeper/internal/server/dto"
	"go-keeper/pkg/logging"
)

const (
	invalidUserID = -1
)

type AuthRepository struct {
	storage DBStorage
	logger  *logging.ZapLogger
}

func NewAuthRepository(storage DBStorage, logger *logging.ZapLogger) *AuthRepository {
	return &AuthRepository{
		storage: storage,
		logger:  logger,
	}
}

//go:embed sql/auth/insert_user.sql
var insertUserQuery string

func (db *AuthRepository) InsertUser(ctx context.Context, creds dto.Creds) (userID int, err error) {
	err = db.storage.QueryValue(ctx, insertUserQuery, []any{creds.Username, creds.Password}, []any{&userID})
	if err != nil {
		return invalidUserID, handleSQLError(err)
	}
	return userID, nil
}

//go:embed sql/auth/validate_user.sql
var validateUserQuery string

func (db *AuthRepository) ValidateUser(ctx context.Context, creds dto.Creds) (userID int, err error) {
	result := struct {
		userID          int
		passwordMatches bool
	}{}
	err = db.storage.QueryValue(
		ctx,
		validateUserQuery,
		[]any{creds.Username, creds.Password},
		[]any{&result.userID, &result.passwordMatches},
	)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return invalidUserID, data.ErrInvalidLogin
		default:
			return invalidUserID, fmt.Errorf("failed to validate user: %w", err)
		}
	}
	if !result.passwordMatches {
		return invalidUserID, data.ErrInvalidPassword
	}
	return result.userID, nil
}
