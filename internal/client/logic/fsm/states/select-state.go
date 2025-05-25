package states

import (
	"context"
	"fmt"
	"go-keeper/internal/client/logic/fsm"
)

var _ fsm.State = (*SelectState)(nil)

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
	cmds := []string{
		"list",
		"load",
		"store",
		"sync",
		"logout",
		"quit",
	}

	err = s.dc.Commands.WriteWithLabel("available commands", cmds)
	if err != nil {
		return nil, err
	}

	for {
		cmd, err := s.dc.Commands.ReadWithLabel("enter command", ctx)
		if err != nil {
			return nil, err
		}

		switch cmd {
		case "list":
			return NewListState(s.dc), nil
		case "load":
		case "store":
			return NewStoreState(s.dc), nil
		case "sync":
			return NewSyncState(s.dc), nil
		case "logout":
			return NewLogoutState(s.dc), nil
		case "quit":
			return nil, nil
		default:
			fmt.Printf("unknown command: %s\n", cmd)
		}
	}
}
