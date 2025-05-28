package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	errorOpeningLogFile  = "error opening log file: %v"
	errorWritingLogFile  = "error writing log file: %v"
	errorReadingLogFile  = "error reading log file: %v"
	errorFormatTimestamp = "error format timestamp must start -> [TIMESTAMP]"
	errorParseTimestamp  = "error parse timestamp: %v"
	errorFormatMessage   = "error format message must be -> LEVEL: MESSAGE"
	errorLogLevelUnknown = "error LogLevel unknown: %s"
)

type FileLogger struct {
	mu       sync.Mutex
	logs     []logEntry
	filePath string
	isSave   bool
}

func NewFileLogger(filePath string) *FileLogger {
	return &FileLogger{
		logs:     make([]logEntry, 0),
		filePath: filePath,
		isSave:   false,
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

func (l *FileLogger) WriteToFile() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if err := l.createDir(); err != nil {
		return l.myError(err)
	}

	file, err := os.OpenFile(l.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		myErr := fmt.Errorf(errorOpeningLogFile, err)
		return l.myError(myErr)
	}
	os.Truncate(l.filePath, 0)

	defer file.Close()

	if _, err := file.WriteString("=== System Logs ===\n"); err != nil {
		myErr := fmt.Errorf(errorWritingLogFile, err)
		return l.myError(myErr)
	}

	for _, entry := range l.logs {
		logLine := entry.format()
		if _, err := file.WriteString(logLine + "\n"); err != nil {
			myErr := fmt.Errorf(errorWritingLogFile, err)
			return l.myError(myErr)
		}
	}
	l.logs = make([]logEntry, 0)
	l.isSave = true
	return nil
}

func (l *FileLogger) ReadFromFile() ([]logEntry, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if !l.isSave {
		return l.logs, nil
	}

	file, err := os.ReadFile(l.filePath)
	if err != nil {
		myErr := fmt.Errorf(errorReadingLogFile, err)
		return nil, l.myError(myErr)
	}

	if len(file) == 0 {
		return make([]logEntry, 0), nil
	}

	lines := strings.Split(string(file), "\n")
	logs := make([]logEntry, 0)

	// Ignoramos la primera línea que contiene "=== System Logs ==="
	for i := 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		entry, err := l.parseLogLine(line)
		if err != nil {
			return nil, l.myError(err)
		}
		logs = append(logs, entry)
	}

	l.logs = make([]logEntry, 0)
	l.logs = append(l.logs, logs...)

	return logs, nil
}

func (l *FileLogger) parseLogLine(line string) (logEntry, error) {
	// Formato esperado: [TIMESTAMP] LEVEL: MESSAGE/ERROR
	entry := logEntry{}

	timestampStart := strings.Index(line, "[")
	timestampEnd := strings.Index(line, "]")
	if timestampStart == -1 || timestampEnd == -1 {
		myErr := fmt.Errorf("%s", errorFormatTimestamp)
		return entry, l.myError(myErr)
	}

	timestamp, err := time.Parse(LogTimeFormat,
		strings.TrimSpace(line[timestampStart+1:timestampEnd]))

	if err != nil {
		myErr := fmt.Errorf(errorParseTimestamp, err)
		return entry, l.myError(myErr)
	}
	entry.Timestamp = timestamp

	remainder := strings.TrimSpace(line[timestampEnd+1:])

	parts := strings.SplitN(remainder, ":", 2)
	if len(parts) != 2 {
		myErr := fmt.Errorf("%s", errorFormatMessage)
		return entry, l.myError(myErr)
	}

	level := strings.TrimSpace(parts[0])
	switch level {
	case "ERROR":
		entry.Level = ERROR
		entry.Error = fmt.Errorf(strings.TrimSpace(parts[1]))
	case "INFO":
		entry.Level = INFO
		entry.Message = strings.TrimSpace(parts[1])
	default:
		myErr := fmt.Errorf(errorLogLevelUnknown, level)
		return entry, l.myError(myErr)
	}

	return entry, nil
}

func (l *FileLogger) createDir() error {
	dirPath := filepath.Dir(l.filePath)
	_, err := os.Stat(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(dirPath, 0755); err != nil {
				return l.myError(err)
			}
		}
		return err
	}
	return nil
}

func (l *FileLogger) myError(err error) error {
	var log = logEntry{
		Level: ERROR,
		Error: err,
	}
	log.Timestamp = time.Now()
	l.logs = append(l.logs, log)
	return err
}
