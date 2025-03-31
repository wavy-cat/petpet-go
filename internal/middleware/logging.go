package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func RequestLogger(logger *zap.Logger, service string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestId := uuid.New()
			r = r.WithContext(context.WithValue(r.Context(), "logger", logger))
			r = r.WithContext(context.WithValue(r.Context(), "requestId", requestId))
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			ww.Header().Add("X-RequestID", requestId.String())

			t1 := time.Now()
			defer func() {
				logger.Info("HTTP request",
					zap.String("service", service),
					zap.Dict("request", zap.String("url", r.URL.String()), zap.String("method", r.Method), zap.String("proto", r.Proto), zap.String("requestId", requestId.String()), zap.String("userAgent", r.UserAgent())),
					zap.Dict("response", zap.Int("status", ww.Status()), zap.Int("contentLength", ww.BytesWritten()), zap.Duration("elapsed", time.Since(t1))),
				)
			}()

			next.ServeHTTP(ww, r)
		})
	}
}
