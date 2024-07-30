package main

import (
	"errors"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	c := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec

	for key, val := range env {
		if val.NeedRemove {
			os.Unsetenv(key)
			continue
		}
		os.Setenv(key, val.Value)
	}
	c.Env = os.Environ()

	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if err := c.Run(); err != nil {
		var e *exec.ExitError

		if errors.As(err, &e) {
			return e.ExitCode()
		}
	}

	return 0
}
