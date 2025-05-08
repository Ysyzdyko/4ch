package http

import (
	"1337b04rd/internal/application/domain"
	ports "1337b04rd/internal/ports/left"
	"context"
	"log/slog"
	"net/http"
	"time"
)

// Определяем свой собственный тип для ключа контекста
type ContextKey string

const (
	// Используем наш новый тип для ключа
	SessionKey ContextKey = "session"
)

type Middleware func(http.Handler) http.Handler

func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

func WithSession(sessionService ports.SessionPort) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var session *domain.Session
			var err error

			cookie, err := r.Cookie("session_id")
			if err == nil && cookie.Value != "" {
				session, err = sessionService.GetSessionByID(context.Background(), cookie.Value)
			}

			if err != nil || session == nil || !session.IsActive {
				session, err = sessionService.CreateSession(context.Background())
				if err != nil {
					http.Error(w, "failed to create session", http.StatusInternalServerError)
					slog.Error("Failed to create session: " + err.Error())
					return
				}

				// Устанавливаем новую cookie
				http.SetCookie(w, &http.Cookie{
					Name:     "session_id",
					Value:    session.ID,
					Path:     "/",
					Expires:  time.Now().Add(7 * 24 * time.Hour),
					HttpOnly: true,
					SameSite: http.SameSiteLaxMode,
				})
			}

			// Используем наш новый тип contextKey
			ctx := context.WithValue(context.Background(), SessionKey, session)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
