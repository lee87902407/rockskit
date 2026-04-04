package cgo

/*
#include <stdlib.h>
#include <rocksdb/c.h>
*/
import "C"

import (
	"fmt"
	"unsafe"
)

type PinnableSlice struct {
	ptr   *C.rocksdb_pinnableslice_t
	value []byte
}

func (s *PinnableSlice) Close() {
	if s == nil || s.ptr == nil {
		return
	}
	C.rocksdb_pinnableslice_destroy(s.ptr)
	s.ptr = nil
	s.value = nil
}

func errClosedDB(action string) error {
	return fmt.Errorf("%s on closed rocksdb", action)
}

func newPinnableSlice(ptr *C.rocksdb_pinnableslice_t, data *C.char, length C.size_t) *PinnableSlice {
	if ptr == nil {
		return nil
	}
	value := []byte{}
	if data != nil && length > 0 {
		value = unsafe.Slice((*byte)(unsafe.Pointer(data)), int(length))
	}
	return &PinnableSlice{ptr: ptr, value: value}
}
