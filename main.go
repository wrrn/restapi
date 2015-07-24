package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/login/", func(rw http.ResponseWriter, r *http.Request) {
	})

	http.HandleFunc("/logout/", func(rw http.ResponseWriter, r *http.Request) {
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
