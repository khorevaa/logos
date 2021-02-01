package logos

import (
	"github.com/khorevaa/logos/appender"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var StackTraceLevelEnabler = zap.NewAtomicLevelAt(zapcore.ErrorLevel)

func newZapCore(config map[string]zapcore.Level, appenders map[string]*appender.Appender) zapcore.Core {

	zcs := make([]zapcore.Core, 0)

	for name, level := range config {

		if level == OffLevel {
			continue
		}

		if a, ok := appenders[name]; ok {
			zcs = append(zcs, zapcore.NewCore(a.Encoder, a.Writer, level))
		}

	}
	if len(zcs) == 0 {
		return zapcore.NewNopCore()
	}

	return zapcore.NewTee(zcs...)
}

func newZapLogger(name string, core zapcore.Core, option ...zap.Option) *zap.Logger {

	return zap.New(core, option...).Named(name)
}
