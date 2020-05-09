package mux

import (
	"context"
	"net/http"
)

type ContextKey string

const (
	ParamsKey = ContextKey("params")
	PathKey   = ContextKey("path")
)

type httpRouter struct {
	tree            map[string]*trie
	notFoundHandler http.Handler
}

func (hr *httpRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	trie := hr.tree[r.Method]

	if trie == nil {
		hr.notFoundHandler.ServeHTTP(w, r)
		return
	}

	value, params := trie.handlerForPath(r.URL.Path)

	if value == nil {
		hr.notFoundHandler.ServeHTTP(w, r)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), PathKey, value.path))
	if len(params) > 0 {
		r = r.WithContext(context.WithValue(r.Context(), ParamsKey, params))
	}

	value.handler.ServeHTTP(w, r)
}

type Config struct {
	RootGroup            *Group
	NotFoundHandler      http.Handler
	CommonPrefix         string
	PathSegmentValidator func(string) bool
	GlobalMiddlewares    []Middleware
}

func NewRouter(cfg *Config) (http.Handler, []*Endpoint) {
	if cfg == nil || cfg.NotFoundHandler == nil || cfg.RootGroup == nil {
		return nil, nil
	}

	router := &httpRouter{
		tree:            make(map[string]*trie),
		notFoundHandler: cfg.NotFoundHandler,
	}

	var badEndpoints []*Endpoint

	walkGroup(cfg.RootGroup, cfg.GlobalMiddlewares, func(mws []Middleware, e *Endpoint) {
		if httpMethods[e.Method] {
			tree := router.tree[e.Method]

			if tree == nil {
				tree = newPathTrie(cfg.CommonPrefix, cfg.PathSegmentValidator)

				router.tree[e.Method] = tree
			}

			if tree.insertPathHandler(e.Path, chainHandler(e.Handler, mws...)) {
				return
			}
		}

		badEndpoints = append(badEndpoints, e)
	})

	return router, badEndpoints
}
