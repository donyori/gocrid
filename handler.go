package gocrid

import (
	"errors"
	"net/http"
)

type Handler interface {
	Handle(*http.Request, *Context)
}

type HandlerFunc func(*http.Request, *Context)

var ErrNilHandler error = errors.New("gocrid: Handler is nil")

func (f HandlerFunc) Handle(r *http.Request, ctx *Context) {
	f(r, ctx)
}
