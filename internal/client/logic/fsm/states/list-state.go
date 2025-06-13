package states

import (
	"context"
	"encoding/json"
	"fmt"
	"go-keeper/internal/client/logic/fsm"
	"go-keeper/internal/client/logic/services"
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
	err = ListEntries(s.dc)
	if err != nil {
		return NewErrorState(s.dc, err), nil
	}

	return NewSelectState(s.dc), nil
}

func ListEntries(dc DependenciesContainer) error {
	entries, err := dc.StorageService.List()
	if err != nil {
		return err
	}

	fmt.Println("\nStored entries:")
	for guid, entry := range entries {
		fmt.Println(FormatEntry(guid, entry))
	}
	fmt.Println()

	return nil
}

func FormatEntry(guid string, entryMeta services.EntryMeta) string {
	metadataJSON, err := json.Marshal(entryMeta.Metadata)
	if err != nil {
		metadataJSON = []byte("{ FAILED TO PARSE }")
	}
	return fmt.Sprintf(
		"--- %s: %s (stored on server: %v)",
		guid,
		metadataJSON,
		entryMeta.StoredOnServer,
	)
}
