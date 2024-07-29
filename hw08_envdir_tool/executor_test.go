package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("Sets env", func(t *testing.T) {
		env := Environment{
			"BAR": EnvValue{
				Value:      "bar",
				NeedRemove: false,
			},
			"UNSET": EnvValue{
				Value:      "",
				NeedRemove: true,
			},
			"ADDED": EnvValue{
				Value:      "from original env",
				NeedRemove: false,
			},
		}

		os.Setenv("BAR", "foo")
		code := RunCmd([]string{"echo"}, env)

		require.Equal(t, "bar", os.Getenv("BAR"), "Sets variable")
		require.Equal(t, "from original env", os.Getenv("ADDED"), "Doesn't reset original env")

		_, u := os.LookupEnv("UNSET")
		require.False(t, u, "Removes variable")

		require.Equal(t, 0, code)
	})
	t.Run("Returns exit code", func(t *testing.T) {
		env := Environment{}
		code := RunCmd([]string{"false"}, env)

		require.Equal(t, 1, code)
	})
}
