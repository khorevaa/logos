package json

import (
	"github.com/khorevaa/logos/appender"
	ec "github.com/khorevaa/logos/encoder/common"
	"github.com/khorevaa/logos/internal/common"
	"go.uber.org/zap/zapcore"
)

var defaultConfig = ec.JsonEncoderConfig{
	TimeKey:       "ts",
	LevelKey:      "level",
	NameKey:       "logger",
	CallerKey:     "caller",
	MessageKey:    "msg",
	StacktraceKey: "stacktrace",
	LineEnding:    "\n",
	TimeEncoder:   "ISO8601",
}

func init() {
	appender.RegisterEncoderType("json", func(cfg *common.Config) (zapcore.Encoder, error) {
		config := defaultConfig
		if cfg != nil {
			if err := cfg.Unpack(&config); err != nil {
				return nil, err
			}
		}

		encoderConfig := zapcore.EncoderConfig{
			TimeKey:        config.TimeKey,
			LevelKey:       config.LevelKey,
			NameKey:        config.NameKey,
			CallerKey:      config.CallerKey,
			MessageKey:     config.MessageKey,
			StacktraceKey:  config.StacktraceKey,
			LineEnding:     config.LineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.EpochTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}

		te, err := ec.GetTimeEncoder(config.TimeEncoder)
		if err != nil {
			return nil, err
		}
		encoderConfig.EncodeTime = te

		return zapcore.NewJSONEncoder(encoderConfig), nil
	})
}
