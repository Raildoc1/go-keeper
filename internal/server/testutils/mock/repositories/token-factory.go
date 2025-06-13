package repositories

type TokenFactoryMock struct{}

func NewTokenFactoryMock() *TokenFactoryMock {
	return &TokenFactoryMock{}
}

func (t *TokenFactoryMock) Generate(extraClaims map[string]string) (string, error) {
	return "test_token", nil
}
