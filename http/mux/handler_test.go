package mux

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/FantLab/go-kit/assert"
)

func Test_Router(t *testing.T) {
	t.Run("bad config", func(t *testing.T) {
		{
			handler, badEndpoints := NewRouter(nil)
			assert.True(t, handler == nil && badEndpoints == nil)
		}
		{
			handler, badEndpoints := NewRouter(&Config{})
			assert.True(t, handler == nil && badEndpoints == nil)
		}
	})

	router, invalidEndpoints := makeTestRouter()

	t.Run("invalid endpoints", func(t *testing.T) {
		assert.True(t, len(invalidEndpoints) == 1 && invalidEndpoints[0].Path == "/blog_topics/:id/message")
	})

	ts := httptest.NewServer(router)
	defer ts.Close()

	t.Run("override", func(t *testing.T) {
		{
			resp := sendTestRequest(http.MethodGet, ts.URL+"/v1/work/1/info", false)
			assert.True(t, resp == "work 11")
		}
	})

	t.Run("success", func(t *testing.T) {
		{
			resp := sendTestRequest(http.MethodPost, ts.URL+"/v1/auth/login", false)
			assert.True(t, resp == "login")
		}
		{
			resp := sendTestRequest(http.MethodGet, ts.URL+"/v1/forums/qwe", false)
			assert.True(t, resp == "forum qwe")
		}
		{
			resp := sendTestRequest(http.MethodPut, ts.URL+"/v1/topics/kek/subscription", true)
			assert.True(t, resp == "topic subscription kek")
		}
	})

	t.Run("fail", func(t *testing.T) {
		{
			resp := sendTestRequest(http.MethodGet, ts.URL+"/v1/auth/login", false)
			assert.True(t, resp == "not found")
		}
		{
			resp := sendTestRequest(http.MethodDelete, ts.URL+"/v1/auth/login", false)
			assert.True(t, resp == "not found")
		}
		{
			resp := sendTestRequest(http.MethodDelete, ts.URL+"/v1/forums/qwe", false)
			assert.True(t, resp == "not found")
		}
		{
			resp := sendTestRequest(http.MethodPut, ts.URL+"/v1/topics/kek/subscription", false)
			assert.True(t, resp == "bad auth")
		}
	})
}

func checkTestAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("auth") != "1" {
			w.WriteHeader(400)
			_, _ = w.Write([]byte("bad auth"))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func sendTestRequest(method, url string, auth bool) string {
	req, err := http.NewRequest(method, url, nil)
	if auth {
		req.Header.Add("auth", "1")
	}
	if err != nil {
		return ""
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return ""
	}
	body, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	return string(body)
}

func makeTestRouter() (http.Handler, []*Endpoint) {
	getPathParamsFromContext := func(ctx context.Context, valueKey string) (value string, exists bool) {
		if values, ok := ctx.Value(ParamsKey).(map[string]string); ok {
			if values != nil {
				value, exists = values[valueKey]
			}
		}
		return
	}

	cfg := &Config{
		RootGroup: new(Group),
		NotFoundHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("not found"))
		}),
		CommonPrefix: "v1",
		PathSegmentValidator: func(s string) bool {
			for _, r := range s {
				if r < 'a' || r > 'z' {
					return false
				}
			}
			return true
		},
		GlobalMiddlewares: nil,
	}

	cfg.RootGroup.Subgroup(func(g *Group) {
		g.Endpoint(http.MethodPost, "/auth/login", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			_, _ = getPathParamsFromContext(r.Context(), "test")
			_, _ = w.Write([]byte("login"))
		}))
		g.Endpoint(http.MethodGet, "/forums", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			_, _ = w.Write([]byte("forums"))
		}))
		g.Endpoint(http.MethodGet, "/forums/:forum_id", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			fid, _ := getPathParamsFromContext(r.Context(), "forum_id")
			_, _ = w.Write([]byte("forum " + fid))
		}))
		g.Endpoint(http.MethodGet, "/topics/:topic_id", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			tid, _ := getPathParamsFromContext(r.Context(), "topic_id")
			_, _ = w.Write([]byte("topic " + tid))
		}))
		g.Endpoint(http.MethodGet, "/work/:work_id/info", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			wid, _ := getPathParamsFromContext(r.Context(), "work_id")
			_, _ = w.Write([]byte("work " + wid))
		}))
		// override previous
		g.Endpoint(http.MethodGet, "/work/:work_id_2/info", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			wid, _ := getPathParamsFromContext(r.Context(), "work_id")
			wid2, _ := getPathParamsFromContext(r.Context(), "work_id_2")
			_, _ = w.Write([]byte("work " + wid + wid2))
		}))

		g.Subgroup(func(g *Group) {
			g.Middleware(checkTestAuthMiddleware)

			g.Endpoint(http.MethodPost, "/auth/refresh", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				_, _ = w.Write([]byte("refreshed"))
			}))
			g.Endpoint(http.MethodGet, "/work/:work_id/userclassification", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				wid, _ := getPathParamsFromContext(r.Context(), "work_id")
				_, _ = w.Write([]byte("userclassification for work " + wid))
			}))
			g.Endpoint(http.MethodPost, "/blog_topics/:id/message", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// bad endpoint
			}))
			g.Endpoint(http.MethodPut, "/topics/:topic_id/subscription", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				tid, _ := getPathParamsFromContext(r.Context(), "topic_id")
				_, _ = w.Write([]byte("topic subscription " + tid))
			}))
		})
	})

	return NewRouter(cfg)
}
