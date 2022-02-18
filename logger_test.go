package logos

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_With(t *testing.T) {

	tests := []struct {
		name    string
		logName string
		fields  []Field
		want    string
	}{
		{
			"simple",
			"testLog",
			[]Field{
				String("string", "value"),
			},

			"logg",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := New(tt.logName)
			logWith := log.Named("named").With(tt.fields...)
			logWith.Info("info msg")

			logWith.Debug("Debug before")
			SetLevel(tt.logName, DebugLevel, "CONSOLE")
			logWith.Debug("Debug after")
			log.Debug("log debug")

			assert.Equalf(t, tt.want, "", "With(%v)", tt.fields)
		})
	}
}
