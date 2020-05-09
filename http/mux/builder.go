package mux

import "net/http"

func (g *Group) Middleware(fn func(http.Handler) http.Handler) {
	g.Middlewares = append(g.Middlewares, fn)
}

func (g *Group) Endpoint(method, path string, handler http.Handler) {
	g.Endpoints = append(g.Endpoints, &Endpoint{
		Method:  method,
		Path:    path,
		Handler: handler,
	})
}

func (g *Group) Subgroup(fn func(g *Group)) {
	sg := new(Group)
	fn(sg)
	g.Subgroups = append(g.Subgroups, sg)
}
