package console

import (
	"fmt"
	"time"

	"github.com/khorevaa/logos/appender"
	"github.com/khorevaa/logos/internal/common"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

var defaultTimestampFormat = "2006-01-02T15:04:05.000Z0700"
var baseTimestamp = time.Now()

var defaultConfig = Config{
	ConsoleSeparator: " ",
	TimestampFormat:  defaultTimestampFormat,
	LineEnding:       "\n",
}

func init() {
	appender.RegisterEncoderType("console", func(cfg *common.Config) (zapcore.Encoder, error) {
		config := defaultConfig
		if cfg != nil {
			if err := cfg.Unpack(&config); err != nil {
				return nil, err
			}
		}

		encoderConfig := EncoderConfig{
			DisableColors:            config.DisableColors,
			ForceColors:              config.ForceColors,
			DisableNaming:            config.DisableNaming,
			DisableTimestamp:         config.DisableTimestamp,
			ConsoleSeparator:         config.ConsoleSeparator,
			LineEnding:               config.LineEnding,
			TimestampFormat:          config.TimestampFormat,
			UseTimePassedAsTimestamp: config.UseTimePassedAsTimestamp,
			UseUppercaseLevel:        config.UseUppercaseLevel,
		}

		encoderConfig.Schema = defaultScheme

		if config.ColorSchema != nil {
			encoderConfig.Schema = config.ColorSchema.Parse()
		}
		en := NewEncoder(encoderConfig)
		return en, nil

	})
}

// NewEncoder initializes a a bol.com tailored Encoder
func NewEncoder(cfg EncoderConfig) *Encoder {
	return &Encoder{
		buf:           bufferpool.Get(),
		EncoderConfig: cfg,
	}
}

// Encoder is a bol.com tailored zap encoder for
// writing human readable logs to the console
type Encoder struct {
	buf *buffer.Buffer
	EncoderConfig
}

// Clone implements the Clone method of the zapcore Encoder interface
func (e *Encoder) Clone() zapcore.Encoder {
	clone := e.clone()
	_, _ = clone.buf.Write(e.buf.Bytes())
	return clone
}

func (e *Encoder) clone() *Encoder {
	clone := get()
	clone.EncoderConfig = e.EncoderConfig
	clone.buf = bufferpool.Get()
	return clone
}

// EncodeEntry implements the EncodeEntry method of the zapcore Encoder interface
func (e *Encoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {

	line := bufferpool.Get()

	lvlColor := getLevelColor(ent.Level, e.Schema)

	e.appendTimeInfo(line, ent)

	e.colorizeText(line, ent.Level.CapitalString(), lvlColor)

	// pp.Println(e.DisableNaming, ent.LoggerName)
	if !e.DisableNaming && len(ent.LoggerName) > 0 {

		e.addSeparatorIfNecessary(line)
		e.colorizeText(line, ent.LoggerName, e.Schema.LogNaming)

	}

	if ent.Caller.Defined {

		e.addSeparatorIfNecessary(line)
		e.colorizeText(line, ent.Caller.TrimmedPath(), e.Schema.Nil)

	}

	if len(ent.Message) > 0 {
		e.addSeparatorIfNecessary(line)
		e.colorizeText(line, ent.Message, lvlColor)
	}

	if e.buf.Len() > 0 {
		line.Write(e.buf.Bytes())
	}

	// Add any structured context.
	e.writeContext(lvlColor, line, fields)

	// If there's no stacktrace key, honor that; this allows users to force
	// single-line output.
	if ent.Stack != "" {
		line.AppendByte('\n')
		line.AppendString(ent.Stack)
	}

	line.AppendString(zapcore.DefaultLineEnding)

	return line, nil

}

// appendTimeInfo appends the time related info
// appends nothing on DisableTimestamp
// appends [seconds] on UseTimePassedAsTimestamp
// appends formatted TimestampFormat else
func (e *Encoder) appendTimeInfo(buf *buffer.Buffer, entry zapcore.Entry) {
	if !e.DisableTimestamp {

		if e.UseTimePassedAsTimestamp {

			e.colorizeText(buf, fmt.Sprintf("[%04d]", int(entry.Time.Sub(baseTimestamp)/time.Second)), e.Schema.Timestamp)

		} else {

			if e.TimestampFormat == "" {
				e.colorizeText(buf, entry.Time.Format(defaultTimestampFormat), e.Schema.Timestamp)
			} else {
				e.colorizeText(buf, entry.Time.Format(e.TimestampFormat), e.Schema.Timestamp)
			}

		}

		e.addSeparatorIfNecessary(buf)
	}

}

func (e *Encoder) writeContext(defColor uint16, out *buffer.Buffer, extra []zapcore.Field) {

	if len(extra) == 0 {
		return
	}

	enc := getColoredEncoder(defColor, e.Schema, e.DisableColors)
	defer putColoredEncoder(enc)

	addFields(enc, extra)

	if enc.buf.Len() > 0 {
		out.Write(enc.buf.Bytes())
	}

}

func (e *Encoder) addSeparatorIfNecessary(line *buffer.Buffer) {
	if line.Len() > 0 {
		line.AppendString(e.ConsoleSeparator)
	}
}

func (e *Encoder) colorizeText(in *buffer.Buffer, text string, color uint16) {

	if e.DisableColors {
		in.AppendString(text)
		return
	}

	colorizeTextW(in, text, color)
}

func (e *Encoder) addKey(key string) {

	e.buf.Write([]byte(e.ConsoleSeparator))
	if e.DisableColors {
		e.buf.AppendString(key)
	} else {
		colorizeTextW(e.buf, key, e.Schema.FieldName)
	}
	e.buf.AppendByte('=')

}

func (e *Encoder) encodeTime(val time.Time) string {

	timestampFormat := e.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = defaultTimestampFormat
	}
	return val.Format(timestampFormat)

}

func addFields(enc zapcore.ObjectEncoder, fields []zapcore.Field) {
	for i := range fields {
		fields[i].AddTo(enc)
	}
}
