package main

import (
	"fmt"
	"io"
	"os"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Logger  LoggerConf
	DB      DBConf
	Storage string
	Server  Server
}

type LoggerConf struct {
	Level string
}

type DBConf struct {
	User     string
	Password string
	Name     string
}

type Server struct {
	Host string
	Port string
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
