package benchmarks

import (
	"errors"
	"fmt"
	"github.com/khorevaa/logos"
	"time"

	"go.uber.org/multierr"
	"go.uber.org/zap"

	"go.uber.org/zap/zapcore"
)

var (
	errExample = errors.New("fail")

	_messages   = fakeMessages(1000)
	_tenInts    = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	_tenStrings = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	_tenTimes   = []time.Time{
		time.Unix(0, 0),
		time.Unix(1, 0),
		time.Unix(2, 0),
		time.Unix(3, 0),
		time.Unix(4, 0),
		time.Unix(5, 0),
		time.Unix(6, 0),
		time.Unix(7, 0),
		time.Unix(8, 0),
		time.Unix(9, 0),
	}
	_oneUser = &user{
		Name:      "Jane Doe",
		Email:     "jane@test.com",
		CreatedAt: time.Date(1980, 1, 1, 12, 0, 0, 0, time.UTC),
	}
	_tenUsers = users{
		_oneUser,
		_oneUser,
		_oneUser,
		_oneUser,
		_oneUser,
		_oneUser,
		_oneUser,
		_oneUser,
		_oneUser,
		_oneUser,
	}
)

func init() {
	const newConfig = `
appenders:
  console:
    - name: CONSOLE
      no_color: true
      target: discard
      encoder:
        console:
           disable_colors: true
           disable_timestamp: true
loggers:
  root:
    level: error
    appender_refs:
      - CONSOLE
  loggers:
    logger:
      name: test
      level: error
`
	err := logos.InitWithConfigContent(newConfig)
	if err != nil {
		panic(err)
	}
}

func fakeMessages(n int) []string {
	messages := make([]string, n)
	for i := range messages {
		messages[i] = fmt.Sprintf("Test logging, but use a somewhat realistic message length. (#%v)", i)
	}
	return messages
}

func getMessage(iter int) string {
	return _messages[iter%1000]
}

type users []*user

func (uu users) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	var err error
	for i := range uu {
		err = multierr.Append(err, arr.AppendObject(uu[i]))
	}
	return err
}

type user struct {
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *user) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("name", u.Name)
	enc.AddString("email", u.Email)
	enc.AddInt64("createdAt", u.CreatedAt.UnixNano())
	return nil
}

func newZapLogger(lvl zapcore.Level) logos.Logger {
	log := logos.New("test")
	//log.SetLLevel(lvl)
	return log

}

func fakeFields() []logos.Field {
	return []logos.Field{
		zap.Int("int", _tenInts[0]),
		zap.Ints("ints", _tenInts),
		zap.String("string", _tenStrings[0]),
		zap.Strings("strings", _tenStrings),
		zap.Time("time", _tenTimes[0]),
		zap.Times("times", _tenTimes),
		zap.Object("user1", _oneUser),
		zap.Object("user2", _oneUser),
		zap.Array("users", _tenUsers),
		zap.Error(errExample),
	}
}

func fakeSugarFields() []interface{} {
	return []interface{}{
		"int", _tenInts[0],
		"ints", _tenInts,
		"string", _tenStrings[0],
		"strings", _tenStrings,
		"time", _tenTimes[0],
		"times", _tenTimes,
		"user1", _oneUser,
		"user2", _oneUser,
		"users", _tenUsers,
		"error", errExample,
	}
}

func fakeFmtArgs() []interface{} {
	// Need to keep this a function instead of a package-global var so that we
	// pay the cast-to-interface{} penalty on each call.
	return []interface{}{
		_tenInts[0],
		_tenInts,
		_tenStrings[0],
		_tenStrings,
		_tenTimes[0],
		_tenTimes,
		_oneUser,
		_oneUser,
		_tenUsers,
		errExample,
	}
}
