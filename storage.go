package main

import (
	"sync"
)

// Storage represents our thread-safe key-value store
type Storage struct {
	mu   sync.RWMutex
	data map[string]string
}

// NewStorage creates a new Storage instance
func NewStorage() *Storage {
	return &Storage{
		data: make(map[string]string),
	}
}

// Set stores a key-value pair
func (s *Storage) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

// Get retrieves a value by key
func (s *Storage) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, exists := s.data[key]
	return value, exists
}

// Del removes one or more key-value pairs and returns the number of keys that were deleted
func (s *Storage) Del(keys ...string) int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	var deleted int64
	for _, key := range keys {
		if _, exists := s.data[key]; exists {
			delete(s.data, key)
			deleted++
		}
	}
	return deleted
}

// Len returns the number of stored key-value pairs
func (s *Storage) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.data)
}
