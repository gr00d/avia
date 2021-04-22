package zaplogger

import (
	"aviasales/pkg/logger/logimpl"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	atom zap.AtomicLevel
	z    *zap.SugaredLogger
}

var _ logimpl.Implementation = &Logger{}
var _ logimpl.Configurable = &Logger{}

func (z *Logger) Log(lvl logimpl.Level, msg string, kvs ...interface{}) {
	z.log(lvl, msg, kvs...)
}

func (z *Logger) log(lvl logimpl.Level, msg string, kvs ...interface{}) {
	switch lvl {
	case logimpl.DebugLevel:
		z.z.Debugw(msg, kvs...)
	case logimpl.InfoLevel:
		z.z.Infow(msg, kvs...)
	case logimpl.WarnLevel:
		z.z.Warnw(msg, kvs...)
	case logimpl.ErrorLevel:
		z.z.Errorw(msg, kvs...)
	case logimpl.DPanicLevel:
		z.z.DPanicw(msg, kvs...)
	case logimpl.PanicLevel:
		z.z.Panicw(msg, kvs...)
	case logimpl.FatalLevel:
		z.z.Fatalw(msg, kvs...)
	}
}
func (z *Logger) Debug(msg string, kvs ...interface{}) {
	z.log(logimpl.DebugLevel, msg, kvs...)
}
func (z *Logger) Info(msg string, kvs ...interface{}) {
	z.log(logimpl.InfoLevel, msg, kvs...)
}
func (z *Logger) Warn(msg string, kvs ...interface{}) {
	z.log(logimpl.WarnLevel, msg, kvs...)
}
func (z *Logger) Error(msg string, kvs ...interface{}) {
	z.log(logimpl.ErrorLevel, msg, kvs...)
}
func (z *Logger) DPanic(msg string, kvs ...interface{}) {
	z.log(logimpl.DPanicLevel, msg, kvs...)
}
func (z *Logger) Panic(msg string, kvs ...interface{}) {
	z.log(logimpl.PanicLevel, msg, kvs...)
}
func (z *Logger) Fatal(msg string, kvs ...interface{}) {
	z.log(logimpl.FatalLevel, msg, kvs...)
}

func (z *Logger) With(kvs ...interface{}) logimpl.Implementation {
	return &Logger{z: z.z.With(kvs...)}
}

func (z *Logger) SetLevel(lvl logimpl.Level) {
	z.atom.SetLevel(lvl)
}

func (z *Logger) WithLevel(lvl logimpl.Level) logimpl.Implementation {
	atom := zap.NewAtomicLevelAt(lvl)

	l := z.z.Desugar().WithOptions(zap.WrapCore(func(pc zapcore.Core) zapcore.Core {
		return zapCoreWrapper{core: pc, lvl: atom}
	})).Sugar()
	return &Logger{
		atom: atom,
		z:    l,
	}
}

func (z *Logger) Sync() {
	_ = z.z.Sync()
}

func NewProduction() (*Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	l, err := cfg.Build(zap.AddCallerSkip(3))
	if err != nil {
		return nil, err
	}
	return &Logger{
		atom: cfg.Level,
		z:    l.Sugar(),
	}, nil
}

func NewDevelopment() (*Logger, error) {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	l, err := cfg.Build(zap.AddCallerSkip(3))
	if err != nil {
		return nil, err
	}
	return &Logger{
		atom: cfg.Level,
		z:    l.Sugar(),
	}, nil
}
