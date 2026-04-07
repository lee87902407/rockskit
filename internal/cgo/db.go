package cgo

/*
#cgo CFLAGS: -I${SRCDIR}/../../c/include
#include <stdlib.h>
#include <rocksdb/c.h>
#include "rockskit.h"
*/
import "C"

import (
	"fmt"
	"unsafe"
)

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

func (db *DB) Delete(opts *WriteOptions, key []byte) error {
	if db == nil || db.ptr == nil {
		return errClosedDB("delete")
	}
	writeOpts := ensureWriteOptions(opts)
	var errptr *C.char
	C.rocksdb_delete(
		db.ptr,
		writeOpts.ptr,
		bytesPtr(key),
		C.size_t(len(key)),
		&errptr,
	)
	return errFromC(errptr, "delete rocksdb key")
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

type ColumnFamilyHandle struct {
	ptr *C.rocksdb_column_family_handle_t
}

func (h *ColumnFamilyHandle) Close() {
	if h == nil || h.ptr == nil {
		return
	}
	C.rocksdb_column_family_handle_destroy(h.ptr)
	h.ptr = nil
}

func (db *DB) CreateColumnFamily(opts *Options, name string) (*ColumnFamilyHandle, error) {
	if db == nil || db.ptr == nil {
		return nil, errClosedDB("create column family")
	}
	if opts == nil || opts.ptr == nil {
		return nil, fmt.Errorf("create column family: options is nil")
	}
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	var errptr *C.char
	handle := C.rocksdb_create_column_family(db.ptr, opts.ptr, cName, &errptr)
	if err := errFromC(errptr, "create column family"); err != nil {
		return nil, err
	}
	return &ColumnFamilyHandle{ptr: handle}, nil
}

func (db *DB) DropColumnFamily(handle *ColumnFamilyHandle) error {
	if db == nil || db.ptr == nil {
		return errClosedDB("drop column family")
	}
	if handle == nil || handle.ptr == nil {
		return nil
	}
	var errptr *C.char
	C.rocksdb_drop_column_family(db.ptr, handle.ptr, &errptr)
	return errFromC(errptr, "drop column family")
}

func (db *DB) GetDefaultColumnFamily() *ColumnFamilyHandle {
	if db == nil || db.ptr == nil {
		return nil
	}
	handle := C.rocksdb_get_default_column_family_handle(db.ptr)
	if handle == nil {
		return nil
	}
	return &ColumnFamilyHandle{ptr: handle}
}

func (db *DB) DeleteCF(handle *ColumnFamilyHandle, opts *WriteOptions, key []byte) error {
	if db == nil || db.ptr == nil {
		return errClosedDB("delete cf")
	}
	if handle == nil || handle.ptr == nil {
		return fmt.Errorf("delete cf: column family handle is nil")
	}
	writeOpts := ensureWriteOptions(opts)
	var errptr *C.char
	C.rocksdb_delete_cf(
		db.ptr,
		writeOpts.ptr,
		handle.ptr,
		bytesPtr(key),
		C.size_t(len(key)),
		&errptr,
	)
	return errFromC(errptr, "delete cf")
}

func (db *DB) PutCF(handle *ColumnFamilyHandle, opts *WriteOptions, key, value []byte) error {
	if db == nil || db.ptr == nil {
		return errClosedDB("put cf")
	}
	if handle == nil || handle.ptr == nil {
		return fmt.Errorf("put cf: column family handle is nil")
	}
	writeOpts := ensureWriteOptions(opts)
	var errptr *C.char
	C.rocksdb_put_cf(
		db.ptr,
		writeOpts.ptr,
		handle.ptr,
		bytesPtr(key),
		C.size_t(len(key)),
		bytesPtr(value),
		C.size_t(len(value)),
		&errptr,
	)
	return errFromC(errptr, "put cf")
}

func (db *DB) GetPinnedCF(handle *ColumnFamilyHandle, opts *ReadOptions, key []byte) (*PinnableSlice, error) {
	if db == nil || db.ptr == nil {
		return nil, errClosedDB("get pinned cf")
	}
	if handle == nil || handle.ptr == nil {
		return nil, fmt.Errorf("get pinned cf: column family handle is nil")
	}
	readOpts := ensureReadOptions(opts)
	var errptr *C.char
	var data *C.char
	var length C.size_t
	ptr := C.rockskit_get_pinned_cf(
		db.ptr,
		handle.ptr,
		readOpts.ptr,
		bytesPtr(key),
		C.size_t(len(key)),
		(**C.char)(unsafe.Pointer(&data)),
		&length,
		&errptr,
	)
	if err := errFromC(errptr, "get pinned cf"); err != nil {
		return nil, err
	}
	if ptr == nil {
		return nil, nil
	}
	return newPinnableSlice(ptr, data, length), nil
}

func (db *DB) MultiGet(opts *ReadOptions, keys [][]byte) ([][]byte, []error) {
	if db == nil || db.ptr == nil {
		errs := make([]error, len(keys))
		for i := range errs {
			errs[i] = errClosedDB("multi get")
		}
		return nil, errs
	}
	n := len(keys)
	if n == 0 {
		return nil, nil
	}
	readOpts := ensureReadOptions(opts)

	cKeys := make([]*C.char, n)
	cKeySizes := make([]C.size_t, n)
	for i, k := range keys {
		cKeys[i] = bytesPtr(k)
		cKeySizes[i] = C.size_t(len(k))
	}
	cVals := make([]*C.char, n)
	cValSizes := make([]C.size_t, n)
	cErrs := make([]*C.char, n)

	C.rocksdb_multi_get(
		db.ptr,
		readOpts.ptr,
		C.size_t(n),
		(**C.char)(unsafe.Pointer(&cKeys[0])),
		(*C.size_t)(unsafe.Pointer(&cKeySizes[0])),
		(**C.char)(unsafe.Pointer(&cVals[0])),
		(*C.size_t)(unsafe.Pointer(&cValSizes[0])),
		(**C.char)(unsafe.Pointer(&cErrs[0])),
	)

	vals := make([][]byte, n)
	errs := make([]error, n)
	for i := 0; i < n; i++ {
		if cErrs[i] != nil {
			errs[i] = errFromC(cErrs[i], fmt.Sprintf("multi get key %d", i))
			continue
		}
		if cVals[i] != nil {
			vals[i] = C.GoBytes(unsafe.Pointer(cVals[i]), C.int(cValSizes[i]))
			C.rocksdb_free(unsafe.Pointer(cVals[i]))
		}
	}
	return vals, errs
}

func (db *DB) MultiGetCF(handles []*ColumnFamilyHandle, opts *ReadOptions, keys [][]byte) ([][]byte, []error) {
	if db == nil || db.ptr == nil {
		errs := make([]error, len(keys))
		for i := range errs {
			errs[i] = errClosedDB("multi get cf")
		}
		return nil, errs
	}
	if len(handles) != len(keys) {
		errs := make([]error, len(keys))
		for i := range errs {
			errs[i] = fmt.Errorf("multi get cf: handles count %d != keys count %d", len(handles), len(keys))
		}
		return nil, errs
	}
	n := len(keys)
	if n == 0 {
		return nil, nil
	}
	readOpts := ensureReadOptions(opts)

	cCFs := make([]*C.rocksdb_column_family_handle_t, n)
	cKeys := make([]*C.char, n)
	cKeySizes := make([]C.size_t, n)
	for i, k := range keys {
		if handles[i] == nil || handles[i].ptr == nil {
			return nil, []error{fmt.Errorf("multi get cf: column family handle %d is nil", i)}
		}
		cCFs[i] = handles[i].ptr
		cKeys[i] = bytesPtr(k)
		cKeySizes[i] = C.size_t(len(k))
	}
	cVals := make([]*C.char, n)
	cValSizes := make([]C.size_t, n)
	cErrs := make([]*C.char, n)

	C.rocksdb_multi_get_cf(
		db.ptr,
		readOpts.ptr,
		(**C.rocksdb_column_family_handle_t)(unsafe.Pointer(&cCFs[0])),
		C.size_t(n),
		(**C.char)(unsafe.Pointer(&cKeys[0])),
		(*C.size_t)(unsafe.Pointer(&cKeySizes[0])),
		(**C.char)(unsafe.Pointer(&cVals[0])),
		(*C.size_t)(unsafe.Pointer(&cValSizes[0])),
		(**C.char)(unsafe.Pointer(&cErrs[0])),
	)

	vals := make([][]byte, n)
	errs := make([]error, n)
	for i := 0; i < n; i++ {
		if cErrs[i] != nil {
			errs[i] = errFromC(cErrs[i], fmt.Sprintf("multi get cf key %d", i))
			continue
		}
		if cVals[i] != nil {
			vals[i] = C.GoBytes(unsafe.Pointer(cVals[i]), C.int(cValSizes[i]))
			C.rocksdb_free(unsafe.Pointer(cVals[i]))
		}
	}
	return vals, errs
}

func (db *DB) PutList(opts *WriteOptions, keys, values [][]byte) error {
	if db == nil || db.ptr == nil {
		return errClosedDB("put list")
	}
	n := len(keys)
	if n == 0 {
		return nil
	}
	if len(values) != n {
		return fmt.Errorf("put list: keys count %d != values count %d", n, len(values))
	}
	writeOpts := ensureWriteOptions(opts)

	cKeys := make([]*C.char, n)
	cKeySizes := make([]C.size_t, n)
	cVals := make([]*C.char, n)
	cValSizes := make([]C.size_t, n)
	for i := 0; i < n; i++ {
		cKeys[i] = bytesPtr(keys[i])
		cKeySizes[i] = C.size_t(len(keys[i]))
		cVals[i] = bytesPtr(values[i])
		cValSizes[i] = C.size_t(len(values[i]))
	}

	var errptr *C.char
	C.rockskit_put_list(
		db.ptr,
		writeOpts.ptr,
		C.size_t(n),
		(**C.char)(unsafe.Pointer(&cKeys[0])),
		(*C.size_t)(unsafe.Pointer(&cKeySizes[0])),
		(**C.char)(unsafe.Pointer(&cVals[0])),
		(*C.size_t)(unsafe.Pointer(&cValSizes[0])),
		&errptr,
	)
	return errFromC(errptr, "put list")
}

func (db *DB) PutListCF(handle *ColumnFamilyHandle, opts *WriteOptions, keys, values [][]byte) error {
	if db == nil || db.ptr == nil {
		return errClosedDB("put list cf")
	}
	if handle == nil || handle.ptr == nil {
		return fmt.Errorf("put list cf: column family handle is nil")
	}
	n := len(keys)
	if n == 0 {
		return nil
	}
	if len(values) != n {
		return fmt.Errorf("put list cf: keys count %d != values count %d", n, len(values))
	}
	writeOpts := ensureWriteOptions(opts)

	cKeys := make([]*C.char, n)
	cKeySizes := make([]C.size_t, n)
	cVals := make([]*C.char, n)
	cValSizes := make([]C.size_t, n)
	for i := 0; i < n; i++ {
		cKeys[i] = bytesPtr(keys[i])
		cKeySizes[i] = C.size_t(len(keys[i]))
		cVals[i] = bytesPtr(values[i])
		cValSizes[i] = C.size_t(len(values[i]))
	}

	var errptr *C.char
	C.rockskit_put_list_cf(
		db.ptr,
		writeOpts.ptr,
		handle.ptr,
		C.size_t(n),
		(**C.char)(unsafe.Pointer(&cKeys[0])),
		(*C.size_t)(unsafe.Pointer(&cKeySizes[0])),
		(**C.char)(unsafe.Pointer(&cVals[0])),
		(*C.size_t)(unsafe.Pointer(&cValSizes[0])),
		&errptr,
	)
	return errFromC(errptr, "put list cf")
}

func (db *DB) DeleteList(opts *WriteOptions, keys [][]byte) error {
	if db == nil || db.ptr == nil {
		return errClosedDB("delete list")
	}
	n := len(keys)
	if n == 0 {
		return nil
	}
	writeOpts := ensureWriteOptions(opts)

	cKeys := make([]*C.char, n)
	cKeySizes := make([]C.size_t, n)
	for i, k := range keys {
		cKeys[i] = bytesPtr(k)
		cKeySizes[i] = C.size_t(len(k))
	}

	var errptr *C.char
	C.rockskit_delete_list(
		db.ptr,
		writeOpts.ptr,
		C.size_t(n),
		(**C.char)(unsafe.Pointer(&cKeys[0])),
		(*C.size_t)(unsafe.Pointer(&cKeySizes[0])),
		&errptr,
	)
	return errFromC(errptr, "delete list")
}

func (db *DB) DeleteListCF(handle *ColumnFamilyHandle, opts *WriteOptions, keys [][]byte) error {
	if db == nil || db.ptr == nil {
		return errClosedDB("delete list cf")
	}
	if handle == nil || handle.ptr == nil {
		return fmt.Errorf("delete list cf: column family handle is nil")
	}
	n := len(keys)
	if n == 0 {
		return nil
	}
	writeOpts := ensureWriteOptions(opts)

	cKeys := make([]*C.char, n)
	cKeySizes := make([]C.size_t, n)
	for i, k := range keys {
		cKeys[i] = bytesPtr(k)
		cKeySizes[i] = C.size_t(len(k))
	}

	var errptr *C.char
	C.rockskit_delete_list_cf(
		db.ptr,
		writeOpts.ptr,
		handle.ptr,
		C.size_t(n),
		(**C.char)(unsafe.Pointer(&cKeys[0])),
		(*C.size_t)(unsafe.Pointer(&cKeySizes[0])),
		&errptr,
	)
	return errFromC(errptr, "delete list cf")
}

func OpenWithColumnFamilies(path string, cfg *Config, createIfMissing bool) (*DB, map[string]*ColumnFamilyHandle, error) {
	artifacts, err := buildArtifacts(cfg, createIfMissing, false)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		if err != nil {
			artifacts.Close()
		}
	}()

	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	var cfCount C.size_t
	var listErrptr *C.char
	cfNamesC := C.rocksdb_list_column_families(artifacts.options.ptr, cPath, &cfCount, &listErrptr)
	if listErrptr != nil {
		return nil, nil, errFromC(listErrptr, "list column families")
	}
	if cfCount == 0 {
		C.rocksdb_list_column_families_destroy(cfNamesC, cfCount)
		return nil, nil, fmt.Errorf("list column families: no column families found")
	}

	names := make([]string, int(cfCount))
	cNamePtrs := make([]*C.char, int(cfCount))
	cfNamesSlice := unsafe.Slice(cfNamesC, int(cfCount))
	for i := 0; i < int(cfCount); i++ {
		names[i] = C.GoString(cfNamesSlice[i])
		cNamePtrs[i] = C.CString(names[i])
	}
	C.rocksdb_list_column_families_destroy(cfNamesC, cfCount)
	for i := range cNamePtrs {
		defer C.free(unsafe.Pointer(cNamePtrs[i]))
	}

	cfOpts := make([]*C.rocksdb_options_t, int(cfCount))
	for i := range cfOpts {
		cfOpts[i] = artifacts.options.ptr
	}

	cfHandles := make([]*C.rocksdb_column_family_handle_t, int(cfCount))

	var openErrptr *C.char
	dbPtr := C.rocksdb_open_column_families(
		artifacts.options.ptr,
		cPath,
		C.int(cfCount),
		(**C.char)(unsafe.Pointer(&cNamePtrs[0])),
		(**C.rocksdb_options_t)(unsafe.Pointer(&cfOpts[0])),
		(**C.rocksdb_column_family_handle_t)(unsafe.Pointer(&cfHandles[0])),
		&openErrptr,
	)
	if err := errFromC(openErrptr, "open column families"); err != nil {
		return nil, nil, err
	}

	cfMap := make(map[string]*ColumnFamilyHandle, int(cfCount))
	for i, name := range names {
		cfMap[name] = &ColumnFamilyHandle{ptr: cfHandles[i]}
	}

	return &DB{
		ptr:       dbPtr,
		artifacts: artifacts,
	}, cfMap, nil
}
