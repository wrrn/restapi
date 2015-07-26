package main

import (
	"database/sql"
	"log"
	"net/http"
	"strings"

	_ "github.com/lib/pq"
	"github.com/warrenharper/restapi/auth"
)

func SetupDB() *sql.DB {
	db, err := sql.Open("postgres", "user=tenable password=insecure dbname=restapi")
	if err != nil {
		log.Fatal(err)
	}
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM configuration")
	db.Exec("DELETE FROM sessions")
	return db
}

var (
	authentication *auth.Auth = &auth.Auth{SetupDB()}
)

func main() {
	authentication.RegisterUser(auth.User{Username: "john_doe", Password: "password"})

	// Login
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if strings.ToUpper(r.Method) != "POST" {
			http.Error(w, "", http.StatusMethodNotAllowed)
		}
		authentication.HandleLogin(w, r)
		return
	})

	//Logout
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		if strings.ToUpper(r.Method) != "POST" {
			http.Error(w, "", http.StatusMethodNotAllowed)
		}
		authentication.HandleLogout(w, r)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
