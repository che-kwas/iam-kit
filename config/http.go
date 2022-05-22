package config

import "time"

const (
	DefaultHTTPMode        = "release"
	DefaultHTTPAddr        = "0.0.0.0:8000"
	DefaultHTTPPingTimeout = time.Duration(10 * time.Second)
	DefaultHealthz         = true
	DefaultMetrics         = false
	DefaultProfiling       = false
)

// HTTPOptions defines options for building an HTTPServer.
type HTTPOptions struct {
	Mode        string
	Addr        string
	Middlewares []string
	PingTimeout time.Duration `mapstructure:"ping-timeout"`
	Healthz     bool
	Metrics     bool
	Profiling   bool
}

// DefaultHTTPOptions news an HTTPOptions.
func DefaultHTTPOptions() *HTTPOptions {
	return &HTTPOptions{
		Mode:        DefaultHTTPMode,
		Addr:        DefaultHTTPAddr,
		Middlewares: []string{},
		PingTimeout: DefaultHTTPPingTimeout,
		Healthz:     DefaultHealthz,
		Metrics:     DefaultMetrics,
		Profiling:   DefaultProfiling,
	}
}
