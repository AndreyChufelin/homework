package main

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

var ErrInvalidDir = errors.New("invalid directory")

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, ErrInvalidDir
	}

	env := make(Environment, len(files))
	for _, f := range files {
		fInfo, err := f.Info()
		if err != nil {
			return nil, err
		}

		name := fInfo.Name()
		name = strings.ReplaceAll(name, "=", "")

		if fInfo.Size() <= 0 {
			env[name] = EnvValue{"", true}
			continue
		}

		file, err := os.Open(filepath.Join(dir, fInfo.Name()))
		if err != nil {
			return nil, err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)

		scanner.Scan()
		v := scanner.Bytes()
		v = bytes.ReplaceAll(v, []byte("\x00"), []byte("\n"))
		value := strings.TrimRight(string(v), " \t\n")

		if err := scanner.Err(); err != nil {
			return nil, err
		}

		env[name] = EnvValue{value, false}
	}

	return env, nil
}
