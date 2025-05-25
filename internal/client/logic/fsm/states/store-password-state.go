package states

import (
	"context"
	"encoding/json"
	"go-keeper/internal/client/logic/fsm"
	"go-keeper/internal/client/logic/services"
)

var _ fsm.State = (*StoreCredsState)(nil)

type creds struct {
	Login    string
	Password string
}

type StoreCredsState struct {
	dc DependenciesContainer
}

func NewStoreCredsState(dc DependenciesContainer) *StoreCredsState {
	return &StoreCredsState{
		dc: dc,
	}
}

func (s *StoreCredsState) OnEnter() {}
func (s *StoreCredsState) OnLeave() {}

func (s *StoreCredsState) Process(ctx context.Context) (next fsm.State, err error) {
	login, err := s.dc.Commands.ReadWithLabel("type login to store", ctx)
	if err != nil {
		return nil, err
	}
	password, err := s.dc.Commands.ReadWithLabel("type password to store", ctx)
	if err != nil {
		return nil, err
	}

	c := creds{
		Login:    login,
		Password: password,
	}

	credsJSON, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}

	err = s.dc.StorageService.Store(services.Entry{
		Data: credsJSON,
		Metadata: map[string]string{
			"type": "password",
		},
		StoredOnServer: false,
	})

	if err != nil {
		return nil, err
	}

	return NewSelectState(s.dc), nil
}
