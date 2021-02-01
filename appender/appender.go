package appender

import (
	"fmt"
	"github.com/khorevaa/logos/appender/console"
	"github.com/khorevaa/logos/appender/file"
	"github.com/khorevaa/logos/appender/gelfudp"
	"github.com/khorevaa/logos/appender/rollingfile"
	"github.com/khorevaa/logos/internal/common"
	"go.uber.org/zap/zapcore"
)

var (
	writers  = map[string]WriterFactory{}
	encoders = map[string]EncoderFactory{}
)

type WriterFactory func(config *common.Config) (zapcore.WriteSyncer, error)
type EncoderFactory func(*common.Config) (zapcore.Encoder, error)

type Appender struct {
	Writer  zapcore.WriteSyncer
	Encoder zapcore.Encoder
}

func init() {
	RegisterWriterType("console", console.New)
	RegisterWriterType("file", file.New)
	RegisterWriterType("rolling_file", rollingfile.New)
	RegisterWriterType("gelf_udp", gelfudp.New)
}

func CreateAppender(writerType string, config *common.Config) (*Appender, error) {
	w, err := NewWriter(writerType, config)
	if err != nil {
		return nil, err
	}
	encoderConfig, err := config.Child("encoder", -1)
	if err != nil {
		return nil, err
	}
	ec := EncoderConfig{}
	if err := encoderConfig.Unpack(&ec); err != nil {
		return nil, err
	}
	e, err := CreateEncoder(ec)
	if err != nil {
		return nil, err
	}
	return &Appender{w, e}, nil
}

func RegisterWriterType(name string, f WriterFactory) {
	if writers[name] != nil {
		panic(fmt.Errorf("writer type  '%v' exists already", name))
	}
	writers[name] = f
}

func NewWriter(name string, config *common.Config) (zapcore.WriteSyncer, error) {
	factory := writers[name]
	if factory == nil {
		return nil, fmt.Errorf("writer type %v undefined", name)
	}
	return factory(config)
}

type EncoderConfig struct {
	Namespace common.ConfigNamespace `logos-config:",inline"`
}

func RegisterEncoderType(name string, gen EncoderFactory) {
	if _, exists := encoders[name]; exists {
		panic(fmt.Sprintf("encoder %q already registered", name))
	}
	encoders[name] = gen
}

func CreateEncoder(cfg EncoderConfig) (zapcore.Encoder, error) {
	// default to json encoder
	encoder := "json"
	if name := cfg.Namespace.Name(); name != "" {
		encoder = name
	}

	factory := encoders[encoder]
	if factory == nil {
		return nil, fmt.Errorf("'%v' encoder is not available", encoder)
	}
	return factory(cfg.Namespace.Config())
}
