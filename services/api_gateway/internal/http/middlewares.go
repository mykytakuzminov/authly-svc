package http

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	c "github.com/mykytakuzminov/ridely-svc/shared/context"
	log "github.com/mykytakuzminov/ridely-svc/shared/logging"
	authpb "github.com/mykytakuzminov/ridely-svc/shared/proto/auth"
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

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func AuthMiddleware(client authpb.AuthServiceClient, logger *zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			traceID := c.GetTraceID(r.Context())

			accessToken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
			if accessToken == "" {
				responseUnauthorized(w, "missing token")
				return
			}

			claims, err := client.Validate(r.Context(), &authpb.ValidateRequest{
				AccessToken: accessToken,
			})
			if err != nil {
				log.Declined(logger, traceID, "token validation", err)
				responseUnauthorized(w, "invalid token")
				return
			}

			userID, err := uuid.Parse(claims.UserId)
			if err != nil {
				log.Failed(logger, traceID, "user id parsing", err)
				responseUnauthorized(w, "invalid token")
				return
			}

			ctx := context.WithValue(r.Context(), c.UserIDKey, userID)
			ctx = context.WithValue(ctx, c.UserRoleKey, claims.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
