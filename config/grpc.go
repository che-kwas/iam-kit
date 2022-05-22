package config

const (
	DefaultGRPCAddr       = "0.0.0.0:8001"
	DefaultGRPCMaxMsgSize = 4 * 1024 * 1024
)

// GRPCOptions defines options for building a GRPCServer.
type GRPCOptions struct {
	Addr       string
	MaxMsgSize int `mapstructure:"max-msg-size"`
}

// DefaultGRPCOptions news an GRPCOptions.
func DefaultGRPCOptions() *GRPCOptions {
	return &GRPCOptions{Addr: DefaultGRPCAddr, MaxMsgSize: DefaultGRPCMaxMsgSize}
}
