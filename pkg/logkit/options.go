package logkit

import (
	"os"
	"runtime"

	"golang.org/x/exp/slog"
)

func WithLevel(l slog.Leveler) HandlerOption {
	return hanlderOptionFunc(func(ho *slog.HandlerOptions) {
		ho.Level = l
	})
}

func WithCaller() HandlerOption {
	return hanlderOptionFunc(func(ho *slog.HandlerOptions) {
		ho.AddSource = true
	})
}

// CheckWriteHook is a custom action that may be executed after an entry is
// written.
//
// Register one on a CheckedEntry with the After method.
//
//	if ce := logger.Check(...); ce != nil {
//	  ce = ce.After(hook)
//	  ce.Write(...)
//	}
//
// You can configure the hook for Fatal log statements at the logger level with
// the zap.WithFatalHook option.
type CheckWriteHook interface {
	// OnWrite is invoked with the CheckedEntry that was written and a list
	// of fields added with that entry.
	//
	// The list of fields DOES NOT include fields that were already added
	// to the logger with the With method.
	OnWrite(r *slog.Record)
}

// CheckWriteAction indicates what action to take after a log entry is
// processed. Actions are ordered in increasing severity.
type CheckWriteAction uint8

const (
	// WriteThenNoop indicates that nothing special needs to be done. It's the
	// default behavior.
	WriteThenNoop CheckWriteAction = iota
	// WriteThenGoexit runs runtime.Goexit after Write.
	WriteThenGoexit
	// WriteThenPanic causes a panic after Write.
	WriteThenPanic
	// WriteThenFatal causes an os.Exit(1) after Write.
	WriteThenFatal
)

// OnWrite implements the OnWrite method to keep CheckWriteAction compatible
// with the new CheckWriteHook interface which deprecates CheckWriteAction.
func (a CheckWriteAction) OnWrite(r *slog.Record) {
	switch a {
	case WriteThenGoexit:
		runtime.Goexit()
	case WriteThenPanic:
		r.Attrs(func(a slog.Attr) {
			if a.Key == ErrorKey {
				panic(a.Value)
			}
		})
		panic(r.Message)
	case WriteThenFatal:
		os.Exit(1)
	}
}

func WithOnFatal(action CheckWriteAction) LoggerOption {
	return loggerOptionFunc(func(log *slog.Logger) {
		type Interface interface {
			SetOnFatal(onFatal CheckWriteHook)
		}
		if h, ok := log.Handler().(Interface); ok {
			h.SetOnFatal(action)
		}
	})
}

// Hooks registers functions which will be called each time the Logger writes
// out an Record. Repeated use of Hooks is additive.
func Hooks(hooks ...func(slog.Record) error) LoggerOption {
	type Interface interface {
		AppendHook(...func(slog.Record) error)
	}
	return loggerOptionFunc(func(log *slog.Logger) {
		if h, ok := log.Handler().(Interface); ok {
			h.AppendHook(hooks...)
		}
	})
}
