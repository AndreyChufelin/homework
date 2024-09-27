package internalhttp

import (
	"net/http"
	"strings"
	"time"

	"github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/internal/logger"
)

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (sr *statusRecorder) WriteHeader(code int) {
	sr.statusCode = code
	sr.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(logger logger.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		addr := strings.Split(req.RemoteAddr, ":")
		if len(addr) > 0 {
			addr = addr[:len(addr)-1]
		}
		ip := strings.Join(addr, "")
		date := time.Now()
		sr := &statusRecorder{ResponseWriter: res, statusCode: http.StatusOK}

		next.ServeHTTP(sr, req)

		latency := time.Since(date)

		logger.Info("HTTP request handled",
			"ip", ip,
			"date", date.Format(time.RFC822Z),
			"method", req.Method,
			"path", req.URL.Path,
			"version", req.Proto,
			"statusCode", sr.statusCode,
			"latency", latency.Milliseconds(),
			"userAgent", req.UserAgent(),
		)
	})
}
