package zaplogger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapCoreWrapper struct {
	core zapcore.Core
	lvl  zap.AtomicLevel
}

func (c zapCoreWrapper) Enabled(l zapcore.Level) bool {
	return c.lvl.Enabled(l)
}
func (c zapCoreWrapper) With(ff []zapcore.Field) zapcore.Core {
	return zapCoreWrapper{
		lvl:  c.lvl,
		core: c.core.With(ff),
	}
}

// nolint:gocritic
func (c zapCoreWrapper) Check(e zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(e.Level) {
		return ce.AddCore(e, c)
	}
	return ce
}

// nolint:gocritic
func (c zapCoreWrapper) Write(e zapcore.Entry, ff []zapcore.Field) error {
	return c.core.Write(e, ff)
}

func (c zapCoreWrapper) Sync() error {
	return c.core.Sync()
}
