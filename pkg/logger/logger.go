package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Stderr ..
const Stderr = "stderr"

// Config represents the logger configuration
type Config struct {
	Name           string
	Level          zapcore.Level
	WithCaller     bool
	WithStacktrace bool
	OutputPaths    []string
	Debug          bool
}

// New will create a logger according to the given Config (or based on the default configuration if no Config is provided)
func New(config ...Config) *zap.Logger {
	var conf Config
	if len(config) > 0 {
		conf = config[0]
	}

	cfg := zap.Config{
		Encoding:          "json",
		Level:             zap.NewAtomicLevelAt(conf.Level),
		ErrorOutputPaths:  []string{"stderr"},
		DisableCaller:     !conf.WithCaller,
		DisableStacktrace: !conf.WithStacktrace,
		EncoderConfig: zapcore.EncoderConfig{
			NameKey:       "logger",
			MessageKey:    "msg",
			LevelKey:      "level",
			EncodeLevel:   zapcore.CapitalLevelEncoder,
			TimeKey:       "time",
			EncodeTime:    zapcore.ISO8601TimeEncoder,
			CallerKey:     "caller",
			EncodeCaller:  zapcore.ShortCallerEncoder,
			StacktraceKey: "stacktrace",
		},
	}

	if len(conf.OutputPaths) > 0 {
		cfg.OutputPaths = make([]string, 0)

		for _, p := range conf.OutputPaths {
			if len(p) == 0 {
				continue
			}

			if p != Stderr {
				// Create log files and directories
				dir := filepath.Dir(p)
				if _, err := os.Stat(dir); os.IsNotExist(err) {
					err := os.MkdirAll(dir, os.ModePerm)
					if err != nil {
						return nil
					}
				}
			}

			cfg.OutputPaths = append(cfg.OutputPaths, p)
		}
	}

	// Log to stderr by default
	if len(cfg.OutputPaths) == 0 {
		cfg.OutputPaths = []string{"stderr"}
	}

	logger, err := cfg.Build()
	if err != nil {
		fmt.Printf("Couldn't create logger: %s", err.Error())
		return nil
	}

	if len(conf.Name) > 0 {
		logger = logger.Named(conf.Name)
	}

	return logger
}
