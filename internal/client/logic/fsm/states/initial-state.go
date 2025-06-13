package states

import (
	"context"
	"go-keeper/internal/client/logic/fsm"
)

var _ fsm.State = (*InitialState)(nil)

type InitialState struct {
	dc DependenciesContainer
}

func NewInitialState(dc DependenciesContainer) *InitialState {
	return &InitialState{
		dc: dc,
	}
}

func (s *InitialState) OnEnter() {}
func (s *InitialState) OnLeave() {}

func (s *InitialState) Process(ctx context.Context) (next fsm.State, err error) {
	if s.dc.TokenRepository.HasToken() {
		return NewSelectState(s.dc), nil
	} else {
		return NewAuthState(s.dc), nil
	}
}
