package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"calendar/internal/domain"
	"calendar/internal/service"
)

type Handler struct {
	svc *service.CalendarService
}

func NewHandler(svc *service.CalendarService) *Handler {
	return &Handler{svc: svc}
}

func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func errorResponse(w http.ResponseWriter, status int, msg string) {
	jsonResponse(w, status, map[string]string{"error": msg})
}

func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorResponse(w, http.StatusMethodNotAllowed, "Только POST")
		return
	}

	var e domain.Event
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		errorResponse(w, http.StatusBadRequest, "неверный формат JSON")
		return
	}

	if err := h.svc.CreateEvent(e); err != nil {
		errorResponse(w, http.StatusServiceUnavailable, err.Error())
		return
	}
	jsonResponse(w, http.StatusOK, map[string]string{"result": "событие создано"})
}

func (h *Handler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	var e domain.Event
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		errorResponse(w, http.StatusBadRequest, "неверный формат JSON")
		return
	}

	if err := h.svc.UpdateEvent(e); err != nil {
		errorResponse(w, http.StatusServiceUnavailable, err.Error())
		return
	}
	jsonResponse(w, http.StatusOK, map[string]string{"result": "событие обновлено"})
}

func (h *Handler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	var req map[string]string
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, "неверный формат JSON")
		return
	}

	id := req["id"]
	if err := h.svc.DeleteEvent(id); err != nil {
		errorResponse(w, http.StatusServiceUnavailable, err.Error())
		return
	}
	jsonResponse(w, http.StatusOK, map[string]string{"result": "событие удалено"})
}

func (h *Handler) getEvents(w http.ResponseWriter, r *http.Request, period string) {
	if r.Method != http.MethodGet {
		errorResponse(w, http.StatusMethodNotAllowed, "Только GET")
		return
	}

	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "некорректный или отсутствующий user_id")
		return
	}

	events, err := h.svc.GetEvents(userID, period)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if events == nil {
		events = []domain.Event{}
	}
	jsonResponse(w, http.StatusOK, map[string]interface{}{"result": events})
}

func (h *Handler) GetEventsDay(w http.ResponseWriter, r *http.Request)   { h.getEvents(w, r, "day") }
func (h *Handler) GetEventsWeek(w http.ResponseWriter, r *http.Request)  { h.getEvents(w, r, "week") }
func (h *Handler) GetEventsMonth(w http.ResponseWriter, r *http.Request) { h.getEvents(w, r, "month") }
