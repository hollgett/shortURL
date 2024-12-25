package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWrite struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (l *loggingResponseWrite) Write(b []byte) (int, error) {
	size, err := l.ResponseWriter.Write(b)
	l.responseData.size += size
	return size, err
}

func (l *loggingResponseWrite) WriteHeader(statusCode int) {
	l.ResponseWriter.WriteHeader(statusCode)
	l.responseData.status = statusCode
}

var Log *zap.Logger = zap.NewNop()

func InitLogger() error {
	cfg := zap.NewDevelopmentConfig()

	log, err := cfg.Build()
	if err != nil {
		return err
	}
	Log = log
	return nil
}

func RequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		Log.Info(
			"Request start",
			zap.String("uri", r.RequestURI),
			zap.String("method", r.Method),
		)

		next.ServeHTTP(w, r)
		duration := time.Since(start)
		Log.Info(
			"Request complete",
			zap.String("uri", r.RequestURI),
			zap.String("method", r.Method),
			zap.Duration("duration", duration),
		)
	})
}

func ResponseMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		respData := responseData{}

		logResp := loggingResponseWrite{
			ResponseWriter: w,
			responseData:   &respData,
		}

		next.ServeHTTP(&logResp, r)

		Log.Info(
			"Response complete",
			zap.Int("status", respData.status),
			zap.Int("size", respData.size),
		)
	}

	return http.HandlerFunc(fn)
}
