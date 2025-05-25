package states

import (
	"context"
	"encoding/json"
	"go-keeper/internal/client/logic/fsm"
	"go-keeper/internal/client/logic/services"
)

var _ fsm.State = (*StoreCardState)(nil)

type card struct {
	Number     string
	CVC        string
	HolderName string
}

type StoreCardState struct {
	dc DependenciesContainer
}

func NewStoreCardState(dc DependenciesContainer) *StoreCardState {
	return &StoreCardState{
		dc: dc,
	}
}

func (s *StoreCardState) OnEnter() {}
func (s *StoreCardState) OnLeave() {}

func (s *StoreCardState) Process(ctx context.Context) (next fsm.State, err error) {
	number, err := s.dc.Commands.ReadWithLabel("type card number", ctx)
	if err != nil {
		return nil, err
	}
	cvc, err := s.dc.Commands.ReadWithLabel("type card security code", ctx)
	if err != nil {
		return nil, err
	}
	holderName, err := s.dc.Commands.ReadWithLabel("type card holder name", ctx)
	if err != nil {
		return nil, err
	}

	cardData := card{
		Number:     number,
		CVC:        cvc,
		HolderName: holderName,
	}

	cardJSON, err := json.Marshal(cardData)
	if err != nil {
		return nil, err
	}

	err = s.dc.StorageService.Store(services.Entry{
		Data: cardJSON,
		Metadata: map[string]string{
			"type": "card",
		},
		StoredOnServer: false,
	})

	if err != nil {
		return nil, err
	}

	return NewSelectState(s.dc), nil
}
