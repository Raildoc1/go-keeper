package repository

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
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

//go:embed sql/storage/store.sql
var storeQuery string

func (db *StorageRepository) Store(ctx context.Context, userID int, guid string, entry dto.Entry) error {
	metadata, err := json.Marshal(entry.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}
	_, err = db.storage.Exec(ctx, storeQuery, userID, guid, metadata, entry.Data)
	if err != nil {
		return handleSQLError(err)
	}
	return nil
}

//go:embed sql/storage/load.sql
var loadQuery string

func (db *StorageRepository) Load(ctx context.Context, userID int, guid string) (dto.Entry, error) {
	var metadataJSON string
	var data []byte
	err := db.storage.QueryValue(ctx, loadQuery, []any{guid, userID}, []any{&metadataJSON, &data})
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

//go:embed sql/storage/load-all.sql
var loadAllQuery string

func (db *StorageRepository) LoadAll(ctx context.Context, userID int) (map[string]dto.Entry, error) {
	rows, err := db.storage.Query(ctx, loadAllQuery, userID)
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

//go:embed sql/storage/delete.sql
var deleteQuery string

func (db *StorageRepository) Delete(ctx context.Context, userID int, guid string) error {
	return errors.New("unimplemented") // todo
}
