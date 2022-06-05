// Package logger is the logger builder.
package logger

import (
	"context"
	"fmt"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Defines common log fields.
const (
	KeyRequestID string = "requestID"
	KeyUsername  string = "username"
)

// LogOptions defines options for building a logger.
type LogOptions struct {
	Name          string
	Development   bool
	DisableCaller bool `mapstructure:"disable-caller"`
}

type Logger struct {
	*zap.SugaredLogger
}

// NewLogger creates a logger.
func NewLogger() *Logger {
	opts, _ := getLogOpts()

	return newLoggerWithOpts(opts)
}

// NewGormLogger creates a gorm logger.
func NewGormLogger() *Logger {
	opts, _ := getLogOpts()
	opts.Name = "gorm"
	opts.Development = false

	return newLoggerWithOpts(opts)
}

// NewGinLogger creates a gin logger.
func NewGinLogger() *Logger {
	opts, _ := getLogOpts()
	opts.Name = "gin"
	opts.DisableCaller = true

	return newLoggerWithOpts(opts)
}

// X adds requestID and username fields to the logging context.
func (l *Logger) X(ctx context.Context) *Logger {
	if requestID := ctx.Value(KeyRequestID); requestID != nil {
		l = l.With(KeyRequestID, requestID)
	}
	if username := ctx.Value(KeyUsername); username != nil {
		l = l.With(KeyUsername, username)
	}

	return l
}

// With adds a variadic number of fields to the logging context.
func (l *Logger) With(fields ...interface{}) *Logger {
	return &Logger{l.SugaredLogger.With(fields...)}
}

// Named adds a new path segment to the logger's name. Segments are joined by
// periods. By default, Loggers are unnamed.
func (l *Logger) Named(s string) *Logger {
	return &Logger{l.SugaredLogger.Named(s)}
}

// Printf logs a message at level Print on the compatibleLogger.
func (l *Logger) Printf(format string, args ...interface{}) {
	l.Info(fmt.Sprintf(format, args...))
}

func newLoggerWithOpts(opts *LogOptions) *Logger {
	var cfg zap.Config
	if opts.Development {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}
	cfg.DisableCaller = opts.DisableCaller

	zaplog, _ := cfg.Build()
	logger := &Logger{zaplog.Sugar().Named(opts.Name)}
	logger.Debugf("new logger with options: %+v", opts)

	return logger
}

func getLogOpts() (*LogOptions, error) {
	opts := &LogOptions{Development: true}
	if err := viper.UnmarshalKey("log", opts); err != nil {
		return nil, err
	}
	return opts, nil
}
