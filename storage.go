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

// Del removes a key-value pair
func (s *Storage) Del(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, exists := s.data[key]
	if exists {
		delete(s.data, key)
	}
	return exists
}

// Len returns the number of stored key-value pairs
func (s *Storage) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.data)
}
