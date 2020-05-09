package mux

import (
	"net/http"
	"testing"

	"github.com/FantLab/go-kit/assert"
)

func Test_ChainHandler(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		var calls []string

		makeMW := func(s string) Middleware {
			return func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					calls = append(calls, s)
					next.ServeHTTP(w, r)
				})
			}
		}

		emptyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

		chainHandler(emptyHandler, makeMW("1"), makeMW("2"), makeMW("3"), makeMW("4")).ServeHTTP(nil, nil)

		assert.DeepEqual(t, calls, []string{"1", "2", "3", "4"})
	})

	t.Run("empty", func(t *testing.T) {
		assert.True(t, chainHandler(nil, nil) == nil)
	})
}

func Test_WalkGroup(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		walkGroup(nil, nil, nil)
		assert.True(t, true)
	})

	t.Run("simple", func(t *testing.T) {
		g := new(Group)

		g.Subgroup(func(g *Group) {
			g.Middleware(nil)
			g.Endpoint("", "1", nil)

			g.Subgroup(func(g *Group) {
				g.Middleware(nil)
				g.Endpoint("", "2", nil)

				g.Subgroup(func(g *Group) {
					g.Middleware(nil)
					g.Endpoint("", "3", nil)

					g.Subgroup(func(g *Group) {
						g.Middleware(nil)
						g.Endpoint("", "4", nil)
					})
				})
			})
		})

		var count int

		walkGroup(g, nil, func(mws []Middleware, e *Endpoint) {
			switch e.Path {
			case "1":
				assert.True(t, len(mws) == 1)
			case "2":
				assert.True(t, len(mws) == 2)
			case "3":
				assert.True(t, len(mws) == 3)
			case "4":
				assert.True(t, len(mws) == 4)
			}
			count++
		})

		assert.True(t, count == 4)
	})
}
