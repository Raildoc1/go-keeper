package states

import (
	"context"
	"go-keeper/internal/client/logic/fsm"
	"go-keeper/internal/client/logic/services"
)

var _ fsm.State = (*StoreTextState)(nil)

type StoreTextState struct {
	dc DependenciesContainer
}

func NewStoreTextState(dc DependenciesContainer) *StoreTextState {
	return &StoreTextState{
		dc: dc,
	}
}

func (s *StoreTextState) OnEnter() {}
func (s *StoreTextState) OnLeave() {}

func (s *StoreTextState) Process(ctx context.Context) (next fsm.State, err error) {
	text, err := s.dc.Commands.ReadWithLabel("type text to store", ctx)
	if err != nil {
		return nil, err
	}

	err = s.dc.StorageService.Store(services.Entry{
		Data: []byte(text),
		Metadata: map[string]string{
			"type": "text",
		},
		StoredOnServer: false,
	})

	if err != nil {
		return nil, err
	}

	return NewSelectState(s.dc), nil
}
