package logos

import (
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _ Logger = (*warpLogger)(nil)
var _ SugaredLogger = (*warpLogger)(nil)

func newLogger(name string, logger *zap.Logger) *warpLogger {

	l := &warpLogger{
		Name:      name,
		defLogger: logger,
		mu:        sync.RWMutex{},
		emitLevel: InfoLevel,
	}

	l.SetUsedAt(time.Now())

	return l
}

type warpLogger struct {
	Name string

	defLogger     *zap.Logger
	sugaredLogger *zap.SugaredLogger
	mu            sync.RWMutex

	emitLevel zapcore.Level
	_usedAt   uint32 // atomic
	_locked   uint32
	_lockWait sync.WaitGroup
}

func (log *warpLogger) copy() *warpLogger {

	return &warpLogger{
		Name:      log.Name,
		defLogger: log.defLogger.WithOptions(),
		mu:        sync.RWMutex{},
		emitLevel: log.emitLevel,
	}
}

func (log *warpLogger) Sugar() SugaredLogger {

	log.initSugaredLogger()
	return log
}

func (log *warpLogger) Named(s string) Logger {

	copyLog := log.copy()

	if copyLog.Name == "" {
		copyLog.Name = s
	} else {
		copyLog.Name = strings.Join([]string{copyLog.Name, s}, ".")
	}

	copyLog.defLogger = copyLog.defLogger.Named(s)
	return copyLog
}

func (log *warpLogger) With(fields ...Field) Logger {
	log.checkLock()
	logger := log.copy()
	logger.defLogger = logger.defLogger.WithOptions(zap.Fields(fields...))
	return logger
}

func (log *warpLogger) Debug(msg string, fields ...Field) {
	log.checkLock()
	log.defLogger.Debug(msg, fields...)
}

func (log *warpLogger) Info(msg string, fields ...Field) {

	log.checkLock()
	log.defLogger.Info(msg, fields...)
}

func (log *warpLogger) Warn(msg string, fields ...Field) {
	log.checkLock()
	log.defLogger.Warn(msg, fields...)
}

func (log *warpLogger) Error(msg string, fields ...Field) {
	log.checkLock()
	log.defLogger.Error(msg, fields...)
}

func (log *warpLogger) Fatal(msg string, fields ...Field) {
	log.checkLock()
	log.defLogger.Fatal(msg, fields...)
}

func (log *warpLogger) Panic(msg string, fields ...Field) {
	log.checkLock()
	log.defLogger.Panic(msg, fields...)
}

func (log *warpLogger) DPanic(msg string, fields ...Field) {
	log.checkLock()
	log.defLogger.DPanic(msg, fields...)
}

func (log *warpLogger) Debugf(format string, args ...interface{}) {
	log.checkLock()
	log.sugaredLogger.Debugf(format, args...)
}

func (log *warpLogger) Infof(format string, args ...interface{}) {
	log.checkLock()
	log.sugaredLogger.Infof(format, args...)
}

func (log *warpLogger) Warnf(format string, args ...interface{}) {
	log.checkLock()
	log.sugaredLogger.Warnf(format, args...)
}

func (log *warpLogger) Errorf(format string, args ...interface{}) {
	log.checkLock()
	log.sugaredLogger.Errorf(format, args...)
}

func (log *warpLogger) Fatalf(format string, args ...interface{}) {
	log.checkLock()
	log.sugaredLogger.Fatalf(format, args...)
}

func (log *warpLogger) Panicf(format string, args ...interface{}) {
	log.checkLock()
	log.sugaredLogger.Panicf(format, args...)
}

func (log *warpLogger) DPanicf(format string, args ...interface{}) {
	log.checkLock()
	log.sugaredLogger.DPanicf(format, args...)
}

func (log *warpLogger) Debugw(msg string, keysAndValues ...interface{}) {
	log.checkLock()
	log.sugaredLogger.Debugw(msg, keysAndValues...)
}

func (log *warpLogger) Infow(msg string, keysAndValues ...interface{}) {
	log.checkLock()
	log.sugaredLogger.Infow(msg, keysAndValues...)
}

func (log *warpLogger) Warnw(msg string, keysAndValues ...interface{}) {
	log.checkLock()
	log.sugaredLogger.Warnw(msg, keysAndValues...)
}

func (log *warpLogger) Errorw(msg string, keysAndValues ...interface{}) {
	log.checkLock()
	log.sugaredLogger.Errorw(msg, keysAndValues...)
}

func (log *warpLogger) Fatalw(msg string, keysAndValues ...interface{}) {
	log.checkLock()
	log.sugaredLogger.Fatalw(msg, keysAndValues...)
}

func (log *warpLogger) Panicw(msg string, keysAndValues ...interface{}) {
	log.checkLock()
	log.sugaredLogger.Panicw(msg, keysAndValues...)
}

func (log *warpLogger) DPanicw(msg string, keysAndValues ...interface{}) {
	log.checkLock()
	log.sugaredLogger.DPanicw(msg, keysAndValues...)
}

func (log *warpLogger) Sync() error {
	return log.defLogger.Sync()
}

func (log *warpLogger) Desugar() Logger {
	return log
}

func (log *warpLogger) UsedAt() time.Time {
	unix := atomic.LoadUint32(&log._usedAt)
	return time.Unix(int64(unix), 0)
}

func (log *warpLogger) SetUsedAt(tm time.Time) {
	atomic.StoreUint32(&log._usedAt, uint32(tm.Unix()))
}

func (log *warpLogger) initSugaredLogger() {

	if log.sugaredLogger == nil {
		log.sugaredLogger = log.defLogger.Sugar()
	}

}

func (log *warpLogger) updateLogger(logger *zap.Logger) {

	_ = log.Sync()

	log.defLogger = logger

	if log.sugaredLogger != nil {
		log.sugaredLogger = log.defLogger.Sugar()
	}
}

func (log *warpLogger) lock() {

	atomic.StoreUint32(&log._locked, 1)
}

func (log *warpLogger) unlock() {

	atomic.StoreUint32(&log._locked, 0)
	log._lockWait.Done()
	log.mu.Unlock()

}

func (log *warpLogger) checkLock() {

	if log.locked() {
		log._lockWait.Wait()
	}

}

func (log *warpLogger) locked() bool {
	return atomic.LoadUint32(&log._locked) == 1
}
