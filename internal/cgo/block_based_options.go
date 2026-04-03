package cgo

/*
#cgo CFLAGS: -I${SRCDIR}/../../native/darwin_arm64/include -I${SRCDIR}/../../native/linux_amd64/include -I${SRCDIR}/../../native/linux_arm64/include
#include <rocksdb/c.h>
*/
import "C"

type FilterPolicy struct {
	ptr *C.rocksdb_filterpolicy_t
}

func NewBloomFilterPolicy(bitsPerKey float64) *FilterPolicy {
	return &FilterPolicy{ptr: C.rocksdb_filterpolicy_create_bloom(C.double(bitsPerKey))}
}

func (p *FilterPolicy) Close() {
	if p == nil || p.ptr == nil {
		return
	}
	C.rocksdb_filterpolicy_destroy(p.ptr)
	p.ptr = nil
}

type BlockBasedOptions struct {
	ptr    *C.rocksdb_block_based_table_options_t
	cache  *LRUCache
	filter *FilterPolicy
}

func NewBlockBasedOptions() *BlockBasedOptions {
	return &BlockBasedOptions{ptr: C.rocksdb_block_based_options_create()}
}

func (o *BlockBasedOptions) SetBlockSize(size uint64) {
	if o == nil || o.ptr == nil {
		return
	}
	C.rocksdb_block_based_options_set_block_size(o.ptr, C.size_t(size))
}

func (o *BlockBasedOptions) SetNoBlockCache(enabled bool) {
	if o == nil || o.ptr == nil {
		return
	}
	C.rocksdb_block_based_options_set_no_block_cache(o.ptr, boolToChar(enabled))
}

func (o *BlockBasedOptions) SetBlockCache(cache *LRUCache) {
	if o == nil || o.ptr == nil || cache == nil || cache.ptr == nil {
		return
	}
	o.cache = cache
	C.rocksdb_block_based_options_set_block_cache(o.ptr, cache.ptr)
}

func (o *BlockBasedOptions) SetFilterPolicy(policy *FilterPolicy) {
	if o == nil || o.ptr == nil || policy == nil || policy.ptr == nil {
		return
	}
	o.filter = policy
	C.rocksdb_block_based_options_set_filter_policy(o.ptr, policy.ptr)
}

func (o *BlockBasedOptions) SetPartitionFilters(enabled bool) {
	if o == nil || o.ptr == nil {
		return
	}
	C.rocksdb_block_based_options_set_partition_filters(o.ptr, boolToChar(enabled))
}

func (o *BlockBasedOptions) SetOptimizeFiltersForMemory(enabled bool) {
	if o == nil || o.ptr == nil {
		return
	}
	C.rocksdb_block_based_options_set_optimize_filters_for_memory(o.ptr, boolToChar(enabled))
}

func (o *BlockBasedOptions) SetWholeKeyFiltering(enabled bool) {
	if o == nil || o.ptr == nil {
		return
	}
	C.rocksdb_block_based_options_set_whole_key_filtering(o.ptr, boolToChar(enabled))
}

func (o *BlockBasedOptions) SetIndexType(indexType int) {
	if o == nil || o.ptr == nil {
		return
	}
	C.rocksdb_block_based_options_set_index_type(o.ptr, C.int(indexType))
}

func (o *BlockBasedOptions) SetDataBlockIndexType(indexType int) {
	if o == nil || o.ptr == nil {
		return
	}
	C.rocksdb_block_based_options_set_data_block_index_type(o.ptr, C.int(indexType))
}

func (o *BlockBasedOptions) SetDataBlockHashRatio(ratio float64) {
	if o == nil || o.ptr == nil || ratio <= 0 {
		return
	}
	C.rocksdb_block_based_options_set_data_block_hash_ratio(o.ptr, C.double(ratio))
}

func (o *BlockBasedOptions) SetCacheIndexAndFilterBlocks(enabled bool) {
	if o == nil || o.ptr == nil {
		return
	}
	C.rocksdb_block_based_options_set_cache_index_and_filter_blocks(o.ptr, boolToChar(enabled))
}

func (o *BlockBasedOptions) SetCacheIndexAndFilterBlocksWithHighPriority(enabled bool) {
	if o == nil || o.ptr == nil {
		return
	}
	C.rocksdb_block_based_options_set_cache_index_and_filter_blocks_with_high_priority(o.ptr, boolToChar(enabled))
}

func (o *BlockBasedOptions) SetPinL0FilterAndIndexBlocksInCache(enabled bool) {
	if o == nil || o.ptr == nil {
		return
	}
	C.rocksdb_block_based_options_set_pin_l0_filter_and_index_blocks_in_cache(o.ptr, boolToChar(enabled))
}

func (o *BlockBasedOptions) SetPinTopLevelIndexAndFilter(enabled bool) {
	if o == nil || o.ptr == nil {
		return
	}
	C.rocksdb_block_based_options_set_pin_top_level_index_and_filter(o.ptr, boolToChar(enabled))
}

func (o *BlockBasedOptions) SetTopLevelIndexPinningTier(tier int) {
	if o == nil || o.ptr == nil {
		return
	}
	C.rocksdb_block_based_options_set_top_level_index_pinning_tier(o.ptr, C.int(tier))
}

func (o *BlockBasedOptions) SetPartitionPinningTier(tier int) {
	if o == nil || o.ptr == nil {
		return
	}
	C.rocksdb_block_based_options_set_partition_pinning_tier(o.ptr, C.int(tier))
}

func (o *BlockBasedOptions) SetUnpartitionedPinningTier(tier int) {
	if o == nil || o.ptr == nil {
		return
	}
	C.rocksdb_block_based_options_set_unpartitioned_pinning_tier(o.ptr, C.int(tier))
}

func (o *BlockBasedOptions) Close() {
	if o == nil || o.ptr == nil {
		return
	}
	C.rocksdb_block_based_options_destroy(o.ptr)
	o.ptr = nil
	if o.filter != nil {
		o.filter.ptr = nil
	}
	o.cache = nil
	o.filter = nil
}
