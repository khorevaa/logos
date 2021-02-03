package logos

type loggerFields interface {
	// Debug uses fmt.Sprint to construct and log a message.
	Debug(msg string, fields ...Field)

	// Info uses fmt.Sprint to construct and log a message.
	Info(msg string, fields ...Field)

	// Warn uses fmt.Sprint to construct and log a message.
	Warn(msg string, fields ...Field)

	// Error uses fmt.Sprint to construct and log a message.
	Error(msg string, fields ...Field)

	// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit(1).
	Fatal(msg string, fields ...Field)

	// Panic uses fmt.Sprint to construct and log a message, then panics.
	Panic(msg string, fields ...Field)

	// DPanic uses fmt.Sprint to construct and log a message. In development, the
	// logger then panics.
	DPanic(msg string, fields ...Field)
}

type loggerFormat interface {
	// Debugf uses fmt.Sprintf to construct and log a message.
	Debugf(format string, args ...interface{})

	// Infof uses fmt.Sprintf to log a templated message.
	Infof(format string, args ...interface{})

	// Warnf uses fmt.Sprintf to log a templated message.
	Warnf(format string, args ...interface{})

	// Errorf uses fmt.Sprintf to log a templated message.
	Errorf(format string, args ...interface{})

	// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit(1).
	Fatalf(format string, args ...interface{})

	// Panicf uses fmt.Sprintf to log a templated message, then panics.
	Panicf(format string, args ...interface{})

	// DPanicf uses fmt.Sprintf to log a templated message. In development, the
	// logger then panics.
	DPanicf(format string, args ...interface{})
}

type Logger interface {
	loggerFields

	Sync() error

	Sugar() SugaredLogger

	Job(name string, kvs ...map[string]interface{}) *Job

	EventEmitter() Emitter
}

type SugaredLogger interface {
	loggerFields
	loggerFormat

	// Debugw logs a message with some additional context. The additional context
	// is added in the form of key-value pairs. The optimal way to write the value
	// to the log message will be inferred by the value's type. To explicitly
	// specify a type you can pass a Field such as logos.Stringer.
	Debugw(msg string, keysAndValues ...interface{})

	// Infow logs a message with some additional context. The additional context
	// is added in the form of key-value pairs. The optimal way to write the value
	// to the log message will be inferred by the value's type. To explicitly
	// specify a type you can pass a Field such as logos.Stringer.
	Infow(msg string, keysAndValues ...interface{})

	// Warnw logs a message with some additional context. The additional context
	// is added in the form of key-value pairs. The optimal way to write the value
	// to the log message will be inferred by the value's type. To explicitly
	// specify a type you can pass a Field such as logos.Stringer.
	Warnw(msg string, keysAndValues ...interface{})

	// Errorw logs a message with some additional context. The additional context
	// is added in the form of key-value pairs. The optimal way to write the value
	// to the log message will be inferred by the value's type. To explicitly
	// specify a type you can pass a Field such as logos.Stringer.
	Errorw(msg string, keysAndValues ...interface{})

	// Fatalw logs a message with some additional context, then calls os.Exit(1).
	// The additional context is added in the form of key-value pairs. The optimal
	// way to write the value to the log message will be inferred by the value's
	// type. To explicitly specify a type you can pass a Field such as
	// logos.Stringer.
	Fatalw(msg string, keysAndValues ...interface{})

	// Panicw logs a message with some additional context, then panics. The
	// additional context is added in the form of key-value pairs. The optimal way
	// to write the value to the log message will be inferred by the value's type.
	// To explicitly specify a type you can pass a Field such as logos.Stringer.
	Panicw(msg string, keysAndValues ...interface{})

	// DPanicw logs a message with some additional context. The logger panics only
	// in Development mode.  The additional context is added in the form of
	// key-value pairs. The optimal way to write the value to the log message will
	// be inferred by the value's type. To explicitly specify a type you can pass a
	// Field such as logos.Stringer.
	DPanicw(msg string, keysAndValues ...interface{})

	Sync() error

	Desugar() Logger

	Job(name string, kvs ...map[string]interface{}) *Job

	EventEmitter() Emitter
}
