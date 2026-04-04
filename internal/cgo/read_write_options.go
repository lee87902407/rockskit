package cgo

/*
#include <rocksdb/c.h>
*/
import "C"

type ReadOptions struct {
	ptr *C.rocksdb_readoptions_t
}

func NewReadOptions() *ReadOptions {
	return &ReadOptions{ptr: C.rocksdb_readoptions_create()}
}

func (o *ReadOptions) Close() {
	if o == nil || o.ptr == nil {
		return
	}
	C.rocksdb_readoptions_destroy(o.ptr)
	o.ptr = nil
}

func (o *ReadOptions) SetFillCache(enabled bool) {
	if o == nil || o.ptr == nil {
		return
	}
	C.rocksdb_readoptions_set_fill_cache(o.ptr, boolToChar(enabled))
}

func (o *ReadOptions) SetVerifyChecksums(enabled bool) {
	if o == nil || o.ptr == nil {
		return
	}
	C.rocksdb_readoptions_set_verify_checksums(o.ptr, boolToChar(enabled))
}

func (o *ReadOptions) SetPinData(enabled bool) {
	if o == nil || o.ptr == nil {
		return
	}
	C.rocksdb_readoptions_set_pin_data(o.ptr, boolToChar(enabled))
}

func ensureReadOptions(opts *ReadOptions) *ReadOptions {
	if opts != nil && opts.ptr != nil {
		return opts
	}
	return defaultReadOptions
}

type WriteOptions struct {
	ptr *C.rocksdb_writeoptions_t
}

func NewWriteOptions() *WriteOptions {
	return &WriteOptions{ptr: C.rocksdb_writeoptions_create()}
}

func (o *WriteOptions) Close() {
	if o == nil || o.ptr == nil {
		return
	}
	C.rocksdb_writeoptions_destroy(o.ptr)
	o.ptr = nil
}

func (o *WriteOptions) DisableWAL(disabled bool) {
	if o == nil || o.ptr == nil {
		return
	}
	C.rocksdb_writeoptions_disable_WAL(o.ptr, C.int(boolToInt(disabled)))
}

func (o *WriteOptions) SetSync(enabled bool) {
	if o == nil || o.ptr == nil {
		return
	}
	C.rocksdb_writeoptions_set_sync(o.ptr, boolToChar(enabled))
}

func ensureWriteOptions(opts *WriteOptions) *WriteOptions {
	if opts != nil && opts.ptr != nil {
		return opts
	}
	return defaultWriteOptions
}

var (
	defaultReadOptions  = NewReadOptions()
	defaultWriteOptions = NewWriteOptions()
)

func boolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}
