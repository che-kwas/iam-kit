// Package logger is the logger builder.
package logger

import (
	"context"
	"fmt"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Defines common log fields.
const (
	KeyRequestID string = "requestID"
	KeyUsername  string = "username"
)

// Defines config default values.
const (
	confKey = "log"

	defaultLevel      = "debug"
	defaultEncoding   = "console"
	defaultOutputPath = "/var/log/iam-apiserver/iam-apiserver.log"
	defaultMaxSize    = 100
	defaultMaxAge     = 30
)

// LogOptions defines options for building a logger.
type LogOptions struct {
	Level      string
	Encoding   string
	OutputPath string `mapstructure:"output-path"`
	MaxSize    int    `mapstructure:"max-size"`
	MaxAge     int    `mapstructure:"max-age"`
	MaxBackups int    `mapstructure:"max-backups"`
}

type Logger struct {
	*zap.SugaredLogger
}

// NewLogger creates a logger.
func NewLogger() *Logger {
	return newLoggerWithLevel("")
}

// NewInfoLogger creates a logger with INFO level.
func NewInfoLogger() *Logger {
	return newLoggerWithLevel("info")
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

func newLoggerWithLevel(level string) *Logger {
	opts, err := getLogOpts()
	if err != nil {
		panic(err)
	}
	if level != "" {
		opts.Level = level
	}

	encoder := newEncoder(opts.Encoding)
	ws := newWriteSyncer(opts)
	lv := newLevel(opts.Level)
	core := zapcore.NewCore(encoder, ws, lv)
	logger := &Logger{zap.New(core, zap.AddCaller()).Sugar()}
	logger.Debugf("new logger with options: %+v", opts)

	return logger
}

func getLogOpts() (*LogOptions, error) {
	opts := &LogOptions{
		Level:      defaultLevel,
		Encoding:   defaultEncoding,
		OutputPath: defaultOutputPath,
		MaxSize:    defaultMaxSize,
		MaxAge:     defaultMaxAge,
	}

	if err := viper.UnmarshalKey(confKey, opts); err != nil {
		return nil, err
	}
	return opts, nil
}

func newEncoder(encoding string) zapcore.Encoder {
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncodeLevel = zapcore.CapitalLevelEncoder

	if encoding == "console" {
		return zapcore.NewConsoleEncoder(cfg)
	}

	return zapcore.NewJSONEncoder(cfg)
}

func newWriteSyncer(opts *LogOptions) zapcore.WriteSyncer {
	logger := &lumberjack.Logger{
		Filename:   opts.OutputPath,
		MaxSize:    opts.MaxSize,
		MaxAge:     opts.MaxAge,
		MaxBackups: opts.MaxBackups,
	}

	return zapcore.AddSync(logger)
}

func newLevel(l string) zapcore.Level {
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(l)); err != nil {
		level = zapcore.InfoLevel
	}

	return level
}
