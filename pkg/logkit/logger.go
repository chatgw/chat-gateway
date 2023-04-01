package logkit

import (
	"errors"
	"os"

	"golang.org/x/exp/slog"
)

var (
	// Log is global logger
	Log *slog.Logger

	// customTimeFormat is custom Time format
	customTimeFormat string

	ErrAlreadyInitialized = errors.New("logger already initialized")

	ErrorKey = "error"
)

// Init initializes log by input parameters
// lvl - global log level: Debug(-1), Info(0), Warn(1), Error(2), DPanic(3), Panic(4), Fatal(5)
// timeFormat - custom time format for logger of empty string to use default
func Init(opts *Options) (*slog.Logger, error) {
	if Log != nil {
		Log.Error("Logger already initialized once, No need to do it multiple times", ErrAlreadyInitialized)
		return nil, ErrAlreadyInitialized
	}
	var err error
	Log, err = New(opts)
	return Log, err
}

// New initializes log by input parameters
// lvl - global log level: Debug(-1), Info(0), Warn(1), Error(2), DPanic(3), Panic(4), Fatal(5)
// timeFormat - custom time format for logger of empty string to use default
func New(opts *Options) (*slog.Logger, error) {
	// Configure console output.
	logTimeFormat := ""
	var useCustomTimeFormat bool

	if len(logTimeFormat) > 0 {
		customTimeFormat = logTimeFormat
		useCustomTimeFormat = true
	}

	hopts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			return a
		},
	}
	if opts != nil {
		opts.ApplyHanlder(hopts)
	}

	// High-priority output should also go to standard error, and low-priority
	// output should also go to standard out.
	// It is useful for Kubernetes deployment.
	// Kubernetes interprets os.Stdout log items as INFO and os.Stderr log items
	// as ERROR by default.
	highPriority := LevelEnablerFunc(func(lvl slog.Level) bool {
		return lvl >= slog.LevelError
	})
	lowPriority := LevelEnablerFunc(func(lvl slog.Level) bool {
		return lvl >= hopts.Level.Level() && lvl < slog.LevelError
	})

	// Join the outputs, encoders, and level-handling functions into
	// slog.
	handler := NewTeeHandler(
		NewLevelHandler(highPriority, hopts.NewTextHandler(os.Stderr)),
		NewLevelHandler(lowPriority, hopts.NewTextHandler(os.Stdout)),
	)
	handler.hopts = hopts

	// From a slog.Logger, it's easy to construct a Logger.
	log := slog.New(handler)
	if opts != nil {
		opts.ApplyLogger(log)
	}

	if !useCustomTimeFormat {
		log.Warn("time format for logger is not provided - use slog default")
	}

	return log, nil
}
