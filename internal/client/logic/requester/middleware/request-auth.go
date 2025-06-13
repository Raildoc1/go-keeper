package middleware

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

var _ resty.RequestMiddleware = (&AuthMiddleware{}).Execute

type TokenRepository interface {
	GetToken() (string, error)
}

type AuthMiddleware struct {
	tokenRepository TokenRepository
}

func NewAuthMiddleware(tokenRepository TokenRepository) *AuthMiddleware {
	return &AuthMiddleware{
		tokenRepository: tokenRepository,
	}
}

func (m *AuthMiddleware) Execute(_ *resty.Client, r *resty.Request) error {
	tkn, err := m.tokenRepository.GetToken()
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}
	r.SetHeader("Authorization", fmt.Sprintf("Bearer %s", tkn))
	return nil
}
