package states

import (
	"go-keeper/internal/client/logic/commands"
	"go-keeper/internal/client/logic/config"
	"go-keeper/internal/client/logic/services"
)

type TokenRepository interface {
	HasToken() bool
	GetToken() (string, error)
	SetToken(token string) error
}

type AuthService interface {
	Register(username, password string) (tkn string, err error)
	Login(username, password string) (tkn string, err error)
}

type StorageService interface {
	List() (map[string]services.EntryMeta, error)
	Store(entry services.Entry) error
}

type DependenciesContainer struct {
	Config          config.Config
	TokenRepository TokenRepository
	Commands        *commands.Commands
	AuthService     AuthService
	StorageService  StorageService
}
