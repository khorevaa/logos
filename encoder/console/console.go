package console

import (
	"fmt"
	"github.com/khorevaa/logos/appender"
	"github.com/khorevaa/logos/internal/common"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
	"time"
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
func (e Encoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {

	line := bufferpool.Get()

	lvlColor := getLevelColor(ent.Level, e.Schema)

	e.appendTimeInfo(line, ent)

	line.AppendString(e.colorizeText(ent.Level.CapitalString(), lvlColor))

	//pp.Println(e.DisableNaming, ent.LoggerName)
	if !e.DisableNaming && len(ent.LoggerName) > 0 {

		e.addSeparatorIfNecessary(line)
		line.AppendString(e.colorizeText(ent.LoggerName, e.Schema.LogNaming))

	}

	if ent.Caller.Defined {

		e.addSeparatorIfNecessary(line)
		line.AppendString(e.colorizeText(ent.Caller.TrimmedPath(), e.Schema.Nil))

	}

	if len(ent.Message) > 0 {
		e.addSeparatorIfNecessary(line)
		line.AppendString(e.colorizeText(ent.Message, lvlColor))
	}
	//e.addSeparatorIfNecessary(line)
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
		var timeInfo string
		if e.UseTimePassedAsTimestamp {
			timeInfo = fmt.Sprintf("[%04d]", int(entry.Time.Sub(baseTimestamp)/time.Second))
		} else {
			timestampFormat := e.TimestampFormat
			if timestampFormat == "" {
				timestampFormat = defaultTimestampFormat
			}
			timeInfo = entry.Time.Format(timestampFormat)
		}

		buf.AppendString(e.colorizeText(timeInfo, e.Schema.Timestamp))
		e.addSeparatorIfNecessary(buf)
	}

}

func (e *Encoder) writeContext(defColor uint16, out *buffer.Buffer, extra []zapcore.Field) {

	if len(extra) == 0 {
		return
	}

	var enc zapcore.ObjectEncoder
	if !e.DisableColors {
		enc = getColoredEncoder(defColor, e.Schema, e.DisableColors)
		defer putColoredEncoder(enc.(*coloredEncoder))
	} else {
		enc = e
	}

	addFields(enc, extra)
	var buf *buffer.Buffer
	switch t := enc.(type) {
	case *Encoder:
		buf = t.buf
	case *coloredEncoder:
		buf = t.buf
	}

	if buf.Len() > 0 {
		out.Write(buf.Bytes())
	}

}

func (e *Encoder) addSeparatorIfNecessary(line *buffer.Buffer) {
	if line.Len() > 0 {
		line.AppendString(e.ConsoleSeparator)
	}
}

func (e *Encoder) colorizeText(text string, color uint16) string {

	if e.DisableColors {
		return text
	}

	return colorizeText(text, color)

}

func (e *Encoder) addKey(key string) {

	if !e.DisableColors {
		key = colorizeText(key, e.Schema.FieldName)
	}

	e.buf.Write([]byte(e.ConsoleSeparator))
	e.buf.AppendString(key)
	e.buf.AppendByte('=')

}

func (e *Encoder) encodeTime(val time.Time) string {

	timestampFormat := e.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = defaultTimestampFormat
	}
	timeInfo := val.Format(timestampFormat)

	return timeInfo

}

func addFields(enc zapcore.ObjectEncoder, fields []zapcore.Field) {
	for i := range fields {
		fields[i].AddTo(enc)
	}
}
