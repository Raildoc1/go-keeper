package client

import (
	"go-keeper/internal/client/config"
	"go-keeper/internal/client/logic/fsm"
	"go-keeper/internal/client/logic/states"
)

type Client struct {
	sm *fsm.StateMachine
}

func New(config config.Config) *Client {
	dc := states.DependenciesContainer{
		Config: config.LogicConfig,
	}
	sm := fsm.NewStateMachine(states.NewConnectingState(dc))
	return &Client{
		sm: sm,
	}
}

func (c *Client) Run() error {
	c.sm.Process()
	return nil
}

func (c *Client) Stop() error {
	return nil
}
