package rocksdb

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

const markerFileName = "store.json"

type DB struct {
	path string

	mu     sync.Mutex
	closed bool
}

func Create(path string) (*DB, error) {
	cleanPath := filepath.Clean(path)

	if _, err := os.Stat(cleanPath); err == nil {
		return nil, fmt.Errorf("create rocksdb at %q: already exists", cleanPath)
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("stat rocksdb path %q: %w", cleanPath, err)
	}

	if err := os.MkdirAll(cleanPath, 0o755); err != nil {
		return nil, fmt.Errorf("create rocksdb dir %q: %w", cleanPath, err)
	}
	if err := os.WriteFile(filepath.Join(cleanPath, markerFileName), []byte("{}\n"), 0o644); err != nil {
		return nil, fmt.Errorf("initialize rocksdb at %q: %w", cleanPath, err)
	}

	return &DB{path: cleanPath}, nil
}

func Open(path string) (*DB, error) {
	cleanPath := filepath.Clean(path)

	info, err := os.Stat(cleanPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("open rocksdb at %q: not found", cleanPath)
		}
		return nil, fmt.Errorf("stat rocksdb path %q: %w", cleanPath, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("open rocksdb at %q: not a directory", cleanPath)
	}
	if _, err := os.Stat(filepath.Join(cleanPath, markerFileName)); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("open rocksdb at %q: not initialized", cleanPath)
		}
		return nil, fmt.Errorf("check rocksdb marker at %q: %w", cleanPath, err)
	}

	return &DB{path: cleanPath}, nil
}

func (db *DB) Close() error {
	if db == nil {
		return nil
	}

	db.mu.Lock()
	defer db.mu.Unlock()
	db.closed = true
	return nil
}
