package common

import (
	"fmt"
	"go.uber.org/zap/zapcore"
)

// Config is used to pass encoding parameters to New.
type JsonEncoderConfig struct {
	TimeKey       string `logos-config:"time_key"`
	LevelKey      string `logos-config:"level_key"`
	NameKey       string `logos-config:"name_key"`
	CallerKey     string `logos-config:"caller_key"`
	MessageKey    string `logos-config:"message_key"`
	StacktraceKey string `logos-config:"stacktrace_key"`
	LineEnding    string `logos-config:"line_ending"`
	TimeEncoder   string `logos-config:"time_encoder" logos-validate:"logos.oneof=epoch epoch_millis epoch_nanos ISO8601"`
}

func GetTimeEncoder(name string) (zapcore.TimeEncoder, error) {
	switch name {
	case "epoch":
		return zapcore.EpochTimeEncoder, nil
	case "epoch_millis":
		return zapcore.EpochMillisTimeEncoder, nil
	case "epoch_nanos":
		return zapcore.EpochNanosTimeEncoder, nil
	case "ISO8601":
		return zapcore.ISO8601TimeEncoder, nil
	default:
		return nil, fmt.Errorf("no such TimeEncoder %q", name)
	}
}
