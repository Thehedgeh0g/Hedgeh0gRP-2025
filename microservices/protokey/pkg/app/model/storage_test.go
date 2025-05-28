package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStorage_SetAndGet(t *testing.T) {
	s := NewStorage()

	key := "testKey"
	val := "123"

	if err := s.Set(key, val); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	got, err := s.Get(key)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if got != val {
		t.Errorf("Expected %s, got %s", val, got)
	}
}

func TestStorage_GetNonexistentKey(t *testing.T) {
	s := NewStorage()

	_, err := s.Get("noKey")
	assert.Error(t, err)
}

func TestStorage_KeysWithPrefix(t *testing.T) {
	s := NewStorage()

	err := s.Set("apple", "1")
	assert.NoError(t, err)
	err = s.Set("apricot", "2")
	assert.NoError(t, err)
	err = s.Set("banana", "3")
	assert.NoError(t, err)

	keys, err := s.Keys("ap")
	if err != nil {
		t.Fatalf("Keys failed: %v", err)
	}

	expected := map[string]bool{
		"apple":   true,
		"apricot": true,
	}

	if len(keys) != 2 {
		t.Errorf("Expected 2 keys, got %d", len(keys))
	}
	for _, k := range keys {
		if !expected[k] {
			t.Errorf("Unexpected key: %s", k)
		}
	}
}
