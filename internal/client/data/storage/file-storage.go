package storage

import (
	"encoding/json"
	"errors"
	"os"
)

var (
	ErrNotFound = errors.New("item not found")
)

type FileStorage struct {
	data map[string]string
}

func NewFileStorage(path string) (*FileStorage, error) {
	data, err := Recover(path)
	if err != nil {
		return nil, err
	}
	return &FileStorage{
		data: data,
	}, nil
}

func Recover(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]string), nil
		}
		return nil, err
	}
	defer file.Close()
	var result map[string]string
	err = json.NewDecoder(file).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (f *FileStorage) Has(key string) bool {
	_, ok := f.data[key]
	return ok
}

func (s *FileStorage) set(key string, value string) {
	s.data[key] = value
}

func (s *FileStorage) get(key string) (string, error) {
	res, ok := s.data[key]
	if !ok {
		return "", ErrNotFound
	}
	return res, nil
}

func (s *FileStorage) Save(path string) error {
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
