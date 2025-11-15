package models

import "time"

type Event struct {
	ID          int       `json:"id"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
}

func NewEvent(
	startTime, endTime time.Time,
	title, description string,
) Event {
	event := Event{
		StartTime:   startTime,
		EndTime:     endTime,
		Title:       title,
		Description: description,
	}

	return event
}
