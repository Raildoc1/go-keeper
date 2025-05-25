package states

import (
	"context"
	"fmt"
	"go-keeper/internal/client/logic/fsm"
)

var _ fsm.State = (*AuthState)(nil)

type AuthState struct {
	dc DependenciesContainer
}

func NewAuthState(dc DependenciesContainer) *AuthState {
	return &AuthState{
		dc: dc,
	}
}

func (s *AuthState) OnEnter() {}
func (s *AuthState) OnLeave() {}

func (s *AuthState) Process(ctx context.Context) (next fsm.State, err error) {
	cmds := []string{
		"register",
		"login",
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
		case "register":
			login, password, err := s.GetCreds(ctx)
			if err != nil {
				return nil, err
			}
			tkn, err := s.dc.AuthService.Register(login, password)
			if err != nil {
				return nil, err
			}
			s.dc.TokenRepository.SetToken(tkn)
			fmt.Printf("token successfully received (%s)\n", tkn)
			return NewSelectState(s.dc), nil
		case "login":
			login, password, err := s.GetCreds(ctx)
			if err != nil {
				return nil, err
			}
			tkn, err := s.dc.AuthService.Login(login, password)
			if err != nil {
				return nil, err
			}
			s.dc.TokenRepository.SetToken(tkn)
			fmt.Printf("token successfully received (%s)\n", tkn)
			return NewSelectState(s.dc), nil
		case "quit":
			return nil, nil
		default:
			fmt.Printf("unknown command: %s\n", cmd)
		}
	}
}

func (s *AuthState) GetCreds(ctx context.Context) (login string, password string, err error) {
	login, err = s.dc.Commands.ReadWithLabel("enter login", ctx)
	if err != nil {
		return "", "", fmt.Errorf("failed to read login: %w", err)
	}
	password, err = s.dc.Commands.ReadWithLabel("enter password", ctx)
	if err != nil {
		return "", "", fmt.Errorf("failed to read password: %w", err)
	}
	return login, password, nil
}
