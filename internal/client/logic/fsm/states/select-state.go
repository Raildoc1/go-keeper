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
		"quit",
	}

	err = s.dc.Commands.Write(cmds)
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
			entries, err := s.dc.StorageService.List()
			if err != nil {
				return nil, err
			}
			for _, entry := range entries {
				fmt.Println(entry)
			}
		case "load":
		case "store":
		case "quit":
			return nil, nil
		default:
			fmt.Printf("unknown command: %s\n", cmd)
		}
	}
}
