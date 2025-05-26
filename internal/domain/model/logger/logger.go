package logger

import (
	"fmt"
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
