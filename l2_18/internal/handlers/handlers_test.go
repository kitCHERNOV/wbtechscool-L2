package handlers

import (
	"bytes"
	"calendar/internal/logger"
	"calendar/internal/models"
	"calendar/internal/storage"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func setupTestHandler() *Handler {
	logger, _ := logger.NewLogger()
	eventStorage := storage.NewEventStorage(logger)
	return NewHandler(eventStorage, logger)
}

func TestCreateEventHandler(t *testing.T) {
	handler := setupTestHandler()
	
	event := models.Event{
		ID:          1,
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(time.Hour),
		Title:       "Test Event",
		Description: "Test Description",
	}
	
	body, _ := json.Marshal(event)
	req := httptest.NewRequest(http.MethodPost, "/create_event", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	
	handler.CreateEventHandler(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestCreateEventHandlerWrongMethod(t *testing.T) {
	handler := setupTestHandler()
	
	req := httptest.NewRequest(http.MethodGet, "/create_event", nil)
	w := httptest.NewRecorder()
	
	handler.CreateEventHandler(w, req)
	
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}

func TestCreateEventHandlerInvalidJSON(t *testing.T) {
	handler := setupTestHandler()
	
	invalidJSON := []byte(`{"invalid json`)
	req := httptest.NewRequest(http.MethodPost, "/create_event", bytes.NewBuffer(invalidJSON))
	w := httptest.NewRecorder()
	
	handler.CreateEventHandler(w, req)
	
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestUpdateEventHandler(t *testing.T) {
	handler := setupTestHandler()
	
	// Сначала создаем событие
	event := models.Event{
		ID:        1,
		StartTime: time.Now(),
		EndTime:   time.Now().Add(time.Hour),
		Title:     "Original Title",
	}
	handler.storage.AddEvent(event)
	
	// Обновляем событие
	updatedEvent := models.Event{
		ID:        1,
		StartTime: time.Now(),
		EndTime:   time.Now().Add(2 * time.Hour),
		Title:     "Updated Title",
	}
	
	body, _ := json.Marshal(updatedEvent)
	req := httptest.NewRequest(http.MethodPost, "/update_event", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	
	handler.UpdateEventHandler(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestDeleteEventHandler(t *testing.T) {
	handler := setupTestHandler()
	
	// Создаем событие
	event := models.Event{
		ID:        1,
		StartTime: time.Now(),
		EndTime:   time.Now().Add(time.Hour),
		Title:     "Event to Delete",
	}
	handler.storage.AddEvent(event)
	
	// Удаляем событие
	body, _ := json.Marshal(event)
	req := httptest.NewRequest(http.MethodPost, "/delete_event", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	
	handler.DeleteEventHandler(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestEventsForDayHandler(t *testing.T) {
	handler := setupTestHandler()
	
	// Добавляем событие
	targetDate := time.Date(2025, 11, 14, 10, 0, 0, 0, time.UTC)
	event := models.Event{
		ID:        1,
		StartTime: targetDate,
		EndTime:   targetDate.Add(time.Hour),
		Title:     "Test Event",
	}
	handler.storage.AddEvent(event)
	
	// Формируем запрос с query параметром
	req := httptest.NewRequest(http.MethodGet, "/events_for_day?date=2025-11-14", nil)
	w := httptest.NewRecorder()
	
	handler.EventsForDayHandler(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestEventsForDayHandlerInvalidDate(t *testing.T) {
	handler := setupTestHandler()
	
	req := httptest.NewRequest(http.MethodGet, "/events_for_day?date=invalid-date", nil)
	w := httptest.NewRecorder()
	
	handler.EventsForDayHandler(w, req)
	
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestEventsForWeekHandler(t *testing.T) {
	handler := setupTestHandler()
	
	// Добавляем событие на неделю
	targetDate := time.Date(2025, 11, 14, 10, 0, 0, 0, time.UTC)
	event := models.Event{
		ID:        1,
		StartTime: targetDate,
		EndTime:   targetDate.Add(time.Hour),
		Title:     "Week Event",
	}
	handler.storage.AddEvent(event)
	
	req := httptest.NewRequest(http.MethodGet, "/events_for_week?date=2025-11-14", nil)
	w := httptest.NewRecorder()
	
	handler.EventsForWeekHandler(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestEventsForMonthHandler(t *testing.T) {
	handler := setupTestHandler()
	
	// Добавляем событие на месяц
	targetDate := time.Date(2025, 11, 14, 10, 0, 0, 0, time.UTC)
	event := models.Event{
		ID:        1,
		StartTime: targetDate,
		EndTime:   targetDate.Add(time.Hour),
		Title:     "Month Event",
	}
	handler.storage.AddEvent(event)
	
	req := httptest.NewRequest(http.MethodGet, "/events_for_month?date=2025-11-14", nil)
	w := httptest.NewRecorder()
	
	handler.EventsForMonthHandler(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}
