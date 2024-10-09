package loggerslog

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	t.Run("info", func(t *testing.T) {
		var buffer bytes.Buffer
		l, err := New(&buffer, "Info")
		require.NoError(t, err)

		l.Info("test", slog.Any("key", "value"))
		s := strings.Split(buffer.String(), " ")

		require.Equal(t, "msg=test", s[2])
		require.Equal(t, "key=value\n", s[3])
	})
	t.Run("error", func(t *testing.T) {
		var buffer bytes.Buffer
		l, err := New(&buffer, "Info")
		require.NoError(t, err)

		l.Error("test", slog.Any("key", "value"))
		s := strings.Split(buffer.String(), " ")

		require.Equal(t, "msg=test", s[2])
		require.Equal(t, "key=value\n", s[3])
	})
}
