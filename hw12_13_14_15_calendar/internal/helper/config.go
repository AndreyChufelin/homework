package helper

import (
	"fmt"
	"io"
	"os"
	"regexp"

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

	tomlStr := string(b)
	re := regexp.MustCompile(`\$\{(\w+)\}`)

	tomlStr = re.ReplaceAllStringFunc(tomlStr, func(match string) string {
		varName := match[2 : len(match)-1]
		if value, exists := os.LookupEnv(varName); exists {
			return value
		}
		return match
	})

	err = toml.Unmarshal([]byte(tomlStr), &config)
	if err != nil {
		return nil, fmt.Errorf("NewConfig: failed unmarshal toml: %w", err)
	}

	return &config, nil
}
