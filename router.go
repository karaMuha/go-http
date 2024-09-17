package http

type Router struct {
	routes map[string]handler
}

type handler func(res *HttpResponse, req *HttpRequest)

func NewRouter() *Router {
	return &Router{
		routes: make(map[string]handler),
	}
}

func (r *Router) HandleFunc(target string, handleFunc func(*HttpResponse, *HttpRequest)) {
	// TODO: check if given string is a valid target
	r.routes[target] = handleFunc
}
