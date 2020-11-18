package httplog

import (
	"net/http"
	"strings"
	"time"

	"github.com/apex/log"
)

// New middleware wrapping `h`.
func New(h http.Handler, opts ...LoggerOption) *Logger {
	logger := &Logger{Handler: h}
	for _, opt := range opts {
		opt(logger)
	}
	return logger
}

// LoggerOption configures a Logger.
type LoggerOption func(*Logger)

// Logger middleware wrapping Handler.
type Logger struct {
	http.Handler
	excludePathPrefixes []string
}

// ExcludePathPrefix configures the logger to not log requests with the given path prefix
func ExcludePathPrefix(pathPrefix string) LoggerOption {
	return func(l *Logger) {
		l.excludePathPrefixes = append(l.excludePathPrefixes, pathPrefix)
	}
}

// ServeHTTP implementation.
func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, prefix := range l.excludePathPrefixes {
		if strings.HasPrefix(r.URL.Path, prefix) {
			return
		}
	}

	start := time.Now()

	res := makeLoggingResponseWriter(w)

	ctx := log.
		FromContext(r.Context()).
		WithFields(log.Fields{
			"url":        r.RequestURI,
			"method":     r.Method,
			"remoteAddr": r.RemoteAddr,
		})

	ctx.Info("request")
	l.Handler.ServeHTTP(res, r)

	ctx = ctx.WithFields(log.Fields{
		"status":   res.Status(),
		"size":     res.Size(),
		"duration": ms(time.Since(start)),
	})

	switch {
	case res.Status() >= 500:
		ctx.Error("response")
	case res.Status() >= 400:
		ctx.Warn("response")
	default:
		ctx.Info("response")
	}
}

// ms returns the duration in milliseconds.
func ms(d time.Duration) int {
	return int(d / time.Millisecond)
}
