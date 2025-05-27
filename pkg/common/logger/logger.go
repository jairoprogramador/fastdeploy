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

type logEntry struct {
	Level     LogLevel
	Message   string
	Error     error
	Timestamp time.Time
}

func (l *logEntry) format() string {
	var levelStr string
	switch l.Level {
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
		l.Timestamp.Format("2006-01-02 15:04:05"),
		levelStr,
		l.Message)

	if l.Error != nil {
		logLine += fmt.Sprintf(" - Error: %v", l.Error)
	}

	return logLine
}
