package states

import (
	"context"
	"errors"
	"fmt"
	"go-keeper/internal/client/logic/fsm"
	"go-keeper/internal/client/logic/services"
)

var _ fsm.State = (*ErrorState)(nil)

type ErrorState struct {
	dc  DependenciesContainer
	err error
}

func NewErrorState(dc DependenciesContainer, err error) *ErrorState {
	return &ErrorState{
		dc:  dc,
		err: err,
	}
}

func (s *ErrorState) OnEnter() {}
func (s *ErrorState) OnLeave() {}

func (s *ErrorState) Process(ctx context.Context) (next fsm.State, err error) {
	if errors.Is(s.err, services.ErrTokenExpired) {
		fmt.Println("Token expired")
		return NewLogoutState(s.dc), nil
	}

	if errors.Is(s.err, services.ErrInvalidCreds) {
		fmt.Println("Invalid creds")
		return NewAuthState(s.dc), nil
	}

	return nil, err
}
