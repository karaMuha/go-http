package http

import "net"

type Router struct {
	routes map[string]handler
}

type handler func(conn net.Conn, r *HttpRequest)

func NewRouter() Router {
	return Router{
		routes: make(map[string]handler),
	}
}

func (r *Router) HandleFunc(target string, handleFunc handler) {
	// TODO: check if given string is a valid target
	r.routes[target] = handleFunc
}
