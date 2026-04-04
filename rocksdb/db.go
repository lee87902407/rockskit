package rocksdb

import (
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"

	native "github.com/lee87902407/rockskit/internal/cgo"
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

func OpenWithColumnFamilies(path string, cfg *Config) (*DB, map[string]*ColumnFamilyHandle, error) {
	cleanPath := filepath.Clean(path)
	cfg = cloneConfig(cfg)
	if err := ValidateConfig(cfg); err != nil {
		return nil, nil, err
	}
	hdl, rawCFMap, err := native.OpenWithColumnFamilies(cleanPath, cfg, true)
	if err != nil {
		return nil, nil, fmt.Errorf("open with column families at %q: %w", cleanPath, err)
	}
	cfMap := make(map[string]*ColumnFamilyHandle, len(rawCFMap))
	for name, handle := range rawCFMap {
		cfMap[name] = &ColumnFamilyHandle{hdl: handle}
	}
	return &DB{path: cleanPath, hdl: hdl}, cfMap, nil
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

type ColumnFamilyHandle struct {
	hdl *native.ColumnFamilyHandle
}

func (db *DB) CreateColumnFamily(name string) (*ColumnFamilyHandle, error) {
	if db == nil || db.closed.Load() {
		return nil, fmt.Errorf("create column family %q: db is closed", name)
	}
	opts := native.NewOptions()
	defer opts.Close()
	hdl, err := db.hdl.CreateColumnFamily(opts, name)
	if err != nil {
		return nil, fmt.Errorf("create column family %q: %w", name, err)
	}
	return &ColumnFamilyHandle{hdl: hdl}, nil
}

func (db *DB) DropColumnFamily(cf *ColumnFamilyHandle) error {
	if db == nil || db.closed.Load() {
		return fmt.Errorf("drop column family: db is closed")
	}
	if cf == nil || cf.hdl == nil {
		return nil
	}
	if err := db.hdl.DropColumnFamily(cf.hdl); err != nil {
		return fmt.Errorf("drop column family: %w", err)
	}
	cf.hdl.Close()
	cf.hdl = nil
	return nil
}

func (db *DB) GetDefaultColumnFamily() *ColumnFamilyHandle {
	if db == nil || db.closed.Load() {
		return nil
	}
	hdl := db.hdl.GetDefaultColumnFamily()
	if hdl == nil {
		return nil
	}
	return &ColumnFamilyHandle{hdl: hdl}
}

func (cf *ColumnFamilyHandle) Close() {
	if cf == nil || cf.hdl == nil {
		return
	}
	cf.hdl.Close()
	cf.hdl = nil
}
