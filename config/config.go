package config

import (
	"github.com/khorevaa/logos/internal/common"
)

const DefaultConfig = `
appenders:
  console:
    - name: CONSOLE
      target: stdout
      encoder:
        console:
loggers:
  root:
    level: info
    appender_refs:
      - CONSOLE

scan: false
scan_period: 1m
`

type Config struct {
	Appenders map[string][]*common.Config `logos-config:"appenders"`
	Loggers   Loggers                     `logos-config:"loggers"`
}

type ScanConfig struct {
	Scan       bool   `logos-config:"scan"`
	ScanPeriod string `logos-config:"scan_period"`
}

type Loggers struct {
	Root   RootLogger     `logos-config:"root"`
	Logger []LoggerConfig `logos-config:"logger"`
}

type RootLogger struct {
	Level          string           `logos-config:"level"`
	AppenderRefs   []string         `logos-config:"appender_refs"`
	AppenderConfig []AppenderConfig `logos-config:"appenders"`
}

type LoggerConfig struct {
	Name           string           `logos-config:"name" logos-validate:"required"`
	Level          string           `logos-config:"level"`
	AddCaller      bool             `logos-config:"add_caller"`
	TraceLevel     string           `logos-config:"trace_level"`
	AppenderRefs   []string         `logos-config:"appender_refs"`
	AppenderConfig []AppenderConfig `logos-config:"appenders"`
}

type AppenderConfig struct {
	Name  string `logos-config:"name" logos-validate:"required"`
	Level string `logos-config:"level"`
}
