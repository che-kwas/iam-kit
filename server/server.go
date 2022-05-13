package server

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

	err error
}

func NewServer(name string) (*Server, error) {
	server := &Server{
		name: name,
		gs:   shutdown.New(10 * time.Second),
	}

	server.buildHTTPServer().buildGRPCServer()
	return server, server.err
}

func (s *Server) buildHTTPServer() *Server {
	s.httpServer, s.err = NewHTTPServerBuilder().Build()
	s.gs.AddShutdownCallback(shutdown.ShutdownFunc(s.httpServer.Shutdown))

	return s
}

func (s *Server) buildGRPCServer() *Server {
	if s.err != nil {
		return s
	}

	s.grpcServer, s.err = NewGRPCServerBuilder().Build()
	s.gs.AddShutdownCallback(shutdown.ShutdownFunc(s.grpcServer.Shutdown))

	return s
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
