package handler

import (
	"encoding/json"
	"errors"
	"habit-tracker/internal/repository"
	"net/http"
)

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
		if errors.Is(err, repository.ErrUserExists) {
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
		if errors.Is(err, repository.ErrUserNotFound) {
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

func (h *Handler) getHabits(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserIDFromContext(r)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	habits, err := h.services.Users.GetHabits(userID)
	if err != nil {
		if errors.Is(err, repository.ErrHabitsNotFound) {
			http.Error(w, "habits not found", http.StatusNotFound)
			return
		}
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(habits)
}
