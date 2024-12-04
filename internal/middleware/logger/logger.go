package logger

import (
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"
)

var Log *zap.Logger = zap.NewNop()

type (
	responseData struct {
		statusCode int
		size       int
	}

	LoggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func Initialize(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}

	cfg := zap.NewProductionConfig()

	cfg.Level = lvl

	zl, err := cfg.Build()
	if err != nil {
		return err
	}

	Log = zl
	return nil
}

func (lrw *LoggingResponseWriter) Write(b []byte) (int, error) {
	size, err := lrw.ResponseWriter.Write(b)
	lrw.responseData.size += size
	return size, err
}

func (lrw *LoggingResponseWriter) WriteHeader(statusCode int) {
	lrw.ResponseWriter.WriteHeader(statusCode)
	lrw.responseData.statusCode = statusCode
}

func RequestLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		method := r.Method
		uri := r.RequestURI
		resData := &responseData{
			size:       0,
			statusCode: 0,
		}

		lwr := LoggingResponseWriter{
			ResponseWriter: w,
			responseData:   resData,
		}

		h.ServeHTTP(&lwr, r)

		duration := time.Since(start)
		Log.Info("request",
			zap.String("method", method),
			zap.String("uri", uri),
			zap.String("duration", duration.String()),
		)
		Log.Info("response",
			zap.String("statusCode", strconv.Itoa(resData.statusCode)),
			zap.String("response size", strconv.Itoa(resData.size)),
		)
	})
}
