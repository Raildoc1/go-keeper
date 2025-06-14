package client

import (
	"context"
	"go-keeper/internal/client/logic/commands"
	"go-keeper/internal/client/logic/fsm"
	"go-keeper/internal/client/logic/fsm/states"
)

type Client struct {
	sm *fsm.StateMachine
}

func New(
	tokenRepository states.TokenRepository,
	cmds *commands.Commands,
	authService states.AuthService,
	storageService states.StorageService,
) *Client {
	dc := states.DependenciesContainer{
		TokenRepository: tokenRepository,
		Commands:        cmds,
		AuthService:     authService,
		StorageService:  storageService,
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
