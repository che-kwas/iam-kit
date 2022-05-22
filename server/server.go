// Package server is a server builder.
package server // import "github.com/che-kwas/iam-kit/server"

import (
	"context"
	"time"

	"github.com/che-kwas/iam-kit/config"
	"github.com/che-kwas/iam-kit/shutdown"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

type Servable interface {
	Run() error
	Shutdown(ctx context.Context) error
}

type Server struct {
	name       string
	httpServer *HTTPServer
	grpcServer *GRPCServer
	gs         *shutdown.GracefulShutdown
}

func NewServer(
	name string,
	httpOpts *config.HTTPOptions,
	grpcOpts *config.GRPCOptions,
) (*Server, error) {
	gs := shutdown.New(10 * time.Second)
	httpServer := NewHTTPServer(httpOpts)
	grpcServer := NewGRPCServer(grpcOpts)
	gs.AddShutdownCallback(shutdown.ShutdownFunc(httpServer.Shutdown))
	gs.AddShutdownCallback(shutdown.ShutdownFunc(grpcServer.Shutdown))

	server := &Server{
		name:       name,
		httpServer: httpServer,
		grpcServer: grpcServer,
		gs:         gs,
	}

	return server, nil
}

func (s *Server) Run() error {
	var eg errgroup.Group
	eg.Go(s.grpcServer.Run)
	eg.Go(s.httpServer.Run)

	s.gs.Start()
	return eg.Wait()
}

func (s *Server) InitRouter(initFunc func(g *gin.Engine)) {
	initFunc(s.httpServer.Engine)
}
