package auth

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/wrrn/token"
)

var (
	InvalidTokenErr = errors.New("Invalid Token")
	MissingTokenErr = errors.New("Token missing from request")
)

type Auth struct {
	RowQueryer
}

// RowQueryer is implemented by the *sql.DB struct
type RowQueryer interface {
	QueryRow(string, ...interface{}) *sql.Row
}

func (a Auth) ValidToken(t string) bool {
	var count int
	err := a.QueryRow(`SELECT COUNT(*)
                       FROM token INNER JOIN user_account
                       ON token.user_id = user_account.id
                       WHERE token.token = ?`, t).Scan(&count)
	return err == nil && count > 0
}

func (a Auth) ValidateTokens(h http.Handler) http.Handler {
	return token.ValidateTokens(a, h)
}
