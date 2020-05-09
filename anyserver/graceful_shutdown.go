package anyserver

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Server struct {
	Start           func() error
	Stop            func(context.Context) error
	SetupError      error
	ShutdownTimeout time.Duration
	DisposeBag      []func() error
}

func RunWithGracefulShutdown(errorFunc func(error), servers ...*Server) {
	ctx, cancel := context.WithCancel(context.Background())
	{
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			<-quit
			signal.Stop(quit)
			cancel()
		}()
	}
	runServers(servers, ctx.Done(), errorFunc)
}

func runServers(servers []*Server, quit <-chan struct{}, errorFunc func(error)) {
	wg := new(sync.WaitGroup)
	for _, server := range servers {
		if server == nil {
			continue
		}
		wg.Add(1)
		go func(server *Server) {
			runServer(server, quit, wg.Done, errorFunc)
		}(server)
	}
	wg.Wait()
}

func runServer(server *Server, quit <-chan struct{}, finishFunc func(), errorFunc func(error)) {
	defer func() {
		for _, fn := range server.DisposeBag {
			if err := fn(); err != nil {
				errorFunc(err)
			}
		}
	}()

	if server.SetupError != nil {
		errorFunc(server.SetupError)

		finishFunc()

		return
	}

	if server.Start == nil || server.Stop == nil {
		finishFunc()

		return
	}

	fail := make(chan struct{})

	go func() {
		defer finishFunc()

		select {
		case <-quit:
			break
		case <-fail:
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), server.ShutdownTimeout)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			errorFunc(err)
		}
	}()

	if err := server.Start(); err != nil {
		errorFunc(err)

		fail <- struct{}{}
	}
}
