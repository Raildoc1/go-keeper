package repositories

import "go-keeper/internal/server/services"

var _ services.TokenFactory = (*TokenFactoryMock)(nil)

type TokenFactoryMock struct{}

func NewTokenFactoryMock() *TokenFactoryMock {
	return &TokenFactoryMock{}
}

func (t *TokenFactoryMock) Generate(extraClaims map[string]string) (string, error) {
	return "test_token", nil
}
