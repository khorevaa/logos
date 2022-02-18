package logos

import (
	"github.com/khorevaa/logos/appender"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type loggerConfig struct {
	Name string

	// Global core config
	Level         zap.AtomicLevel
	AddCaller     bool
	AddStacktrace zapcore.LevelEnabler

	Parent      *loggerConfig
	coreConfigs map[string]zap.AtomicLevel
}

func (l *loggerConfig) updateConfigLevel(appenderName string, level zapcore.Level) {

	if atomicLevel, ok := l.coreConfigs[appenderName]; ok {
		atomicLevel.SetLevel(level)
	}
}

func (l *loggerConfig) CreateLogger(appenders map[string]*appender.Appender) *warpLogger {

	zc := newZapCore(l.coreConfigs, appenders)
	zl := newZapLogger(l.Name, zc, zap.WithCaller(l.AddCaller), zap.AddStacktrace(l.AddStacktrace), zap.AddCallerSkip(1))
	return newLogger(l.Name, zl)

}

func (l *loggerConfig) UpdateLogger(logger *warpLogger, appenders map[string]*appender.Appender) {

	zc := newZapCore(l.coreConfigs, appenders)

	newLogger := zap.New(zc, zap.WithCaller(l.AddCaller), zap.AddStacktrace(l.AddStacktrace), zap.AddCallerSkip(1))

	if len(l.Name) > 0 {
		newLogger = newLogger.Named(l.Name)
	}

	logger.updateLogger(newLogger)

}

func (l *loggerConfig) copy(name string) *loggerConfig {

	log := &loggerConfig{
		Name:        name,
		Level:       l.Level,
		Parent:      l.Parent,
		coreConfigs: make(map[string]zap.AtomicLevel),
	}

	copyMapConfig(log.coreConfigs, l.coreConfigs)

	return log

}

func copyMapConfig(dst map[string]zap.AtomicLevel, src map[string]zap.AtomicLevel) {

	if len(src) == 0 {
		return
	}

	if dst == nil {
		dst = make(map[string]zap.AtomicLevel, len(src))
	}

	for name, level := range src {

		dst[name] = zap.NewAtomicLevelAt(level.Level())
	}

}
