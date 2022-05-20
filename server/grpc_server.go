package server

import (
	"context"
	"log"
	"net"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

const (
	ConfigGRPCKey         = "grpc"
	DefaultGRPCAddr       = "0.0.0.0:8001"
	DefaultGRPCMaxMsgSize = 4 * 1024 * 1024
)

// GRPCServerBuilder defines options for building a GRPCServer.
type GRPCServerBuilder struct {
	Addr       string
	MaxMsgSize int `mapstructure:"max-msg-size"`

	err error
}

// NewGRPCServerBuilder is used to build an GRPCServer.
func NewGRPCServerBuilder() *GRPCServerBuilder {
	b := &GRPCServerBuilder{Addr: DefaultGRPCAddr, MaxMsgSize: DefaultGRPCMaxMsgSize}
	b.err = viper.UnmarshalKey(ConfigGRPCKey, b)

	return b
}

// Build builds an GRPCServer.
func (b *GRPCServerBuilder) Build() (*GRPCServer, error) {
	if b.err != nil {
		return nil, b.err
	}

	opts := []grpc.ServerOption{grpc.MaxRecvMsgSize(b.MaxMsgSize)}
	server := grpc.NewServer(opts...)

	return &GRPCServer{Server: server, addr: b.Addr}, nil
}

type GRPCServer struct {
	*grpc.Server
	addr string
}

var _ Servable = &GRPCServer{}

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
