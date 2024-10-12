//go:generate protoc --proto_path=../../api --go_out=../../api --go_opt=paths=source_relative --go-grpc_out=../../api --go-grpc_opt=paths=source_relative ../../api/EventService.proto

package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/internal/app"
	"github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/internal/helper"
	loggerslog "github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/internal/logger/slog"
	internalgrpc "github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/internal/server/http"
	_ "github.com/lib/pq"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := LoadConfig(configFile)
	if err != nil {
		log.Fatalf("failed to read config from %q: %v", configFile, err)
	}

	logg, err := loggerslog.New(os.Stderr, config.Logger.Level)
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	storage, closeStorage, err := helper.InitStorage(ctx, helper.DBConfig{
		User:     config.DB.User,
		Password: config.DB.Password,
		Name:     config.DB.Name,
		Host:     config.DB.Host,
		Port:     config.DB.Port,
	}, config.Storage)
	if err != nil {
		logg.Error("failed to run database", "err", err)
		cancel()
	}
	defer closeStorage()

	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(logg, calendar, config.Server.Host, config.Server.Port)
	grpc := internalgrpc.NewServer(logg, calendar, config.GRPC.Host, config.GRPC.Port)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
		grpc.Stop()
	}()

	logg.Info("calendar is running...")

	go func() {
		if err := grpc.Start(); err != nil {
			logg.Error("failed to start grpc server: " + err.Error())
			cancel()
			os.Exit(1)
		}
	}()

	if err := server.Start(); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
