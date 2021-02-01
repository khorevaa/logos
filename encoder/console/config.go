package console

import "strings"

// Schema is the color schema for the default log parts/levels
type ColorSchemaConfig struct {
	Timestamp   string `logos-config:"timestamp"`
	Naming      string `logos-config:"naming"`
	InfoLevel   string `logos-config:"info_level"`
	WarnLevel   string `logos-config:"warn_level"`
	ErrorLevel  string `logos-config:"error_level"`
	FatalLevel  string `logos-config:"fatal_level"`
	PanicLevel  string `logos-config:"panic_level"`
	DPanicLevel string `logos-config:"dpanic_level"`
	DebugLevel  string `logos-config:"debug_level"`
}

func (c ColorSchemaConfig) Parse() ColorScheme {

	scheme := ColorScheme{}

	scheme.Timestamp = parseFieldColor(c.Timestamp)
	scheme.LogNaming = parseFieldColor(c.Naming)
	scheme.InfoLevel = parseFieldColor(c.InfoLevel)
	scheme.WarnLevel = parseFieldColor(c.WarnLevel)
	scheme.ErrorLevel = parseFieldColor(c.ErrorLevel)
	scheme.FatalLevel = parseFieldColor(c.FatalLevel)
	scheme.PanicLevel = parseFieldColor(c.PanicLevel)
	scheme.DPanicLevel = parseFieldColor(c.DPanicLevel)
	scheme.DebugLevel = parseFieldColor(c.DebugLevel)

	scheme.fixColors()
	return scheme
}

func parseFieldColor(colorString string) uint16 {

	colors := strings.Split(colorString, ",")

	if len(colors) == 0 {
		return 0
	}

	color := colors[0]
	bold := uint16(0)
	bgColor := ""

	if arr := strings.Split(color, "+"); len(arr) > 1 {
		color = strings.TrimSpace(arr[0])

		if strings.TrimSpace(arr[1]) == "b" {
			bold = Bold
		}

	}
	if len(colors) > 1 {
		bgColor = strings.TrimSpace(colors[1])
	}

	return getColor(color) | bold | getBgColor(bgColor)

}

// Config is used to pass encoding parameters to New.
type Config struct {

	// color schema for messages
	ColorSchema *ColorSchemaConfig `logos-config:"color_scheme"`
	// no colors
	DisableColors bool `logos-config:"disable_colors"`
	// no check for TTY terminal
	ForceColors bool `logos-config:"force_colors"`
	// false -> name passed, true -> github.com/khorevaa/logos
	DisableNaming bool `logos-config:"disable_naming"`
	// no timestamp
	DisableTimestamp bool `logos-config:"disable_timestamp"`
	// console separator default space
	ConsoleSeparator string `logos-config:"console_separator"`
	// false -> time passed, true -> timestamp
	UseTimePassedAsTimestamp bool `logos-config:"pass_timestamp"`
	// false -> info, true -> INFO
	UseUppercaseLevel bool `logos-config:"uppercase_level"`

	TimestampFormat string `logos-config:"timestamp_format"`

	LineEnding string `logos-config:"line_ending"`
}

// EncoderConfig is used to pass encoding parameters to New.
type EncoderConfig struct {

	// no colors
	DisableColors bool
	// no check for TTY terminal
	ForceColors bool
	// false -> time passed, true -> timestamp
	UseTimePassedAsTimestamp bool
	// false -> info, true -> INFO
	UseUppercaseLevel bool
	// false -> name passed, true -> github.com/khorevaa/logos
	DisableNaming bool
	// no timestamp
	DisableTimestamp bool
	// console separator default space
	ConsoleSeparator string
	// line end for log
	LineEnding string
	// time format string
	TimestampFormat string
	// color schema for messages
	Schema ColorScheme
}
