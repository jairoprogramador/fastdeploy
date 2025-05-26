package logger

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type FileLogger struct {
	mu       sync.Mutex
	logs     []LogEntry
	filePath string
}

func NewFileLogger(filePath string) *FileLogger {
	return &FileLogger{
		logs:     make([]LogEntry, 0),
		filePath: filePath,
	}
}

func (l *FileLogger) Error(err error) {
	if err == nil {
		return
	}
	var log = LogEntry{
		Level: ERROR,
		Error: err,
	}
	l.Log(log)
}

func (l *FileLogger) Info(message string) {
	if message == "" {
		return
	}
	var log = LogEntry{
		Level:   INFO,
		Message: message,
	}
	l.Log(log)
}

func (l *FileLogger) Log(entry LogEntry) {
	l.mu.Lock()
	defer l.mu.Unlock()

	entry.Timestamp = time.Now()
	l.logs = append(l.logs, entry)
}

func (l *FileLogger) GetLogs() []LogEntry {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.logs
}

func (l *FileLogger) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logs = make([]LogEntry, 0)
	os.Truncate(l.filePath, 0)
}

func (l *FileLogger) WriteToFile() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	file, err := os.OpenFile(l.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error abriendo archivo de log: %v", err)
	}
	defer file.Close()

	if _, err := file.WriteString("=== System Logs ===\n"); err != nil {
		return fmt.Errorf("error escribiendo en archivo de log: %v", err)
	}

	for _, entry := range l.logs {
		logLine := formatLogEntry(entry)
		if _, err := file.WriteString(logLine + "\n"); err != nil {
			return fmt.Errorf("error escribiendo en archivo de log: %v", err)
		}
	}
	return nil
}
