package anyserver

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/FantLab/go-kit/assert"
)

func serverFromTestHTTP(ts *httptest.Server, shutdownTimeout time.Duration) *Server {
	return &Server{
		Start: func() error {
			ts.Start()
			return nil
		},
		Stop: func(ctx context.Context) error {
			ts.CloseClientConnections()
			ts.Close()
			return nil
		},
		ShutdownTimeout: shutdownTimeout,
	}
}

func Test_runServers(t *testing.T) {
	t.Run("start error", func(t *testing.T) {
		server := &Server{
			Start: func() error {
				return errors.New("start error")
			},
			Stop: func(ctx context.Context) error {
				return nil
			},
		}

		var errToCheck error

		runServers([]*Server{server}, context.Background().Done(), func(err error) {
			errToCheck = err
		})

		assert.True(t, errToCheck.Error() == "start error")
	})

	t.Run("stop error", func(t *testing.T) {
		server := &Server{
			Start: func() error {
				return nil
			},
			Stop: func(ctx context.Context) error {
				return errors.New("stop error")
			},
		}

		ctx, cancel := context.WithCancel(context.Background())

		go cancel()

		var errToCheck error

		runServers([]*Server{server}, ctx.Done(), func(err error) {
			errToCheck = err
		})

		assert.True(t, errToCheck.Error() == "stop error")
	})

	t.Run("setup error", func(t *testing.T) {
		server := &Server{
			SetupError: errors.New("setup error"),
		}

		var errToCheck error

		runServers([]*Server{server}, context.Background().Done(), func(err error) {
			errToCheck = err
		})

		assert.True(t, errToCheck.Error() == "setup error")
	})

	t.Run("dispose", func(t *testing.T) {
		server := &Server{
			DisposeBag: []func() error{
				func() error {
					return errors.New("in dispose")
				},
			},
		}

		var errToCheck error

		runServers([]*Server{server}, context.Background().Done(), func(err error) {
			errToCheck = err
		})

		assert.True(t, errToCheck.Error() == "in dispose")
	})

	t.Run("single server", func(t *testing.T) {
		var x uint32

		ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			select {
			case <-r.Context().Done():
				return
			case <-time.After(2 * time.Second):
				atomic.StoreUint32(&x, 1)
				return
			}
		}))

		server := serverFromTestHTTP(ts, 100*time.Millisecond)

		ctx, cancel := context.WithCancel(context.Background())

		go func() {
			time.Sleep(50 * time.Millisecond)
			cancel()
		}()

		go func() {
			time.Sleep(20 * time.Millisecond)
			_, _ = http.Get(ts.URL)
		}()

		runServers([]*Server{server}, ctx.Done(), func(err error) {})

		assert.True(t, atomic.LoadUint32(&x) == 0)
	})

	t.Run("multiple servers", func(t *testing.T) {
		var x, y uint32

		ts1 := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			atomic.StoreUint32(&x, 10)
		}))

		ts2 := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			atomic.StoreUint32(&y, 20)
		}))

		s1 := serverFromTestHTTP(ts1, 100*time.Millisecond)
		s2 := serverFromTestHTTP(ts2, 100*time.Millisecond)

		ctx, cancel := context.WithCancel(context.Background())

		go func() {
			time.Sleep(50 * time.Millisecond)
			cancel()
		}()

		go func() {
			time.Sleep(20 * time.Millisecond)
			_, _ = http.Get(ts1.URL)
			_, _ = http.Get(ts2.URL)
		}()

		runServers([]*Server{s1, nil, s2}, ctx.Done(), func(err error) {})

		assert.True(t, atomic.LoadUint32(&x) == 10)
		assert.True(t, atomic.LoadUint32(&y) == 20)
	})
}
