package main

import (
	"fmt"
	"io"
	"os"

	"github.com/pelletier/go-toml/v2"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger LoggerConf
	// TODO
}

type LoggerConf struct {
	Level string
	// TODO
}

func LoadConfig(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return Config{}, fmt.Errorf("NewConfig: failed opening file: %w", err)
	}
	defer file.Close()

	var config Config

	b, err := io.ReadAll(file)
	if err != nil {
		return Config{}, fmt.Errorf("NewConfig: failed reading file: %w", err)
	}

	err = toml.Unmarshal(b, &config)
	if err != nil {
		return Config{}, fmt.Errorf("NewConfig: failed unmarshal toml: %w", err)
	}

	return config, nil
}

// TODO
