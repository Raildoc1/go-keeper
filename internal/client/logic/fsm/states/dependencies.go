package states

import (
	"go-keeper/internal/client/logic/commands"
	"go-keeper/internal/client/logic/services"
)

type TokenRepository interface {
	HasToken() bool
	GetToken() (string, error)
	SetToken(token string) error
	Reset() error
}

type AuthService interface {
	Register(username, password string) (tkn string, err error)
	Login(username, password string) (tkn string, err error)
}

type StorageService interface {
	List() (map[string]services.EntryMeta, error)
	Store(entry services.Entry) error
	Load(guid string) (services.Entry, error)
	Sync() error
}

type DependenciesContainer struct {
	TokenRepository TokenRepository
	Commands        *commands.Commands
	AuthService     AuthService
	StorageService  StorageService
}
