package services

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go-keeper/internal/client/data/repositories"
	"go-keeper/internal/client/logic/requester"
	"go-keeper/internal/common/protocol"
	"net/http"
)

var (
	ErrTokenExpired  = errors.New("token expired")
	ErrEntryNotFound = errors.New("not found")
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
	Get(guid string) (repositories.Entry, error)
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

func (s *StorageService) Load(guid string) (Entry, error) {
	entry, err := s.dataRepository.Get(guid)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return Entry{}, ErrEntryNotFound
		}
		return Entry{}, err
	}
	return Entry{
		Metadata:       entry.Metadata,
		StoredOnServer: entry.StoredOnServer,
		Data:           entry.Data,
	}, nil
}

func (s *StorageService) Sync() error {
	remoteEntries, statusCode, err := requester.Get[map[string]protocol.Entry](s.req, "/api/user/loadall")
	if err != nil {
		if errors.Is(err, requester.ErrUnexpectedStatusCode) {
			err = s.statusCodeToError(statusCode)
		}
		return fmt.Errorf("failed to get remote entries: %w", err)
	}

	localEntries, err := s.dataRepository.GetAll()
	if err != nil {
		return fmt.Errorf("failed to get local entries: %w", err)
	}

	for guid, localEntry := range localEntries {
		if _, ok := remoteEntries[guid]; !ok {
			if localEntry.StoredOnServer {
				// entry was deleted from server
				delete(localEntries, guid)
			} else {
				err = s.upload(guid, localEntry)
				if err != nil {
					return err
				}
				localEntry.StoredOnServer = true
				localEntries[guid] = localEntry
			}
		}
	}

	for guid, remoteEntry := range remoteEntries {
		localEntries[guid] = repositories.Entry{
			Metadata:       remoteEntry.Metadata,
			Data:           remoteEntry.Data,
			StoredOnServer: true,
		}
	}

	err = s.dataRepository.SetAll(localEntries)
	if err != nil {
		return err
	}

	return nil
}

func (s *StorageService) upload(guid string, entry repositories.Entry) error {
	protocolEntry := protocol.StoreRequest{
		GUID: guid,
		Entry: protocol.Entry{
			Data:     entry.Data,
			Metadata: entry.Metadata,
		},
	}

	resp, err := s.req.Post("/api/user/store", protocolEntry)
	if err != nil {
		return fmt.Errorf("failed to upload to remote server: %w", err)
	}

	err = s.statusCodeToError(resp.StatusCode())
	if err != nil {
		return fmt.Errorf("failed to upload entry: %w", err)
	}

	return nil
}

func (s *StorageService) statusCodeToError(statusCode int) error {
	switch statusCode {
	case http.StatusOK:
		return nil
	case http.StatusUnauthorized:
		return ErrTokenExpired
	default:
		return fmt.Errorf("unexpected status code: %v", statusCode)
	}
}
