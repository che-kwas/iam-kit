// Package server is a server builder.
package server // import "github.com/che-kwas/iam-kit/server"

import (
	"context"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/che-kwas/iam-kit/shutdown"
)

// Servable defines the behavior of a server.
type Servable interface {
	Run() error
	Shutdown(context.Context) error
}

// Server defines the structure of a server.
type Server struct {
	Name       string
	HTTPServer *HTTPServer
	GRPCServer *GRPCServer
	gs         *shutdown.GracefulShutdown

	err error
}

// Option configures the server.
type Option interface {
	apply(*Server)
}

// OptionFunc wraps a func so it satisfies the Option interface.
type OptionFunc func(*Server)

func (o OptionFunc) apply(s *Server) {
	o(s)
}

func NewServer(name string, opts ...Option) (*Server, error) {
	gs := shutdown.New(10 * time.Second)

	httpS, err := NewHTTPServer()
	if err != nil {
		return nil, err
	}
	gs.AddShutdownCallback(shutdown.ShutdownFunc(httpS.Shutdown))

	s := &Server{
		Name:       name,
		HTTPServer: httpS,
		gs:         gs,
	}

	for _, opt := range opts {
		opt.apply(s)
	}

	return s, s.err
}

func WithGRPC() Option {
	return OptionFunc(func(s *Server) {
		if s.err != nil {
			return
		}

		s.GRPCServer, s.err = NewGRPCServer()
		if s.err == nil {
			s.gs.AddShutdownCallback(shutdown.ShutdownFunc(s.GRPCServer.Shutdown))
		}
	})
}

func WithShutdown(sd shutdown.ShutdownFunc) Option {
	return OptionFunc(func(s *Server) {
		if s.err != nil {
			return
		}

		s.gs.AddShutdownCallback(sd)
	})
}

func (s *Server) Run() error {
	var eg errgroup.Group
	eg.Go(s.HTTPServer.Run)
	if s.GRPCServer != nil {
		eg.Go(s.GRPCServer.Run)
	}

	s.gs.Start()
	return eg.Wait()
}
