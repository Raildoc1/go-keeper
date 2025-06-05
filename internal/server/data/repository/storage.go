package repository

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"go-keeper/internal/server/dto"
	"go-keeper/pkg/logging"
)

type StorageRepository struct {
	storage DBStorage
	logger  *logging.ZapLogger
}

func NewStorageRepository(storage DBStorage, logger *logging.ZapLogger) *StorageRepository {
	return &StorageRepository{
		storage: storage,
		logger:  logger,
	}
}

func (db *StorageRepository) Store(ctx context.Context, userID int, guid string, entry dto.Entry) error {
	metadata, err := json.Marshal(entry.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	query, args, err := sq.Insert(dataDB).Columns(data_owner, data_guid, data_metadata, data_data).
		Values(userID, guid, string(metadata), entry.Data).
		Suffix(fmt.Sprintf(`RETURNING %s`, data_ID)).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = db.storage.Exec(ctx, query, args...)
	if err != nil {
		return handleSQLError(err)
	}
	return nil
}

func (db *StorageRepository) Load(ctx context.Context, userID int, guid string) (dto.Entry, error) {
	query, args, err := sq.Select(data_metadata, data_data).From(dataDB).
		Where(sq.And{sq.Eq{data_guid: guid}, sq.Eq{data_owner: userID}}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return dto.Entry{}, fmt.Errorf("failed to build query: %w", err)
	}

	var metadataJSON string
	var data []byte
	err = db.storage.QueryValue(ctx, query, args, []any{&metadataJSON, &data})
	if err != nil {
		return dto.Entry{}, handleSQLError(err)
	}

	var metadata map[string]string
	err = json.Unmarshal([]byte(metadataJSON), &metadata)
	if err != nil {
		return dto.Entry{}, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	return dto.Entry{
		Metadata: metadata,
		Data:     data,
	}, nil
}

func (db *StorageRepository) LoadAll(ctx context.Context, userID int) (map[string]dto.Entry, error) {
	query, args, err := sq.Select(data_guid, data_metadata, data_data).From(dataDB).
		Where(sq.Eq{data_owner: userID}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := db.storage.Query(ctx, query, args...)
	if err != nil {
		return nil, handleSQLError(err)
	}
	defer rows.Close()

	result := make(map[string]dto.Entry)

	if err = rows.Err(); err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return result, nil
		default:
			return nil, handleSQLError(err)
		}
	}

	for rows.Next() {
		var guid string
		var metadataJSON string
		var data []byte
		err := rows.Scan(
			&guid,
			&metadataJSON,
			&data,
		)
		if err != nil {
			return nil, handleSQLError(err)
		}
		var metadata map[string]string
		err = json.Unmarshal([]byte(metadataJSON), &metadata)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}

		result[guid] = dto.Entry{
			Metadata: metadata,
			Data:     data,
		}
	}

	return result, nil
}
