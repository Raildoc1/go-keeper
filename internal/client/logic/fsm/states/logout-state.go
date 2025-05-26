package states

import (
	"context"
	"go-keeper/internal/client/logic/fsm"
)

var _ fsm.State = (*LogoutState)(nil)

type LogoutState struct {
	dc DependenciesContainer
}

func NewLogoutState(dc DependenciesContainer) *LogoutState {
	return &LogoutState{
		dc: dc,
	}
}

func (s *LogoutState) OnEnter() {}
func (s *LogoutState) OnLeave() {}

func (s *LogoutState) Process(ctx context.Context) (next fsm.State, err error) {

	err = s.dc.TokenRepository.Reset()
	if err != nil {
		return NewErrorState(s.dc, err), nil
	}

	return NewAuthState(s.dc), nil
}
