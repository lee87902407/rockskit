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

func (db *DB) Put(key, value []byte) error {
	if db == nil || db.closed.Load() {
		return fmt.Errorf("put: db is closed")
	}
	if err := db.hdl.Put(nil, key, value); err != nil {
		return fmt.Errorf("put: %w", err)
	}
	return nil
}

func (db *DB) Delete(key []byte) error {
	if db == nil || db.closed.Load() {
		return fmt.Errorf("delete: db is closed")
	}
	if err := db.hdl.Delete(nil, key); err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}

func (db *DB) PutCF(cf *ColumnFamilyHandle, key, value []byte) error {
	if db == nil || db.closed.Load() {
		return fmt.Errorf("put cf: db is closed")
	}
	if cf == nil || cf.hdl == nil {
		return fmt.Errorf("put cf: column family handle is nil")
	}
	if err := db.hdl.PutCF(cf.hdl, nil, key, value); err != nil {
		return fmt.Errorf("put cf: %w", err)
	}
	return nil
}

func (db *DB) DeleteCF(cf *ColumnFamilyHandle, key []byte) error {
	if db == nil || db.closed.Load() {
		return fmt.Errorf("delete cf: db is closed")
	}
	if cf == nil || cf.hdl == nil {
		return fmt.Errorf("delete cf: column family handle is nil")
	}
	if err := db.hdl.DeleteCF(cf.hdl, nil, key); err != nil {
		return fmt.Errorf("delete cf: %w", err)
	}
	return nil
}

func (db *DB) MultiGet(keys [][]byte) ([][]byte, []error) {
	if db == nil || db.closed.Load() {
		errs := make([]error, len(keys))
		for i := range errs {
			errs[i] = fmt.Errorf("multi get: db is closed")
		}
		return nil, errs
	}
	vals, errs := db.hdl.MultiGet(nil, keys)
	for i, e := range errs {
		if e != nil {
			errs[i] = fmt.Errorf("multi get: %w", e)
		}
	}
	return vals, errs
}

func (db *DB) MultiGetCF(cfs []*ColumnFamilyHandle, keys [][]byte) ([][]byte, []error) {
	if db == nil || db.closed.Load() {
		errs := make([]error, len(keys))
		for i := range errs {
			errs[i] = fmt.Errorf("multi get cf: db is closed")
		}
		return nil, errs
	}
	hdls := make([]*native.ColumnFamilyHandle, len(cfs))
	for i, cf := range cfs {
		if cf == nil || cf.hdl == nil {
			errs := make([]error, len(keys))
			errs[i] = fmt.Errorf("multi get cf: column family handle %d is nil", i)
			return nil, errs
		}
		hdls[i] = cf.hdl
	}
	vals, errs := db.hdl.MultiGetCF(hdls, nil, keys)
	for i, e := range errs {
		if e != nil {
			errs[i] = fmt.Errorf("multi get cf: %w", e)
		}
	}
	return vals, errs
}

func (db *DB) PutList(keys, values [][]byte) error {
	if db == nil || db.closed.Load() {
		return fmt.Errorf("put list: db is closed")
	}
	if err := db.hdl.PutList(nil, keys, values); err != nil {
		return fmt.Errorf("put list: %w", err)
	}
	return nil
}

func (db *DB) PutListCF(cf *ColumnFamilyHandle, keys, values [][]byte) error {
	if db == nil || db.closed.Load() {
		return fmt.Errorf("put list cf: db is closed")
	}
	if cf == nil || cf.hdl == nil {
		return fmt.Errorf("put list cf: column family handle is nil")
	}
	if err := db.hdl.PutListCF(cf.hdl, nil, keys, values); err != nil {
		return fmt.Errorf("put list cf: %w", err)
	}
	return nil
}

func (db *DB) DeleteList(keys [][]byte) error {
	if db == nil || db.closed.Load() {
		return fmt.Errorf("delete list: db is closed")
	}
	if err := db.hdl.DeleteList(nil, keys); err != nil {
		return fmt.Errorf("delete list: %w", err)
	}
	return nil
}

func (db *DB) DeleteListCF(cf *ColumnFamilyHandle, keys [][]byte) error {
	if db == nil || db.closed.Load() {
		return fmt.Errorf("delete list cf: db is closed")
	}
	if cf == nil || cf.hdl == nil {
		return fmt.Errorf("delete list cf: column family handle is nil")
	}
	if err := db.hdl.DeleteListCF(cf.hdl, nil, keys); err != nil {
		return fmt.Errorf("delete list cf: %w", err)
	}
	return nil
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
