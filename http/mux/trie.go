package mux

import (
	"net/http"
	"strings"
)

type pathHandler struct {
	path    string
	handler http.Handler
}

type node struct {
	value    *pathHandler
	keys     map[string]struct{}
	anyChild *node
	children map[string]*node
}

func (n *node) insertPathHandler(path []string, value *pathHandler) {
	if len(path) == 0 {
		n.value = value

		return
	}

	wildcard, name := parseSegment(path[0])

	var child *node

	if wildcard {
		child = n.anyChild

		if child == nil {
			child = new(node)

			n.anyChild = child
		}

		if child.keys == nil {
			child.keys = make(map[string]struct{})
		}

		child.keys[name] = struct{}{}
	} else {
		if n.children != nil {
			child = n.children[name]
		}

		if child == nil {
			child = new(node)

			if n.children == nil {
				n.children = make(map[string]*node)
			}

			n.children[name] = child
		}
	}

	child.insertPathHandler(path[1:], value)
}

func (n *node) handlerForPath(path []string, saveParam func(key, value string)) *pathHandler {
	if len(path) == 0 {
		return n.value
	}

	name, path := path[0], path[1:]

	if child := n.children[name]; child != nil {
		if value := child.handlerForPath(path, saveParam); value != nil {
			return value
		}
	}

	if n.anyChild != nil {
		if value := n.anyChild.handlerForPath(path, saveParam); value != nil {
			for key := range n.anyChild.keys {
				saveParam(key, name)
			}
			return value
		}
	}

	return nil
}

func parseSegment(s string) (bool, string) {
	if s[0] == ':' {
		return true, s[1:]
	}
	return false, s
}

type trie struct {
	maxDepth         int
	prefix           string
	segmentValidator func(string) bool
	root             *node
}

func (t *trie) insertPathHandler(path string, handler http.Handler) bool {
	if handler == nil {
		return false
	}

	segments := strings.FieldsFunc(path, func(r rune) bool {
		return r == '/'
	})

	if t.segmentValidator != nil {
		for _, segment := range segments {
			wildcard, name := parseSegment(segment)
			if wildcard {
				continue
			}
			if !t.segmentValidator(name) {
				return false
			}
		}
	}

	t.root.insertPathHandler(segments, &pathHandler{
		path:    path,
		handler: handler,
	})

	if len(segments) > t.maxDepth {
		t.maxDepth = len(segments)
	}

	return true
}

func (t *trie) handlerForPath(path string) (*pathHandler, map[string]string) {
	segments := strings.FieldsFunc(path, func(r rune) bool {
		return r == '/'
	})

	if t.prefix != "" {
		if len(segments) == 0 || t.prefix != segments[0] {
			return nil, nil
		}
		segments = segments[1:]
	}

	if len(segments) > t.maxDepth {
		return nil, nil
	}

	var params map[string]string

	handler := t.root.handlerForPath(segments, func(key, value string) {
		if params == nil {
			params = make(map[string]string)
		}
		params[key] = value
	})

	return handler, params
}

func newPathTrie(prefix string, segmentValidator func(string) bool) *trie {
	return &trie{
		prefix:           prefix,
		segmentValidator: segmentValidator,
		root:             &node{},
	}
}
