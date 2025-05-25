package services

import (
	"context"
	"go-keeper/internal/server/dto"
)

type StorageRepository interface {
	Store(ctx context.Context, userID int, guid string, entry dto.Entry) error
	Load(ctx context.Context, userID int, guid string) (dto.Entry, error)
	Delete(ctx context.Context, userID int, guid string) error
	LoadAll(ctx context.Context, userID int) (map[string]dto.Entry, error)
}

type StorageService struct {
	repository StorageRepository
}

func NewStorageService(repository StorageRepository) *StorageService {
	return &StorageService{
		repository: repository,
	}
}

func (s *StorageService) Store(ctx context.Context, userID int, guid string, entry dto.Entry) error {
	return s.repository.Store(ctx, userID, guid, entry)
}

func (s *StorageService) Load(ctx context.Context, userID int, guid string) (dto.Entry, error) {
	return s.repository.Load(ctx, userID, guid)
}

func (s *StorageService) Delete(ctx context.Context, userID int, guid string) error {
	return s.repository.Delete(ctx, userID, guid)
}

func (s *StorageService) LoadAll(ctx context.Context, userID int) (map[string]dto.Entry, error) {
	return s.repository.LoadAll(ctx, userID)
}
