package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/wrrn/restapi/auth"
	"github.com/wrrn/restapi/configuration"
	"github.com/wrrn/restapi/configuration/confighandler"
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
		db                         = SetupDB()
		configHandler http.Handler = confighandler.Handler{
			configuration.ConfigurationController{db},
		}
		port       port = 80
	)

	flag.Var(&port, "port", "Port for server to listen on")
	flag.Parse()
	configHandler = auth.Auth{db}.ValidateTokens(configHandler)

	mux := http.NewServeMux()

	mux.Handle("/configurations/", http.StripPrefix("/configurations", configHandler))

	log.Fatal(http.ListenAndServe(":"+port.String(), mux))

}
