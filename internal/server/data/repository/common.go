package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go-keeper/internal/server/data"
)

const (
	usersDB        = "users"
	users_ID       = "id"
	users_login    = "login"
	users_password = "password"

	dataDB        = "data"
	data_ID       = "id"
	data_owner    = "owner"
	data_guid     = "guid"
	data_metadata = "metadata"
	data_data     = "data"
)

type DBStorage interface {
	Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, query string, args ...any) (pgx.Row, error)
	Query(ctx context.Context, query string, args ...any) (pgx.Rows, error)
	QueryValue(ctx context.Context, query string, args []any, dest []any) error
}

func handleSQLError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return data.ErrUniqueConstraintViolation
		}
	}
	return err
}
