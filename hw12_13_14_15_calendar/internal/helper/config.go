package helper

import (
	"fmt"
	"io"
	"os"

	"github.com/pelletier/go-toml/v2"
)

func NewConfig[C interface{}](path string) (*C, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("NewConfig: failed opening file: %w", err)
	}
	defer file.Close()

	var config C

	b, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("NewConfig: failed reading file: %w", err)
	}

	err = toml.Unmarshal(b, &config)
	if err != nil {
		return nil, fmt.Errorf("NewConfig: failed unmarshal toml: %w", err)
	}

	return &config, nil
}
