package console

import (
	"encoding/base64"
	"fmt"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
	"strconv"
	"sync"
	"time"
)

var poolColoredEncoder = sync.Pool{
	New: func() interface{} {
		return &coloredEncoder{
			buf: bufferpool.Get(),
		}
	},
}

func getColoredEncoder(lvlColor uint16, scheme ColorScheme) *coloredEncoder {
	enc := poolColoredEncoder.Get().(*coloredEncoder)
	enc.buf = bufferpool.Get()
	enc.scheme = scheme
	enc.entLevelColor = lvlColor
	return enc
}

func getLevelColor(level zapcore.Level, scheme ColorScheme) uint16 {
	switch level {
	case zapcore.DebugLevel:
		return scheme.DebugLevel
	case zapcore.InfoLevel:
		return scheme.InfoLevel
	case zapcore.ErrorLevel:
		return scheme.ErrorLevel
	case zapcore.PanicLevel:
		return scheme.PanicLevel
	case zapcore.WarnLevel:
		return scheme.WarnLevel
	case zapcore.DPanicLevel:
		return scheme.DPanicLevel
	case zapcore.FatalLevel:
		return scheme.FatalLevel
	default:
		return scheme.String
	}
}

func putColoredEncoder(enc *coloredEncoder) {
	enc.scheme = defaultScheme
	enc.buf.Free()
	enc.entLevelColor = NoColor
	enc.EncodeDuration = nil
	enc.EncodeTime = nil
	poolColoredEncoder.Put(enc)
}

type coloredEncoder struct {
	buf          *buffer.Buffer
	disableColor bool
	scheme       ColorScheme

	entLevelColor uint16

	EncodeDuration zapcore.DurationEncoder
	EncodeTime     zapcore.TimeEncoder
}

func (e *coloredEncoder) addKey(key string) {
	e.buf.AppendByte(' ')
	e.appendColoredString(key, e.scheme.FieldName)
	e.buf.AppendByte('=')
}

func (e *coloredEncoder) appendColoredString(val string, color uint16) {

	if e.disableColor {
		e.buf.AppendString(val)
		return
	}

	e.buf.AppendString(colorizeText(val, color))

}

func (e *coloredEncoder) AppendBool(val bool) {
	e.addElementSeparator()
	e.appendColoredString(strconv.FormatBool(val), e.scheme.Bool)
}

func (e *coloredEncoder) AppendByteString(bstr []byte) {
	e.addElementSeparator()
	e.buf.AppendString(string(bstr))
}

func (e *coloredEncoder) AppendComplex128(val complex128) {
	r, i := float64(real(val)), float64(imag(val))

	str := fmt.Sprintf("%s%s",
		strconv.FormatFloat(r, 'f', -1, 64),
		strconv.FormatFloat(i, 'f', -1, 64))
	e.addElementSeparator()
	e.appendColoredString(str, e.scheme.Float)

}

func (e *coloredEncoder) AppendComplex64(val complex64) {
	e.AppendComplex128(complex128(val))
}

func (e *coloredEncoder) AppendFloat64(val float64) {
	e.addElementSeparator()
	e.appendColoredString(strconv.FormatFloat(val, 'f', -1, 64), e.scheme.Float)
}

func (e *coloredEncoder) AppendFloat32(val float32) {
	e.AppendFloat64(float64(val))
}

func (e *coloredEncoder) AppendInt(val int) {
	e.AppendInt64(int64(val))
}

func (e *coloredEncoder) AppendInt64(val int64) {
	e.addElementSeparator()
	e.appendColoredString(strconv.FormatInt(val, 10), e.scheme.Integer)
}

func (e *coloredEncoder) AppendInt32(val int32) {
	e.AppendInt64(int64(val))
}

func (e *coloredEncoder) AppendInt16(val int16) {
	e.AppendInt64(int64(val))
}

func (e *coloredEncoder) AppendInt8(val int8) {
	e.AppendInt64(int64(val))
}

func (e *coloredEncoder) AppendString(str string) {
	e.addElementSeparator()
	e.appendColoredString(str, e.scheme.String)
}

func (e *coloredEncoder) AppendUint(val uint) {
	e.AppendUint64(uint64(val))
}

func (e *coloredEncoder) AppendUint64(val uint64) {
	e.addElementSeparator()
	e.appendColoredString(strconv.FormatUint(val, 64), e.scheme.PointerAddress)
}

func (e *coloredEncoder) AppendUint32(val uint32) {
	e.AppendUint64(uint64(val))
}

func (e *coloredEncoder) AppendUint16(val uint16) {
	e.AppendUint64(uint64(val))
}

func (e *coloredEncoder) AppendUint8(val uint8) {
	e.AppendUint64(uint64(val))
}

func (e *coloredEncoder) AppendUintptr(val uintptr) {
	e.AppendUint64(uint64(val))
}

func (e *coloredEncoder) AddArray(key string, marshaler zapcore.ArrayMarshaler) error {
	e.addKey(key)
	return e.AppendArray(marshaler)
}

func (e *coloredEncoder) AddObject(key string, obj zapcore.ObjectMarshaler) error {
	e.addKey(key)
	return e.AppendObject(obj)
}

