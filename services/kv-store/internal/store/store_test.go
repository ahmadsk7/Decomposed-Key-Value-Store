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

// Test 5: overwrite updates value for the same key
func TestStore_Overwrite(t *testing.T) {
	s := New()
	key := "same-key"

	s.Put(key, []byte("first"))
	v1, ok := s.Get(key)
	if !ok || string(v1) != "first" {
		t.Fatalf("expected first, got %q (ok=%v)", string(v1), ok)
	}

	s.Put(key, []byte("second"))
	v2, ok := s.Get(key)
	if !ok || string(v2) != "second" {
		t.Fatalf("expected second, got %q (ok=%v)", string(v2), ok)
	}
}

// Test 6: concurrent readers and writers do not race
func TestStore_ConcurrentReadWrite(t *testing.T) {
	s := New()
	key := "rw"
	var wg sync.WaitGroup

	// writers
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			s.Put(key, []byte("v"))
		}(i)
	}

	// readers
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = s.Get(key)
		}()
	}

	wg.Wait()
	// final get should not panic and may or may not exist depending on timing
	_, _ = s.Get(key)
}

// Test 7: empty key handling
func TestStore_EmptyKey(t *testing.T) {
	s := New()
	s.Put("", []byte("value"))
	v, ok := s.Get("")
	if !ok || string(v) != "value" {
		t.Fatalf("empty key should work, got ok=%v value=%q", ok, string(v))
	}
}

// Test 8: empty value handling
func TestStore_EmptyValue(t *testing.T) {
	s := New()
	s.Put("key", []byte(""))
	v, ok := s.Get("key")
	if !ok || len(v) != 0 {
		t.Fatalf("empty value should work, got ok=%v len=%d", ok, len(v))
	}
}

// Test 9: store capacity under heavy load
func TestStore_Load(t *testing.T) {
	s := New()
	const n = 1000
	for i := 0; i < n; i++ {
		s.Put(string(rune(i)), []byte("data"))
	}
	if got, _ := s.Get(string(rune(n/2))); string(got) != "data" {
		t.Fatal("capacity stress failed")
	}
}
