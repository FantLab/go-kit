package mux

import "net/http"

type Endpoint struct {
	Method  string
	Path    string
	Handler http.Handler
}

type Middleware func(http.Handler) http.Handler

type Group struct {
	Middlewares []Middleware
	Endpoints   []*Endpoint
	Subgroups   []*Group
}
