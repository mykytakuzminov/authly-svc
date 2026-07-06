package http

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"

	c "github.com/mykytakuzminov/ridely-svc/shared/context"
)

func TimeoutMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func TraceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := uuid.New()

		w.Header().Set("X-Trace-ID", traceID.String())
		ctx := context.WithValue(r.Context(), c.TraceIDKey, traceID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
