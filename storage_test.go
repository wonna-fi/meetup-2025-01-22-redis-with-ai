package main

import (
	"fmt"
	"sync"
	"testing"
)

func TestStorage(t *testing.T) {
	s := NewStorage()

	// Test Set and Get
	s.Set("key1", "value1")
	value, exists := s.Get("key1")
	if !exists {
		t.Error("Expected key1 to exist")
	}
	if value != "value1" {
		t.Errorf("Expected value1, got %s", value)
	}

	// Test Get non-existent key
	_, exists = s.Get("nonexistent")
	if exists {
		t.Error("Expected nonexistent key to not exist")
	}

	// Test Del
	existed := s.Del("key1")
	if !existed {
		t.Error("Expected key1 to exist before deletion")
	}
	_, exists = s.Get("key1")
	if exists {
		t.Error("Expected key1 to be deleted")
	}

	// Test Del non-existent key
	existed = s.Del("nonexistent")
	if existed {
		t.Error("Expected Del of nonexistent key to return false")
	}

	// Test Len
	if s.Len() != 0 {
		t.Errorf("Expected length 0, got %d", s.Len())
	}
	s.Set("key1", "value1")
	s.Set("key2", "value2")
	if s.Len() != 2 {
		t.Errorf("Expected length 2, got %d", s.Len())
	}
}

func TestStorageConcurrent(t *testing.T) {
	s := NewStorage()
	var wg sync.WaitGroup

	// Test concurrent writes
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", i)
			value := fmt.Sprintf("value%d", i)
			s.Set(key, value)
		}(i)
	}
	wg.Wait()

	if s.Len() != 100 {
		t.Errorf("Expected 100 items, got %d", s.Len())
	}

	// Test concurrent reads and writes
	wg = sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(2)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", i)
			value := fmt.Sprintf("newvalue%d", i)
			s.Set(key, value)
		}(i)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", i)
			_, _ = s.Get(key)
		}(i)
	}
	wg.Wait()

	// Test concurrent deletes
	wg = sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", i)
			s.Del(key)
		}(i)
	}
	wg.Wait()

	if s.Len() != 0 {
		t.Errorf("Expected 0 items after deletion, got %d", s.Len())
	}
}

func TestStorageRaceConditions(t *testing.T) {
	s := NewStorage()
	var wg sync.WaitGroup

	// Perform all operations concurrently on the same key
	key := "testkey"
	operations := 100

	for i := 0; i < operations; i++ {
		wg.Add(3)

		// Writer
		go func(i int) {
			defer wg.Done()
			s.Set(key, fmt.Sprintf("value%d", i))
		}(i)

		// Reader
		go func() {
			defer wg.Done()
			_, _ = s.Get(key)
		}()

		// Deleter
		go func() {
			defer wg.Done()
			s.Del(key)
		}()
	}

	wg.Wait()
	// We can't assert the final state as it depends on the order of operations,
	// but we can verify that no panics occurred
}
