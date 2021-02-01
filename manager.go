package logos

import (
	"fmt"
	"github.com/khorevaa/logos/appender"
	config2 "github.com/khorevaa/logos/config"
	"github.com/khorevaa/logos/internal/common"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
)

const (
	rootLoggerName    = "root"
	packageSepBySlash = '/'
)

type logManager struct {
	getLoggerLocker sync.RWMutex
	loggerConfigs   sync.Map
	coreLoggers     sync.Map //

	appenders map[string]*appender.Appender

	rootLevel        zapcore.Level
	rootLogger       *warpLogger
	rootLoggerConfig *loggerConfig

	cancelRedirectStdLog func()
}

func newLogManager(rawConfig *common.Config) (*logManager, error) {

	config := config2.Config{}
	err := rawConfig.Unpack(&config)
	if err != nil {
		return nil, err
	}

	m := logManager{
		loggerConfigs: sync.Map{},
		coreLoggers:   sync.Map{},
		appenders:     map[string]*appender.Appender{},
	}

	for appenderType, appenderConfigs := range config.Appenders {
		for _, appenderConfig := range appenderConfigs {
			createAppender, err := appender.CreateAppender(appenderType, appenderConfig)
			if err != nil {
				return nil, err
			}
			name, err := appenderConfig.Name()
			if err != nil {
				return nil, err
			}
			m.appenders[name] = createAppender
		}
	}

	err = m.newRootLoggerFromCfg(config.Loggers.Root)

	if err != nil {
		return nil, err
	}

	// loggerConfigs
	for _, lc := range config.Loggers.Logger {
		l, err := m.newLoggerFromCfg(lc)
		if err != nil {
			return nil, err
		}

		m.loggerConfigs.Store(lc.Name, l)

		//if _, loaded := m.loggerConfigs.LoadOrStore(lc.Name, l); loaded {
		//	return nil, fmt.Errorf("duplicated logger %q", lc.Name)
		//}
	}

	return &m, nil

}

func (m *logManager) NewLogger(name string) Logger {
	return m.getLogger(name)
}

func (m *logManager) SetLevel(name string, level zapcore.Level) {

	// TODO Чтото сделать
	//return m.getLogger(name)
}

func (m *logManager) getLogger(name string, lock ...bool) *warpLogger {

	if len(name) == 0 {
		return m.rootLogger
	}

	if len(lock) > 0 && lock[0] || len(lock) == 0 {
		m.getLoggerLocker.Lock()
		defer m.getLoggerLocker.Unlock()
	}

	if core, ok := m.coreLoggers.Load(name); ok {
		return core.(*warpLogger)
	}

	if cfg, ok := m.loggerConfigs.Load(name); ok {

		logConfig := cfg.(*loggerConfig)

		core := logConfig.CreateLogger(m.appenders)
		m.coreLoggers.Store(logConfig.Name, core)

		return core
	}

	logConfig := m.newCoreLoggerConfig(name)
	m.loggerConfigs.Store(name, logConfig)

	core := logConfig.CreateLogger(m.appenders)
	m.coreLoggers.Store(logConfig.Name, core)

	return core

}

func (m *logManager) getParent(name string) *loggerConfig {

	parent := m.getRootLoggerConfig()
	for i, c := range name {
		// Search for package separator character
		if c == packageSepBySlash {
			parentName := name[0:i]
			if parentName != "" {
				parent = m.loadCoreLoggerConfig(parentName, parent)
			}
		}
	}

	return parent
}

func (m *logManager) getRootLoggerConfig() *loggerConfig {

	if m.rootLoggerConfig != nil {
		return m.rootLoggerConfig
	}

	name := rootLoggerName

	if log, ok := m.loggerConfigs.Load(name); ok {
		return log.(*loggerConfig)
	}

	log := &loggerConfig{
		Name:          name,
		Level:         m.rootLevel,
		coreConfigs:   make(map[string]zapcore.Level),
		AddStacktrace: StackTraceLevelEnabler,
	}
	log.coreConfigs["CONSOLE"] = InfoLevel

	m.rootLoggerConfig = log
	m.loggerConfigs.Store(name, log)
	m.rootLogger = log.CreateLogger(m.appenders)

	return m.rootLoggerConfig
}

func (m *logManager) newLoggerFromCfg(loggerCfg config2.LoggerConfig) (*loggerConfig, error) {

	name := loggerCfg.Name
	levelName := loggerCfg.Level
	appenderConfigs := loggerCfg.AppenderConfig
	appenders := loggerCfg.AppenderRefs

	level, err := createLevel(levelName)
	if err != nil {
		return nil, err
	}

	log := m.newCoreLoggerConfig(name)
	log.Level = level

	if len(appenders) > 0 {
		log.coreConfigs = make(map[string]zapcore.Level, len(appenders))
		for _, appenderName := range appenders {
			log.coreConfigs[appenderName] = level
		}
	}

	for _, appenderConfig := range appenderConfigs {

		if len(appenderConfig.Level) > 0 {
			appenderLevel, err := createLevel(appenderConfig.Level)
			if err != nil {
				debugf("creating appender level <%s> error: %s", appenderConfig.Level, err)
				continue
			}
			log.coreConfigs[appenderConfig.Name] = appenderLevel
		}

	}

	log.AddCaller = loggerCfg.AddCaller
	log.AddStacktrace = StackTraceLevelEnabler

	if tLevel, err := createLevel(loggerCfg.TraceLevel); len(loggerCfg.TraceLevel) > 0 && err == nil {
		log.AddStacktrace = zap.NewAtomicLevelAt(tLevel)
	}

	return log, nil
}

