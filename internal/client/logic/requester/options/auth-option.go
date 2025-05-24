package options

var _ Option = (*AuthOption)(nil)

type TokenRepository interface {
	GetToken() string
}

type AuthOption struct {
	tokenRepository TokenRepository
}

func NewAuthOption(tokenRepository TokenRepository) *AuthOption {
	return &AuthOption{
		tokenRepository: tokenRepository,
	}
}

func (o *AuthOption) GetHeaders() map[string]string {
	return map[string]string{
		"Authorization": "Bearer " + o.tokenRepository.GetToken(),
	}
}
