package gocrid

var DefaultManager *Manager = NewManager(nil, nil)

func Handle(pattern string, handler Handler) {
	DefaultManager.Handle(pattern, handler)
}

func HandleFunc(pattern string, handler func(*Context, error)) {
	DefaultManager.HandleFunc(pattern, handler)
}
