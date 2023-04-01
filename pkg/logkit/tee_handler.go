package logkit

import (
	"context"

	"golang.org/x/exp/slog"
)

// A TeeHandler wraps a Handler with an Enabled method
// that returns false for levels below a minimum.
type TeeHandler struct {
	hopts    *slog.HandlerOptions
	handlers []slog.Handler
	hooks    []func(slog.Record) error
	onFatal  CheckWriteHook // default is WriteThenFatal
}

// NewLevelHandler returns a TeeHandler with the given level.
// All methods except Enabled delegate to h.
func NewTeeHandler(handlers ...slog.Handler) *TeeHandler {
	return &TeeHandler{
		handlers: handlers,
		hooks:    make([]func(slog.Record) error, 0),
		onFatal:  WriteThenNoop,
	}
}

// Enabled implements Handler.Enabled by reporting whether
// level is at least as large as h's level.
func (h *TeeHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for i := 0; i < len(h.handlers); i++ {
		if h.handlers[i].Enabled(ctx, level) {
			return true
		}
	}
	return false
}

// Handle implements Handler.Handle.
func (h *TeeHandler) Handle(ctx context.Context, r slog.Record) (err error) {
	for _, handler := range h.handlers {
		if !handler.Enabled(ctx, r.Level) {
			continue
		}
		if err = handler.Handle(ctx, r); err != nil {
			return
		}
	}
	for i := 0; i < len(h.hooks); i++ {
		if err = h.hooks[i](r); err != nil {
			return
		}
	}
	if r.Level > slog.LevelError {
		h.onFatal.OnWrite(&r)
	}
	return nil
}

// WithAttrs implements Handler.WithAttrs.
func (h *TeeHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, len(h.handlers))
	for i := 0; i < len(h.handlers); i++ {
		handlers[i] = h.handlers[i].WithAttrs(attrs)
	}
	nh := NewTeeHandler(handlers...)
	nh.AppendHook(h.hooks...)
	return nh
}

// WithGroup implements Handler.WithGroup.
func (h *TeeHandler) WithGroup(name string) slog.Handler {
	handlers := make([]slog.Handler, len(h.handlers))
	for i := 0; i < len(h.handlers); i++ {
		handlers[i] = h.handlers[i].WithGroup(name)
	}
	nh := NewTeeHandler(handlers...)
	nh.AppendHook(h.hooks...)
	return nh
}

// Handlers returns the Handlers wrapped by h.
func (h *TeeHandler) Handlers() []slog.Handler {
	return h.handlers
}

// AppendHandlers returns the Handlers wrapped by h.
func (h *TeeHandler) AppendHandlers(handlers ...slog.Handler) {
	h.handlers = append(h.handlers, handlers...)
}

// AppendHook append hooks to hander.
func (h *TeeHandler) AppendHook(hooks ...func(slog.Record) error) {
	h.hooks = append(h.hooks, hooks...)
}

// SetOnFatal set on fatal action.
func (h *TeeHandler) SetOnFatal(onFatal CheckWriteHook) {
	h.onFatal = onFatal
}

func (h *TeeHandler) HandlerOptions() *slog.HandlerOptions {
	return h.hopts
}
