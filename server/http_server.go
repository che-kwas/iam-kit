package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"

	"github.com/che-kwas/iam-kit/httputil"
)

const (
	ConfKeyHTTP = "http"

	DefaultHTTPMode        = "release"
	DefaultHTTPAddr        = "0.0.0.0:8000"
	DefaultHTTPPingTimeout = time.Duration(10 * time.Second)
	DefaultHealthz         = true
	DefaultMetrics         = false
	DefaultProfiling       = false

	RouterVersion   = "/version"
	RouterHealthz   = "/healthz"
	RouterMetrics   = "/metrics"
	RouterProfiling = "/debug/pprof"
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

// HTTPServer is both a HTTPServer and a gin.Engine.
type HTTPServer struct {
	*gin.Engine

	addr        string
	server      *http.Server
	middlewares []string
	healthz     bool
	metrics     bool
	profiling   bool
	pingTimeout time.Duration
}

var _ Servable = &HTTPServer{}

// NewHTTPServer builds an HTTPServer.
func NewHTTPServer() (*HTTPServer, error) {
	opts, err := getHTTPOptions()
	if err != nil {
		return nil, err
	}

	gin.SetMode(opts.Mode)

	s := &HTTPServer{
		Engine:      gin.New(),
		addr:        opts.Addr,
		middlewares: opts.Middlewares,
		healthz:     opts.Healthz,
		metrics:     opts.Metrics,
		profiling:   opts.Profiling,
		pingTimeout: opts.PingTimeout,
	}

	s.setupMiddlewares()
	s.setupAPIs()

	return s, nil
}

// Run runs the HTTP server and conducts a self health check.
func (s *HTTPServer) Run() error {
	s.server = &http.Server{
		Addr:    s.addr,
		Handler: s,
	}

	var eg errgroup.Group
	eg.Go(func() error {
		log.Printf("[HTTP] server start to listening on %s", s.addr)

		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}

		log.Printf("[HTTP] server on %s stopped", s.addr)
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
			httputil.WriteResponse(c, nil, map[string]string{"status": "OK"})
		})
	}

	// TODO: setup metrics, profiling, and version
}

func (s *HTTPServer) ping(ctx context.Context) error {
	url := fmt.Sprintf("http://%s%s", s.addr, RouterHealthz)
	url = strings.Replace(url, "0.0.0.0", "127.0.0.1", 1)

	for {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return err
		}

		resp, err := http.DefaultClient.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			log.Print("Server self health check success.")
			resp.Body.Close()
			return nil
		}

		log.Print("Waiting for the router, retry in 1 second.")
		time.Sleep(1 * time.Second)

		select {
		case <-ctx.Done():
			return errors.New("self health check of the http server failed")
		default:
		}
	}
}

func getHTTPOptions() (*HTTPOptions, error) {
	opts := &HTTPOptions{
		Mode:        DefaultHTTPMode,
		Addr:        DefaultHTTPAddr,
		Middlewares: []string{},
		PingTimeout: DefaultHTTPPingTimeout,
		Healthz:     DefaultHealthz,
		Metrics:     DefaultMetrics,
		Profiling:   DefaultProfiling,
	}

	if err := viper.UnmarshalKey(ConfKeyHTTP, opts); err != nil {
		return nil, err
	}
	return opts, nil
}
