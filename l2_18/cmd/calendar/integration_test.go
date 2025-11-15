package main

import (
	"bytes"
	"calendar/internal/handlers"
	"calendar/internal/logger"
	"calendar/internal/models"
	"calendar/internal/storage"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

// setupSimpleTestServer создает тестовый сервер БЕЗ middleware
// Это упрощенная версия для тестирования только handlers
func setupSimpleTestServer(t *testing.T) (*handlers.Handler, *logger.Logger) {
	// Создаем директорию для логов если её нет
	err := os.MkdirAll("../../internal/logs", 0755)
	if err != nil {
		t.Logf("Warning: Failed to create logs directory: %v", err)
	}

	// Создаем logger
	logger, err := logger.NewLogger()
	if err != nil {
		t.Logf("Warning: Failed to create logger: %v, using nil logger", err)
		// Продолжаем с nil logger для демонстрации
	}

	eventStorage := storage.NewEventStorage(logger)
	handler := handlers.NewHandler(eventStorage, logger)

	return handler, logger
}

// TestSimple_CreateAndRetrieveEvent простой тест создания и получения события
func TestSimple_CreateAndRetrieveEvent(t *testing.T) {
	handler, logger := setupSimpleTestServer(t)
	if logger != nil {
		defer logger.Close()
	}

	// Создаем событие
	event := models.Event{
		ID:          1,
		StartTime:   time.Date(2025, 11, 14, 10, 0, 0, 0, time.UTC),
		EndTime:     time.Date(2025, 11, 14, 11, 0, 0, 0, time.UTC),
		Title:       "Integration Test Event",
		Description: "Testing full cycle",
	}

	// Отправляем POST запрос на создание
	body, _ := json.Marshal(event)
	req := httptest.NewRequest(http.MethodPost, "/create_event", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.CreateEventHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
		t.Logf("Response body: %s", w.Body.String())
	}

	// Получаем события на день
	req = httptest.NewRequest(http.MethodGet, "/events_for_day?date=2025-11-14", nil)
	w = httptest.NewRecorder()

	handler.EventsForDayHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
		t.Logf("Response body: %s", w.Body.String())
		return
	}

	var events []models.Event
	if err := json.NewDecoder(w.Body).Decode(&events); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(events))
	}

	if len(events) > 0 && events[0].Title != "Integration Test Event" {
		t.Errorf("Expected title 'Integration Test Event', got '%s'", events[0].Title)
	}
}

// TestSimple_UpdateEvent тест обновления события
func TestSimple_UpdateEvent(t *testing.T) {
	handler, logger := setupSimpleTestServer(t)
	if logger != nil {
		defer logger.Close()
	}

	// Создаем событие
	event := models.Event{
		ID:        2,
		StartTime: time.Date(2025, 11, 14, 14, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2025, 11, 14, 15, 0, 0, 0, time.UTC),
		Title:     "Original Title",
	}

	body, _ := json.Marshal(event)
	req := httptest.NewRequest(http.MethodPost, "/create_event", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	handler.CreateEventHandler(w, req)

	// Обновляем событие
	updatedEvent := models.Event{
		ID:        2,
		StartTime: time.Date(2025, 11, 14, 14, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2025, 11, 14, 16, 0, 0, 0, time.UTC),
		Title:     "Updated Title",
	}

	body, _ = json.Marshal(updatedEvent)
	req = httptest.NewRequest(http.MethodPost, "/update_event", bytes.NewBuffer(body))
	w = httptest.NewRecorder()

	handler.UpdateEventHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
		t.Logf("Response body: %s", w.Body.String())
	}

	// Проверяем что событие обновилось
	req = httptest.NewRequest(http.MethodGet, "/events_for_day?date=2025-11-14", nil)
	w = httptest.NewRecorder()
	handler.EventsForDayHandler(w, req)

	var events []models.Event
	json.NewDecoder(w.Body).Decode(&events)

	found := false
	for _, e := range events {
		if e.ID == 2 && e.Title == "Updated Title" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Updated event not found")
	}
}

