package states

import (
	"context"
	"fmt"
	"go-keeper/internal/client/logic/fsm"
)

var _ fsm.State = (*StoreState)(nil)

type StoreState struct {
	dc DependenciesContainer
}

func NewStoreState(dc DependenciesContainer) *StoreState {
	return &StoreState{
		dc: dc,
	}
}

func (s *StoreState) OnEnter() {}
func (s *StoreState) OnLeave() {}

func (s *StoreState) Process(ctx context.Context) (next fsm.State, err error) {

	cmds := []string{
		"password",
		"text",
		"file",
		"card",
		"back",
	}

	err = s.dc.Commands.WriteWithLabel("available types", cmds)
	if err != nil {
		return nil, err
	}

	for {
		cmd, err := s.dc.Commands.ReadWithLabel("choose type", ctx)
		if err != nil {
			return nil, err
		}

		switch cmd {
		case "password":
			return NewStoreCredsState(s.dc), nil
		case "text":
			return NewStoreTextState(s.dc), nil
		case "card":
			return NewStoreCardState(s.dc), nil
		case "back":
			return NewSelectState(s.dc), nil
		default:
			fmt.Printf("unknown command: %s\n", cmd)
		}
	}
}
