package storage

import (
	"calendar/internal/logger"
	"calendar/internal/models"
	"errors"
	"sync"
	"time"
)

type EventStorage struct {
	events map[int]models.Event
	mu     sync.RWMutex
	logger *logger.Logger
}

func NewEventStorage(logger *logger.Logger) *EventStorage {
	return &EventStorage{
		events: make(map[int]models.Event, 10),
		logger: logger,
	}
}

func (s *EventStorage) AddEvent(event models.Event) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    s.logger.Info("adding event", "event_id", event.ID, "title", event.Title)

    // Валидация, бизнес-правила
    if event.Title == "" {
        return errors.New("title is required")
    }
    
    // s.events = append(s.events, event)
	s.events[event.ID] = event
    return nil
}

func (s *EventStorage) UpdateEvent(event models.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.logger.Info("update event", "event id", event, "title", event.Title)

	if event.Title == "" {
		return errors.New("title is required")
	}

	if event.StartTime.IsZero() || event.EndTime.IsZero() {
		return errors.New("time is not set")
	}

	s.events[event.ID] = event
	return nil
}

func (s *EventStorage) DeleteEvent(eventID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.logger.Info("delete event", "event id", eventID, "title", s.events[eventID].Title)

	delete(s.events, eventID)

	return nil
}

func (s *EventStorage) GetEventsForDay(date time.Time) []models.Event {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]models.Event, 0)
	for _, e := range s.events {
		if isSameDay(e.StartTime, date) {
			result = append(result, e)
		}
	}
	return result
}


func isSameDay(d1, d2 time.Time) bool {
	y1, m1, day1 := d1.Date()
	y2, m2, day2 := d2.Date()
	return y1 == y2 && m1 == m2 && day1 == day2
}

func (s *EventStorage) GetEventsForWeek(date time.Time) []models.Event {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	// Get week days (0 = Sunday, 1 = Monday)
	weekday := int(date.Weekday())
	
	// Transform to russian week format: Monday = 1, Sunday = 7
	if weekday == 0 {
		weekday = 7
	}
	
	// start time of week (monday 00:00:00)
	daysFromMonday := weekday - 1
	weekStart := date.AddDate(0, 0, -daysFromMonday)
	weekStart = time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, date.Location())
	
	// end of week (next monday 00:00:00)
	weekEnd := weekStart.AddDate(0, 0, 7)
	
	result := make([]models.Event, 0)
	for _, e := range s.events {
		// check the event is inclded?
		if !e.StartTime.Before(weekStart) && e.StartTime.Before(weekEnd) {
			result = append(result, e)
		}
	}
	
	return result
}

func (s *EventStorage) GetEventsForMonth(date time.Time) []models.Event {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    // Первый день месяца
    monthStart := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
    
    // Первый день следующего месяца
    monthEnd := monthStart.AddDate(0, 1, 0)
    
    result := make([]models.Event, 0)
    for _, e := range s.events {
        if !e.StartTime.Before(monthStart) && e.StartTime.Before(monthEnd) {
            result = append(result, e)
        }
    }
    
    return result
}
