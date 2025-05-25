package repositories

import "go-keeper/internal/client/data/storage"

const (
	DataKey = "data"
)

type Entry struct {
	Metadata       map[string]string
	Data           []byte
	StoredOnServer bool
}

type DataRepository struct {
	storage storage.Storage
}

func NewDataRepository(storage storage.Storage) *DataRepository {
	return &DataRepository{
		storage: storage,
	}
}

func (r *DataRepository) GetAll() (map[string]Entry, error) {
	return r.getAllInternal()
}

func (r *DataRepository) getAllInternal() (map[string]Entry, error) {
	if !r.storage.Has(DataKey) {
		return make(map[string]Entry), nil
	}
	res, err := storage.Get[map[string]Entry](r.storage, DataKey)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *DataRepository) SetAll(data map[string]Entry) error {
	return r.setAllInternal(data)
}
func (r *DataRepository) setAllInternal(data map[string]Entry) error {
	err := storage.Set(r.storage, DataKey, data)
	if err != nil {
		return err
	}
	return nil
}

func (r *DataRepository) Set(guid string, value Entry) error {
	current, err := r.getAllInternal()
	if err != nil {
		return err
	}
	current[guid] = value
	err = r.setAllInternal(current)
	if err != nil {
		return err
	}
	return nil
}
