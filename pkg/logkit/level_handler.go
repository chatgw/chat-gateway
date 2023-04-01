package logkit

import (
	"context"

	"golang.org/x/exp/slog"
)

type LevelEnablerFunc func(slog.Level) bool

// Enabled calls the wrapped function.
func (f LevelEnablerFunc) Enabled(ctx context.Context, lvl slog.Level) bool { return f(lvl) }

// Enabled implements Handler.Enabled by reporting whether
// level is at least as large as h's level.
type LevelEnabler interface {
	Enabled(context.Context, slog.Level) bool
}

// A LevelHandler wraps a Handler with an Enabled method
// that returns false for levels below a minimum.
type LevelHandler struct {
	LevelEnabler

	handler slog.Handler
}

// NewLevelHandler returns a LevelHandler with the given level.
// All methods except Enabled delegate to h.
func NewLevelHandler(level LevelEnabler, h slog.Handler) *LevelHandler {
	// Optimization: avoid chains of LevelHandlers.
	if lh, ok := h.(*LevelHandler); ok {
		h = lh.Handler()
	}
	return &LevelHandler{level, h}
}

// Handle implements Handler.Handle.
func (h *LevelHandler) Handle(ctx context.Context, r slog.Record) error {
	return h.handler.Handle(ctx, r)
}

// WithAttrs implements Handler.WithAttrs.
func (h *LevelHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return NewLevelHandler(h.LevelEnabler, h.handler.WithAttrs(attrs))
}

// WithGroup implements Handler.WithGroup.
func (h *LevelHandler) WithGroup(name string) slog.Handler {
	return NewLevelHandler(h.LevelEnabler, h.handler.WithGroup(name))
}

// Handler returns the Handler wrapped by h.
func (h *LevelHandler) Handler() slog.Handler {
	return h.handler
}
