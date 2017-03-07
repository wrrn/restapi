package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/wrrn/restapi/auth"
	"github.com/wrrn/restapi/configuration"
	"github.com/wrrn/restapi/configuration/confighandler"
)

func setupDB(driver, dataSource string) *sql.DB {
	db, err := sql.Open(driver, dataSource)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func main() {
	var (
		driver     string
		dataSource string
		port       port = 80
	)

	flag.StringVar(&driver, "db-driver", "", "The driver to use: sqlite3 or postgres")
	flag.StringVar(&dataSource, "db-source", "", "Database source configurations. (Could contain username, password, and database name)")
	flag.Var(&port, "port", "Port for server to listen on")
	flag.Parse()
	db := setupDB(driver, dataSource)

	var configHandler http.Handler = confighandler.Handler{
		configuration.ConfigurationController{db},
	}

	configHandler = auth.Auth{db}.ValidateTokens(configHandler)

	mux := http.NewServeMux()

	mux.Handle("/configurations/", http.StripPrefix("/configurations", configHandler))

	log.Fatal(http.ListenAndServe(":"+port.String(), mux))

}
