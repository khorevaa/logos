package console

import (
	"encoding/base64"
	"fmt"
	"time"

	"go.uber.org/zap/zapcore"
)

// This file contains the methods that implement
// the zapcore ConsoleEncoder interfaces

func (e *Encoder) AppendBool(val bool) {
	e.buf.AppendBool(val)
}

func (e *Encoder) AppendByteString(bstr []byte) {
	e.buf.AppendString(string(bstr))
}

func (e *Encoder) AppendComplex128(val complex128) {
	r, i := float64(real(val)), float64(imag(val))
	e.buf.AppendFloat(r, 64)
	e.buf.AppendFloat(i, 64)
}

func (e *Encoder) AppendComplex64(val complex64) {
	e.AppendComplex128(complex128(val))
}

func (e *Encoder) AppendFloat64(val float64) {
	e.buf.AppendFloat(val, 64)
}

func (e *Encoder) AppendFloat32(val float32) {
	e.buf.AppendFloat(float64(val), 32)
}

func (e *Encoder) AppendInt(val int) {
	e.buf.AppendInt(int64(val))
}

func (e *Encoder) AppendInt64(val int64) {
	e.buf.AppendInt(val)
}

func (e *Encoder) AppendInt32(val int32) {
	e.buf.AppendInt(int64(val))
}

func (e *Encoder) AppendInt16(val int16) {
	e.buf.AppendInt(int64(val))
}

func (e *Encoder) AppendInt8(val int8) {
	e.buf.AppendInt(int64(val))
}

func (e *Encoder) AppendString(str string) {
	e.buf.AppendString(str)
}

func (e *Encoder) AppendUint(val uint) {
	e.buf.AppendUint(uint64(val))
}

func (e *Encoder) AppendUint64(val uint64) {
	e.buf.AppendUint(val)
}

func (e *Encoder) AppendUint32(val uint32) {
	e.buf.AppendUint(uint64(val))
}

func (e *Encoder) AppendUint16(val uint16) {
	e.buf.AppendUint(uint64(val))
}

func (e *Encoder) AppendUint8(val uint8) {
	e.buf.AppendUint(uint64(val))
}

func (e *Encoder) AppendUintptr(val uintptr) {
	e.AppendUint64(uint64(val))
}

func (e *Encoder) AddArray(key string, marshaler zapcore.ArrayMarshaler) error {
	e.addKey(key)
	return marshaler.MarshalLogArray(e)
}

func (e *Encoder) AddObject(key string, obj zapcore.ObjectMarshaler) error {
	e.addKey(key)
	return obj.MarshalLogObject(e)
}

func (e *Encoder) AddBinary(key string, val []byte) {
	e.AddString(key, base64.StdEncoding.EncodeToString(val))
}

func (e *Encoder) AddByteString(key string, val []byte) {
	e.AddString(key, base64.StdEncoding.EncodeToString(val))
}

func (e *Encoder) AddBool(key string, val bool) {
	e.addKey(key)
	e.AppendBool(val)
}

func (e *Encoder) AddComplex128(key string, val complex128) {
	e.addKey(key)
	e.AppendComplex128(val)
}

func (e *Encoder) AddComplex64(key string, val complex64) {
	e.addKey(key)
	e.AppendComplex128(complex128(val))
}

func (e *Encoder) AddDuration(key string, val time.Duration) {
	e.addKey(key)
	e.AppendDuration(val)
}

func (e *Encoder) AddFloat64(key string, val float64) {
	e.addKey(key)
	e.AppendFloat64(val)
}

func (e *Encoder) AddFloat32(key string, val float32) {
	e.AddFloat64(key, float64(val))
}

func (e *Encoder) AddInt(key string, val int) {
	e.AddInt64(key, int64(val))
}

func (e *Encoder) AddInt64(key string, val int64) {
	e.addKey(key)
	e.AppendInt64(val)
}

func (e *Encoder) AddInt32(key string, val int32) {
	e.AddInt64(key, int64(val))
}

func (e *Encoder) AddInt16(key string, val int16) {
	e.AddInt64(key, int64(val))
}

func (e *Encoder) AddInt8(key string, val int8) {
	e.AddInt64(key, int64(val))
}

func (e *Encoder) AddString(key string, val string) {
	e.addKey(key)
	e.AppendString(val)
}

func (e *Encoder) AddTime(key string, val time.Time) {
	e.addKey(key)
	e.AppendTime(val)
}

func (e *Encoder) AddUint(key string, val uint) {
	e.AddUint64(key, uint64(val))
}

func (e *Encoder) AddUint64(key string, val uint64) {
	e.addKey(key)
	e.AppendUint64(val)
}

func (e *Encoder) AddUint32(key string, val uint32) {
	e.AddUint64(key, uint64(val))
}

func (e *Encoder) AddUint16(key string, val uint16) {
	e.AddUint64(key, uint64(val))
}

func (e *Encoder) AddUint8(key string, val uint8) {
	e.AddUint64(key, uint64(val))
}

func (e *Encoder) AddUintptr(key string, val uintptr) {
	e.AddUint64(key, uint64(val))
}

func (e *Encoder) AddReflected(key string, val interface{}) error {
	e.addKey(key)
	v, ok := val.(string)
	if !ok {
		v = fmt.Sprintf("%v", val)
	}
	e.AppendString(v)
	return nil
}

func (e *Encoder) OpenNamespace(_ string) {
	// no-op -
	// namespaces do not really visually apply to console logs
}

func (e *Encoder) AppendDuration(val time.Duration) {
	//cur := e.buf.Len()

	//e.EncoderConfig.TimestampFormat.EncodeDuration(val, e)
	//if cur == e.buf.Len() {
	e.AppendInt64(int64(val))
	//}
}

func (e *Encoder) AppendTime(val time.Time) {

	cur := e.buf.Len()
	encodeTime := e.encodeTime(val)

	if len(encodeTime) > 0 {
		e.buf.AppendString(encodeTime)
	}

	if cur == e.buf.Len() {
		e.AppendInt64(val.UnixNano())
	}
}

func (e *Encoder) AppendArray(arr zapcore.ArrayMarshaler) error {
	e.buf.AppendByte('[')
	err := arr.MarshalLogArray(e)
	e.buf.AppendByte(']')
	return err
}

func (e *Encoder) AppendObject(obj zapcore.ObjectMarshaler) error {
	e.buf.AppendByte('{')
	err := obj.MarshalLogObject(e)
	e.buf.AppendByte('}')
	return err
}

func (e *Encoder) AppendReflected(val interface{}) error {
	v, ok := val.(string)
	if !ok {
		v = fmt.Sprintf("%v", val)
	}
	e.AppendString(v)
	return nil
}
