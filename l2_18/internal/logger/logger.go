package logger

import (
	"fmt"
	"os"
	"time"
	"sync"
	"log/slog"
)

type Logger struct {
	file   *os.File
	mu     sync.Mutex
	logger *slog.Logger
}

func NewLogger() (*Logger, error) {
	const pathOfLogFile = "interanl/logs/wgetmanager.log"
	file, err := os.OpenFile(pathOfLogFile, os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0777)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	var logger = &Logger{
		file: file,
	}
	
	logger.logger = slog.New(slog.NewJSONHandler(file, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))


	return logger, nil
}

// func (l *Logger) Write(p []byte) (int, error) {
//     // время и пробел
//     ts := time.Now().Format("15.04 02.01.2006")
//     // уберём конечный \n у входа, чтобы не плодить пустые строки
//     if len(p) > 0 && p[len(p)-1] == '\n' {
//         p = p[:len(p)-1]
//     }
//     line := fmt.Sprintf("%s %s\n", ts, p)
//     return l.file.Write([]byte(line))
// }


func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// Write реализует io.Writer для совместимости
func (l *Logger) Write(p []byte) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	ts := time.Now().Format("2006-01-02 15:04:05")
	
	if len(p) > 0 && p[len(p)-1] == '\n' {
		p = p[:len(p)-1]
	}
	
	line := fmt.Sprintf("[%s] %s\n", ts, p)
	return l.file.Write([]byte(line))
}

// Методы для уровней логирования
func (l *Logger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}

func (l *Logger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}