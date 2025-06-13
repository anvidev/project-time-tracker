package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/anvidev/project-time-tracker/internal/contextkeys"
	"github.com/anvidev/project-time-tracker/internal/store/sessions"
)

func (api *api) bearerAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			api.unauthorizedError(w, r, fmt.Errorf("missing authorization header"))
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 && parts[0] != "Bearer" {
			api.unauthorizedError(w, r, fmt.Errorf("invalid authorization header"))
			return
		}

		ctx := r.Context()
		token := parts[1]

		session, err := api.store.Sessions.Validate(ctx, token)
		if err != nil {
			switch err {
			case sessions.ErrSessionExpired:
				api.unauthorizedError(w, r, fmt.Errorf("access expired"))
			default:
				api.unauthorizedError(w, r, fmt.Errorf("access denied"))
			}
			return
		}

		ctx = context.WithValue(ctx, contextkeys.SessionToken, session.Token)
		ctx = context.WithValue(ctx, contextkeys.UserId, session.UserId)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserId(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(contextkeys.UserId).(int64)
	return userID, ok
}

func getSessionToken(ctx context.Context) (int64, bool) {
	token, ok := ctx.Value(contextkeys.SessionToken).(int64)
	return token, ok
}
