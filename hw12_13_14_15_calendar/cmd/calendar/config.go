package main

import "github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/internal/helper"

type Config struct {
	Logger  LoggerConf
	DB      DBConf
	Storage string
	Server  Server
	GRPC    GRPC
}

type LoggerConf struct {
	Level string
}

type DBConf struct {
	User     string
	Password string
	Name     string
	Host     string
	Port     string
}

type Server struct {
	Host string
	Port string
}

type GRPC struct {
	Host string
	Port string
}

func LoadConfig(path string) (Config, error) {
	config, err := helper.NewConfig[Config](path)
	if err != nil {
		return Config{}, err
	}

	return *config, nil
}
