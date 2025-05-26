package states

import (
	"context"
	"go-keeper/internal/client/logic/fsm"
)

var _ fsm.State = (*LoadState)(nil)

type LoadState struct {
	dc DependenciesContainer
}

func NewLoadState(dc DependenciesContainer) *LoadState {
	return &LoadState{
		dc: dc,
	}
}

func (s *LoadState) OnEnter() {}
func (s *LoadState) OnLeave() {}

func (s *LoadState) Process(ctx context.Context) (next fsm.State, err error) {
	err = ListEntries(s.dc)
	if err != nil {
		return NewErrorState(s.dc, err), nil
	}

	return NewSelectState(s.dc), nil
}
