package logger

import (
	"fmt"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	SUCCESS
	WARNING
	ERROR
	UNKNOWN
)

type logEntry struct {
	Level   LogLevel
	Message string
	Error   error
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

	logLine := fmt.Sprintf("%s: %s",
		levelStr,
		l.Message)

	if l.Error != nil {
		logLine += fmt.Sprintf("%s", l.Error)
	}

	return logLine
}
