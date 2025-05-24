package data

type TokenRepository struct {
	tkn string
}

func NewTokenRepository() *TokenRepository {
	return &TokenRepository{
		tkn: "",
	}
}

func (t *TokenRepository) HasToken() bool {
	return t.tkn != ""
}

func (t *TokenRepository) GetToken() string {
	return t.tkn
}

func (t *TokenRepository) SetToken(token string) {
	t.tkn = token
}