// TestSimple_DeleteEvent тест удаления события
func TestSimple_DeleteEvent(t *testing.T) {
	handler, logger := setupSimpleTestServer(t)
	if logger != nil {
		defer logger.Close()
	}

	// Создаем событие
	event := models.Event{
		ID:        3,
		StartTime: time.Date(2025, 11, 14, 16, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2025, 11, 14, 17, 0, 0, 0, time.UTC),
		Title:     "Event to Delete",
	}

	body, _ := json.Marshal(event)
	req := httptest.NewRequest(http.MethodPost, "/create_event", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	handler.CreateEventHandler(w, req)

	// Удаляем событие
	body, _ = json.Marshal(event)
	req = httptest.NewRequest(http.MethodPost, "/delete_event", bytes.NewBuffer(body))
	w = httptest.NewRecorder()

	handler.DeleteEventHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Проверяем что событие удалено
	req = httptest.NewRequest(http.MethodGet, "/events_for_day?date=2025-11-14", nil)
	w = httptest.NewRecorder()
	handler.EventsForDayHandler(w, req)

	var events []models.Event
	json.NewDecoder(w.Body).Decode(&events)

	for _, e := range events {
		if e.ID == 3 {
			t.Error("Event should be deleted but still exists")
		}
	}
}

// TestSimple_WeekEvents тест получения событий на неделю
func TestSimple_WeekEvents(t *testing.T) {
	handler, logger := setupSimpleTestServer(t)
	if logger != nil {
		defer logger.Close()
	}

	// Создаем события на разные дни недели
	monday := time.Date(2025, 11, 10, 10, 0, 0, 0, time.UTC)
	friday := time.Date(2025, 11, 14, 10, 0, 0, 0, time.UTC)
	nextMonday := time.Date(2025, 11, 17, 10, 0, 0, 0, time.UTC)

	events := []models.Event{
		{ID: 10, StartTime: monday, EndTime: monday.Add(time.Hour), Title: "Monday"},
		{ID: 11, StartTime: friday, EndTime: friday.Add(time.Hour), Title: "Friday"},
		{ID: 12, StartTime: nextMonday, EndTime: nextMonday.Add(time.Hour), Title: "Next Monday"},
	}

	for _, e := range events {
		body, _ := json.Marshal(e)
		req := httptest.NewRequest(http.MethodPost, "/create_event", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		handler.CreateEventHandler(w, req)
	}

	// Получаем события на неделю (с пятницы)
	req := httptest.NewRequest(http.MethodGet, "/events_for_week?date=2025-11-14", nil)
	w := httptest.NewRecorder()
	handler.EventsForWeekHandler(w, req)

	var weekEvents []models.Event
	json.NewDecoder(w.Body).Decode(&weekEvents)

	// Должно быть 2 события (Monday и Friday), но не Next Monday
	if len(weekEvents) != 2 {
		t.Errorf("Expected 2 events for the week, got %d", len(weekEvents))
	}
}

// TestSimple_MonthEvents тест получения событий на месяц
func TestSimple_MonthEvents(t *testing.T) {
	handler, logger := setupSimpleTestServer(t)
	if logger != nil {
		defer logger.Close()
	}

	// События в ноябре и декабре
	novEvent1 := models.Event{
		ID:        20,
		StartTime: time.Date(2025, 11, 5, 10, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2025, 11, 5, 11, 0, 0, 0, time.UTC),
		Title:     "November Event 1",
	}
	novEvent2 := models.Event{
		ID:        21,
		StartTime: time.Date(2025, 11, 25, 10, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2025, 11, 25, 11, 0, 0, 0, time.UTC),
		Title:     "November Event 2",
	}
	decEvent := models.Event{
		ID:        22,
		StartTime: time.Date(2025, 12, 5, 10, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2025, 12, 5, 11, 0, 0, 0, time.UTC),
		Title:     "December Event",
	}

	for _, e := range []models.Event{novEvent1, novEvent2, decEvent} {
		body, _ := json.Marshal(e)
		req := httptest.NewRequest(http.MethodPost, "/create_event", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		handler.CreateEventHandler(w, req)
	}

	// Получаем события на ноябрь
	req := httptest.NewRequest(http.MethodGet, "/events_for_month?date=2025-11-14", nil)
	w := httptest.NewRecorder()
	handler.EventsForMonthHandler(w, req)

	var monthEvents []models.Event
	json.NewDecoder(w.Body).Decode(&monthEvents)

	if len(monthEvents) != 2 {
		t.Errorf("Expected 2 events for November, got %d", len(monthEvents))
	}
}

// TestSimple_InvalidRequests тест обработки неверных запросов
func TestSimple_InvalidRequests(t *testing.T) {
	handler, logger := setupSimpleTestServer(t)
	if logger != nil {
		defer logger.Close()
	}

	tests := []struct {
		name           string
		method         string
		url            string
		body           string
		expectedStatus int
		handlerFunc    http.HandlerFunc
	}{
		{
			name:           "Wrong method for create",
			method:         "GET",
			url:            "/create_event",
			body:           "",
			expectedStatus: http.StatusMethodNotAllowed,
			handlerFunc:    handler.CreateEventHandler,
		},
		{
			name:           "Invalid JSON",
			method:         "POST",
			url:            "/create_event",
			body:           "invalid json",
			expectedStatus: http.StatusBadRequest,
			handlerFunc:    handler.CreateEventHandler,
		},
		{
			name:           "Invalid date format",
			method:         "GET",
			url:            "/events_for_day?date=invalid",
			body:           "",
			expectedStatus: http.StatusBadRequest,
			handlerFunc:    handler.EventsForDayHandler,
		},
		{
			name:           "Empty title",
			method:         "POST",
			url:            "/create_event",
			body:           `{"id":100,"start_time":"2025-11-14T10:00:00Z","end_time":"2025-11-14T11:00:00Z","title":""}`,
			expectedStatus: http.StatusServiceUnavailable,
			handlerFunc:    handler.CreateEventHandler,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.url, bytes.NewBufferString(tt.body))
			w := httptest.NewRecorder()

			tt.handlerFunc(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
				t.Logf("Response body: %s", w.Body.String())
			}
		})
	}
}
