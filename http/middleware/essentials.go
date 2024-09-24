package middleware

import (
	"context"
	"net/http"
)

type Essentials struct {
	Objects map[string]any
	Next    http.Handler
}

func (mw *Essentials) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for key, val := range mw.Objects {
		r = r.WithContext(context.WithValue(r.Context(), key, val))
	}
	mw.Next.ServeHTTP(w, r)
}
