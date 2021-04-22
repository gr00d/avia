/*Package logimpl describes implementation that gets
passed around by `logger` via contexts.*/
package logimpl

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Level = zapcore.Level

var (
	DebugLevel Level = zap.DebugLevel
	InfoLevel  Level = zap.InfoLevel
	WarnLevel  Level = zap.WarnLevel
	ErrorLevel Level = zap.ErrorLevel
	// DPanicLevel logs are particularly important errors. In development the
	// logger panics after writing the message.
	DPanicLevel Level = zap.DPanicLevel
	PanicLevel  Level = zap.PanicLevel
	FatalLevel  Level = zap.FatalLevel
)

// Implementation writes logs somewhere.
type Implementation interface {
	Debug(msg string, kvs ...interface{})
	Info(msg string, kvs ...interface{})
	Warn(msg string, kvs ...interface{})
	Error(msg string, kvs ...interface{})
	DPanic(msg string, kvs ...interface{})
	Panic(msg string, kvs ...interface{})
	Fatal(msg string, kvs ...interface{})

	With(kvs ...interface{}) Implementation

	Log(lvl Level, msg string, kvs ...interface{})

	// WithLevel creates a child logger that has a different logging level
	// than its parent.
	//
	// Child still has any data embedded via With().
	WithLevel(Level) Implementation

	// Sync flushes buffered logs, if any.
	Sync()
}

// Configurable is a logger that configures itself.
type Configurable interface {
	Implementation
	SetLevel(Level)
}
