package http

import "net"

type Router struct {
	routes map[string]handler
}

type handler func(conn net.Conn, r *Request)

func NewRouter() Router {
	return Router{
		routes: make(map[string]handler),
	}
}

func (r *Router) HandleFunc(target string, handleFunc handler) {
	r.routes[target] = handleFunc
}
