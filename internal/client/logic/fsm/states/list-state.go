package states

import (
	"context"
	"fmt"
	"go-keeper/internal/client/logic/fsm"
)

var _ fsm.State = (*ListState)(nil)

type ListState struct {
	dc DependenciesContainer
}

func NewListState(dc DependenciesContainer) *ListState {
	return &ListState{
		dc: dc,
	}
}

func (s *ListState) OnEnter() {}
func (s *ListState) OnLeave() {}

func (s *ListState) Process(ctx context.Context) (next fsm.State, err error) {

	entries, err := s.dc.StorageService.List()
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		fmt.Println(entry)
	}

	return NewSelectState(s.dc), nil
}
