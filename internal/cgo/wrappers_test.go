package cgo

import "testing"

func TestWrapperCloseMethodsAreIdempotent(t *testing.T) {
	var (
		allocator    *MemoryAllocator
		cacheOptions *LRUCacheOptions
		options      *Options
		readOptions  *ReadOptions
		writeOptions *WriteOptions
		filterPolicy *FilterPolicy
		blockOptions *BlockBasedOptions
		cache        *LRUCache
		rateLimiter  *RateLimiter
		db           *DB
		pinnable     *PinnableSlice
	)

	allocator.Close()
	cacheOptions.Close()
	options.Close()
	readOptions.Close()
	writeOptions.Close()
	filterPolicy.Close()
	blockOptions.Close()
	cache.Close()
	rateLimiter.Close()
	db.Close()
	pinnable.Close()
}
