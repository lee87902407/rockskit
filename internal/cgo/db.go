package cgo

/*
#include <stdlib.h>
#include <rocksdb/c.h>
*/
import "C"

import "unsafe"

type DB struct {
	ptr *C.rocksdb_t
	*artifacts
}

func Create(path string, cfg *Config) (*DB, error) {
	return open(path, cfg, true, true)
}

func Open(path string, cfg *Config) (*DB, error) {
	return open(path, cfg, false, false)
}

func open(path string, cfg *Config, createIfMissing bool, errorIfExists bool) (*DB, error) {
	artifacts, err := buildArtifacts(cfg, createIfMissing, errorIfExists)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			artifacts.Close()
		}
	}()

	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	var errptr *C.char
	db := C.rocksdb_open(artifacts.options.ptr, cPath, &errptr)
	if err := errFromC(errptr, "open rocksdb"); err != nil {
		return nil, err
	}
	return &DB{
		ptr:       db,
		artifacts: artifacts,
	}, nil
}

func (db *DB) Close() {
	if db == nil || db.ptr == nil {
		return
	}
	C.rocksdb_close(db.ptr)
	db.ptr = nil
	if db.artifacts != nil {
		db.artifacts.Close()
		db.artifacts = nil
	}
}

type artifacts struct {
	options      *Options
	blockOptions *BlockBasedOptions
	filterPolicy *FilterPolicy
	cache        *LRUCache
	rateLimiter  *RateLimiter
}

func (a *artifacts) Close() {
	if a == nil {
		return
	}
	if a.options != nil {
		a.options.Close()
		a.options = nil
	}
	if a.blockOptions != nil {
		a.blockOptions.Close()
		a.blockOptions = nil
	}
	a.filterPolicy = nil
	if a.cache != nil {
		a.cache.Close()
		a.cache = nil
	}
	if a.rateLimiter != nil {
		a.rateLimiter.Close()
		a.rateLimiter = nil
	}
}
