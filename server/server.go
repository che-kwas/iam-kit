// Package server is a server builder.
package server // import "github.com/che-kwas/iam-kit/server"

import (
	"context"
	"time"

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

	// err is used to eliminate error checking hell
	err error
}

func NewServer(name string) (*Server, error) {
	server := &Server{
		name: name,
		gs:   shutdown.New(10 * time.Second),
	}
	server.buildHTTP().buildGRPC()

	return server, server.err
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

func (s *Server) buildHTTP() *Server {
	s.httpServer, s.err = NewHTTPServer()
	s.gs.AddShutdownCallback(shutdown.ShutdownFunc(s.httpServer.Shutdown))

	return s
}

func (s *Server) buildGRPC() *Server {
	if s.err != nil {
		return s
	}

	s.grpcServer, s.err = NewGRPCServer()
	s.gs.AddShutdownCallback(shutdown.ShutdownFunc(s.grpcServer.Shutdown))

	return s
}
