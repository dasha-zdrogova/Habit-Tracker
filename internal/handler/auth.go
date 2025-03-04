package handler

import (
	"context"
	"encoding/json"
	"errors"
	"habit-tracker/internal/repository"
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

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	var user registerRequest
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if len(user.Username) < 2 {
		http.Error(w, "username too short", http.StatusBadRequest)
		return
	}

	if len(user.Username) > 50 {
		http.Error(w, "username too long", http.StatusBadRequest)
		return
	}

	if len(user.Password) < 6 {
		http.Error(w, "password too short", http.StatusBadRequest)
		return
	}

	if err := h.services.Users.Register(user.Username, user.Password); err != nil {
		if err == repository.ErrUserExists {
			http.Error(w, "username already taken", http.StatusConflict)
			return
		}
		http.Error(w, "internal service error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var user loginRequest
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// TODO: проверять на поля на пустоту, если что bad request
	userID, err := h.services.Users.Login(user.Username, user.Password)
	if err != nil {
		if err == repository.ErrUserNotFound {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		http.Error(w, "internal service error", http.StatusInternalServerError)
		return
	}

	token := h.tokenManager.GenerateToken()
	h.tokenManager.AddToken(token, userID)

	json.NewEncoder(w).Encode(authResponse{
		Token: token,
	})
}

func (*Handler) getUserIDFromContext(r *http.Request) (int, error) {
	userID, ok := r.Context().Value(userIDKey{}).(int)
	if !ok {
		return -1, errors.New("no user id in context")
	}
	return userID, nil
}
