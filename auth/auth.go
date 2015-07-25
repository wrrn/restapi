package auth

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"log"
	"math/big"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

const (
	CookieName       = "RESTAPI"
	randMax    int64 = 54050505434503053
)

type User struct {
	id       int
	Username string `json:"username""`
	Password string `json:"password"`
}
type Auth struct {
	*sql.DB
}

func (a Auth) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var (
		user User
		err  error
	)

	json.NewDecoder(r.Body).Decode(&user)
	user, err = a.login(user.Username, user.Password)
	if err != nil {
		Unauthorized(w)
		return
	}
	sessionID, err := a.createSession(user)

	if err != nil {
		log.Println(err)
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	cookie := generateCookie(sessionID)
	http.SetCookie(w, cookie)

	if _, err := w.Write([]byte("Authorized")); err != nil {
		log.Println(err)
		http.Error(w, "Server Error", http.StatusInternalServerError)
	}

}

func Unauthorized(w http.ResponseWriter) {
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}

func generateCookie(sessionID string) *http.Cookie {
	return &http.Cookie{
		Name:  CookieName,
		Value: sessionID,
	}

}

func (a Auth) login(username, password string) (user User, err error) {
	user.Username = username
	if err = a.DB.QueryRow("SELECT id, password FROM users where username = $1", username).Scan(&user.id, &user.Password); err != nil {
		return user, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return user, err

}

func (a Auth) createSession(user User) (sessionID string, err error) {

	sessNum, err := rand.Int(rand.Reader, big.NewInt(randMax))
	if err == nil {
		sessionID = sessNum.String()
		_, err = a.DB.Exec("INSERT INTO sessions(session_id, user_id) VALUES($1, $2)", sessionID, user.id)
	}
	return string(sessionID), err
}
