package logkit

import (
	"golang.org/x/exp/slog"
)

// WrapOptions adds Option's to a test Logger built by NewLogger.
func WrapOptions(opts ...any) *Options {
	ret := &Options{
		o:  make([]LoggerOption, 0),
		ho: make([]HandlerOption, 0),
	}
	for _, opt := range opts {
		switch o := opt.(type) {
		case LoggerOption:
			ret.o = append(ret.o, o)
		case HandlerOption:
			ret.ho = append(ret.ho, o)
		default:
			panic("unknown logkit option")
		}
	}
	return ret
}

type Options struct {
	o  []LoggerOption
	ho []HandlerOption
}

func (opts Options) ApplyHanlder(h *slog.HandlerOptions) {
	for _, fn := range opts.ho {
		fn.apply(h)
	}
}

func (opts Options) ApplyLogger(logger *slog.Logger) {
	for _, fn := range opts.o {
		fn.apply(logger)
	}
}

// HandlerOption configures the test logger built by NewLogger.
type HandlerOption interface {
	apply(*slog.HandlerOptions)
}

type hanlderOptionFunc func(*slog.HandlerOptions)

func (f hanlderOptionFunc) apply(opts *slog.HandlerOptions) {
	f(opts)
}

// An LoggerOption configures a Logger.
type LoggerOption interface {
	apply(*slog.Logger)
}

// loggerOptionFunc wraps a func so it satisfies the Option interface.
type loggerOptionFunc func(*slog.Logger)

func (f loggerOptionFunc) apply(log *slog.Logger) {
	f(log)
}
