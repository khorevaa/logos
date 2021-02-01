package logos

import "errors"

var (
	ErrEnvConfigNotSet = errors.New("environment variable 'LOGOS_CONFIG' is not set")
)
