package main

import (
	"calendar/internal/calendar/middlewares"
	config "calendar/internal/config"
	"calendar/internal/handlers"
	"calendar/internal/logger"
	"calendar/internal/storage"
	"log"
	"net/http"
	"time"
)

func main() {
	// Настройка конфигурации
	cfg := config.MustLoad()

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
	//addr := ":8080"
	log.Printf("Starting server on %s", cfg.Address)

	server := &http.Server{
		Addr:         cfg.Address,
		Handler:      middleware,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
