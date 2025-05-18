package services

import (
	"context"
	"go-keeper/internal/server/dto"
)

type StorageService struct{}

func NewStorageService() *StorageService {
	return &StorageService{}
}

func (s *StorageService) Store(ctx context.Context, entry dto.Entry) error {
	return nil
}

func (s *StorageService) Load(ctx context.Context, id string) (dto.Entry, error) {
	return dto.Entry{}, nil
}

func (s *StorageService) Delete(ctx context.Context, id string) error {
	return nil
}

func (s *StorageService) LoadAll(ctx context.Context) ([]dto.Entry, error) {
	return nil, nil
}
