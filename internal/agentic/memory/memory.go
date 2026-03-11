package memory

import (
	"context"
	"time"
)

// Memory provides shared scratchpad for agents.
type Memory struct {
	store map[string]interface{}
}

// New creates a Memory store.
func New() *Memory {
	return &Memory{store: make(map[string]interface{})}
}

// Set stores a value with a key.
func (m *Memory) Set(ctx context.Context, key string, value interface{}) {
	m.store[key] = value
}

// Get retrieves a value by key.
func (m *Memory) Get(ctx context.Context, key string) (interface{}, bool) {
	val, ok := m.store[key]
	return val, ok
}

// List returns all keys.
func (m *Memory) List(ctx context.Context) []string {
	var keys []string
	for k := range m.store {
		keys = append(keys, k)
	}
	return keys
}
