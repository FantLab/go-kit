package mux

import "net/http"

var httpMethods = map[string]bool{
	http.MethodGet:     true,
	http.MethodHead:    true,
	http.MethodPost:    true,
	http.MethodPut:     true,
	http.MethodPatch:   true,
	http.MethodDelete:  true,
	http.MethodConnect: true,
	http.MethodOptions: true,
	http.MethodTrace:   true,
}

func walkGroup(g *Group, mws []Middleware, fn func(mws []Middleware, endpoint *Endpoint)) {
	if g == nil {
		return
	}
	mws = append(mws, g.Middlewares...)
	for _, endpoint := range g.Endpoints {
		if endpoint != nil {
			fn(mws, endpoint)
		}
	}
	for _, sg := range g.Subgroups {
		walkGroup(sg, mws, fn)
	}
}

func chainHandler(endpoint http.Handler, mws ...Middleware) http.Handler {
	n := len(mws)
	if n == 0 || endpoint == nil {
		return endpoint
	}
	h := mws[n-1](endpoint)
	for i := n - 2; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}
