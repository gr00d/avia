package logger

import (
	"aviasales/pkg/logger/logimpl"
	"aviasales/pkg/logger/zaplogger"
	"context"
)

type Level = logimpl.Level

var (
	DebugLevel Level = logimpl.DebugLevel
	InfoLevel  Level = logimpl.InfoLevel
	WarnLevel  Level = logimpl.WarnLevel
	ErrorLevel Level = logimpl.ErrorLevel
	// DPanicLevel logs are particularly important errors. In development the
	// logger panics after writing the message.
	DPanicLevel Level = logimpl.DPanicLevel
	PanicLevel  Level = logimpl.PanicLevel
	FatalLevel  Level = logimpl.FatalLevel
)

// With adds keys and values to the child context.
//
// Any log entry made via child context will have these
// kv pairs in it.
func With(ctx context.Context, kvs ...interface{}) context.Context {
	return withLogger(ctx, fromCtxOrDefault(ctx).With(kvs...))
}

func Debug(ctx context.Context, msg string, kvs ...interface{}) {
	fromCtxOrDefault(ctx).Debug(msg, kvs...)
}
func Info(ctx context.Context, msg string, kvs ...interface{}) {
	fromCtxOrDefault(ctx).Info(msg, kvs...)
}
func Warn(ctx context.Context, msg string, kvs ...interface{}) {
	fromCtxOrDefault(ctx).Warn(msg, kvs...)
}
func WarnE(ctx context.Context, msg string, err error, kvs ...interface{}) {
	fromCtxOrDefault(ctx).Warn(msg, append([]interface{}{"err", err}, kvs...)...)
}
func Error(ctx context.Context, msg string, err error, kvs ...interface{}) {
	fromCtxOrDefault(ctx).Error(msg, append([]interface{}{"err", err}, kvs...)...)
}
func Errorn(ctx context.Context, msg string, kvs ...interface{}) {
	fromCtxOrDefault(ctx).Error(msg, kvs...)
}
func DPanic(ctx context.Context, msg string, kvs ...interface{}) {
	fromCtxOrDefault(ctx).DPanic(msg, kvs...)
}
func Panic(ctx context.Context, msg string, kvs ...interface{}) {
	fromCtxOrDefault(ctx).Panic(msg, kvs...)
}
func Fatal(ctx context.Context, msg string, kvs ...interface{}) {
	fromCtxOrDefault(ctx).Fatal(msg, kvs...)
}
func FatalE(ctx context.Context, msg string, err error, kvs ...interface{}) {
	fromCtxOrDefault(ctx).Fatal(msg, append([]interface{}{"err", err}, kvs...)...)
}

func Log(ctx context.Context, lvl Level, msg string, kvs ...interface{}) {
	fromCtxOrDefault(ctx).Log(lvl, msg, kvs...)
}

// WithLevel overrides log level for logs made via child context.
// Only child context's loglevel will be changed; parent will stay
// as is.
func WithLevel(ctx context.Context, lvl Level) context.Context {
	// panic here means SetLogger(ctx,WithLevel(ctx,lvl)) was called before
	return withLogger(ctx, fromCtxOrDefault(ctx).WithLevel(lvl))
}

// SetLevel sets logging level for root logger.
func SetLevel(lvl Level) {
	defaultLogger.SetLevel(lvl)
}

var defaultLogger logimpl.Configurable

// SetGlobalLogger overrides global root logger.
func SetGlobalLogger(l logimpl.Configurable) {
	defaultLogger = l
}

// nolint:gochecknoinits
func init() {
	zl, err := zaplogger.NewDevelopment()
	if err != nil {
		panic(err)
	}
	SetGlobalLogger(zl)
}
