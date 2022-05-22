package server

import (
	"context"
	"log"
	"net"

	"github.com/che-kwas/iam-kit/config"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	*grpc.Server
	addr string
}

var _ Servable = &GRPCServer{}

// NewGRPCServer builds an GRPCServer.
func NewGRPCServer(opts *config.GRPCOptions) *GRPCServer {
	grpcOpts := []grpc.ServerOption{grpc.MaxRecvMsgSize(opts.MaxMsgSize)}
	server := grpc.NewServer(grpcOpts...)

	return &GRPCServer{Server: server, addr: opts.Addr}
}

// Run runs the HTTP server and conducts a self health check.
func (s *GRPCServer) Run() error {
	log.Printf("[gRPC] server start to listening on %s", s.addr)

	listen, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	return s.Serve(listen)
}

// Shutdown shuts down the HTTP server.
func (s *GRPCServer) Shutdown(_ context.Context) error {
	s.GracefulStop()
	return nil
}
