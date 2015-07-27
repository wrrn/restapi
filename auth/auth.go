package auth

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

const (
	CookieName       = "RESTAPI"
	secret           = "Not a secret"
	randMax    int64 = 54050505434503053
)

var (
	InvalidSessionErr = errors.New("Invalid Session")
)

type User struct {
	// id refers to the ID that is stored in the database
	id       int    `json:-`
	Username string `json:"username"`
	Password string `json:"password"`
}
type Auth struct {
	*sql.DB
}

// HandleLogin checks decodes the request and creates a session for valid
// credentials. If the users credentials are correct and session could be
// created than a 200 code with a message of "Authorized" will be returned.
// If the credentials are bad 401 code with a message of "Unauthorized" will
// be returned. A 501 error  with a message of "Server Error" will be returned
// if a session cannot be created or the body of the response cannot be written.
func (a Auth) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var (
		user User
		err  error
	)

	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Format Error", http.StatusBadRequest)
		return
	}
	user, err = a.login(user.Username, user.Password)
	if err != nil {
		Unauthorized(w)
		return
	}
	sessionID, err := a.createSession(user)

	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	cookie := generateCookie(sessionID)
	http.SetCookie(w, cookie)

	if _, err := w.Write([]byte("Authorized")); err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
	}

}

// Handlelogout will write a 200 code with a message of success to the response.
// If the response cannot be written to, a 500 code with the message "Server Error" will be sent
func (a Auth) HandleLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(CookieName)
	if err == nil {
		a.revokeSession(cookie.Value)
	}

	if _, err := w.Write([]byte("Success")); err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
	}

}

// CheckSession checks the request to verify that the value of cookie with the
// name "RESTAPI" matches a session id  stored in the database
func (a Auth) CheckSession(r *http.Request) (user User, err error) {
	cookie, err := r.Cookie(CookieName)
	if err != nil {
		return user, err
	}
	err = a.DB.QueryRow("SELECT id, username FROM users INNER JOIN sessions ON users.id = sessions.user_id WHERE sessions.session_id = $1", cookie.Value).Scan(&user.id, &user.Username)

	return user, err
}

// VerifySessions will return a handler that will verify that a session
// exists before allowing the handler in the arugment to be called.
// If a session does not exist sends a 403 code.
func (a Auth) VerifySessions(h http.Handler) http.Handler {
	return sessionsHandler{
		Handler: h,
		Auth:    a,
	}
}

// RegisterUser register a user and stores them in the database.
func (a Auth) RegisterUser(user User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = a.DB.Exec("INSERT INTO users(username, password) VALUES($1, $2)", user.Username, hashedPassword)
	return err
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

// generateCookie returns a cookie whose name is "RESTAPI" and whose value is
// the value of the argument.
func generateCookie(sessionID string) *http.Cookie {
	return &http.Cookie{
		Name:  CookieName,
		Value: sessionID,
	}

}

// login verifies that the credentials are valid and returns a populated user
// if they are valid with a nil error. If the error is non-nil credential
// validation failed
func (a Auth) login(username, password string) (user User, err error) {
	user.Username = username
	if err = a.DB.QueryRow("SELECT id, password FROM users where username = $1", username).Scan(&user.id, &user.Password); err != nil {
		return user, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return user, err

}

// createSession creates a session id and adds to the database, and returns
// the created session id.
func (a Auth) createSession(user User) (sessionID string, err error) {

	sessNum, err := rand.Int(rand.Reader, big.NewInt(randMax))
	if err == nil {
		sessionID = sessNum.String()
		_, err = a.DB.Exec("INSERT INTO sessions(session_id, user_id) VALUES($1, $2)", sessionID, user.id)
	}
	return sessionID, err
}

// revokeSession removes the argument from the database.
func (a Auth) revokeSession(sessionID string) error {
	_, err := a.DB.Exec("DELETE FROM sessions where session_id = $1", sessionID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return err
}
