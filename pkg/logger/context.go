package logger

import (
	"aviasales/pkg/logger/logimpl"
	"context"
)

type ctxKeyType string

var ctxKey ctxKeyType = "logger_override"

// withLogger sets up a logger for context; any global function will extract
// it from context for logging.
func withLogger(ctx context.Context, l logimpl.Implementation) context.Context {
	return context.WithValue(ctx, ctxKey, l)
}

func fromCtxOrDefault(ctx context.Context) logimpl.Implementation {
	v, ok := ctx.Value(ctxKey).(logimpl.Implementation)
	if !ok {
		return defaultLogger
	}
	return v
}
