package config

import (
	"github.com/khorevaa/logos/appender/console"
	common2 "github.com/khorevaa/logos/encoder/common"
	"github.com/khorevaa/logos/internal/common"

	"testing"
)

type ConfigNamespace struct {
	name   string
	config console.Config
}

func TestConfigFrom(t *testing.T) {

	tests := []struct {
		name   string
		config interface{}
		text   []string
	}{
		{
			name: "simple",
			config: struct {
				Appenders map[string][]interface{} `logos-config:"appenders"`
				Loggers   Loggers                  `logos-config:"loggerConfigs"`
			}{
				Appenders: map[string][]interface{}{
					"console": {struct {
						Name    string      `logos-config:"name"`
						Target  string      `logos-config:"target"`
						Encoder interface{} `logos-config:"encoder"`
					}{
						Name:   "CONSOLE",
						Target: "stderr",
						Encoder: struct {
							Json common2.JsonEncoderConfig `logos-config:"json"`
						}{
							Json: common2.JsonEncoderConfig{
								TimeEncoder: "ISO8601",
							},
						},
					},
					},
				},
				Loggers: Loggers{
					Root: RootLogger{
						Level:        "error",
						AppenderRefs: []string{"CONSOLE"},
						//AppenderConfig: nil,
					},
				}},

			text: []string{"hello world", "hello"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			cfg, _ := common.NewConfigFrom(tt.config)

			defConfig, err := common.NewConfigFrom(DefaultConfig)

			cfg, err = common.MergeConfigs(defConfig, cfg)
			t.Error(err)

		})
	}
}
