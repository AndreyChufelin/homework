package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	loggerslog "github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/internal/logger/slog"
	"github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/internal/queue"
	"github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/internal/sender"
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
	}
	defer q.Stop()

	consumer := queue.NewConsumer("notification_queue", q.Conn)
	err = consumer.Start()
	if err != nil {
		logg.Error("failed start consumer", "err", err)
	}
	defer consumer.Stop()

	go func() {
		<-ctx.Done()
		consumer.Stop()
	}()

	s := sender.NewSender(consumer, logg)
	err = s.Start()
	if err != nil {
		logg.Error("sender failed", "err", err)
	}
}
