package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/internal/helper"
	loggerslog "github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/internal/logger/slog"
	"github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/internal/queue"
	"github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/internal/sender"
	_ "github.com/lib/pq"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	config, err := LoadConfig(configFile)
	if err != nil {
		log.Fatalf("failed to read config from %q: %v", configFile, err)
	}

	logg, err := loggerslog.New(os.Stderr, "INFO")
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	q := queue.NewQueue(config.Queue.User, config.Queue.Password, config.Queue.Host, config.Queue.Port)
	err = q.Start()
	if err != nil {
		logg.Error("failed start queue", "err", err)
		cancel()
	}
	defer q.Stop()

	storage, closeStorage, err := helper.InitStorage(ctx, helper.DBConfig{
		User:     config.DB.User,
		Password: config.DB.Password,
		Name:     config.DB.Name,
		Host:     config.DB.Host,
		Port:     config.DB.Port,
	}, "sql")
	if err != nil {
		logg.Error("failed to run database", "err", err)
		cancel()
	}
	defer closeStorage()

	consumer := queue.NewConsumer("notification_queue", q.Conn)
	err = consumer.Start()
	if err != nil {
		logg.Error("failed start consumer", "err", err)
		cancel()
	}
	defer consumer.Stop()

	go func() {
		<-ctx.Done()
		consumer.Stop()
	}()

	s := sender.NewSender(consumer, logg, storage)
	err = s.Start()
	if err != nil {
		logg.Error("sender failed", "err", err)
	}
}
