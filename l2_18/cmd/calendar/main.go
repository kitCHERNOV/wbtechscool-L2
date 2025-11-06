package main

import (
	"calendar/internal/calendar/middlewares"
	"calendar/internal/logger"
	"encoding/json"
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
    mux := http.NewServeMux()
	
    // Регистрация обработчиков
    mux.HandleFunc("/create_event", createEventHandler)
    mux.HandleFunc("/update_event", updateEventHandler)
    mux.HandleFunc("/delete_event", deleteEventHandler)
    mux.HandleFunc("/events_for_day", eventsForDayHandler)
    mux.HandleFunc("/events_for_week", eventsForWeekHandler)
    mux.HandleFunc("/events_for_month", eventsForMonthHandler)

    // Middleware для логирования
	logger := logger.NewLogger()
    handler := middlewares.LoggingMiddleware(mux, logger)

    // Запуск сервера
    addr := ":8080"
    log.Printf("Starting server on %s", addr)
    
    server := &http.Server{
        Addr:         addr,
        Handler:      handler,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
    }

    if err := server.ListenAndServe(); err != nil {
        log.Fatal(err)
    }
}

// Заглушки обработчиков (реализуй сам)
func createEventHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    // TODO: реализовать
    respondJSON(w, http.StatusOK, map[string]string{"status": "created"})
}

func updateEventHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    // TODO: реализовать
    respondJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

func deleteEventHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    // TODO: реализовать
    respondJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func eventsForDayHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    // TODO: реализовать
    respondJSON(w, http.StatusOK, []string{})
}

func eventsForWeekHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    // TODO: реализовать
    respondJSON(w, http.StatusOK, []string{})
}

func eventsForMonthHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    // TODO: реализовать
    respondJSON(w, http.StatusOK, []string{})
}

// Вспомогательная функция для JSON ответов
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}