package testutils

import (
	"os"
	"strconv"
	"testing"
)

func IsCI() bool {
	env, ok := os.LookupEnv("CI")
	if !ok {
		return false
	}

	isCI, err := strconv.ParseBool(env)
	if err != nil {
		return false
	}

	return isCI
}

func SkipIfCI(t *testing.T, args ...any) {
	if IsCI() {
		t.Skip(args...)
	}
}
