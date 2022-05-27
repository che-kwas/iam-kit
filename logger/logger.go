// Package logger is the logger builder.
package logger

import (
	"context"
	"log"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Defines common log fields.
const (
	KeyRequestID   string = "requestID"
	KeyUsername    string = "username"
	KeyWatcherName string = "watcher"
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

type key int

const ctxKey key = 0

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
	opts, err := getLogOpts()
	if err != nil {
		log.Fatal("get log options error: ", err)
	}

	encoder := newEncoder(opts.Encoding)
	ws := newWriteSyncer(opts)
	level := newLevel(opts.Level)
	core := zapcore.NewCore(encoder, ws, level)
	logger := &Logger{zap.New(core, zap.AddCaller()).Sugar()}
	logger.Debugf("NewLogger, opts: %+v", opts)

	return logger
}

// FromContext returns the logger in the context.
func FromContext(ctx context.Context) *Logger {
	if ctx != nil {
		if logger := ctx.Value(ctxKey); logger != nil {
			return logger.(*Logger)
		}
	}

	return NewLogger()
}

// WithContext puts the logger into the context.
func (l *Logger) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKey, l)
}

// X adds iam-specific fields to the loggin context.
func (l *Logger) X(ctx context.Context) *Logger {
	if requestID := ctx.Value(KeyRequestID); requestID != nil {
		l.SugaredLogger = l.With(KeyRequestID, requestID)
	}
	if username := ctx.Value(KeyUsername); username != nil {
		l.SugaredLogger = l.With(KeyUsername, username)
	}
	if watcherName := ctx.Value(KeyWatcherName); watcherName != nil {
		l.SugaredLogger = l.With(KeyWatcherName, watcherName)
	}

	return l
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
