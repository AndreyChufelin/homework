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
	"github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/internal/storage/memory"
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

	logFile, err := os.OpenFile("./logs/calendar.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer logFile.Close()

	logg, err := logger.New(logFile, config.Logger.Level)
	if err != nil {
		log.Fatalf("failed to create logger: %v", err) //nolint:gocritic
	}

	storage := memorystorage.New()
	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(logg, calendar)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
