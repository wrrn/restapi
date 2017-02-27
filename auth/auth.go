package auth

import (
	"errors"
	"net/http"
)

var (
	InvalidTokenErr = errors.New("Invalid Token")
	MissingTokenErr = errors.New("Token missing from request")
)

// VerifyTokens will return a handler that will verify that a session
// exists before allowing the handler in the arugment to be called.
// If the token is not valid it responds with a 401 code.
func VerifyTokens(v TokenValidator, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if token := getToken(r); !v.ValidToken(token) {
			Unauthorized(w)
			return
		}

		h.ServeHTTP(w, r)
	})
}

// Unauthorized is just a convience function that allows us to write a
// status code of 401 and a message of "Unauthorized" to the response
func Unauthorized(w http.ResponseWriter) {
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}

// Forbidden is just a convience function that allows us to write a
// status code of 403 and a message of "Forbidden" to the response
func Forbidden(w http.ResponseWriter) {
	http.Error(w, "Forbidden", http.StatusForbidden)
}

// GetToken
func getToken(r *http.Request) (token string) {
	token, _, _ = r.BasicAuth()
	return token
}
