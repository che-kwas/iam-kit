package server

import (
	"context"
	"net"

	"github.com/spf13/viper"
	"google.golang.org/grpc"

	"github.com/che-kwas/iam-kit/logger"
)

const (
	confKeyGRPC = "grpc"

	defaultGRPCAddr       = "0.0.0.0:8001"
	defaultGRPCMaxMsgSize = 4 * 1024 * 1024
)

// GRPCOptions defines options for building a GRPCServer.
type GRPCOptions struct {
	Addr       string
	MaxMsgSize int `mapstructure:"max-msg-size"`
}

type GRPCServer struct {
	*grpc.Server
	addr string
}

var _ Servable = &GRPCServer{}

// NewGRPCServer builds an GRPCServer.
func NewGRPCServer() (*GRPCServer, error) {
	opts, err := getGRPCOptions()
	if err != nil {
		return nil, err
	}
	logger.L().Debugf("New grpc server with options: %+v", opts)

	grpcOpts := []grpc.ServerOption{grpc.MaxRecvMsgSize(opts.MaxMsgSize)}
	server := grpc.NewServer(grpcOpts...)

	return &GRPCServer{Server: server, addr: opts.Addr}, nil
}

// Run runs the HTTP server and conducts a self health check.
func (s *GRPCServer) Run() error {
	logger.L().Infof("[gRPC] server start to listening on %s", s.addr)

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

func getGRPCOptions() (*GRPCOptions, error) {
	opts := &GRPCOptions{Addr: defaultGRPCAddr, MaxMsgSize: defaultGRPCMaxMsgSize}

	if err := viper.UnmarshalKey(confKeyGRPC, opts); err != nil {
		return nil, err
	}
	return opts, nil
}
