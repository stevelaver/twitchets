package testutils

import (
	"os"
	"strconv"
	"testing"
)

func SkipIfCI(t *testing.T, args ...any) {
	env, ok := os.LookupEnv("CI")
	if !ok {
		return
	}

	isCI, err := strconv.ParseBool(env)
	if err != nil {
		return
	}

	if isCI {
		t.Skip(args...)
	}
}
