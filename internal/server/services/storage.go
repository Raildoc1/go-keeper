package services

import (
	"context"
	"go-keeper/internal/server/dto"
)

type StorageRepository interface {
	Store(ctx context.Context, userID int, entry dto.Entry) error
	Load(ctx context.Context, entryId, userID int) (dto.Entry, error)
	Delete(ctx context.Context, userID int, id int) error
	LoadAll(ctx context.Context, userID int) (map[int]dto.Entry, error)
}

type StorageService struct {
	repository StorageRepository
}

func NewStorageService(repository StorageRepository) *StorageService {
	return &StorageService{
		repository: repository,
	}
}

func (s *StorageService) Store(ctx context.Context, userID int, entry dto.Entry) error {
	return s.repository.Store(ctx, userID, entry)
}

func (s *StorageService) Load(ctx context.Context, userID int, id int) (dto.Entry, error) {
	return s.repository.Load(ctx, userID, id)
}

func (s *StorageService) Delete(ctx context.Context, userID int, id int) error {
	return s.repository.Delete(ctx, userID, id)
}

func (s *StorageService) LoadAll(ctx context.Context, userID int) (map[int]dto.Entry, error) {
	return s.repository.LoadAll(ctx, userID)
}
