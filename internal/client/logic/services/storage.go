package services

import (
	"github.com/google/uuid"
	"go-keeper/internal/client/data/repositories"
	"go-keeper/internal/client/logic/requester"
)

type EntryMeta struct {
	Metadata       map[string]string
	StoredOnServer bool
}

type Entry struct {
	Metadata       map[string]string
	Data           []byte
	StoredOnServer bool
}

type DataRepository interface {
	GetAll() (map[string]repositories.Entry, error)
	SetAll(data map[string]repositories.Entry) error
	Set(guid string, value repositories.Entry) error
}

type StorageService struct {
	dataRepository DataRepository
	req            *requester.Requester
}

func NewStorageService(dataRepository DataRepository, req *requester.Requester) *StorageService {
	return &StorageService{
		dataRepository: dataRepository,
		req:            req,
	}
}

func (s *StorageService) List() (map[string]EntryMeta, error) {
	data, err := s.dataRepository.GetAll()
	if err != nil {
		return nil, err
	}
	res := make(map[string]EntryMeta, len(data))
	for guid, entry := range data {
		res[guid] = EntryMeta{
			Metadata:       entry.Metadata,
			StoredOnServer: entry.StoredOnServer,
		}
	}
	return res, nil
}

func (s *StorageService) Store(entry Entry) error {
	guid := uuid.New().String()
	err := s.dataRepository.Set(
		guid,
		repositories.Entry{
			Metadata:       entry.Metadata,
			StoredOnServer: entry.StoredOnServer,
			Data:           entry.Data,
		},
	)
	if err != nil {
		return err
	}
	return nil
}
