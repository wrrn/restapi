package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/warrenharper/restapi/auth"
	"github.com/warrenharper/restapi/configuration"
	"github.com/warrenharper/restapi/utils/request"
	"github.com/warrenharper/restapi/utils/response"
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

func main() {
	var (
		db                        = SetupDB()
		authentication *auth.Auth = &auth.Auth{db}
		configHandler             = configuration.ConfigurationHandler{
			configuration.ConfigurationController{db},
		}
		_ = configHandler
	)
	authentication.RegisterUser(auth.User{Username: "john_doe", Password: "password"})

	mux := http.NewServeMux()

	// Login
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if !request.Is(r, "POST") {
			response.MethodNotAllowed(w)
			return
		}
		authentication.HandleLogin(w, r)
		return
	})

	//Logout
	mux.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		if !request.Is(r, "POST") {
			response.MethodNotAllowed(w)
			return
		}
		authentication.HandleLogout(w, r)
	})

	mux.Handle("/configurations/", http.StripPrefix("/configurations", configHandler))

	log.Fatal(http.ListenAndServe(":8080", mux))

}
