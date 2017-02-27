package auth

type ValidatorFunc func(token string) bool

type TokenValidator interface {
	ValidToken(token string) bool
}

func (v ValidatorFunc) ValidToken(token string) bool {
	return v(token)
}
