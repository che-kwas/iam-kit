package server

import (
	"context"
	"net"

	"github.com/marmotedu/iam/pkg/log"
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
	Address    string
	MaxMsgSize int `mapstructure:"max-msg-size"`

	err error
}

// NewGRPCServerBuilder is used to build an GRPCServer.
func NewGRPCServerBuilder() *GRPCServerBuilder {
	b := &GRPCServerBuilder{Address: DefaultGRPCAddr, MaxMsgSize: DefaultGRPCMaxMsgSize}
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

	return &GRPCServer{Server: server, address: b.Address}, nil
}

type GRPCServer struct {
	*grpc.Server
	address string
}

var _ Servable = &GRPCServer{}

// Run runs the HTTP server and conducts a self health check.
func (s *GRPCServer) Run() error {
	log.Infof("[gRPC] server start to listening on %s", s.address)

	listen, err := net.Listen("tcp", s.address)
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
