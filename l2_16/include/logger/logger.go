package logger

import (
	"fmt"
	"os"
	"time"
)

type Logger struct {
	Error string
	file *os.File
}

func NewLogger() *Logger {
	const pathOfLogFile = "include/logs/wgetmanager.log"
	file, err := os.OpenFile(pathOfLogFile, os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0777)
	if err != nil {
		fmt.Println("Problems with logging")
	}

	var logger = &Logger{
		file: file,
	}
	
	return logger
}

func (l *Logger) Write(p []byte) (int, error) {
    // время и пробел
    ts := time.Now().Format("15.04 02.01.2006")
    // уберём конечный \n у входа, чтобы не плодить пустые строки
    if len(p) > 0 && p[len(p)-1] == '\n' {
        p = p[:len(p)-1]
    }
    line := fmt.Sprintf("%s %s\n", ts, p)
    return l.file.Write([]byte(line))
}
