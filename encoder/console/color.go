package console

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	// No color
	NoColor uint16 = 1 << 15
)

const (
	// Foreground colors for ColorScheme.
	_ uint16 = iota | NoColor
	Black
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
	bitsForeground       = 0
	maskForeground       = 0xf
	ansiForegroundOffset = 30 - 1
)

const (
	// Background colors for ColorScheme.
	_ uint16 = iota<<bitsBackground | NoColor
	BackgroundBlack
	BackgroundRed
	BackgroundGreen
	BackgroundYellow
	BackgroundBlue
	BackgroundMagenta
	BackgroundCyan
	BackgroundWhite
	bitsBackground       = 4
	maskBackground       = 0xf << bitsBackground
	ansiBackgroundOffset = 40 - 1
)

const (
	// Bold flag for ColorScheme.
	Bold     uint16 = 1<<bitsBold | NoColor
	bitsBold        = 8
	maskBold        = 1 << bitsBold
	ansiBold        = 1
)

// To use with SetColorScheme.
type ColorScheme struct {
	Bool            uint16
	Integer         uint16
	Float           uint16
	String          uint16
	StringQuotation uint16
	EscapedChar     uint16
	FieldName       uint16
	PointerAddress  uint16
	Nil             uint16
	Time            uint16
	StructName      uint16
	ObjectLength    uint16

	LogNaming   uint16
	Timestamp   uint16
	InfoLevel   uint16
	WarnLevel   uint16
	ErrorLevel  uint16
	FatalLevel  uint16
	PanicLevel  uint16
	DPanicLevel uint16
	DebugLevel  uint16
}

var (
	defaultScheme = ColorScheme{
		Bool:            Cyan | Bold,
		Integer:         Blue | Bold,
		Float:           Magenta | Bold,
		String:          Red,
		StringQuotation: Red | Bold,
		EscapedChar:     Magenta | Bold,
		FieldName:       Yellow,
		PointerAddress:  Blue | Bold,
		Nil:             Cyan | Bold,
		Time:            Blue | Bold,
		StructName:      Green,
		ObjectLength:    Blue,

		LogNaming:   Black | Bold | BackgroundWhite,
		Timestamp:   Blue | Bold,
		InfoLevel:   Green,
		WarnLevel:   Yellow | Bold,
		ErrorLevel:  Red,
		FatalLevel:  Red | Bold,
		PanicLevel:  Red | Bold,
		DPanicLevel: Black | Bold | BackgroundRed,
		DebugLevel:  Blue,
	}

	colorMap = map[string]uint16{
		"black":   Black,
		"red":     Red,
		"green":   Green,
		"yellow":  Yellow,
		"blue":    Blue,
		"magenta": Magenta,
		"cyan":    Cyan,
		"white":   White,
	}

	colorBgMap = map[string]uint16{
		"black":   BackgroundBlack,
		"red":     BackgroundRed,
		"green":   BackgroundGreen,
		"yellow":  BackgroundYellow,
		"blue":    BackgroundBlue,
		"magenta": BackgroundMagenta,
		"cyan":    BackgroundCyan,
		"white":   BackgroundWhite,
	}
)

func getColor(color string) uint16 {

	val, ok := colorMap[strings.ToLower(color)]
	if ok {
		return val
	}
	return 0
}

func getBgColor(color string) uint16 {

	val, ok := colorBgMap[strings.ToLower(color)]
	if ok {
		return val
	}
	return 0
}

func (cs *ColorScheme) fixColors() {
	typ := reflect.Indirect(reflect.ValueOf(cs))
	defaultType := reflect.ValueOf(defaultScheme)
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.Uint() == 0 {
			field.SetUint(defaultType.Field(i).Uint())
		}
	}
}

func colorizeText(text string, color uint16) string {
	foreground := color & maskForeground >> bitsForeground
	background := color & maskBackground >> bitsBackground
	bold := color & maskBold

	if foreground == 0 && background == 0 && bold == 0 {
		return text
	}

	modBold := ""
	modForeground := ""
	modBackground := ""

	if bold > 0 {
		modBold = "\033[1m"
	}
	if foreground > 0 {
		modForeground = fmt.Sprintf("\033[%dm", foreground+ansiForegroundOffset)
	}
	if background > 0 {
		modBackground = fmt.Sprintf("\033[%dm", background+ansiBackgroundOffset)
	}

	return fmt.Sprintf("%s%s%s%s\033[0m", modForeground, modBackground, modBold, text)
}
