package cgo

/*
#include <rocksdb/c.h>
*/
import "C"

type RawOptions struct {
	WriteBufferSize                uint64
	Level0FileNumCompactionTrigger int
	TargetFileSizeBase             uint64
	MaxBytesForLevelBase           uint64
	CreateIfMissing                bool
	ErrorIfExists                  bool
	Compression                    int
}

type Options struct {
	ptr *C.rocksdb_options_t
}

func NewOptions(cfg RawOptions) *Options {
	ptr := C.rocksdb_options_create()
	C.rocksdb_options_set_create_if_missing(ptr, boolToChar(cfg.CreateIfMissing))
	C.rocksdb_options_set_error_if_exists(ptr, boolToChar(cfg.ErrorIfExists))
	C.rocksdb_options_set_write_buffer_size(ptr, C.size_t(cfg.WriteBufferSize))
	C.rocksdb_options_set_level0_file_num_compaction_trigger(ptr, C.int(cfg.Level0FileNumCompactionTrigger))
	C.rocksdb_options_set_target_file_size_base(ptr, C.uint64_t(cfg.TargetFileSizeBase))
	C.rocksdb_options_set_max_bytes_for_level_base(ptr, C.uint64_t(cfg.MaxBytesForLevelBase))
	C.rocksdb_options_set_compression(ptr, C.int(cfg.Compression))
	return &Options{ptr: ptr}
}

func (o *Options) Close() {
	if o == nil || o.ptr == nil {
		return
	}
	C.rocksdb_options_destroy(o.ptr)
	o.ptr = nil
}

func boolToChar(value bool) C.uchar {
	if value {
		return 1
	}
	return 0
}
