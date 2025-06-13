package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

var (
	ErrNotFound = errors.New("item not found")
)

type FileStorage struct {
	data map[string]string
	path string
}

func NewFileStorage(path string) (*FileStorage, error) {
	data, err := readFromFile(path)
	if err != nil {
		return nil, err
	}
	return &FileStorage{
		data: data,
		path: path,
	}, nil
}

func readFromFile(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]string), nil
		}
		return nil, err
	}
	defer file.Close()
	return decodeJSON(file)
}

func decodeJSON(r io.Reader) (map[string]string, error) {
	var result map[string]string
	err := json.NewDecoder(r).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *FileStorage) save(path string) error {
	d, err := json.Marshal(s.data)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, d, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (s *FileStorage) Has(key string) bool {
	_, ok := s.data[key]
	return ok
}

func (s *FileStorage) Reset(key string) error {
	delete(s.data, key)
	err := s.save(s.path)
	if err != nil {
		return fmt.Errorf("failed to save to file: %w", err)
	}
	return nil
}

func (s *FileStorage) set(key string, value string) error {
	s.data[key] = value
	err := s.save(s.path)
	if err != nil {
		return fmt.Errorf("failed to save to file: %w", err)
	}
	return nil
}

func (s *FileStorage) get(key string) (string, error) {
	res, ok := s.data[key]
	if !ok {
		return "", ErrNotFound
	}
	return res, nil
}
