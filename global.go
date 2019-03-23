package gocrid

import "net/http"

var DefaultManager *Manager = NewManager(nil, nil)

func Handle(pattern string, handler Handler) error {
	return DefaultManager.Handle(pattern, handler)
}

func HandleFunc(pattern string,
	handler func(*http.Request, *Context)) error {
	return DefaultManager.HandleFunc(pattern, handler)
}

func ParseRequest(w http.ResponseWriter, r *http.Request) (
	context *Context, err error) {
	return DefaultManager.ParseRequest(w, r)
}
