package mux

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/FantLab/go-kit/assert"
)

func Test_Trie(t *testing.T) {
	t.Run("trie", func(t *testing.T) {
		trie := newPathTrie("", func(s string) bool {
			return true
		})

		{
			trie.insertPathHandler("x/y/z", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "1")
			}))

			trie.insertPathHandler("x/y/:z", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "2")
			}))

			trie.insertPathHandler("x/:y/z", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "3")
			}))

			trie.insertPathHandler("x/:y/:z", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "4")
			}))

			trie.insertPathHandler(":x/y/z", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "5")
			}))

			trie.insertPathHandler(":x/y/:z", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "6")
			}))

			trie.insertPathHandler(":x/:y/z", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "7")
			}))

			trie.insertPathHandler(":x/:y/:z", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "8")
			}))
		}

		cases := []func(w *httptest.ResponseRecorder){
			// 1
			func(rr *httptest.ResponseRecorder) {
				value, params := trie.handlerForPath("x/y/z")
				value.handler.ServeHTTP(rr, nil)
				assert.True(t, rr.Body.String() == "1\n")
				assert.True(t, len(params) == 0)
			},
			// 2
			func(rr *httptest.ResponseRecorder) {
				value, params := trie.handlerForPath("x/y/9")
				value.handler.ServeHTTP(rr, nil)
				assert.True(t, rr.Body.String() == "2\n")
				assert.True(t, len(params) == 1 && params["z"] == "9")
			},
			// 3
			func(rr *httptest.ResponseRecorder) {
				value, params := trie.handlerForPath("x/9/z")
				value.handler.ServeHTTP(rr, nil)
				assert.True(t, rr.Body.String() == "3\n")
				assert.True(t, len(params) == 1 && params["y"] == "9")
			},
			// 4
			func(rr *httptest.ResponseRecorder) {
				value, params := trie.handlerForPath("x/9/10")
				value.handler.ServeHTTP(rr, nil)
				assert.True(t, rr.Body.String() == "4\n")
				assert.True(t, len(params) == 2 && params["y"] == "9" && params["z"] == "10")
			},
			// 5
			func(rr *httptest.ResponseRecorder) {
				value, params := trie.handlerForPath("1/y/z")
				value.handler.ServeHTTP(rr, nil)
				assert.True(t, rr.Body.String() == "5\n")
				assert.True(t, len(params) == 1 && params["x"] == "1")
			},
			// 6
			func(rr *httptest.ResponseRecorder) {
				value, params := trie.handlerForPath("1/y/2")
				value.handler.ServeHTTP(rr, nil)
				assert.True(t, rr.Body.String() == "6\n")
				assert.True(t, len(params) == 2 && params["x"] == "1" && params["z"] == "2")
			},
			// 7
			func(rr *httptest.ResponseRecorder) {
				value, params := trie.handlerForPath("1/2/z")
				value.handler.ServeHTTP(rr, nil)
				assert.True(t, rr.Body.String() == "7\n")
				assert.True(t, len(params) == 2 && params["x"] == "1" && params["y"] == "2")
			},
			// 8
			func(rr *httptest.ResponseRecorder) {
				value, params := trie.handlerForPath("1/2/3")
				value.handler.ServeHTTP(rr, nil)
				assert.True(t, rr.Body.String() == "8\n")
				assert.True(t, len(params) == 3 && params["x"] == "1" && params["y"] == "2" && params["z"] == "3")
			},
		}

		for i := 0; i < 10; i++ {
			rand.Seed(time.Now().UnixNano())
			rand.Shuffle(len(cases), func(i, j int) { cases[i], cases[j] = cases[j], cases[i] })
			for _, fn := range cases {
				fn(httptest.NewRecorder())
			}
		}
	})
}

func Test_EdgeCases(t *testing.T) {
	emptyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	t.Run("insert nil handler", func(t *testing.T) {
		trie := newPathTrie("", nil)

		assert.True(t, !trie.insertPathHandler("a/b", nil))
	})

	t.Run("existing handler", func(t *testing.T) {
		trie := newPathTrie("", nil)

		assert.True(t, trie.insertPathHandler("x/y/z", emptyHandler))
		assert.True(t, trie.insertPathHandler("x/y/z", emptyHandler))
	})

	t.Run("not found handler", func(t *testing.T) {
		trie := newPathTrie("", nil)
		assert.True(t, trie.insertPathHandler("x/y/z", emptyHandler))

		{
			value, params := trie.handlerForPath("x/y")
			assert.True(t, value == nil && params == nil)
		}
		{
			value, params := trie.handlerForPath("a/b/c/d")
			assert.True(t, value == nil && params == nil)
		}
	})

	t.Run("invalid segment", func(t *testing.T) {
		trie := newPathTrie("", func(s string) bool {
			for _, r := range s {
				if r < 'a' || r > 'z' {
					return false
				}
			}
			return true
		})

		assert.True(t, !trie.insertPathHandler("x/ y/z", emptyHandler))
		assert.True(t, trie.insertPathHandler("q/w/e", emptyHandler))
		assert.True(t, !trie.insertPathHandler("q/w/e/r/ ", emptyHandler))
	})

	t.Run("prefix", func(t *testing.T) {
		trie := newPathTrie("v1", func(s string) bool {
			return true
		})

		assert.True(t, trie.insertPathHandler("x/y/z", emptyHandler))

		{
			value, params := trie.handlerForPath("x/y/z")
			assert.True(t, value == nil && params == nil)
		}
		{
			value, params := trie.handlerForPath("v1/x/y/z")
			assert.True(t, value != nil && params == nil)
		}
	})

	t.Run("multi param names", func(t *testing.T) {
		trie := newPathTrie("", func(s string) bool {
			return true
		})

		assert.True(t, trie.insertPathHandler("x/:y1/z", emptyHandler))
		assert.True(t, trie.insertPathHandler("x/:y2/z", emptyHandler))

		value, params := trie.handlerForPath("x/y/z")
		assert.True(t, value != nil)
		assert.DeepEqual(t, params, map[string]string{
			"y1": "y",
			"y2": "y",
		})
	})

	t.Run("root handler", func(t *testing.T) {
		{
			trie := newPathTrie("", func(s string) bool {
				return true
			})
			assert.True(t, trie.insertPathHandler("/", emptyHandler))
			value, params := trie.handlerForPath("/")
			assert.True(t, value != nil && params == nil)
			assert.True(t, trie.root.value != nil)
		}

		{
			trie := newPathTrie("v1", func(s string) bool {
				return true
			})
			assert.True(t, trie.insertPathHandler("/", emptyHandler))
			value, params := trie.handlerForPath("v1/")
			assert.True(t, value != nil && params == nil)
			assert.True(t, trie.root.value != nil)
		}
	})
}
