// Package logger is the logger builder.
package logger

import (
	"log"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	ConfKey = "log"

	DefaultLevel      = "debug"
	DefaultEncoding   = "console"
	DefaultOutputPath = "/var/log/iam-apiserver/iam-apiserver.log"
	DefaultMaxSize    = 100
	DefaultMaxAge     = 30
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

// NewLogger creates a zap sugared logger.
func NewLogger() (*zap.SugaredLogger, error) {
	opts, err := getLogOpts()
	if err != nil {
		return nil, err
	}
	log.Printf("NewLogger, opts: %+v", opts)

	encoder := newEncoder(opts.Encoding)
	ws := newWriteSyncer(opts)
	level := newLevel(opts.Level)
	core := zapcore.NewCore(encoder, ws, level)
	logger := zap.New(core, zap.AddCaller())

	return logger.Sugar(), nil
}

func getLogOpts() (*LogOptions, error) {
	opts := &LogOptions{
		Level:      DefaultLevel,
		Encoding:   DefaultEncoding,
		OutputPath: DefaultOutputPath,
		MaxSize:    DefaultMaxSize,
		MaxAge:     DefaultMaxAge,
	}

	if err := viper.UnmarshalKey(ConfKey, opts); err != nil {
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
