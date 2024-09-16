package middleware

import (
	"context"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Logging struct {
	Logger *zap.Logger
	Next   http.Handler
}

func (mw *Logging) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	r = r.WithContext(context.WithValue(r.Context(), "logger", mw.Logger))
	mw.Next.ServeHTTP(w, r)

	duration := time.Since(start)
	mw.Logger.Info("HTTP request",
		zap.String("path", r.URL.Path),
		zap.Duration("duration", duration),
	)
}
