package handlers

import (
	"calendar/internal/logger"
	"calendar/internal/models"
	"calendar/internal/storage"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type Handler struct {
	storage *storage.EventStorage
	logger  *logger.Logger
}

func NewHandler(storage *storage.EventStorage, logger *logger.Logger) *Handler {
	return &Handler{
		storage: storage,
		logger:  logger,
	}
}

// constant answer messages
const (
	notAllowedMethod = "Method not allowed"
)

// Заглушки обработчиков (реализуй сам)
func (h *Handler) CreateEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, notAllowedMethod, http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Incorrect data", http.StatusBadRequest)
		return
	}

	var event models.Event
	err = json.Unmarshal(body, &event)
	if err != nil {
		http.Error(w, "messege struct is not correct", http.StatusBadRequest)
		return
	}

	// call adding event func
	err = h.storage.AddEvent(event)
	if err != nil {
		h.logger.Error(err.Error())
		http.Error(w, "server problem", http.StatusServiceUnavailable)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"status": "event is created",
	})
}

func (h *Handler) UpdateEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Incorrect data", http.StatusBadRequest)
		return
	}

	var event models.Event
	err = json.Unmarshal(body, &event)
	if err != nil {
		http.Error(w, "messege struct is not correct", http.StatusBadRequest)
		return
	}

	// call updating event func
	err = h.storage.UpdateEvent(event)
	if err != nil {
		h.logger.Error(err.Error())
		http.Error(w, "server problem", http.StatusServiceUnavailable)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

func (h *Handler) DeleteEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Incorrect data", http.StatusBadRequest)
		return
	}

	var event models.Event
	err = json.Unmarshal(body, &event)
	if err != nil {
		http.Error(w, "messege struct is not correct", http.StatusBadRequest)
		return
	}

	err = h.storage.DeleteEvent(event.ID)
	if err != nil {
		h.logger.Error(err.Error())
		http.Error(w, "server problem", http.StatusServiceUnavailable)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (h *Handler) EventsForDayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// body, err := io.ReadAll(r.Body)
	// if err != nil {
	//     http.Error(w, "Incorrect data", http.StatusBadRequest)
	//     return
	// }
	dateStr := r.URL.Query().Get("date")

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "invalid date format, use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	events := h.storage.GetEventsForDay(date)

	respondJSON(w, http.StatusOK, events)
}

func (h *Handler) EventsForWeekHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	dateStr := r.URL.Query().Get("date")

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "invalid date format, use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	events := h.storage.GetEventsForWeek(date)
	respondJSON(w, http.StatusOK, events)
}

func (h *Handler) EventsForMonthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	dateStr := r.URL.Query().Get("date")

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "invalid date format, use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	events := h.storage.GetEventsForMonth(date)
	respondJSON(w, http.StatusOK, events)
}

// Вспомогательная функция для JSON ответов
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
