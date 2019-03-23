package main

import (
	"fmt"
	"net/http"

	"github.com/donyori/gocrid"
)

func defaultHandler(r *http.Request, ctx *gocrid.Context) {
	if r.URL.Path != "/" {
		ctx.AfterWriteResp(func(w http.ResponseWriter) bool {
			http.NotFound(w, r)
			return false
		})
		return
	}
	if ctx.IsLogin() {
		id, err := ctx.GetId()
		if err != nil {
			ctx.AfterWriteResp(func(w http.ResponseWriter) bool {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return false
			})
			return
		}
		username, err := ctx.GetUsername()
		if err != nil {
			ctx.AfterWriteResp(func(w http.ResponseWriter) bool {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return false
			})
			return
		}
		host, err := ctx.GetHost()
		if err != nil {
			ctx.AfterWriteResp(func(w http.ResponseWriter) bool {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return false
			})
			return
		}
		ctx.AfterWriteResp(func(w http.ResponseWriter) bool {
			fmt.Fprintln(w, "ID:", id)
			fmt.Fprintln(w, "Username:", username)
			fmt.Fprintln(w, "Host:", host)
			return true
		})
	} else {
		ctx.AfterWriteResp(func(w http.ResponseWriter) bool {
			fmt.Fprintln(w, "Not login.")
			return true
		})
	}
}

func loginHandler(r *http.Request, ctx *gocrid.Context) {
	if ctx.IsLogin() {
		ctx.AfterWriteResp(func(w http.ResponseWriter) bool {
			fmt.Fprintln(w, "Already login.")
			return true
		})
	} else {
		err := r.ParseForm()
		if err != nil {
			ctx.AfterWriteResp(func(w http.ResponseWriter) bool {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return false
			})
			return
		}
		username := r.FormValue("username")
		err = ctx.Login(username)
		if err != nil {
			ctx.AfterWriteResp(func(w http.ResponseWriter) bool {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return false
			})
			return
		}
		ctx.AfterWriteResp(func(w http.ResponseWriter) bool {
			fmt.Fprintln(w, "Login successfully.")
			return true
		})
		id, err := ctx.GetId()
		if err != nil {
			ctx.AfterWriteResp(func(w http.ResponseWriter) bool {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return false
			})
			return
		}
		ctx.AfterWriteResp(func(w http.ResponseWriter) bool {
			fmt.Fprintln(w, "ID:", id)
			return true
		})
		host, err := ctx.GetHost()
		if err != nil {
			ctx.AfterWriteResp(func(w http.ResponseWriter) bool {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return false
			})
			return
		}
		ctx.AfterWriteResp(func(w http.ResponseWriter) bool {
			fmt.Fprintln(w, "Host:", host)
			return true
		})
	}
}

func logoutHandler(r *http.Request, ctx *gocrid.Context) {
	if ctx.IsLogin() {
		err := ctx.Logout()
		if err != nil {
			ctx.AfterWriteResp(func(w http.ResponseWriter) bool {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return false
			})
			return
		}
		ctx.AfterWriteResp(func(w http.ResponseWriter) bool {
			fmt.Fprintln(w, "Logout successfully.")
			return true
		})
	} else {
		ctx.AfterWriteResp(func(w http.ResponseWriter) bool {
			fmt.Fprintln(w, "Not login.")
			return true
		})
	}
}
