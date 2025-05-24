package logger

import (
	"sync"
	"time"
)

/* type UserLoggerInterface interface {
	Log(entry LogEntry)
	GetLogs() []LogEntry
	Clear()
} */

type ConsoleLogger struct {
	mu   sync.Mutex
	logs []LogEntry
}

func NewConsoleLogger() *ConsoleLogger {
	return &ConsoleLogger{
		logs: make([]LogEntry, 0),
	}
}

func (l *ConsoleLogger) Log(entry LogEntry) {
	l.mu.Lock()
	defer l.mu.Unlock()

	entry.Timestamp = time.Now()
	l.logs = append(l.logs, entry)
}

func (l *ConsoleLogger) GetLogs() []LogEntry {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.logs
}

func (l *ConsoleLogger) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logs = make([]LogEntry, 0)
}
