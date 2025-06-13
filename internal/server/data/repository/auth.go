package repository

import (
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
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

func (db *AuthRepository) InsertUser(ctx context.Context, creds dto.Creds) (userID int, err error) {
	query, args, err := sq.Insert(usersDB).Columns(users_login, users_password).
		Values(creds.Username, sq.Expr(`crypt($2, gen_salt('md5'))`, creds.Password)).
		Suffix(fmt.Sprintf(`RETURNING %s`, users_ID)).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return invalidUserID, fmt.Errorf("failed to build query: %w", err)
	}
	err = db.storage.QueryValue(ctx, query, args, []any{&userID})
	if err != nil {
		return invalidUserID, handleSQLError(err)
	}
	return userID, nil
}

func (db *AuthRepository) ValidateUser(ctx context.Context, creds dto.Creds) (userID int, err error) {
	passwordCheck := fmt.Sprintf(`(%s = crypt($2, %s)) AS password_match`, users_password, users_password)
	query, args, err := sq.Select(users_ID, passwordCheck).
		From(usersDB).
		Where(sq.Eq{users_login: creds.Username}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return invalidUserID, fmt.Errorf("failed to build query: %w", err)
	}
	result := struct {
		userID          int
		passwordMatches bool
	}{}
	err = db.storage.QueryValue(
		ctx,
		query,
		append(args, creds.Password),
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
