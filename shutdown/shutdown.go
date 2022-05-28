// Package shutdown assumes the responsibility of graceful shutdown of the server.
package shutdown // import "github.com/che-kwas/iam-kit/shutdown"

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/che-kwas/iam-kit/logger"
)

// ShutdownCallback is an interface you have to implement for callbacks.
type ShutdownCallback interface {
	OnShutdown(ctx context.Context) error
}

// ShutdownFunc is a helper type, so you can easily provide anonymous functions
// as ShutdownCallbacks like this:
//   AddShutdownCallback(shutdown.ShutdownFunc(your_callback))
type ShutdownFunc func(context.Context) error

func (f ShutdownFunc) OnShutdown(ctx context.Context) error {
	return f(ctx)
}

// GracefulShutdown handles ShutdownCallbacks.
type GracefulShutdown struct {
	callbacks []ShutdownCallback
	signals   []os.Signal
	timeout   time.Duration
}

// New creates a GracefulShutdown.
func New(timeout time.Duration) *GracefulShutdown {
	return &GracefulShutdown{
		callbacks: make([]ShutdownCallback, 0, 3),
		signals:   []os.Signal{os.Interrupt, syscall.SIGTERM},
		timeout:   timeout,
	}
}

// Start starts listening to shutdown signals.
func (gs *GracefulShutdown) Start() {
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, gs.signals...)

		<-stop
		gs.Shutdown()
	}()
}

// AddShutdownCallback adds a ShutdownCallback that will be called when
// shutdown is requested.
func (gs *GracefulShutdown) AddShutdownCallback(callback ShutdownCallback) {
	gs.callbacks = append(gs.callbacks, callback)
}

// Shutdown calls all calllbacks, and wait for them to finish.
func (gs *GracefulShutdown) Shutdown() {
	// os.Exit(0) makes Shutdown not testable, so put the logic into
	// a testable function.
	gs.waitCallbacks()
	os.Exit(0)
}

func (gs *GracefulShutdown) waitCallbacks() {
	ctx, cancel := context.WithTimeout(context.Background(), gs.timeout)
	defer cancel()

	var wg sync.WaitGroup
	for _, callback := range gs.callbacks {
		wg.Add(1)
		go func(callback ShutdownCallback) {
			defer wg.Done()
			if err := callback.OnShutdown(ctx); err != nil {
				logger.L().Fatal("Failed to gracefully shutdown: ", err)
			}
		}(callback)
	}

	wg.Wait()
}
