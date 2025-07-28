package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResolveConfigPath_Default(t *testing.T) {
	path := resolveConfigPath("/cwd")
	require.Equal(t, "/cwd/config.yaml", path)
}

func TestResolveConfigPath_Flag(t *testing.T) {
	*configFlag = "/flag/config.yaml"
	path := resolveConfigPath("/cwd")
	require.Equal(t, "/flag/config.yaml", path)
	*configFlag = ""
}

func TestResolveConfigPath_Env(t *testing.T) {
	t.Setenv("TWITCHETS_CONFIG", "/env/config.yaml")
	path := resolveConfigPath("/cwd")
	require.Equal(t, "/env/config.yaml", path)
}

func TestResolveConfigPath_FlagOverridesEnv(t *testing.T) {
	*configFlag = "/flag/config.yaml"
	t.Setenv("TWITCHETS_CONFIG", "/env/config.yaml")
	path := resolveConfigPath("/cwd")
	require.Equal(t, "/flag/config.yaml", path)
	*configFlag = ""
}
