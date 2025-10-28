package vos

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

type LoggerID struct {
	Value string `yaml:"value"`
}

func NewLoggerID() (LoggerID, error) {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return LoggerID{}, fmt.Errorf("could not generate random bytes for logger id: %w", err)
	}
	return LoggerID{Value: hex.EncodeToString(bytes)}, nil
}

func NewLoggerIDFromString(id string) (LoggerID, error) {
	if id == "" {
		return LoggerID{}, fmt.Errorf("logger ID cannot be empty")
	}
	return LoggerID{Value: id}, nil
}

func (e LoggerID) String() string {
	return e.Value
}
