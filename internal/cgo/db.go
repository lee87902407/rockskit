package cgo

/*
#include <stdlib.h>
#include <rocksdb/c.h>
*/
import "C"

import "unsafe"

type DB struct {
	ptr *C.rocksdb_t
}

func Open(path string, opts *Options) (*DB, error) {
	if opts == nil || opts.ptr == nil {
		return nil, errFromC(nil, "open rocksdb with nil options")
	}

	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	var errptr *C.char
	db := C.rocksdb_open(opts.ptr, cPath, &errptr)
	if err := errFromC(errptr, "open rocksdb"); err != nil {
		return nil, err
	}
	return &DB{ptr: db}, nil
}

func (db *DB) Close() {
	if db == nil || db.ptr == nil {
		return
	}
	C.rocksdb_close(db.ptr)
	db.ptr = nil
}
