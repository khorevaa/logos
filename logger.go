package logos

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
	"sync/atomic"
	"time"
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

func (l *warpLogger) Sugar() SugaredLogger {

	l.initSugaredLogger()
	return l
}

func (l *warpLogger) Debug(msg string, fields ...Field) {
	l.checkLock()
	l.defLogger.Debug(msg, fields...)
}

func (l *warpLogger) Info(msg string, fields ...Field) {

	l.checkLock()
	l.defLogger.Info(msg, fields...)
}

func (l *warpLogger) Warn(msg string, fields ...Field) {
	l.checkLock()
	l.defLogger.Warn(msg, fields...)
}

func (l *warpLogger) Error(msg string, fields ...Field) {
	l.checkLock()
	l.defLogger.Error(msg, fields...)
}

func (l *warpLogger) Fatal(msg string, fields ...Field) {
	l.checkLock()
	l.defLogger.Fatal(msg, fields...)
}

func (l *warpLogger) Panic(msg string, fields ...Field) {
	l.checkLock()
	l.defLogger.Panic(msg, fields...)
}

func (l *warpLogger) DPanic(msg string, fields ...Field) {
	l.checkLock()
	l.defLogger.DPanic(msg, fields...)
}

func (l *warpLogger) Debugf(format string, args ...interface{}) {
	l.checkLock()
	l.sugaredLogger.Debugf(format, args...)
}

func (l *warpLogger) Infof(format string, args ...interface{}) {
	l.checkLock()
	l.sugaredLogger.Infof(format, args...)
}

func (l *warpLogger) Warnf(format string, args ...interface{}) {
	l.checkLock()
	l.sugaredLogger.Warnf(format, args...)
}

func (l *warpLogger) Errorf(format string, args ...interface{}) {
	l.checkLock()
	l.sugaredLogger.Errorf(format, args...)
}

func (l *warpLogger) Fatalf(format string, args ...interface{}) {
	l.checkLock()
	l.sugaredLogger.Fatalf(format, args...)
}

func (l *warpLogger) Panicf(format string, args ...interface{}) {
	l.checkLock()
	l.sugaredLogger.Panicf(format, args...)
}

func (l *warpLogger) DPanicf(format string, args ...interface{}) {
	l.checkLock()
	l.sugaredLogger.DPanicf(format, args...)
}

func (l *warpLogger) Debugw(msg string, keysAndValues ...interface{}) {
	l.checkLock()
	l.sugaredLogger.Debugw(msg, keysAndValues...)
}

func (l *warpLogger) Infow(msg string, keysAndValues ...interface{}) {
	l.checkLock()
	l.sugaredLogger.Infow(msg, keysAndValues...)
}

func (l *warpLogger) Warnw(msg string, keysAndValues ...interface{}) {
	l.checkLock()
	l.sugaredLogger.Warnw(msg, keysAndValues...)
}

func (l *warpLogger) Errorw(msg string, keysAndValues ...interface{}) {
	l.checkLock()
	l.sugaredLogger.Errorw(msg, keysAndValues...)
}

func (l *warpLogger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.checkLock()
	l.sugaredLogger.Fatalw(msg, keysAndValues...)
}

func (l *warpLogger) Panicw(msg string, keysAndValues ...interface{}) {
	l.checkLock()
	l.sugaredLogger.Panicw(msg, keysAndValues...)
}

func (l *warpLogger) DPanicw(msg string, keysAndValues ...interface{}) {
	l.checkLock()
	l.sugaredLogger.DPanicw(msg, keysAndValues...)
}

func (l *warpLogger) Sync() error {
	return l.defLogger.Sync()
}

func (l *warpLogger) Desugar() Logger {
	return l
}

func (l *warpLogger) UsedAt() time.Time {
	unix := atomic.LoadUint32(&l._usedAt)
	return time.Unix(int64(unix), 0)
}

func (l *warpLogger) SetUsedAt(tm time.Time) {
	atomic.StoreUint32(&l._usedAt, uint32(tm.Unix()))
}

func (l *warpLogger) initSugaredLogger() {

	if l.sugaredLogger == nil {
		l.sugaredLogger = l.defLogger.Sugar()
	}

}

func (l *warpLogger) updateLogger(logger *zap.Logger) {

	l.lock()
	defer l.unlock()

	_ = l.Sync()

	l.defLogger = logger

	if l.sugaredLogger != nil {
		l.sugaredLogger = l.defLogger.Sugar()
	}
}

func (l *warpLogger) lock() {

	l.mu.Lock()
	atomic.StoreUint32(&l._locked, 1)
	l._lockWait.Add(1)
}

func (l *warpLogger) unlock() {

	atomic.StoreUint32(&l._locked, 0)
	l._lockWait.Done()
	l.mu.Unlock()

}

func (l *warpLogger) checkLock() {

	if l.locked() {
		l._lockWait.Wait()
	}

}

func (l *warpLogger) locked() bool {
	return atomic.LoadUint32(&l._locked) == 1
}
