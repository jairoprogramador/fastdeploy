package model

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// LogLevel define los niveles de log
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	SUCCESS
	ERROR
)

// LogEntry representa una entrada de log
type LogEntry struct {
	Level     LogLevel
	Message   string
	Error     error
	Timestamp time.Time
}

// LoggerInterface define la interfaz para el sistema de logging
type LoggerInterface interface {
	Log(entry LogEntry)
	GetLogs() []LogEntry
	Clear()
	WriteToFile() error
}

// FileLogger implementa LoggerInterface para escribir logs en archivo
type FileLogger struct {
	mu       sync.Mutex
	logs     []LogEntry
	filePath string
}

// NewFileLogger crea una nueva instancia de FileLogger
func NewFileLogger(filePath string) *FileLogger {
	return &FileLogger{
		logs:     make([]LogEntry, 0),
		filePath: filePath,
	}
}

// Log implementa el método Log de LoggerInterface
func (l *FileLogger) Log(entry LogEntry) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Solo guardamos logs de nivel ERROR y SUCCESS
	if entry.Level == ERROR || entry.Level == SUCCESS {
		entry.Timestamp = time.Now()
		l.logs = append(l.logs, entry)
	}
}

// GetLogs implementa el método GetLogs de LoggerInterface
func (l *FileLogger) GetLogs() []LogEntry {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.logs
}

// Clear implementa el método Clear de LoggerInterface
func (l *FileLogger) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logs = make([]LogEntry, 0)

	// Limpiamos el archivo de log
	os.Truncate(l.filePath, 0)
}

// WriteToFile implementa el método WriteToFile de LoggerInterface
func (l *FileLogger) WriteToFile() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	file, err := os.OpenFile(l.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error abriendo archivo de log: %v", err)
	}
	defer file.Close()

	for _, entry := range l.logs {
		var levelStr string
		switch entry.Level {
		case ERROR:
			levelStr = "ERROR"
		case SUCCESS:
			levelStr = "SUCCESS"
		}

		logLine := fmt.Sprintf("[%s] %s: %s",
			entry.Timestamp.Format("2006-01-02 15:04:05"),
			levelStr,
			entry.Message)

		if entry.Error != nil {
			logLine += fmt.Sprintf(" - Error: %v", entry.Error)
		}

		if _, err := file.WriteString(logLine + "\n"); err != nil {
			return fmt.Errorf("error escribiendo en archivo de log: %v", err)
		}
	}

	return nil
}
