package services

import (
	"encoding/json"
	"fmt"
	"go-keeper/internal/client/logic/requester"
	"net/http"
)

type Entry struct {
	ID       int
	Metadata map[string]string
}

type StorageService struct {
	req *requester.Requester
}

func NewStorageService(req *requester.Requester) *StorageService {
	return &StorageService{
		req: req,
	}
}

func (s *StorageService) List() ([]Entry, error) {
	resp, err := s.req.Get("/api/user/list")

	if err != nil {
		return []Entry{}, fmt.Errorf("post request failed: %w", err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		break
	case http.StatusBadRequest:
		return []Entry{}, ErrInvalidInput
	default:
		return []Entry{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	var res []Entry
	err = json.Unmarshal(resp.Body(), &res)
	if err != nil {
		return []Entry{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return res, nil
}
