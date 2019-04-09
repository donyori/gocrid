package main

import (
	"log"
	"net/http"

	"github.com/donyori/gocrid"
)

func main() {
	gocrid.HandleFunc("/login", loginHandler)
	gocrid.HandleFunc("/logout", logoutHandler)
	gocrid.HandleFunc("/", defaultHandler)
	log.Fatal(http.ListenAndServe(":4551", nil))
}
