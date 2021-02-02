package logos

import (
	"github.com/khorevaa/logos/appender"
	"github.com/khorevaa/logos/config"
	"github.com/khorevaa/logos/internal/common"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"strings"
	"testing"
)

func TestInitWithConfigContent(t *testing.T) {
	const newConfig = `
appenders:
  console:
    - name: CONSOLE
      target: discard
      #no_color: true
      encoder:
        console:
          color_scheme:
            #debug_level: "Black,yellow"

  rolling_file:
    - name: ROLL_FILE
      file_name: ./logs/log.log
      max_size: 100
      encoder:
        json:
loggers:
  root:
    level: info
    appender_refs:
      - CONSOLE
  logger:
    - name: github.com
      level: debug
      add_caller: true
      trace_level: error
      appender_refs:
        - ROLL_FILE
      appenders:
        - name: CONSOLE
          level: debug

`
	err := InitWithConfigContent(newConfig)

	//l.SetLLevel(OffLevel)
	assert.Nil(t, err)

	l := New("github.com/logger")
	l.Info("hello")
	l.Debug("world")
	l2 := New("github.com/logger/v1")
	l2.Info("hello world test/logger/v1", zap.String("key", "val"))
	l2.Debug("hello world test/logger/v1")

	//err = InitWithConfigContent(newConfig)
	//assert.NotNil(t, err)

	l.Debug("hello world")
	l2.Debug("hello world test/logger/v1")
	l2.Error("hello world test/logger/v1", zap.Any("interface", []interface{}{1, "2", true}))
	l2.DPanic("hello world test/logger/v1", zap.Any("ints", []int{1, 2, 3321}))
	l2.Warn("hello world test/logger/v1", zap.Int("key", 123), zap.Bool("bool", false), zap.Any("bools", []bool{false, true, true}))

}

func Test_parseConfigFromString(t *testing.T) {

	tests := []struct {
		name    string
		args    []string
		want    string
		wantErr bool
	}{
		{
			"root debug",
			[]string{
				//"appenders.console.1.name=CONSOLE_TEST",
				//"appenders.console.1.target=stdout",
				//"appenders.console.1.no_color=true",
				//"appenders.console.1.encoder.console",
				//"loggers.logger.0.add_caller=true",
				//"loggers.logger.0.level=debug",
				//"loggers.logger.0.name=stdlog",
				//"loggers.root.appender_refs.0=CONSOLE",
				"loggers.root.level=debug",
			},
			`
loggers:
  root:
    level: debug`,
			false,
		},
		{
			"add logger",
			[]string{
				//"appenders.console.1.name=CONSOLE_TEST",
				//"appenders.console.1.target=stdout",
				//"appenders.console.1.no_color=true",
				//"appenders.console.1.encoder.console",
				"loggers.logger.0.add_caller=true",
				"loggers.logger.0.level=debug",
				"loggers.logger.0.name=github.com/khorevaa/logos",
				"loggers.logger.0.appender_refs.0=CONSOLE",
			},
			`
loggers:
  logger:
    - name: github.com/khorevaa/logos
      level: debug
      add_caller: true
      appender_refs:
        - CONSOLE
`,
			false,
		},
		{
			"add appender",
			[]string{
				"appenders.console.0.name=CONSOLE_TEST",
				"appenders.console.0.target=stdout",
				"appenders.console.0.no_color=true",
				"appenders.console.0.encoder.console",
				"loggers.logger.0.add_caller=true",
				"loggers.logger.0.level=debug",
				"loggers.logger.0.name=github.com/khorevaa/logos",
				"loggers.logger.0.appender_refs.0=CONSOLE_TEST",
			},
			`
appenders:
  console:
    - name: CONSOLE_TEST
      target: stdout
      no_color: true
      encoder:
        console:

loggers:
  logger:
    - name: github.com/khorevaa/logos
      level: debug
      add_caller: true
      appender_refs:
        - CONSOLE_TEST
`,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseConfigFromEnvString(strings.Join(tt.args, ";"))

			if err != nil {
				t.Errorf("parseConfigFromEnvString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			var cfg config.Config
			err = got.Unpack(&cfg)

			if (err != nil) != tt.wantErr {
				t.Errorf("parseConfigFromEnvString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			want, err := common.NewConfigFrom(tt.want)
			var wantCfg config.Config
			if err != nil {
				t.Errorf("parseConfigFromEnvString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			err = want.Unpack(&wantCfg)

			if err != nil {
				t.Errorf("parseConfigFromEnvString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(cfg.Appenders) > 0 {

				for appenderType, appenderConfigs := range cfg.Appenders {

					for idx, appenderConfig := range appenderConfigs {
						cfgAppender, err := appender.CreateAppender(appenderType, appenderConfig)
						assert.NoError(t, err)
						wantAppender, err := appender.CreateAppender(appenderType, wantCfg.Appenders[appenderType][idx])
						assert.NoError(t, err)
						assert.EqualValuesf(t, cfgAppender, wantAppender, "parseConfigFromEnvString() got = %v, want %v", cfgAppender, wantAppender)

					}

				}

			} else {
				assert.EqualValuesf(t, cfg.Loggers, wantCfg.Loggers, "parseConfigFromEnvString() got = %v, want %v", cfg, wantCfg)
			}

		})
	}
}