//func (enc *coloredEncoder) AppendArray(arr ArrayMarshaler) error {
//	enc.addElementSeparator()
//	enc.buf.AppendByte('[')
//	err := arr.MarshalLogArray(enc)
//	enc.buf.AppendByte(']')
//	return err
//}
//
//func (enc *jsonEncoder) AppendObject(obj ObjectMarshaler) error {
//	enc.addElementSeparator()
//	enc.buf.AppendByte('{')
//	err := obj.MarshalLogObject(enc)
//	enc.buf.AppendByte('}')
//	return err
//}

func (e *coloredEncoder) AddBinary(key string, val []byte) {
	e.AddString(key, base64.StdEncoding.EncodeToString(val))
}

func (e *coloredEncoder) AddByteString(key string, val []byte) {
	e.AddString(key, base64.StdEncoding.EncodeToString(val))
}

func (e *coloredEncoder) AddBool(key string, val bool) {
	e.addKey(key)
	e.AppendBool(val)
}

func (e *coloredEncoder) AddComplex128(key string, val complex128) {
	e.addKey(key)
	e.AppendComplex128(val)
}

func (e *coloredEncoder) AddComplex64(key string, val complex64) {
	e.addKey(key)
	e.AppendComplex128(complex128(val))
}

func (e *coloredEncoder) AddDuration(key string, val time.Duration) {
	e.addKey(key)
	e.AppendDuration(val)
}

func (e *coloredEncoder) AddFloat64(key string, val float64) {
	e.addKey(key)
	e.AppendFloat64(val)
}

func (e *coloredEncoder) AddFloat32(key string, val float32) {
	e.AddFloat64(key, float64(val))
}

func (e *coloredEncoder) AddInt(key string, val int) {
	e.AddInt64(key, int64(val))
}

func (e *coloredEncoder) AddInt64(key string, val int64) {
	e.addKey(key)
	e.AppendInt64(val)
}

func (e *coloredEncoder) AddInt32(key string, val int32) {
	e.AddInt64(key, int64(val))
}

func (e *coloredEncoder) AddInt16(key string, val int16) {
	e.AddInt64(key, int64(val))
}

func (e *coloredEncoder) AddInt8(key string, val int8) {
	e.AddInt64(key, int64(val))
}

func (e *coloredEncoder) AddString(key string, val string) {
	e.addKey(key)
	e.AppendString(val)
}

func (e *coloredEncoder) AddTime(key string, val time.Time) {
	e.addKey(key)
	e.AppendTime(val)
}

func (e *coloredEncoder) AddUint(key string, val uint) {
	e.AddUint64(key, uint64(val))
}

func (e *coloredEncoder) AddUint64(key string, val uint64) {
	e.addKey(key)
	e.AppendUint64(val)
}

func (e *coloredEncoder) AddUint32(key string, val uint32) {
	e.AddUint64(key, uint64(val))
}

func (e *coloredEncoder) AddUint16(key string, val uint16) {
	e.AddUint64(key, uint64(val))
}

func (e *coloredEncoder) AddUint8(key string, val uint8) {
	e.AddUint64(key, uint64(val))
}

func (e *coloredEncoder) AddUintptr(key string, val uintptr) {
	e.AddUint64(key, uint64(val))
}

func (e *coloredEncoder) AddReflected(key string, val interface{}) error {
	e.addKey(key)
	v, ok := val.(string)
	if !ok {
		v = fmt.Sprintf("%v", val)
	}
	e.AppendString(v)
	return nil
}

func (e *coloredEncoder) OpenNamespace(_ string) {
	// no-op -
	// namespaces do not really visually apply to console logs
}

func (e *coloredEncoder) AppendDuration(val time.Duration) {
	cur := e.buf.Len()
	e.EncodeDuration(val, e)
	if cur == e.buf.Len() {
		e.AppendInt64(int64(val))
	}
}

func (e *coloredEncoder) AppendTime(val time.Time) {
	cur := e.buf.Len()
	e.EncodeTime(val, e)
	if cur == e.buf.Len() {
		e.AppendInt64(val.UnixNano())
	}
}

func (e *coloredEncoder) AppendArray(arr zapcore.ArrayMarshaler) error {
	e.addElementSeparator()
	e.buf.AppendByte('[')
	err := arr.MarshalLogArray(e)
	e.buf.AppendByte(']')
	return err
}

func (e *coloredEncoder) AppendObject(obj zapcore.ObjectMarshaler) error {
	e.addElementSeparator()
	e.buf.AppendByte('{')
	err := obj.MarshalLogObject(e)
	e.buf.AppendByte('}')
	return err
}

func (e *coloredEncoder) AppendReflected(val interface{}) error {
	v, ok := val.(string)
	if !ok {
		v = fmt.Sprintf("%v", val)
	}
	e.addElementSeparator()
	e.AppendString(v)
	return nil
}

func (e *coloredEncoder) addElementSeparator() {
	last := e.buf.Len() - 1
	if last < 0 {
		return
	}
	switch e.buf.Bytes()[last] {
	case '{', '[', ':', ',', ' ', '=':
		return
	default:
		e.buf.AppendByte(',')
	}
}
