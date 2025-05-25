package repositories

import (
	"fmt"
	"go-keeper/internal/client/data/storage"
)

const (
	TokenKey = "token"
)

type TokenRepository struct {
	storage storage.Storage
}

func NewTokenRepository(storage storage.Storage) *TokenRepository {
	return &TokenRepository{
		storage: storage,
	}
}

func (r *TokenRepository) HasToken() bool {
	return r.storage.Has(TokenKey)
}

func (r *TokenRepository) GetToken() (string, error) {
	tkn, err := storage.Get[string](r.storage, TokenKey)
	if err != nil {
		return "", fmt.Errorf("failed to get token: %w", err)
	}
	return tkn, nil
}

func (r *TokenRepository) SetToken(token string) error {
	err := storage.Set(r.storage, TokenKey, token)
	if err != nil {
		return fmt.Errorf("failed to set token: %w", err)
	}
	return nil
}