func debugf(format string, args ...interface{}) {
	if debug {
		fmt.Printf(format, args...)
	}
}

func (m *logManager) loadCoreLoggerConfig(name string, parent *loggerConfig) *loggerConfig {

	if log, ok := m.loggerConfigs.Load(name); ok {
		return log.(*loggerConfig)
	}

	if parent == nil {
		parent = m.rootLoggerConfig
	}

	log := &loggerConfig{
		Name:          name,
		Parent:        parent,
		AddStacktrace: parent.AddStacktrace,
		AddCaller:     parent.AddCaller,
		coreConfigs:   make(map[string]zapcore.Level),
	}

	copyMapConfig(log.coreConfigs, parent.coreConfigs)
	m.loggerConfigs.Store(name, log)
	return log

}

func (m *logManager) newCoreLoggerConfig(name string) *loggerConfig {

	parent := m.getParent(name)
	loggerConfig := m.loadCoreLoggerConfig(name, parent)

	return loggerConfig
}

func (m *logManager) RedirectStdLog() func() {

	stdlog := m.getLogger("stdlog", false)
	m.cancelRedirectStdLog = zap.RedirectStdLog(stdlog.defLogger)
	return m.cancelRedirectStdLog
}

func (m *logManager) CancelRedirectStdLog() {

	if m.cancelRedirectStdLog == nil {
		return
	}

	m.cancelRedirectStdLog()
}

func (m *logManager) Update(rawConfig *common.Config) error {

	nc, err := newLogManager(rawConfig)
	if err != nil {
		return err
	}

	m.getLoggerLocker.Lock()
	defer m.getLoggerLocker.Unlock()

	err = m.Sync()

	if err != nil {
		return err
	}

	m.appenders = nc.appenders
	m.rootLevel = nc.rootLevel
	m.rootLoggerConfig = nc.rootLoggerConfig
	m.rootLoggerConfig.UpdateLogger(m.rootLogger, m.appenders)

	m.loggerConfigs.Range(func(key, value interface{}) bool {
		name := key.(string)
		newLog := nc.newCoreLoggerConfig(name)
		ref := value.(*loggerConfig)
		*ref = *newLog
		return true
	})

	nc.loggerConfigs.Range(func(key, value interface{}) bool {
		if _, found := m.loggerConfigs.Load(key); found {
			return true
		}
		m.loggerConfigs.Store(key, value)
		return true
	})

	m.coreLoggers.Range(func(key, value interface{}) bool {

		newCore := m.newCoreLoggerConfig(key.(string))
		newCore.UpdateLogger(value.(*warpLogger), m.appenders)
		return true
	})

	if m.cancelRedirectStdLog != nil {
		m.cancelRedirectStdLog = m.RedirectStdLog()
	}

	return nil
}

func (m *logManager) Sync() error {
	m.coreLoggers.Range(func(_, value interface{}) bool {
		_ = value.(*warpLogger).Sync()
		return true
	})
	for _, a := range m.appenders {
		_ = a.Writer.Sync()
	}
	return nil
}

func (m *logManager) newRootLoggerFromCfg(root config2.RootLogger) error {

	levelName := root.Level
	appenderConfigs := root.AppenderConfig
	appenders := root.AppenderRefs

	level, err := createLevel(levelName)
	if err != nil {
		return err
	}

	m.rootLevel = level

	rootLoggerConfig := &loggerConfig{
		Name:          rootLoggerName,
		Level:         m.rootLevel,
		coreConfigs:   make(map[string]zapcore.Level),
		AddStacktrace: StackTraceLevelEnabler,
	}

	for _, appenderName := range appenders {
		rootLoggerConfig.coreConfigs[appenderName] = level
	}

	for _, appenderConfig := range appenderConfigs {

		if len(appenderConfig.Level) > 0 {

			appenderLevel, err := createLevel(appenderConfig.Level)
			if err != nil {
				debugf("creating appender level <%s> error: %s", appenderConfig.Level, err)
				continue
			}

			rootLoggerConfig.coreConfigs[appenderConfig.Name] = appenderLevel

		}

	}
	m.rootLoggerConfig = rootLoggerConfig
	m.loggerConfigs.Store(rootLoggerName, m.rootLoggerConfig)
	m.rootLogger = m.rootLoggerConfig.CreateLogger(m.appenders)

	m.coreLoggers.Store(rootLoggerName, m.rootLogger)

	return nil
}

func (m *logManager) UpdateLogger(name string, logger *zap.Logger) {
	core := m.getLogger(name, false)
	core.updateLogger(logger)
}

func createLevel(level string) (zapcore.Level, error) {
	switch level {
	case "off", "OFF", "false":
		return OffLevel, nil
	default:
		var l zapcore.Level
		err := l.UnmarshalText([]byte(level))
		return l, err
	}

}
