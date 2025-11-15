package models

import "time"

type Event struct {
	ID int
	StartTime time.Time
	EndTime time.Time
	Title string
	Description string
}

func NewEvent(
	startTime, endTime time.Time,
	title, description string,
) Event {
	event := Event{
		StartTime: startTime,
		EndTime: endTime,
		Title: title, 
		Description: description,
	}

	return event
} 