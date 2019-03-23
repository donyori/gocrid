package main

import (
	"fmt"
	"net/http"

	"github.com/donyori/gocrid"
)

func httpDefaultHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	ctx, err := gocrid.ParseRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = ctx.Write()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if ctx.IsLogin() {
		id, err := ctx.GetId()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		username, err := ctx.GetUsername()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		host, err := ctx.GetHost()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(w, "ID:", id)
		fmt.Fprintln(w, "Username:", username)
		fmt.Fprintln(w, "Host:", host)
	} else {
		fmt.Fprintln(w, "Not login.")
	}
}

func httpLoginHandler(w http.ResponseWriter, r *http.Request) {
	ctx, err := gocrid.ParseRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if ctx.IsLogin() {
		fmt.Fprintln(w, "Already login.")
	} else {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		username := r.FormValue("username")
		err = ctx.Login(username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = ctx.Write()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(w, "Login successfully.")
		id, err := ctx.GetId()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(w, "ID:", id)
		host, err := ctx.GetHost()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(w, "Host:", host)
	}
}

func httpLogoutHandler(w http.ResponseWriter, r *http.Request) {
	ctx, err := gocrid.ParseRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if ctx.IsLogin() {
		err := ctx.Logout()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = ctx.Write()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(w, "Logout successfully.")
	} else {
		fmt.Fprintln(w, "Not login.")
	}
}
