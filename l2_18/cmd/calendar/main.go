package main

import (
	"calendar/internal/calendar/middlewares"
	"calendar/internal/handlers"
	"calendar/internal/logger"
	"calendar/internal/storage"
	"log"
	"net/http"
	"time"
)

// creating calander http service
// TODO: complete the tasks below
// POST /create_event — создание нового события;
// POST /update_event — обновление существующего;
// POST /delete_event — удаление;
// GET /events_for_day — получить все события на день;
// GET /events_for_week — события на неделю;
// GET /events_for_month — события на месяц.

func main() {
	// Настройка конфигурации

	// инициализация логгера 
	logger, err := logger.NewLogger()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Close()

	logger.Info("aplication starting")



	// база данных событий календаря
	var eventStorage = storage.NewEventStorage(logger)

	// основной обработчик
	handler := handlers.NewHandler(eventStorage, logger)


	// Http сервис
	mux := http.NewServeMux()

	// Регистрация обработчиков
	mux.HandleFunc("/create_event", handler.CreateEventHandler)
	mux.HandleFunc("/update_event", handler.UpdateEventHandler)
	mux.HandleFunc("/delete_event", handler.DeleteEventHandler)
	mux.HandleFunc("/events_for_day", handler.EventsForDayHandler)
	mux.HandleFunc("/events_for_week", handler.EventsForWeekHandler)
	mux.HandleFunc("/events_for_month", handler.EventsForMonthHandler)

	// Middleware для логирования
	middleware := middlewares.LoggingMiddleware(mux, logger)

	// Запуск сервера
	addr := ":8080"
	log.Printf("Starting server on %s", addr)

	server := &http.Server{
		Addr:         addr,
		Handler:      middleware,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
