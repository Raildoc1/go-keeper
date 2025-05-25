package options

import "fmt"

var _ Option = (*AuthOption)(nil)

type TokenRepository interface {
	GetToken() (string, error)
}

type AuthOption struct {
	tokenRepository TokenRepository
}

func NewAuthOption(tokenRepository TokenRepository) *AuthOption {
	return &AuthOption{
		tokenRepository: tokenRepository,
	}
}

func (o *AuthOption) GetHeaders() (map[string]string, error) {
	tkn, err := o.tokenRepository.GetToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	return map[string]string{
		"Authorization": "Bearer " + tkn,
	}, nil
}
