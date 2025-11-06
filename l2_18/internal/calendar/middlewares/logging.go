package middlewares

import (
	"calendar/internal/logger"
	"net/http"
	"time"
    "fmt"
)

// Middleware для логирования запросов
func LoggingMiddleware(next http.Handler, logger *logger.Logger) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        // log.Printf("%s %s", r.Method, r.URL.Path)
		fmt.Fprintf(logger, "%s %s\n", r.Method, r.URL.Path)
        next.ServeHTTP(w, r)
        fmt.Fprintf(logger, "Completed in %v", time.Since(start))
    })
}