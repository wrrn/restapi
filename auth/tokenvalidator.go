package auth

// TokenValidator validates a token.
//
// ValidToken should return true if the token is valid, or
// false if invalid.
type TokenValidator interface {
	ValidToken(token string) bool
}

// ValidatorFunc is an adapter to allow the use of ordinary functions as
// TokenValidator.
type ValidatorFunc func(token string) bool

// ValidToken calls f(token)
func (f ValidatorFunc) ValidToken(token string) bool {
	return f(token)
}
