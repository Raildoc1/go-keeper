package storage

import "encoding/json"

type Storage interface {
	Has(key string) bool
	Reset(key string) error
	set(key string, value string) error
	get(key string) (string, error)
}

func Get[T any](s Storage, key string) (T, error) {
	var zero T
	value, err := s.get(key)
	if err != nil {
		return zero, err
	}
	var result T
	err = json.Unmarshal([]byte(value), &result)
	if err != nil {
		return zero, err
	}
	return result, nil
}

func Set[T any](s Storage, key string, value T) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return s.set(key, string(bytes))
}
