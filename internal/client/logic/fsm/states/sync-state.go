package states

import (
	"context"
	"fmt"
	"go-keeper/internal/client/logic/fsm"
)

var _ fsm.State = (*SyncState)(nil)

type SyncState struct {
	dc DependenciesContainer
}

func NewSyncState(dc DependenciesContainer) *SyncState {
	return &SyncState{
		dc: dc,
	}
}

func (s *SyncState) OnEnter() {}
func (s *SyncState) OnLeave() {}

func (s *SyncState) Process(ctx context.Context) (next fsm.State, err error) {

	fmt.Println("Syncing...")

	err = s.dc.StorageService.Sync()
	if err != nil {
		return nil, err
	}

	fmt.Println("Synced succeeded!")

	return NewSelectState(s.dc), nil
}
