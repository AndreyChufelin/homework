package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("Correct data", func(t *testing.T) {
		env, err := ReadDir("./testdata/env/")

		require.Equal(t, EnvValue{"\"hello\"", false}, env["HELLO"])
		require.Equal(t, EnvValue{"bar", false}, env["BAR"], "Only first line")
		require.Equal(t, EnvValue{"", false}, env["EMPTY"], "Trim")
		require.Equal(t, EnvValue{"   foo\nwith new line", false}, env["FOO"], "Replace 0x00 to \n")
		require.Equal(t, EnvValue{"", true}, env["UNSET"], "Remove when empty file")

		require.NotContains(t, env, "WITH=")
		require.Equal(t, EnvValue{"with =", false}, env["WITH"], "Remove '=' sign from name")

		require.NoError(t, err)
	})
	t.Run("Not existing directory", func(t *testing.T) {
		_, err := ReadDir("./testdata/not")

		require.ErrorIs(t, err, ErrInvalidDir)
	})
}
