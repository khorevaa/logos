# Logos: like log4j, but for golang.

[![go.dev][pkg-img]][pkg] [![goreport][report-img]][report] [![build][build-img]][build] [![coverage][cov-img]][cov] ![stability-stable][stability-img]

## Features

This project is a wrapper around the excellent logging framework zap.

* Dependency
    - `go.uber.org/zap` for logging
    - `github.com/elastic/go-ucfg` for config logging system
* Simple and Clean Interface
* One log manager for all logs
* Hot config update from file or env
* Appenders
    - `Console`, *write to console*
    - `File`, *any log file*
    - `gelfupd`, *greylog logger*
    - `RollingFile`, *rolling file writing & compress*
* Encoders
    - `Console`, *colorful & formatting text for console*
    - `Gelf`, *gelf for greylog*
    - `Json`, *standard json encoder*
* Useful utility function
    - `Setlevel(LogName string, level int)`, *hot update logger level*
    - `UpdateLogger(LogName string, logger *zap.Logger*)`, *hot update core logger*
    - `RedirectStdLog()`, *redirect standard log package*
* High Performance
    - [Significantly faster][high-performance] than all other json loggers.

## Interfaces

### Logger
```go
// DefaultLogger is the global logger.
```

### Json  Writer

To log a machine-friendly, use `json`. [![playground][play-pretty-img]][play-pretty]

```go

```

### Pretty Console Writer

To log a human-friendly, colorized output, use `Console`. [![playground][play-pretty-img]][play-pretty]

```go

```
![Pretty logging][pretty-img]
> Note: pretty logging also works on windows console

### High Performance

A quick and simple benchmark with zap/zerolog, which runs on [github actions][benchmark]:

```go
// go test -v -cpu=4 -run=none -bench=. -benchtime=10s -benchmem log_test.go
package main

import (
	"io/ioutil"
	"testing"

	"github.com/khorevaa/logos"
	"github.com/phuslu/log"
	"github.com/rs/zerolog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var fakeMessage = "Test logging, but use a somewhat realistic message length. "

func BenchmarkLogos(b *testing.B) {

	const newConfig = `
appenders:
  console:
    - name: CONSOLE
      target: discard
      encoder:
        console:
loggers:
  root:
    level: info
    appender_refs:
      - CONSOLE
`
	err := logos.InitWithConfigContent(newConfig)
	if err != nil {
		panic(err)
	}

	logger := logos.New("benchmark")
	for i := 0; i < b.N; i++ {
		logger.Info(fakeMessage, zap.String("foo", "bar"), zap.Int("int", 42))
	}
}

func BenchmarkZap(b *testing.B) {
	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(ioutil.Discard),
		zapcore.InfoLevel,
	))
	for i := 0; i < b.N; i++ {
		logger.Info(fakeMessage, zap.String("foo", "bar"), zap.Int("int", 42))
	}
}

func BenchmarkZeroLog(b *testing.B) {
	logger := zerolog.New(ioutil.Discard).With().Timestamp().Logger()
	for i := 0; i < b.N; i++ {
		logger.Info().Str("foo", "bar").Int("int", 42).Msg(fakeMessage)
	}
}

func BenchmarkPhusLog(b *testing.B) {
	logger := log.Logger{
		TimeFormat: "", // uses rfc3339 by default
		Writer:     log.IOWriter{ioutil.Discard},
	}
	for i := 0; i < b.N; i++ {
		logger.Info().Str("foo", "bar").Int("int", 42).Msg(fakeMessage)
	}
}



```
A Performance result as below, for daily benchmark results see [github actions][benchmark]
```

```


[pkg-img]: http://img.shields.io/badge/godoc-reference-5272B4.svg
[pkg]: https://godoc.org/github.com/khorevaa/logos
[report-img]: https://goreportcard.com/badge/github.com/khorevaa/logos
[report]: https://goreportcard.com/report/github.com/khorevaa/logos
[build-img]: https://github.com/khorevaa/logos/workflows/build/badge.svg
[build]: https://github.com/khorevaa/logos/actions
[cov-img]: http://gocover.io/_badge/github.com/khorevaa/logos
[cov]: https://gocover.io/github.com/khorevaa/logos
[stability-img]: https://img.shields.io/badge/stability-stable-green.svg
[high-performance]: https://github.com/khorevaa/logos#high-performance
[play-simple-img]: https://img.shields.io/badge/playground-NGV25aBKmYH-29BEB0?style=flat&logo=go
[play-simple]: https://play.golang.org/p/NGV25aBKmYH
[play-customize-img]: https://img.shields.io/badge/playground-emTsJJKUGXZ-29BEB0?style=flat&logo=go
[play-customize]: https://play.golang.org/p/emTsJJKUGXZ
[play-file-img]: https://img.shields.io/badge/playground-nS--ILxFyhHM-29BEB0?style=flat&logo=go
[play-file]: https://play.golang.org/p/nS-ILxFyhHM
[play-pretty-img]: https://img.shields.io/badge/playground-SCcXG33esvI-29BEB0?style=flat&logo=go
[play-pretty]: https://play.golang.org/p/SCcXG33esvI
[pretty-img]: https://user-images.githubusercontent.com/195836/101993218-cda82380-3cf3-11eb-9aa2-b8b1c832a72e.png
[play-formatting-img]: https://img.shields.io/badge/playground-UmJmLxYXwRO-29BEB0?style=flat&logo=go
[play-formatting]: https://play.golang.org/p/UmJmLxYXwRO
[play-context-img]: https://img.shields.io/badge/playground-oAVAo302faf-29BEB0?style=flat&logo=go
[play-context]: https://play.golang.org/p/oAVAo302faf
[play-marshal-img]: https://img.shields.io/badge/playground-NxMoqaiVxHM-29BEB0?style=flat&logo=go
[play-marshal]: https://play.golang.org/p/NxMoqaiVxHM
[play-interceptor]: https://play.golang.org/p/upmVP5cO62Y
[play-interceptor-img]: https://img.shields.io/badge/playground-upmVP5cO62Y-29BEB0?style=flat&logo=go
[benchmark]: https://github.com/khorevaa/logos/actions?query=workflow%3Abenchmark
[zerolog]: https://github.com/rs/zerolog
[glog]: https://github.com/golang/glog
[quicktemplate]: https://github.com/valyala/quicktemplate
[gjson]: https://github.com/tidwall/gjson
[zap]: https://github.com/uber-go/zap
[lumberjack]: https://github.com/natefinch/lumberjack