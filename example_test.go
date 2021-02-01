package logos_test

import (
	"github.com/khorevaa/logos"
)

func ExampleNew() {

	log := logos.New("github.com/v8platform/test")
	log.Info("Error")
	log.Sync()

	// Output:
	// 2021-01-28T22:10:37.059+0300	INFO	github.com/v8platform/test	Error	{"url": "url", "attempt": 3, "backoff": 1}

}
