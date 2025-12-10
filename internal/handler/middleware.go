package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/murtazapatel89100/RSS-Aggregartor/internal/auth"
	"github.com/murtazapatel89100/RSS-Aggregartor/internal/database"
)

type contextKey string

const userContextKey contextKey = "user"

func (config ApiConfig) MiddlewareAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apikey, err := auth.GetApiKey(r.Header)
		if err != nil {
			RespondWithError(w, 403, fmt.Sprintf("Authentication error: %v", err))
			return
		}

		user, err := config.DB.GetUserByApiKey(r.Context(), apikey)
		if err != nil {
			RespondWithError(w, 403, "Invalid API key")
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserFromContext(r *http.Request) (database.User, bool) {
	user, ok := r.Context().Value(userContextKey).(database.User)
	return user, ok
}
