package middleware

import (
	"context"
	"net/http"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type ctxKeyLogger int

// LoggerKey is the key that holds the zap logger-presets.
const LoggerKey ctxKeyLogger = 0

type ctxKeyRequestID int

// RequestIDKey is the key that holds the unique request ID in a request context.
const RequestIDKey ctxKeyRequestID = 0

// RequestIDHeader is the name of the HTTP Header which contains the request id
const RequestIDHeader = "X-Request-ID"

// LoggerFromContext returns the request logger stored by Logger.
func LoggerFromContext(ctx context.Context) (*zap.Logger, bool) {
	logger, ok := ctx.Value(LoggerKey).(*zap.Logger)
	return logger, ok
}

func Logger(logger *zap.Logger, service string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			requestID := r.Header.Get(RequestIDHeader)
			if requestID == "" {
				var err error
				requestID, err = gonanoid.New()
				if err != nil {
					logger.Fatal("Failed to generate request ID", zap.Error(err))
				}
			}

			logger := logger.
				With(zap.String("requestId", requestID)).
				With(zap.String("service", service))

			ctx = context.WithValue(ctx, RequestIDKey, requestID)
			ctx = context.WithValue(ctx, LoggerKey, logger)

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			ww.Header().Add(RequestIDHeader, requestID)

			t1 := time.Now()
			defer func() {
				logger.Info("HTTP request",
					zap.Dict("request",
						zap.String("url", r.URL.String()),
						zap.String("method", r.Method),
						zap.String("proto", r.Proto),
						zap.String("userAgent", r.UserAgent())),
					zap.Dict("response",
						zap.Int("status", ww.Status()),
						zap.Int("contentLength", ww.BytesWritten()),
						zap.Duration("elapsed", time.Since(t1))),
				)
			}()

			next.ServeHTTP(ww, r.WithContext(ctx))
		})
	}
}
