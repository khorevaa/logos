package logos

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
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
loggerConfigs:
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
