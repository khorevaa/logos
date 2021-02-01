package file

import (
	"github.com/khorevaa/logos/internal/common"
	"go.uber.org/zap/zapcore"
	"os"
)

type File struct {
	*os.File
}

type Config struct {
	FileName string `logos-config:"file_name" logos-validate:"required"`
}

var (
	defaultConfig = Config{}
)

func DefaultConfig() Config {
	return defaultConfig
}

func New(v *common.Config) (zapcore.WriteSyncer, error) {
	cfg := DefaultConfig()
	if err := v.Unpack(&cfg); err != nil {
		return nil, err
	}
	f, err := os.OpenFile(cfg.FileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	return &File{f}, nil
}
