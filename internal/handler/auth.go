package handler

import (
	"context"
	"errors"
	"net/http"
)

type userIDKey struct{}

type registerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type authResponse struct {
	Token string `json:"token"`
}

func (h *Handler) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "missing auth token", http.StatusUnauthorized)
			return
		}

		userID, err := h.tokenManager.ValidateToken(token)
		if err != nil {
			http.Error(w, "invalid or expired auth token", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), userIDKey{}, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (*Handler) getUserIDFromContext(r *http.Request) (int, error) {
	userID, ok := r.Context().Value(userIDKey{}).(int)
	if !ok {
		return -1, errors.New("no user id in context")
	}
	return userID, nil
}
