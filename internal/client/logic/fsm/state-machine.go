package fsm

import (
	"context"
)

type State interface {
	OnEnter()
	OnLeave()
	Process(ctx context.Context) (next State, err error)
}

type StateMachine struct {
	currentState State
}

func NewStateMachine(initialState State) *StateMachine {
	initialState.OnEnter()
	return &StateMachine{
		currentState: initialState,
	}
}

func (s *StateMachine) Process(ctx context.Context) error {
	for s.currentState != nil {
		nextState, err := s.currentState.Process(ctx)
		if err != nil {
			return err
		}
		s.currentState.OnLeave()
		s.currentState = nextState
		if nextState != nil {
			s.currentState.OnEnter()
		}
	}
	return nil
}
