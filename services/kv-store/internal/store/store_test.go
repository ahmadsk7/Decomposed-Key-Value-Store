package store

import (
	"sync"
	"testing"
)

// Test 1: put a value then retrieve it
func TestStore_Put(t *testing.T) {
	s := New()
	key := "test-key"
	value := []byte("test-value")

	s.Put(key, value)
	stored, exists := s.Get(key)

	if !exists {
		t.Fatal("value should exist after put")
	}
	if string(stored) != string(value) {
		t.Fatalf("the expected value is %s, but we got %s", string(value), string(stored))
	}
}

// Test 2: get existing value and non-existent key
func TestStore_Get(t *testing.T) {
	s := New()
	key := "test-key"
	value := []byte("test-value")

	s.Put(key, value)
	stored, exists := s.Get(key)

	if !exists {
		t.Fatal("value should exist")
	}
	if string(stored) != string(value) {
		t.Fatalf("expected %s, got %s", string(value), string(stored))
	}

	_, exists = s.Get("non-existent")
	if exists {
		t.Fatal("non-existent key should not exist")
	}
}

// Test 3: delete existing key and non-existent key
func TestStore_Delete(t *testing.T) {
	s := New()
	key := "test-key"
	value := []byte("test-value")

	s.Put(key, value)
	deleted := s.Delete(key)

	if !deleted {
		t.Fatal("delete should return true for existing key")
	}

	_, exists := s.Get(key)
	if exists {
		t.Fatal("key should not exist after delete")
	}

	deleted = s.Delete("non-existent")
	if deleted {
		t.Fatal("delete should return false for non-existent key")
	}
}

// Test 4: concurrent writes to verify thread safety
func TestStore_Concurrent(t *testing.T) {
	s := New()
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			s.Put("key", []byte("value"))
		}(i)
	}

	wg.Wait()
}
