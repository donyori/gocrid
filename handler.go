package gocrid

type Handler interface {
	Handle(*Context, error)
}

type HandlerFunc func(*Context, error)

func (f HandlerFunc) Handle(ctx *Context, err error) {
	f(ctx, err)
}
