package logos

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/khorevaa/logos/config"
	"github.com/khorevaa/logos/internal/common"
	"go.uber.org/zap/zapcore"

	_ "github.com/khorevaa/logos/appender"
	_ "github.com/khorevaa/logos/encoder/common"
	_ "github.com/khorevaa/logos/encoder/console"
	_ "github.com/khorevaa/logos/encoder/gelf"
	_ "github.com/khorevaa/logos/encoder/json"
)

var (
	manager    *logManager
	configFile string
	initLocker sync.Mutex
	debug      bool
)

func resolveConfigFileFromEnv() (string, error) {
	f := os.Getenv("LOGOS_CONFIG_FILE")
	if f == "" {
		return "", errors.New("environment variable 'LOGOS_CONFIG_FILE' is not set")
	}
	return f, nil
}

func resolveConfigFileFromWorkDir() (string, error) {
	matches1, _ := filepath.Glob("logos.yaml")
	matches2, _ := filepath.Glob("logos.yml")
	matches := append(matches1, matches2...)
	switch len(matches) {
	case 0:
		return "", errors.New("no config file found in work dir")
	case 1:
		return matches[0], nil
	default:
		panic(fmt.Errorf("multiple config files found %v", matches))
	}
}

func init() {

	initLocker.Lock()
	defer initLocker.Unlock()

	debug, _ = strconv.ParseBool(os.Getenv("LOGOS_DEBUG"))
	debugf("Logos is debugging on")

	if configFile == "" {
		cf, err := resolveConfigFileFromEnv()
		if err == nil {
			configFile = cf
		}
	}

	if configFile == "" {
		cf, err := resolveConfigFileFromWorkDir()
		if err == nil {
			configFile = cf
		}
	}

	var err error
	var rawConfig *common.Config

	if configFile != "" {
		// load ConfigFile
		configFile, err = filepath.Abs(configFile)
		if err != nil {
			panic(err)
		}

		if debug {
			debugf("logos using config file: <%s>", configFile)
			bs, err := ioutil.ReadFile(configFile)
			if err != nil {
				panic(err)
			}
			debugf(string(bs) + "\n")
		}

		rawConfig, _, err = common.LoadFile(configFile)
	} else {

		debugf("logos using default config:\n" + config.DefaultConfig)
		rawConfig, err = common.NewConfigFrom(config.DefaultConfig)
	}

	if err != nil {
		panic(err)
	}

	envConfig, err := parseConfigFromEnv()
	if err != nil && debug {
		fmt.Printf("logos loading config from env err: %s", err)
	}
	if envConfig != nil {
		err = rawConfig.Merge(envConfig)
		if err != nil {
			fmt.Printf("logos merge configs err: %s", err)
		}
	}

	manager, err = newLogManager(rawConfig)

	if err != nil {
		panic(err)
	}

	manager.RedirectStdLog()

	go func() {
		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
		<-quit
		Sync()
	}()

}
func parseConfigFromEnv() (*common.Config, error) {
	configData := os.Getenv("LOGOS_CONFIG")
	if configData == "" {
		return nil, ErrEnvConfigNotSet
	}
	return parseConfigFromEnvString(configData)
}

func parseConfigFromEnvString(configData string) (*common.Config, error) {

	configData = strings.TrimPrefix(configData, `"`)
	configData = strings.TrimSuffix(configData, `"`)

	newConfig := common.NewConfig()

	data := strings.Split(configData, ";")

	for _, strData := range data {

		strData = strings.TrimSpace(strData)
		pathValue := strings.Split(strData, "=")
		value := ""
		path := pathValue[0]
		if len(pathValue) == 2 {
			value = pathValue[1]
		}

		if len(value) == 0 && !strings.HasSuffix(path, ".") {
			// this is object
			path += "."
		}

		err := newConfig.SetString(path, -1, value)
		if err != nil {
			debugf("error loading config from path %s err <%s>\n", path, err.Error())
		}

	}

	return newConfig, nil
}

func InitWithConfigContent(content string) error {
	initLocker.Lock()
	defer initLocker.Unlock()

	debugf("logos InitWithConfigContent:\n" + content)

	rawConfig, err := common.NewConfigFrom(content)
	if err != nil {
		return err
	}

	err = manager.Update(rawConfig)
	if err != nil {
		return err
	}

	return nil
}

func New(name string) Logger {
	manager.Sync()
	return manager.NewLogger(name)
}
func SetLevel(name string, level zapcore.Level, appender ...string) {
	manager.SetLevel(name, level, appender...)
}

func Sync() {
	_ = manager.Sync()
}

func RedirectStdLog() func() {
	return manager.RedirectStdLog()
}

func CancelRedirectStdLog() {
	manager.CancelRedirectStdLog()
}
