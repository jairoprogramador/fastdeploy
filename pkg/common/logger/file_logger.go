package logger

import (
	"fmt"
	"os"
	"sync"
	"time"
)

const (
	errorOpeningLogFile = "error opening log file: %v"
	errorWritingLogFile = "error writing log file: %v"
)

type FileLogger struct {
	mu       sync.Mutex
	logs     []logEntry
	filePath string
}

func NewFileLogger(filePath string) *FileLogger {
	return &FileLogger{
		logs:     make([]logEntry, 0),
		filePath: filePath,
	}
}

func (l *FileLogger) Error(err error) {
	if err == nil {
		return
	}
	var log = logEntry{
		Level: ERROR,
		Error: err,
	}
	l.Log(log)
}

func (l *FileLogger) Info(message string) {
	if message == "" {
		return
	}
	var log = logEntry{
		Level:   INFO,
		Message: message,
	}
	l.Log(log)
}

func (l *FileLogger) Log(entry logEntry) {
	l.mu.Lock()
	defer l.mu.Unlock()

	entry.Timestamp = time.Now()
	l.logs = append(l.logs, entry)
}

func (l *FileLogger) GetLogs() []logEntry {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.logs
}

func (l *FileLogger) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logs = make([]logEntry, 0)
	os.Truncate(l.filePath, 0)
}

func (l *FileLogger) WriteToFile() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	file, err := os.OpenFile(l.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf(errorOpeningLogFile, err)
	}
	defer file.Close()

	if _, err := file.WriteString("=== System Logs ===\n"); err != nil {
		return fmt.Errorf(errorWritingLogFile, err)
	}

	for _, entry := range l.logs {
		logLine := entry.format()
		if _, err := file.WriteString(logLine + "\n"); err != nil {
			return fmt.Errorf(errorWritingLogFile, err)
		}
	}
	return nil
}
