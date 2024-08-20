package logger

import (
	"fmt"
	"io"
	"log/slog"
)

type Logger struct {
	slog *slog.Logger
}

func New(w io.Writer, level string) (*Logger, error) {
	var lvl slog.Level
	err := lvl.UnmarshalText([]byte(level))
	if err != nil {
		return nil, fmt.Errorf("wrong level value %q: %w", level, err)
	}
	log := slog.New(slog.NewTextHandler(w, &slog.HandlerOptions{Level: lvl}))

	return &Logger{
		slog: log,
	}, nil
}

func (l Logger) Info(msg string, keysAndValues ...interface{}) {
	l.slog.Info(msg, keysAndValues...)
}

func (l Logger) Error(msg string, keysAndValues ...interface{}) {
	l.slog.Error(msg, keysAndValues...)
}

// TODO
