package httpx

import (
	"context"
	"net/http"
	"strings"
)

type ctxKey string

const pathParamsKey ctxKey = "path_params"

type route struct {
	pattern string
	parts   []string
	handler http.HandlerFunc
}

type Router struct {
	routes map[string][]route
}

func NewRouter() *Router {
	return &Router{routes: make(map[string][]route)}
}

func (r *Router) Handle(method, pattern string, h http.HandlerFunc) {
	parts := splitPath(pattern)
	r.routes[method] = append(r.routes[method], route{
		pattern: pattern,
		parts:   parts,
		handler: h,
	})
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	rs := r.routes[req.Method]
	if len(rs) == 0 {
		http.NotFound(w, req)
		return
	}

	pathParts := splitPath(req.URL.Path)

	for _, rt := range rs {
		params, ok := match(rt.parts, pathParts)
		if !ok {
			continue
		}

		ctx := context.WithValue(req.Context(), pathParamsKey, params)
		rt.handler(w, req.WithContext(ctx))
		return
	}
	http.NotFound(w, req)
}

func Param(r *http.Request, name string) string {
	m, _ := r.Context().Value(pathParamsKey).(map[string]string)
	if m == nil {
		return ""
	}
	return m[name]
}

func splitPath(p string) []string {
	p = strings.TrimSpace(p)
	p = strings.Trim(p, "/")
	if p == "" {
		return []string{}
	}
	return strings.Split(p, "/")
}

func match(patternPaths, pathParts []string) (map[string]string, bool) {
	if len(patternPaths) != len(pathParts) {
		return nil, false
	}
	params := map[string]string{}
	for i := range patternPaths {
		pp := patternPaths[i]
		ap := pathParts[i]

		if strings.HasPrefix(pp, "{") && strings.HasSuffix(pp, "}") {
			key := strings.TrimSuffix(strings.TrimPrefix(pp, "{"), "}")
			if key == "" {
				return nil, false
			}
			params[key] = ap
			continue
		}

		if pp != ap {
			return nil, false
		}
	}
	return params, true
}
