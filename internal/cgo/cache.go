package cgo

/*
#include <rocksdb/c.h>
*/
import "C"

import "fmt"
import "runtime"

type MemoryAllocator struct {
	ptr *C.rocksdb_memory_allocator_t
}

func NewJemallocAllocator() (*MemoryAllocator, error) {
	var errptr *C.char
	alloc := C.rocksdb_jemalloc_nodump_allocator_create(&errptr)
	if err := errFromCString(errptr, "create jemalloc allocator"); err != nil {
		return nil, err
	}
	if alloc == nil {
		return nil, fmt.Errorf("create jemalloc allocator returned nil")
	}
	return &MemoryAllocator{ptr: alloc}, nil
}

func (a *MemoryAllocator) Close() {
	if a == nil || a.ptr == nil {
		return
	}
	C.rocksdb_memory_allocator_destroy(a.ptr)
	a.ptr = nil
}

type LRUCacheOptions struct {
	ptr *C.rocksdb_lru_cache_options_t
}

func NewLRUCacheOptions() *LRUCacheOptions {
	return &LRUCacheOptions{ptr: C.rocksdb_lru_cache_options_create()}
}

func (o *LRUCacheOptions) SetCapacity(capacity uint64) {
	if o == nil || o.ptr == nil || capacity == 0 {
		return
	}
	C.rocksdb_lru_cache_options_set_capacity(o.ptr, C.size_t(capacity))
}

func (o *LRUCacheOptions) SetNumShardBits(bits int) {
	if o == nil || o.ptr == nil {
		return
	}
	C.rocksdb_lru_cache_options_set_num_shard_bits(o.ptr, C.int(bits))
}

func (o *LRUCacheOptions) SetMemoryAllocator(allocator *MemoryAllocator) {
	if o == nil || o.ptr == nil || allocator == nil || allocator.ptr == nil {
		return
	}
	C.rocksdb_lru_cache_options_set_memory_allocator(o.ptr, allocator.ptr)
}

func (o *LRUCacheOptions) Close() {
	if o == nil || o.ptr == nil {
		return
	}
	C.rocksdb_lru_cache_options_destroy(o.ptr)
	o.ptr = nil
}

type LRUCache struct {
	ptr       *C.rocksdb_cache_t
	allocator *MemoryAllocator
	options   *LRUCacheOptions
}

func NewLRUCache(capacity uint64, shardBits int) (*LRUCache, error) {
	options := NewLRUCacheOptions()
	options.SetCapacity(capacity)
	options.SetNumShardBits(shardBits)

	var allocator *MemoryAllocator
	if runtime.GOOS == "linux" {
		var err error
		allocator, err = NewJemallocAllocator()
		if err != nil {
			panic(err)
		}
		options.SetMemoryAllocator(allocator)
	}

	cache := C.rocksdb_cache_create_lru_opts(options.ptr)
	if cache == nil {
		options.Close()
		if allocator != nil {
			allocator.Close()
		}
		return nil, fmt.Errorf("create lru cache returned nil")
	}
	return &LRUCache{ptr: cache, allocator: allocator, options: options}, nil
}

func (c *LRUCache) Close() {
	if c == nil || c.ptr == nil {
		return
	}
	C.rocksdb_cache_destroy(c.ptr)
	c.ptr = nil
	if c.options != nil {
		c.options.Close()
		c.options = nil
	}
	if c.allocator != nil {
		c.allocator.Close()
		c.allocator = nil
	}
}
