package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/donyori/gocrid"
)

func main() {
	useHttpHandlers := false
	flag.BoolVar(&useHttpHandlers, "http", false,
		"True if use http handlers instead of gocrid handlers.")
	if !useHttpHandlers {
		gocrid.HandleFunc("/login", loginHandler)
		gocrid.HandleFunc("/logout", logoutHandler)
		gocrid.HandleFunc("/", defaultHandler)
	} else {
		http.HandleFunc("/login", httpLoginHandler)
		http.HandleFunc("/logout", httpLogoutHandler)
		http.HandleFunc("/", httpDefaultHandler)
	}
	log.Fatal(http.ListenAndServe(":4551", nil))
}
