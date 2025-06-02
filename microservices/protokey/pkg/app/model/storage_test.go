package model_test

import (
	"os"
	"testing"
	"time"

	"protokey/pkg/app/model"
)

func cleanupFile(path string) {
	_ = os.Remove(path)
}

func createTestStorage(t *testing.T, path string) *model.Storage {
	// Подменим имя файла через переменную окружения
	t.Setenv("PROTOKEY_DATA_PATH", path)
	return model.NewStorage()
}

func TestSetAndGet(t *testing.T) {
	dataPath := "test_data_set_get.data"
	defer cleanupFile(dataPath)

	store := createTestStorage(t, dataPath)

	err := store.Set("foo", "bar")
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	val, err := store.Get("foo")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if val != "bar" {
		t.Fatalf("Expected 'bar', got '%s'", val)
	}
}

func TestKeyNotFound(t *testing.T) {
	dataPath := "test_data_notfound.data"
	defer cleanupFile(dataPath)

	store := createTestStorage(t, dataPath)

	_, err := store.Get("missing")
	if err != model.ErrKeyNotFound {
		t.Fatalf("Expected ErrKeyNotFound, got %v", err)
	}
}

func TestKeysPrefix(t *testing.T) {
	dataPath := "test_data_keys.data"
	defer cleanupFile(dataPath)

	store := createTestStorage(t, dataPath)

	store.Set("apple", "1")
	store.Set("appetizer", "2")
	store.Set("banana", "3")

	keys, err := store.Keys("app")
	if err != nil {
		t.Fatalf("Keys failed: %v", err)
	}

	expected := map[string]bool{"apple": true, "appetizer": true}
	if len(keys) != 2 {
		t.Fatalf("Expected 2 keys, got %d", len(keys))
	}
	for _, k := range keys {
		if !expected[k] {
			t.Errorf("Unexpected key: %s", k)
		}
	}
}

func TestFlushAndReload(t *testing.T) {
	dataPath := "test_data_reload.data"
	defer cleanupFile(dataPath)

	{
		store := createTestStorage(t, dataPath)
		store.Set("k1", "v1")
		store.Set("k2", "v2")
		time.Sleep(2 * time.Second) // ждём flush
	}

	// Recreate storage (simulate restart)
	store := createTestStorage(t, dataPath)

	v, err := store.Get("k1")
	if err != nil || v != "v1" {
		t.Fatalf("Expected v1, got %s (err: %v)", v, err)
	}

	v2, err := store.Get("k2")
	if err != nil || v2 != "v2" {
		t.Fatalf("Expected v2, got %s (err: %v)", v2, err)
	}
}

func TestInvalidLogLineIsIgnored(t *testing.T) {
	path := "test_data_invalid.data"
	defer cleanupFile(path)

	_ = os.WriteFile(path, []byte("bad json\n"), 0644)

	store := createTestStorage(t, path)

	err := store.Set("key", "val")
	if err != nil {
		t.Fatalf("Set failed after invalid log: %v", err)
	}
}
