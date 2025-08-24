package deployment

import (
	"fmt"
	"sync"
)

type Context interface {
	Get(key string) (string, error)
	Set(key, value string)
}

type DeploymentContext struct {
	mu     sync.RWMutex
	params map[string]string
}

func NewDeploymentContext() Context {
	return &DeploymentContext{
		params: make(map[string]string),
	}
}

func (c *DeploymentContext) Get(key string) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, ok := c.params[key]
	if !ok {
		return "", fmt.Errorf("parámetro no encontrado: %s", key)
	}
	return value, nil
}

func (c *DeploymentContext) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.params[key] = value
}
