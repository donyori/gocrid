package main

import (
	"fmt"
	"net/http"

	"github.com/donyori/gocrid"
)

func defaultHandler(ctx *gocrid.Context, err error) {
	if err != nil {
		ctx.RespA(func(w http.ResponseWriter) bool {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return false
		})
		return
	}
	r := ctx.GetRequest()
	if r.URL.Path != "/" {
		ctx.RespA(func(w http.ResponseWriter) bool {
			http.NotFound(w, r)
			return false
		})
		return
	}
	if ctx.IsLogin() {
		id, err := ctx.GetId()
		if err != nil {
			ctx.RespA(func(w http.ResponseWriter) bool {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return false
			})
			return
		}
		username, err := ctx.GetUsername()
		if err != nil {
			ctx.RespA(func(w http.ResponseWriter) bool {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return false
			})
			return
		}
		host, err := ctx.GetHost()
		if err != nil {
			ctx.RespA(func(w http.ResponseWriter) bool {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return false
			})
			return
		}
		ctx.RespA(func(w http.ResponseWriter) bool {
			fmt.Fprintln(w, "ID:", id)
			fmt.Fprintln(w, "Username:", username)
			fmt.Fprintln(w, "Host:", host)
			return true
		})
	} else {
		ctx.RespA(func(w http.ResponseWriter) bool {
			fmt.Fprintln(w, "Not login.")
			return true
		})
	}
}

func loginHandler(ctx *gocrid.Context, err error) {
	if err != nil {
		ctx.RespA(func(w http.ResponseWriter) bool {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return false
		})
		return
	}
	if ctx.IsLogin() {
		ctx.RespA(func(w http.ResponseWriter) bool {
			fmt.Fprintln(w, "Already login.")
			return true
		})
	} else {
		r := ctx.GetRequest()
		err = r.ParseForm()
		if err != nil {
			ctx.RespA(func(w http.ResponseWriter) bool {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return false
			})
			return
		}
		username := r.FormValue("username")
		err = ctx.Login(username)
		if err != nil {
			ctx.RespA(func(w http.ResponseWriter) bool {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return false
			})
			return
		}
		ctx.RespA(func(w http.ResponseWriter) bool {
			fmt.Fprintln(w, "Login successfully.")
			return true
		})
		id, err := ctx.GetId()
		if err != nil {
			ctx.RespA(func(w http.ResponseWriter) bool {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return false
			})
			return
		}
		ctx.RespA(func(w http.ResponseWriter) bool {
			fmt.Fprintln(w, "ID:", id)
			return true
		})
		host, err := ctx.GetHost()
		if err != nil {
			ctx.RespA(func(w http.ResponseWriter) bool {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return false
			})
			return
		}
		ctx.RespA(func(w http.ResponseWriter) bool {
			fmt.Fprintln(w, "Host:", host)
			return true
		})
	}
}

func logoutHandler(ctx *gocrid.Context, err error) {
	if err != nil {
		ctx.RespA(func(w http.ResponseWriter) bool {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return false
		})
		return
	}
	if ctx.IsLogin() {
		err = ctx.Logout()
		if err != nil {
			ctx.RespA(func(w http.ResponseWriter) bool {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return false
			})
			return
		}
		ctx.RespA(func(w http.ResponseWriter) bool {
			fmt.Fprintln(w, "Logout successfully.")
			return true
		})
	} else {
		ctx.RespA(func(w http.ResponseWriter) bool {
			fmt.Fprintln(w, "Not login.")
			return true
		})
	}
}
