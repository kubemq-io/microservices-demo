package log

import (
	"context"
	"fmt"
	"go.uber.org/zap/zapcore"
	"os"

	"go.uber.org/zap"
)

type key string

var contextKey = key("com.heptio.logger")

func NewLogger() Logger {
	z, err := newZapLogger(2)
	if err != nil {
		fmt.Printf("failed to initialize logger: %v", err)
		os.Exit(1)
	}
	return Wrap(z.Sugar())
}

// Logger is an interface for logging
type Logger interface {
	// Debugf uses fmt.Sprintf to log a templated message.
	Debugf(template string, args ...interface{})

	// Infof uses fmt.Sprintf to log a templated message.
	Infof(template string, args ...interface{})

	// Warnf uses fmt.Sprintf to log a templated message.
	Warnf(template string, args ...interface{})

	// Errorf uses fmt.Sprintf to log a templated message.
	Errorf(template string, args ...interface{})

	Fatalf(template string, args ...interface{})

	With(args ...interface{}) Logger

	WithErr(err error) Logger

	Named(name string) Logger
}

// sugaredLogWrapper adapts a zap.SugaredLogger to the Logger interface
type sugaredLogWrapper struct {
	*zap.SugaredLogger
}

func (s *sugaredLogWrapper) WithErr(err error) Logger {
	return &sugaredLogWrapper{s.SugaredLogger.With("err", err)}
}

func (s *sugaredLogWrapper) With(args ...interface{}) Logger {
	return &sugaredLogWrapper{s.SugaredLogger.With(args...)}
}

func (s *sugaredLogWrapper) Named(name string) Logger {
	return &sugaredLogWrapper{s.SugaredLogger.Named(name)}
}

var _ Logger = (*sugaredLogWrapper)(nil)

// Wrap zap.SugaredLogger as Logger interface
func Wrap(z *zap.SugaredLogger) Logger {
	return &sugaredLogWrapper{z}
}

// NopLogger constructs a nop logger
func NopLogger() Logger {
	return Wrap(zap.NewNop().Sugar())
}

// WithLoggerContext returns a new context with a set logger
func WithLoggerContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, contextKey, logger)
}

// From extracts a logger from the supplied context, or returns a NopLogger if none is found.
func From(ctx context.Context) Logger {
	if ctx == nil {
		return NopLogger()
	}
	v := ctx.Value(contextKey)
	l, ok := v.(Logger)
	if !ok || l == nil {
		return NopLogger()
	}
	return l
}

// Returns a new zap logger, setting level according to the provided
// verbosity level as an offset of the base level, Info.
// i.e. verboseLevel==0, level==Info
//      verboseLevel==1, level==Debug
func newZapLogger(verboseLevel int) (*zap.Logger, error) {
	level := zapcore.InfoLevel - zapcore.Level(verboseLevel)
	if level < zapcore.DebugLevel || level > zapcore.FatalLevel {
		level = zapcore.DebugLevel
	}

	cfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(level),
		Development:      false,
		Encoding:         "console",
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	return cfg.Build()
}
