package main

import (
	"context"
	"log/slog"
)

var _ slog.Handler = (*TeeLogger)(nil)

type TeeLogger struct {
	handlers []slog.Handler
}

func NewTeeLogger(handlers ...slog.Handler) *TeeLogger {
	return &TeeLogger{
		handlers: handlers,
	}
}

func (t *TeeLogger) Enabled(ctx context.Context, level slog.Level) bool {
	for _, handler := range t.handlers {
		if handler.Enabled(ctx, level) {
			return true
		}
	}

	return false
}

func (t *TeeLogger) Handle(ctx context.Context, record slog.Record) error {
	for _, handler := range t.handlers {
		if err := handler.Handle(ctx, record); err != nil {
			return err
		}
	}

	return nil
}

func (t *TeeLogger) WithAttrs(attrs []slog.Attr) slog.Handler {
	var ret = &TeeLogger{
		handlers: make([]slog.Handler, len(t.handlers)),
	}
	for k, handler := range t.handlers {
		ret.handlers[k] = handler.WithAttrs(attrs)
	}

	return ret
}

func (t *TeeLogger) WithGroup(name string) slog.Handler {
	var ret = &TeeLogger{
		handlers: make([]slog.Handler, len(t.handlers)),
	}
	for k, handler := range t.handlers {
		ret.handlers[k] = handler.WithGroup(name)
	}

	return ret
}
