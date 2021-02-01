package console

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/khorevaa/logos/internal/common"
	"github.com/mattn/go-colorable"
	"go.uber.org/zap/zapcore"
)

type Console struct {
	zapcore.WriteSyncer
	colorable bool
}

type Config struct {
	Target  `logos-config:"target" logos-validate:"required,logos.oneof=stderr stdout discard"`
	NoColor bool `logos-config:"no_color"`
}

type Target = string

const (
	Discard Target = "discard"
	Stdout  Target = "stdout"
	Stderr  Target = "stderr"
)

var (
	defaultConfig = Config{
		Target: Stdout,
	}
)

func DefaultConfig() Config {
	return defaultConfig
}

func New(v *common.Config) (zapcore.WriteSyncer, error) {
	cfg := DefaultConfig()
	if err := v.Unpack(&cfg); err != nil {
		return nil, err
	}
	switch cfg.Target {
	case Stdout:
		return NewConsole(cfg, os.Stdout), nil
	case Stderr:
		return NewConsole(cfg, os.Stderr), nil
	case Discard:
		return &Console{zapcore.AddSync(ioutil.Discard), false}, nil
	default:
		return nil, fmt.Errorf("unknown target %q", cfg.Target)
	}

}

func NewConsole(config Config, file *os.File) zapcore.WriteSyncer {

	if config.NoColor {
		return &Console{zapcore.AddSync(colorable.NewNonColorable(file)), false}
	}

	return &Console{zapcore.AddSync(colorable.NewColorable(file)), true}

}
