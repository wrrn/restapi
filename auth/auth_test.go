package auth_test

import (
	"database/sql"
	"log"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/wrrn/restapi/auth"
)

type database struct {
	*sql.DB
}

func TestValidToken(t *testing.T) {
	tests := map[string]struct {
		token   string
		isValid bool
	}{
		"ValidToken": {
			token:   "SECRET",
			isValid: true,
		},
		"InvalidToken": {
			token:   "",
			isValid: false,
		},
	}

	db := newDB()
	defer db.Close()
	db.Exec("DELETE FROM user_account")
	db.Exec("DELETE FROM token")
	if err := db.addUser("joe"); err != nil {
		t.Fatal("Unable to add user:", err)
	}

	if err := db.addToken("joe", "SECRET"); err != nil {
		t.Fatal("Unable to add token:", err)
	}

	for testName, test := range tests {
		authenticator := auth.Auth{db}
		if actual := authenticator.ValidToken(test.token); actual != test.isValid {
			t.Errorf("Test %s failed: Expected: %v, Got: %v\n", testName, test.isValid, actual)
		}
	}

}

func TestInvalidToken(t *testing.T) {

}

func newDB() (db database) {
	d, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal("DB could not be created:", err)
		return
	}

	_, err = d.Exec(`
             CREATE TABLE IF NOT EXISTS user_account (
                 id INTEGER PRIMARY KEY ASC,
                 username TEXT UNIQUE      
             )`)

	if err != nil {
		d.Close()
		log.Fatal("Unable to build the user_account table:", err)
		return
	}

	_, err = d.Exec(`
             CREATE TABLE IF NOT EXISTS token (
                 token TEXT UNIQUE,
                 user_id INTEGER
             );`)

	if err != nil {
		d.Close()
		log.Fatal("Unable to build the token table:", err)
		return
	}

	return database{d}
}

func (db database) addUser(username string) error {
	_, err := db.Exec(`INSERT INTO user_account(username) VALUES(?)`, username)
	return err
}

func (db database) addToken(username, token string) error {
	_, err := db.Exec(`INSERT INTO token(user_id, token) 
                       VALUES( (SELECT id from user_account WHERE username = ?) , ?)`, username, token)
	return err
}
