package rocksdb

import (
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"

	native "github.com/yanjie/rockskit/internal/cgo"
)

type DB struct {
	path string
	hdl  *native.DB

	closed atomic.Bool
}

func Create(path string, cfg *Config) (*DB, error) {
	cleanPath := filepath.Clean(path)
	cfg = cloneConfig(cfg)
	if err := ValidateConfig(cfg); err != nil {
		return nil, err
	}
	hdl, err := native.Create(cleanPath, cfg)
	if err != nil {
		return nil, fmt.Errorf("create rocksdb at %q: %w", cleanPath, err)
	}
	return &DB{path: cleanPath, hdl: hdl}, nil
}

func Open(path string, cfg *Config) (*DB, error) {
	cleanPath := filepath.Clean(path)
	cfg = cloneConfig(cfg)
	if err := ValidateConfig(cfg); err != nil {
		return nil, err
	}
	hdl, err := native.Open(cleanPath, cfg)
	if err != nil {
		return nil, fmt.Errorf("open rocksdb at %q: %w", cleanPath, err)
	}
	return &DB{path: cleanPath, hdl: hdl}, nil
}

func createOrOpen(path string, cfg *Config) (*DB, error) {
	cleanPath := filepath.Clean(path)
	entries, err := os.ReadDir(cleanPath)
	if err == nil && len(entries) > 0 {
		return Open(cleanPath, cfg)
	}
	if err != nil && !os.IsNotExist(err) {
		return Open(cleanPath, cfg)
	}
	return Create(cleanPath, cfg)
}

func (db *DB) Close() error {
	if db == nil {
		return nil
	}
	if !db.closed.CompareAndSwap(false, true) {
		return nil
	}
	if db.hdl != nil {
		db.hdl.Close()
		db.hdl = nil
	}
	return nil
}
