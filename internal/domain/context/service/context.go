package service

import (
	"fmt"
	"sync"
)

type Context interface {
	Get(key string) (string, error)
	Set(key, value string)
	GetAll() map[string]string
	SetAll(data map[string]string)
}

type DataContext struct {
	mu     sync.RWMutex
	params map[string]string
}

func NewDataContext() Context {
	return &DataContext{
		params: make(map[string]string),
	}
}

func (c *DataContext) Get(key string) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, ok := c.params[key]
	if !ok {
		return "", fmt.Errorf("parámetro no encontrado: %s", key)
	}
	return value, nil
}

func (c *DataContext) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.params[key] = value
}

func (c *DataContext) GetAll() map[string]string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.params
}

func (c *DataContext) SetAll(data map[string]string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.params = data
}