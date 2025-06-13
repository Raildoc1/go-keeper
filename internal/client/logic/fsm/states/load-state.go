package states

import (
	"context"
	"errors"
	"fmt"
	"go-keeper/internal/client/logic/fsm"
	"go-keeper/internal/client/logic/services"
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

	guid, err := s.dc.Commands.ReadWithLabel("enter guid to load", ctx)
	if err != nil {
		return NewErrorState(s.dc, err), nil
	}

	entry, err := s.dc.StorageService.Load(guid)
	if err != nil {
		if errors.Is(err, services.ErrEntryNotFound) {
			fmt.Println("No such guid")
			return NewLoadState(s.dc), nil
		}
		return NewErrorState(s.dc, err), nil
	}

	fmt.Println(string(entry.Data))
	fmt.Println()

	return NewSelectState(s.dc), nil
}
