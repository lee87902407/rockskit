package cgo

/*
#cgo CFLAGS: -I${SRCDIR}/../../c/include
#include <stdlib.h>
#include <rocksdb/c.h>
#include "rockskit.h"
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

func (db *DB) Put(opts *WriteOptions, key, value []byte) error {
	if db == nil || db.ptr == nil {
		return errClosedDB("put")
	}
	writeOpts := ensureWriteOptions(opts)
	var errptr *C.char
	C.rocksdb_put(
		db.ptr,
		writeOpts.ptr,
		bytesPtr(key),
		C.size_t(len(key)),
		bytesPtr(value),
		C.size_t(len(value)),
		&errptr,
	)
	return errFromC(errptr, "put rocksdb key")
}

func (db *DB) GetPinned(opts *ReadOptions, key []byte) (*PinnableSlice, error) {
	if db == nil || db.ptr == nil {
		return nil, errClosedDB("get pinned")
	}
	readOpts := ensureReadOptions(opts)
	var errptr *C.char
	var data *C.char
	var length C.size_t
	ptr := C.rockskit_get_pinned(
		db.ptr,
		readOpts.ptr,
		bytesPtr(key),
		C.size_t(len(key)),
		(**C.char)(unsafe.Pointer(&data)),
		&length,
		&errptr,
	)
	if err := errFromC(errptr, "get pinned rocksdb key"); err != nil {
		return nil, err
	}
	if ptr == nil {
		return nil, nil
	}
	return newPinnableSlice(ptr, data, length), nil
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

func bytesPtr(data []byte) *C.char {
	if len(data) == 0 {
		return nil
	}
	return (*C.char)(unsafe.Pointer(&data[0]))
}
