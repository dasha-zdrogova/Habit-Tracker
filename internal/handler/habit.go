package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type createHabitRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (h *Handler) createHabit(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserIDFromContext(r)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	var habit createHabitRequest
	if err := json.NewDecoder(r.Body).Decode(&habit); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if len(habit.Name) < 1 {
		http.Error(w, "The name of the habit should not be empty", http.StatusBadRequest)
		return
	}

	if len(habit.Name) > 50 {
		http.Error(w, "the name of the habit too long", http.StatusBadRequest)
		return
	}

	if len(habit.Description) > 256 {
		http.Error(w, "the description of the habit too long", http.StatusBadRequest)
		return
	}

	err = h.services.Habits.Create(userID, habit.Name, habit.Description)
	if err != nil {
		http.Error(w, "failed to create habit", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) getHabitInfo(w http.ResponseWriter, r *http.Request) {
	ID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid ID format", http.StatusBadRequest)
		return
	}
	habitLogs, err := h.services.Habits.GetInfo(ID)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(habitLogs)
}

func (h *Handler) markHabit(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) deleteHabit(w http.ResponseWriter, r *http.Request) {

}
