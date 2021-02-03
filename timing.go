package logos

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

type Kvs map[string]interface{}

type Emitter interface {
	EmitEvent(job string, event string, kvs map[string]interface{})
	EmitEventErr(job string, event string, err error, kvs map[string]interface{})
	EmitTiming(job string, event string, duration time.Duration, kvs map[string]interface{})
	EmitComplete(job string, status CompletionStatus, duration time.Duration, kvs map[string]interface{})
	EmitGauge(job string, event string, value float64, kvs map[string]interface{})
}

func (l *warpLogger) Job(name string, kvs ...map[string]interface{}) *Job {

	return newJob(name, l, kvs...)

}

func (l *warpLogger) EventEmitter() Emitter {

	return l

}

func newJob(name string, emitter Emitter, kvs ...map[string]interface{}) *Job {

	job := &Job{
		Name:      name,
		emitter:   emitter,
		Start:     time.Now(),
		KeyValues: map[string]interface{}{},
	}

	if len(kvs) > 0 {
		job.KeyValues = kvs[0]
	}

	return job
}

var _ Emitter = (*warpLogger)(nil)

func (e *warpLogger) EmitEvent(job string, event string, kvs map[string]interface{}) {

	var fields []Field
	fields = append(fields, String("job", job), String("event", event))
	kvsToFields(&fields, kvs)

	e.emit(0, e.emitLevel, fields)
}

func (e *warpLogger) EmitEventErr(job string, event string, err error, kvs map[string]interface{}) {

	var fields []Field
	fields = append(fields, String("job", job), String("event", event), Error(err))
	kvsToFields(&fields, kvs)

	e.emit(0, ErrorLevel, fields)

}

func (e *warpLogger) EmitTiming(job string, event string, duration time.Duration, kvs map[string]interface{}) {
	var fields []Field
	fields = append(fields, String("job", job), String("event", event), Duration("duration", duration))
	kvsToFields(&fields, kvs)

	e.emit(0, e.emitLevel, fields)

}

func (e *warpLogger) EmitComplete(job string, status CompletionStatus, duration time.Duration, kvs map[string]interface{}) {
	var fields []Field
	fields = append(fields, String("job", job), String("status", status.String()), Duration("duration", duration))

	kvsToFields(&fields, kvs)

	lvl := e.emitLevel

	switch status {
	case Err, ValidationError, Panic:
		lvl = ErrorLevel
	case Junk:
		lvl = WarnLevel
	}

	e.emit(0, lvl, fields)

}

func (e *warpLogger) EmitGauge(job string, event string, value float64, kvs map[string]interface{}) {
	var fields []Field
	fields = append(fields, String("job", job), String("event", event), Float64("gauge", value))

	kvsToFields(&fields, kvs)
	e.emit(0, e.emitLevel, fields)

}

func kvsToFields(fields *[]Field, kvs map[string]interface{}) {

	for key, value := range kvs {
		*fields = append(*fields, Any(key, value))
	}

	return
}

func (e *warpLogger) emit(callerSkip int, emitLevel zapcore.Level, fields []zapcore.Field) {

	emitter := e.defLogger

	if callerSkip > 0 {
		emitter = e.defLogger.WithOptions(zap.AddCallerSkip(callerSkip))
	}

	switch emitLevel {
	case DebugLevel:
		emitter.Debug("", fields...)
	case InfoLevel:
		emitter.Info("", fields...)
	case WarnLevel:
		emitter.Warn("", fields...)
	case ErrorLevel:
		emitter.Error("", fields...)
	case PanicLevel:
		emitter.Panic("", fields...)
	case DPanicLevel:
		emitter.DPanic("", fields...)
	case FatalLevel:
		emitter.Fatal("", fields...)
	}

}

type CompletionStatus int

const (
	Success CompletionStatus = iota
	ValidationError
	Panic
	Err
	Junk
)

var completionStatusToString = map[CompletionStatus]string{
	Success:         "success",
	ValidationError: "validation_error",
	Panic:           "panic",
	Err:             "error",
	Junk:            "junk",
}

func (cs CompletionStatus) String() string {
	return completionStatusToString[cs]
}

type Job struct {
	Name    string
	emitter Emitter

	Start     time.Time
	KeyValues map[string]interface{}
}

func (j *Job) Event(eventName string) {
	allKvs := j.mergedKeyValues(nil)
	j.emitter.EmitEvent(j.Name, eventName, allKvs)
}

func (j *Job) EventKv(eventName string, kvs map[string]interface{}) {
	allKvs := j.mergedKeyValues(kvs)
	j.emitter.EmitEvent(j.Name, eventName, allKvs)
}

func (j *Job) EventErr(eventName string, err error) error {

	allKvs := j.mergedKeyValues(nil)
	j.emitter.EmitEventErr(j.Name, eventName, err, allKvs)
	return err
}

func (j *Job) EventErrKv(eventName string, err error, kvs map[string]interface{}) error {
	allKvs := j.mergedKeyValues(kvs)
	j.emitter.EmitEventErr(j.Name, eventName, err, allKvs)
	return err
}

func (j *Job) Timing(eventName string, duration time.Duration) {
	allKvs := j.mergedKeyValues(nil)
	j.emitter.EmitTiming(j.Name, eventName, duration, allKvs)
}

func (j *Job) TimingKv(eventName string, duration time.Duration, kvs map[string]interface{}) {
	allKvs := j.mergedKeyValues(kvs)
	j.emitter.EmitTiming(j.Name, eventName, duration, allKvs)
}

func (j *Job) Gauge(eventName string, value float64) {
	allKvs := j.mergedKeyValues(nil)
	j.emitter.EmitGauge(j.Name, eventName, value, allKvs)
}

func (j *Job) GaugeKv(eventName string, value float64, kvs map[string]interface{}) {
	allKvs := j.mergedKeyValues(kvs)
	j.emitter.EmitGauge(j.Name, eventName, value, allKvs)
}

func (j *Job) Complete(status CompletionStatus) {
	allKvs := j.mergedKeyValues(nil)
	j.emitter.EmitComplete(j.Name, status, time.Since(j.Start), allKvs)

}

func (j *Job) CompleteErr(err error) error {

	j.CompleteKv(Err, map[string]interface{}{
		"error": err,
	})
	return err

}

func (j *Job) CompleteKv(status CompletionStatus, kvs map[string]interface{}) {
	allKvs := j.mergedKeyValues(kvs)
	j.emitter.EmitComplete(j.Name, status, time.Since(j.Start), allKvs)
}

func (j *Job) KeyValue(key string, value string) *Job {
	if j.KeyValues == nil {
		j.KeyValues = make(map[string]interface{})
	}
	j.KeyValues[key] = value
	return j
}

func (j *Job) mergedKeyValues(instanceKvs map[string]interface{}) map[string]interface{} {
	var allKvs map[string]interface{}

	// Count how many maps actually have contents in them. If it's 0 or 1, we won't allocate a new map.
	// Also, optimistically set allKvs. We might use it or we might overwrite the value with a newly made map.
	var kvCount = 0
	if len(j.KeyValues) > 0 {
		kvCount += 1
		allKvs = j.KeyValues
	}

	if len(instanceKvs) > 0 {
		kvCount += 1
		allKvs = instanceKvs
	}

	if kvCount > 1 {
		allKvs = make(map[string]interface{})
		for k, v := range j.KeyValues {
			allKvs[k] = v
		}
		for k, v := range instanceKvs {
			allKvs[k] = v
		}
	}

	return allKvs
}
