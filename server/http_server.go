package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"

	"github.com/che-kwas/iam-kit/util"
)

const (
	ConfigHTTPKey = "http"

	DefaultMode        = "release"
	DefaultHTTPAddr    = "0.0.0.0:8000"
	DefaultHealthz     = true
	DefaultMetrics     = false
	DefaultProfiling   = false
	DefaultPingTimeout = "10s"

	RouterVersion   = "/version"
	RouterHealthz   = "/healthz"
	RouterMetrics   = "/metrics"
	RouterProfiling = "/debug/pprof"
)

// HTTPServerBuilder can build an HTTPServer.
type HTTPServerBuilder struct {
	Mode        string
	Address     string
	Middlewares []string
	Healthz     bool
	Metrics     bool
	Profiling   bool
	PingTimeout string `mapstructure:"ping-timeout"`

	err error
}

// NewHTTPServerBuilder is used to build an HTTPServer.
func NewHTTPServerBuilder() *HTTPServerBuilder {
	b := &HTTPServerBuilder{
		Mode:        DefaultMode,
		Address:     DefaultHTTPAddr,
		Middlewares: []string{},
		Healthz:     DefaultHealthz,
		Metrics:     DefaultMetrics,
		Profiling:   DefaultProfiling,
		PingTimeout: DefaultPingTimeout,
	}
	b.err = viper.UnmarshalKey(ConfigHTTPKey, b)

	return b
}

// Build builds an HTTPServer.
func (b *HTTPServerBuilder) Build() (*HTTPServer, error) {
	if b.err != nil {
		return nil, b.err
	}
	var pingTimeout time.Duration
	pingTimeout, b.err = time.ParseDuration(b.PingTimeout)
	if b.err != nil {
		return nil, b.err
	}

	gin.SetMode(b.Mode)

	s := &HTTPServer{
		Engine:      gin.New(),
		address:     b.Address,
		middlewares: b.Middlewares,
		healthz:     b.Healthz,
		metrics:     b.Metrics,
		profiling:   b.Profiling,
		pingTimeout: pingTimeout,
	}

	s.setupMiddlewares()
	s.setupAPIs()

	return s, nil
}

// HTTPServer is both a HTTPServer and a gin.Engine.
type HTTPServer struct {
	*gin.Engine

	address     string
	server      *http.Server
	middlewares []string
	healthz     bool
	metrics     bool
	profiling   bool
	pingTimeout time.Duration
}

var _ Servable = &HTTPServer{}

// Run runs the HTTP server and conducts a self health check.
func (s *HTTPServer) Run() error {
	s.server = &http.Server{
		Addr:    s.address,
		Handler: s,
	}

	var eg errgroup.Group
	eg.Go(func() error {
		log.Infof("[HTTP] server start to listening on %s", s.address)

		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}

		log.Infof("[HTTP] server on %s stopped", s.address)
		return nil
	})

	// self health check
	ctx, cancel := context.WithTimeout(context.Background(), s.pingTimeout)
	defer cancel()
	if s.healthz {
		if err := s.ping(ctx); err != nil {
			return err
		}
	}

	return eg.Wait()
}

// Shutdown shuts down the HTTP server.
func (s *HTTPServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *HTTPServer) setupMiddlewares() {}

func (s *HTTPServer) setupAPIs() {
	if s.healthz {
		s.GET(RouterHealthz, func(c *gin.Context) {
			util.WriteResponse(c, nil, map[string]string{"status": "OK"})
		})
	}

	// TODO: setup metrics, profiling, and version
}

func (s *HTTPServer) ping(ctx context.Context) error {
	url := fmt.Sprintf("http://%s%s", s.address, RouterHealthz)
	url = strings.Replace(url, "0.0.0.0", "127.0.0.1", 1)

	for {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return err
		}

		resp, err := http.DefaultClient.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			log.Debug("Server self health check success.")
			resp.Body.Close()
			return nil
		}

		log.Debug("Waiting for the router, retry in 1 second.")
		time.Sleep(1 * time.Second)

		select {
		case <-ctx.Done():
			return errors.New("self health check of the http server failed")
		default:
		}
	}
}
