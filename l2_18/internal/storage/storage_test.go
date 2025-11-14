package storage

import (
	"calendar/internal/logger"
	"calendar/internal/models"
	"testing"
	"time"
)

func TestAddEvent(t *testing.T) {
	logger, _ := logger.NewLogger()
	defer logger.Close()
	
	storage := NewEventStorage(logger)
	
	event := models.Event{
		ID:          1,
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(time.Hour),
		Title:       "Test Meeting",
		Description: "Test Description",
	}
	
	err := storage.AddEvent(event)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	// Проверяем, что событие добавлено
	if len(storage.events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(storage.events))
	}
}

func TestAddEventEmptyTitle(t *testing.T) {
	logger, _ := logger.NewLogger()
	defer logger.Close()
	
	storage := NewEventStorage(logger)
	
	event := models.Event{
		ID:        1,
		StartTime: time.Now(),
		EndTime:   time.Now().Add(time.Hour),
		Title:     "", // Пустой title
	}
	
	err := storage.AddEvent(event)
	if err == nil {
		t.Error("Expected error for empty title, got nil")
	}
}

func TestUpdateEvent(t *testing.T) {
	logger, _ := logger.NewLogger()
	defer logger.Close()
	
	storage := NewEventStorage(logger)
	
	// Добавляем событие
	event := models.Event{
		ID:          1,
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(time.Hour),
		Title:       "Original Title",
		Description: "Original Description",
	}
	storage.AddEvent(event)
	
	// Обновляем событие
	updatedEvent := models.Event{
		ID:          1,
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(2 * time.Hour),
		Title:       "Updated Title",
		Description: "Updated Description",
	}
	
	err := storage.UpdateEvent(updatedEvent)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	// Проверяем обновление
	if storage.events[1].Title != "Updated Title" {
		t.Errorf("Expected 'Updated Title', got %s", storage.events[1].Title)
	}
}

func TestUpdateEventEmptyTitle(t *testing.T) {
	logger, _ := logger.NewLogger()
	defer logger.Close()
	
	storage := NewEventStorage(logger)
	
	event := models.Event{
		ID:        1,
		StartTime: time.Now(),
		EndTime:   time.Now().Add(time.Hour),
		Title:     "",
	}
	
	err := storage.UpdateEvent(event)
	if err == nil {
		t.Error("Expected error for empty title, got nil")
	}
}

func TestDeleteEvent(t *testing.T) {
	logger, _ := logger.NewLogger()
	defer logger.Close()
	
	storage := NewEventStorage(logger)
	
	// Добавляем событие
	event := models.Event{
		ID:        1,
		StartTime: time.Now(),
		EndTime:   time.Now().Add(time.Hour),
		Title:     "Test Event",
	}
	storage.AddEvent(event)
	
	// Удаляем событие
	err := storage.DeleteEvent(1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	// Проверяем удаление
	if len(storage.events) != 0 {
		t.Errorf("Expected 0 events, got %d", len(storage.events))
	}
}

func TestGetEventsForDay(t *testing.T) {
	logger, _ := logger.NewLogger()
	defer logger.Close()
	
	storage := NewEventStorage(logger)
	
	// Создаем конкретную дату
	targetDate := time.Date(2025, 11, 14, 10, 0, 0, 0, time.UTC)
	otherDate := time.Date(2025, 11, 15, 10, 0, 0, 0, time.UTC)
	
	// Добавляем события на разные дни
	event1 := models.Event{
		ID:        1,
		StartTime: targetDate,
		EndTime:   targetDate.Add(time.Hour),
		Title:     "Event on 14th",
	}
	event2 := models.Event{
		ID:        2,
		StartTime: otherDate,
		EndTime:   otherDate.Add(time.Hour),
		Title:     "Event on 15th",
	}
	event3 := models.Event{
		ID:        3,
		StartTime: targetDate.Add(2 * time.Hour),
		EndTime:   targetDate.Add(3 * time.Hour),
		Title:     "Another event on 14th",
	}
	
	storage.AddEvent(event1)
	storage.AddEvent(event2)
	storage.AddEvent(event3)
	
	// Получаем события на 14 ноября
	events := storage.GetEventsForDay(targetDate)
	
	if len(events) != 2 {
		t.Errorf("Expected 2 events for the day, got %d", len(events))
	}
}

func TestGetEventsForWeek(t *testing.T) {
	logger, _ := logger.NewLogger()
	defer logger.Close()
	
	storage := NewEventStorage(logger)
	
	// Четверг 14 ноября 2025
	thursday := time.Date(2025, 11, 14, 10, 0, 0, 0, time.UTC)
	// Понедельник той же недели (10 ноября)
	monday := time.Date(2025, 11, 10, 10, 0, 0, 0, time.UTC)
	// Следующий понедельник (17 ноября) - уже не в этой неделе
	nextMonday := time.Date(2025, 11, 17, 10, 0, 0, 0, time.UTC)
	
	// События на текущей неделе
	event1 := models.Event{
		ID:        1,
		StartTime: monday,
		EndTime:   monday.Add(time.Hour),
		Title:     "Monday Event",
	}
	event2 := models.Event{
		ID:        2,
		StartTime: thursday,
		EndTime:   thursday.Add(time.Hour),
		Title:     "Thursday Event",
	}
	// Событие на следующей неделе
	event3 := models.Event{
		ID:        3,
		StartTime: nextMonday,
		EndTime:   nextMonday.Add(time.Hour),
		Title:     "Next Monday Event",
	}
	
	storage.AddEvent(event1)
	storage.AddEvent(event2)
	storage.AddEvent(event3)
	
	// Получаем события на неделю с 14 ноября
	events := storage.GetEventsForWeek(thursday)
	
	if len(events) != 2 {
		t.Errorf("Expected 2 events for the week, got %d", len(events))
	}
}

func TestGetEventsForMonth(t *testing.T) {
	logger, _ := logger.NewLogger()
	defer logger.Close()
	
	storage := NewEventStorage(logger)
	
	// Дата в ноябре 2025
	novemberDate := time.Date(2025, 11, 14, 10, 0, 0, 0, time.UTC)
	// Дата в декабре 2025
	decemberDate := time.Date(2025, 12, 5, 10, 0, 0, 0, time.UTC)
	
	// События в ноябре
	event1 := models.Event{
		ID:        1,
		StartTime: time.Date(2025, 11, 1, 10, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2025, 11, 1, 11, 0, 0, 0, time.UTC),
		Title:     "November Start",
	}
	event2 := models.Event{
		ID:        2,
		StartTime: time.Date(2025, 11, 30, 10, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2025, 11, 30, 11, 0, 0, 0, time.UTC),
		Title:     "November End",
	}
	// Событие в декабре
	event3 := models.Event{
		ID:        3,
		StartTime: decemberDate,
		EndTime:   decemberDate.Add(time.Hour),
		Title:     "December Event",
	}
	
	storage.AddEvent(event1)
	storage.AddEvent(event2)
	storage.AddEvent(event3)
	
	// Получаем события на ноябрь
	events := storage.GetEventsForMonth(novemberDate)
	
	if len(events) != 2 {
		t.Errorf("Expected 2 events for November, got %d", len(events))
	}
}
