package logos_test

import (
	"errors"
	"github.com/khorevaa/logos"
	log2 "log"
)

func ExampleNew_simple() {
	rawConfig := `
appenders:
 console:
   - name: CONSOLE
     target: stdout
     no_color: true
     encoder:
       console:
         disable_timestamp: true
         color_scheme:
           info_level: blue+b
           debug_level: green+b

loggers:
 root:
   level: debug
   appender_refs:
     - CONSOLE
`

	err := logos.InitWithConfigContent(rawConfig)
	if err != nil {
		panic(err)
	}

	//logos.CancelRedirectStdLog()

	log2.Println("1")
	log := logos.New("<your-package-name>") // like github.com/khorevaa/logos
	log.Info("This is me first log. Hello world logging systems")
	//cancel()
	//log2.Println("2")

	// Output:
	// INFO stdlog 1
	// INFO <your-package-name> This is me first log. Hello world logging systems
}

func ExampleNew_with_config_json() {

	rawConfig := `
appenders:
  console:
    - name: CONSOLE
      target: stdout
      encoder:
        json:

loggers:
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
      no_color: true
      encoder:
        console:
          disable_timestamp: true
          color_scheme:
            info_level: blue+b
            debug_level: green+b

loggers:
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

	// Output:
	// INFO <your-package-name> This is me first log. Hello world logging systems
	// DEBUG <your-package-name> This is me first error err=log system error
}
