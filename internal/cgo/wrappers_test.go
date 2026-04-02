package cgo

import "testing"

func TestWrapperCloseMethodsAreIdempotent(t *testing.T) {
	var (
		allocator    *MemoryAllocator
		cacheOptions *LRUCacheOptions
		options      *Options
		filterPolicy *FilterPolicy
		blockOptions *BlockBasedOptions
		cache        *LRUCache
		rateLimiter  *RateLimiter
		db           *DB
	)

	allocator.Close()
	cacheOptions.Close()
	options.Close()
	filterPolicy.Close()
	blockOptions.Close()
	cache.Close()
	rateLimiter.Close()
	db.Close()
}
