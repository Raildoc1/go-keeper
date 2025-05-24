package client

import (
	"context"
	"go-keeper/internal/client/config"
	"go-keeper/internal/client/logic/commands"
	"go-keeper/internal/client/logic/fsm"
	"go-keeper/internal/client/logic/fsm/states"
)

type Client struct {
	sm *fsm.StateMachine
}

func New(
	config config.Config,
	tokenRepository states.TokenRepository,
	cmds *commands.Commands,
	authService states.AuthService,
) *Client {
	dc := states.DependenciesContainer{
		Config:          config.LogicConfig,
		TokenRepository: tokenRepository,
		Commands:        cmds,
		AuthService:     authService,
	}
	sm := fsm.NewStateMachine(states.NewInitialState(dc))
	return &Client{
		sm: sm,
	}
}

func (c *Client) Run(ctx context.Context) error {
	return c.sm.Process(ctx)
}

func (c *Client) Stop() error {
	return nil
}
