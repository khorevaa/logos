package logos_test

import (
	"errors"
	"github.com/khorevaa/logos"
)

func ExampleNew_simple() {

	log := logos.New("<your-package-name>") // like github.com/khorevaa/logos
	log.Info("This is me first log. Hello world logging systems")

}

func ExampleNew_with_config_json() {

	rawConfig := `
appenders:
  console:
    - name: CONSOLE
      target: stdout
      encoder:
        json:

loggerConfigs:
  root:
    level: info
    appender_refs:
      - CONSOLE
`

	logos.InitWithConfigContent(rawConfig)

	log := logos.New("<your-package-name>") // like github.com/khorevaa/logos
	log.Info("This is me first log. Hello world logging systems")

}
func ExampleNew_with_color_scheme() {

	rawConfig := `
appenders:
  console:
    - name: CONSOLE
      target: stdout
      encoder:
        console:
          color_scheme:
            info_level: blue+b
            debug_level: green+b

loggerConfigs:
  root:
    level: debug
    appender_refs:
      - CONSOLE
`

	logos.InitWithConfigContent(rawConfig)

	log := logos.New("<your-package-name>") // like github.com/khorevaa/logos
	log.Info("This is me first log. Hello world logging systems")

	err := errors.New("log system error")
	log.Debug("This is me first error", logos.Any("err", err))

}
