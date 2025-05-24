package logger

import (
	"fmt"
	"sync"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	SUCCESS
	WARNING
	ERROR
)

type LogEntry struct {
	Level     LogLevel
	Message   string
	Error     error
	Timestamp time.Time
}

type Logger struct {
	mu            sync.Mutex
	ConsoleLogger *ConsoleLogger
	FileLogger    *FileLogger
}

func NewLogger(filePath string) *Logger {
	return &Logger{
		ConsoleLogger: NewConsoleLogger(),
		FileLogger:    NewFileLogger(filePath),
	}
}

func (l *Logger) NewError(message string) error {
	if message == "" {
		return nil
	}
	err := fmt.Errorf("%s", message)
	l.Error(err)
	return err
}

func (l *Logger) Success(message string) {
	if message == "" {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.ConsoleLogger.Log(LogEntry{
		Level:   SUCCESS,
		Message: message,
	})
}

func (l *Logger) Info(message string) {
	if message == "" {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.ConsoleLogger.Log(LogEntry{
		Level:   INFO,
		Message: message,
	})
}

func (l *Logger) Warning(message string) {
	if message == "" {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.ConsoleLogger.Log(LogEntry{
		Level:   WARNING,
		Message: message,
	})
}

func (l *Logger) Error(err error) {
	if err == nil {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.ConsoleLogger.Log(LogEntry{
		Level: ERROR,
		Error: err,
	})
}

func (l *Logger) Debug(message string) {
	if message == "" {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.ConsoleLogger.Log(LogEntry{
		Level:   DEBUG,
		Message: message,
	})
}

func (l *Logger) SuccessSystem(message string) {
	if message == "" {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.FileLogger.Log(LogEntry{
		Level:   SUCCESS,
		Message: message,
	})
}

func (l *Logger) InfoSystem(message string) {
	if message == "" {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.FileLogger.Log(LogEntry{
		Level:   INFO,
		Message: message,
	})
}

func (l *Logger) WarningSystem(message string) {
	if message == "" {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.FileLogger.Log(LogEntry{
		Level:   WARNING,
		Message: message,
	})
}

func (l *Logger) ErrorSystem(err error) {
	if err == nil {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.FileLogger.Log(LogEntry{
		Level: ERROR,
		Error: err,
	})
}

func (l *Logger) ErrorSystemMessage(message string, err error) {
	if err == nil {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.FileLogger.Log(LogEntry{
		Level:   ERROR,
		Error:   err,
		Message: message,
	})
}

func (l *Logger) DebugSystem(message string) {
	if message == "" {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.FileLogger.Log(LogEntry{
		Level:   DEBUG,
		Message: message,
	})
}

func (l *Logger) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.ConsoleLogger.Clear()
	l.FileLogger.Clear()
}

/* type LoggerInterface interface {
	Log(entry LogEntry)
	GetSystemLogs() []LogEntry
	GetUserLogs() []LogEntry
	Clear()
	WriteToFile() error
} */

func formatLogEntry(entry LogEntry) string {
	var levelStr string
	switch entry.Level {
	case DEBUG:
		levelStr = "DEBUG"
	case INFO:
		levelStr = "INFO"
	case SUCCESS:
		levelStr = "SUCCESS"
	case ERROR:
		levelStr = "ERROR"
	}

	logLine := fmt.Sprintf("[%s] %s: %s",
		entry.Timestamp.Format("2006-01-02 15:04:05"),
		levelStr,
		entry.Message)

	if entry.Error != nil {
		logLine += fmt.Sprintf(" - Error: %v", entry.Error)
	}

	return logLine
}
