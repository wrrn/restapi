package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/warrenharper/restapi/auth"
	"github.com/warrenharper/restapi/configuration"
	"github.com/warrenharper/restapi/configuration/confighandler"
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
		db                          = SetupDB()
		authentication *auth.Auth   = &auth.Auth{db}
		configHandler  http.Handler = confighandler.Handler{
			configuration.ConfigurationController{db},
		}
	)

	configHandler = authentication.VerifySessions(configHandler)

	mux := http.NewServeMux()

	mux.Handle("/configurations/", http.StripPrefix("/configurations", configHandler))

	log.Fatal(http.ListenAndServe(":8080", mux))

}
