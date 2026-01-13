package httpx

import "net/http"

type Router struct {
	routes map[string]map[string]http.HandlerFunc
}

func NewRouter() *Router {
	return &Router{routes: make(map[string]map[string]http.HandlerFunc)}
}

func (r *Router) Handle(method, path string, h http.HandlerFunc) {
	if r.routes[method] == nil {
		r.routes[method] = make(map[string]http.HandlerFunc)
	}
	r.routes[method][path] = h
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if m := r.routes[req.Method]; m != nil {
		if h, ok := m[req.URL.Path]; ok {
			h(w, req)
			return
		}
	}
	http.NotFound(w, req)
}
