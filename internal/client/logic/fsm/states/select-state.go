package states

import (
	"context"
	"go-keeper/internal/client/logic/fsm"
)

var _ fsm.State = (*SelectState)(nil)

const (
	quitOperationName = "quit"
)

type SelectState struct {
	dc DependenciesContainer
}

func NewSelectState(dc DependenciesContainer) *SelectState {
	return &SelectState{
		dc: dc,
	}
}

func (s *SelectState) OnEnter() {}
func (s *SelectState) OnLeave() {}

func (s *SelectState) Process(ctx context.Context) (next fsm.State, err error) {
	return nil, nil
}
