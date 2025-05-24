package states

import (
	"go-keeper/internal/client/logic/commands"
	"go-keeper/internal/client/logic/config"
)

type TokenRepository interface {
	HasToken() bool
	GetToken() string
	SetToken(token string)
}

type AuthService interface {
	Register(username, password string) (tkn string, err error)
	Login(username, password string) (tkn string, err error)
}

type DependenciesContainer struct {
	Config          config.Config
	TokenRepository TokenRepository
	Commands        *commands.Commands
	AuthService     AuthService
}
