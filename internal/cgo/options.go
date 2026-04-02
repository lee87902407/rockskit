package cgo

/*
#cgo CFLAGS: -I${SRCDIR}/../../native/darwin_arm64/include -I${SRCDIR}/../../native/linux_amd64/include -I${SRCDIR}/../../native/linux_arm64/include
#include <stdlib.h>
#include <rocksdb/c.h>
*/
import "C"

type Options struct {
	ptr *C.rocksdb_options_t
}

// NewOptions 创建一个原生 rocksdb_options_t 包装对象。
//
// 示例：
//
//	opts := NewOptions()
//	defer opts.Close()
//	opts.SetCreateIfMissing(true)
//
// 该对象用于在打开数据库前累计数据库级别的配置项。
func NewOptions() *Options {
	return &Options{ptr: C.rocksdb_options_create()}
}

// SetWriteBufferSize 用于设置 memtable 在 flush 前的目标大小。
//
// 示例：调用 SetWriteBufferSize(64 * 1024 * 1024) 表示设置 64MiB 的 memtable。
func (o *Options) SetWriteBufferSize(size uint64) {
	if o == nil || o.ptr == nil || size == 0 {
		return
	}
	C.rocksdb_options_set_write_buffer_size(o.ptr, C.size_t(size))
}

// SetWriteBufferNumber 用于设置写入阻塞前最多允许存在多少个 memtable。
// 示例：调用 SetWriteBufferNumber(6) 表示最多保留 6 个写缓冲。
func (o *Options) SetWriteBufferNumber(n int) {
	if o == nil || o.ptr == nil || n <= 0 {
		return
	}
	C.rocksdb_options_set_max_write_buffer_number(o.ptr, C.int(n))
}

// SetWriteBufferNumberToMerge 用于设置读路径上会合并多少个 immutable memtable。
// 示例：调用 SetWriteBufferNumberToMerge(1) 表示采用更保守的读放大策略。
func (o *Options) SetWriteBufferNumberToMerge(n int) {
	if o == nil || o.ptr == nil || n <= 0 {
		return
	}
	C.rocksdb_options_set_min_write_buffer_number_to_merge(o.ptr, C.int(n))
}

// SetWalSizeLimitMb 用于控制 WAL 轮转阈值，单位为 MB。
// 示例：调用 SetWalSizeLimitMb(512) 表示 WAL 大约达到 512MB 后开始轮转。
func (o *Options) SetWalSizeLimitMb(size uint64) {
	if o == nil || o.ptr == nil {
		return
	}
	C.rocksdb_options_set_WAL_size_limit_MB(o.ptr, C.uint64_t(size))
}

// SetMaxTotalWalSize 用于限制整个数据库允许保留的 WAL 总量。
// 示例：调用 SetMaxTotalWalSize(1 << 30) 表示 WAL 总量上限约为 1GiB。
func (o *Options) SetMaxTotalWalSize(size uint64) {
	if o == nil || o.ptr == nil {
		return
	}
	C.rocksdb_options_set_max_total_wal_size(o.ptr, C.uint64_t(size))
}

// SetLevel0FileNumCompactionTrigger 用于设置 L0 文件达到多少个后触发 compaction。
// 示例：调用 SetLevel0FileNumCompactionTrigger(4) 表示在突发写入下更早开始 compaction。
func (o *Options) SetLevel0FileNumCompactionTrigger(n int) {
	if o == nil || o.ptr == nil || n <= 0 {
		return
	}
	C.rocksdb_options_set_level0_file_num_compaction_trigger(o.ptr, C.int(n))
}

// SetLevel0SlowDownWritesTrigger 用于设置在硬停止前何时开始减速写入。
// 示例：调用 SetLevel0SlowDownWritesTrigger(20) 表示在 L0 文件过多前先减速。
func (o *Options) SetLevel0SlowDownWritesTrigger(n int) {
	if o == nil || o.ptr == nil || n <= 0 {
		return
	}
	C.rocksdb_options_set_level0_slowdown_writes_trigger(o.ptr, C.int(n))
}

// SetLevel0StopWritesTrigger 用于设置 L0 文件过多时何时硬停止写入。
// 示例：调用 SetLevel0StopWritesTrigger(24) 可以防止读放大失控增长。
func (o *Options) SetLevel0StopWritesTrigger(n int) {
	if o == nil || o.ptr == nil || n <= 0 {
		return
	}
	C.rocksdb_options_set_level0_stop_writes_trigger(o.ptr, C.int(n))
}

// SetMaxBytesForLevelBase 用于设置第一个 leveled compaction 层的目标大小。
// 示例：调用 SetMaxBytesForLevelBase(256 << 20) 表示设置 256MiB 的 level base。
func (o *Options) SetMaxBytesForLevelBase(size uint64) {
	if o == nil || o.ptr == nil {
		return
	}
	C.rocksdb_options_set_max_bytes_for_level_base(o.ptr, C.uint64_t(size))
}

// SetTargetFileSizeBase 用于设置第一层 SST 文件的目标大小。
// 示例：调用 SetTargetFileSizeBase(64 << 20) 表示 SST 目标大小为 64MiB。
func (o *Options) SetTargetFileSizeBase(size uint64) {
	if o == nil || o.ptr == nil {
		return
	}
	C.rocksdb_options_set_target_file_size_base(o.ptr, C.uint64_t(size))
}

// SetMaxBackgroundJobs 用于控制 RocksDB 可并发执行多少个 flush/compaction 任务。
// 示例：调用 SetMaxBackgroundJobs(4) 表示最多允许 4 个后台任务。
func (o *Options) SetMaxBackgroundJobs(n int) {
	if o == nil || o.ptr == nil || n <= 0 {
		return
	}
	C.rocksdb_options_set_max_background_jobs(o.ptr, C.int(n))
}

// SetMaxSubCompactions 用于允许单个 compaction 拆分到多个 worker。
// 示例：调用 SetMaxSubCompactions(2) 表示一个 compaction 最多使用两个 worker。
func (o *Options) SetMaxSubCompactions(n uint32) {
	if o == nil || o.ptr == nil || n == 0 {
		return
	}
	C.rocksdb_options_set_max_subcompactions(o.ptr, C.uint32_t(n))
}

// SetMaxLogFileSize 用于限制单个 INFO 日志文件大小。
// 示例：调用 SetMaxLogFileSize(1 << 30) 表示日志接近 1GiB 时轮转。
func (o *Options) SetMaxLogFileSize(size uint64) {
	if o == nil || o.ptr == nil {
		return
	}
	C.rocksdb_options_set_max_log_file_size(o.ptr, C.size_t(size))
}

// SetKeepLogFileNum 用于控制磁盘上保留多少个 INFO 日志。
// 示例：调用 SetKeepLogFileNum(10) 表示保留最新的 10 个日志文件。
func (o *Options) SetKeepLogFileNum(n uint) {
	if o == nil || o.ptr == nil || n == 0 {
		return
	}
	C.rocksdb_options_set_keep_log_file_num(o.ptr, C.size_t(n))
}

// SetRecycleLogFileNum 用于允许 RocksDB 复用旧 WAL 文件而不是创建新文件。
// 示例：调用 SetRecycleLogFileNum(4) 可以减少持续写入下的文件创建抖动。
func (o *Options) SetRecycleLogFileNum(n uint) {
	if o == nil || o.ptr == nil {
		return
	}
	C.rocksdb_options_set_recycle_log_file_num(o.ptr, C.size_t(n))
}

// SetMaxCompactionBytes 用于限制 compaction 的最大输入大小。
// 示例：调用 SetMaxCompactionBytes(0) 表示交给 RocksDB 默认策略处理。
func (o *Options) SetMaxCompactionBytes(size uint64) {
	if o == nil || o.ptr == nil {
		return
	}
	C.rocksdb_options_set_max_compaction_bytes(o.ptr, C.uint64_t(size))
}

// SetPeriodicCompactionSeconds 用于设置 RocksDB 周期性重写旧文件的时间间隔。
// 示例：调用 SetPeriodicCompactionSeconds(86400) 表示每天触发一次周期性 compaction。
func (o *Options) SetPeriodicCompactionSeconds(n int) {
	if o == nil || o.ptr == nil || n < 0 {
		return
	}
	C.rocksdb_options_set_periodic_compaction_seconds(o.ptr, C.uint64_t(n))
}

// SetCompressionParallelThreads 用于设置压缩阶段的辅助线程数量。
// 示例：调用 SetCompressionParallelThreads(4) 可以提升大机器上的压缩吞吐。
func (o *Options) SetCompressionParallelThreads(n int) {
	if o == nil || o.ptr == nil || n <= 0 {
		return
	}
	C.rocksdb_options_set_compression_options_parallel_threads(o.ptr, C.int(n))
}

// SetEnableBlobFiles 用于控制是否开启 BlobDB 风格的大 value 分离。
// 示例：调用 SetEnableBlobFiles(true) 表示大 value 会写入 blob 文件而不是 SST。
func (o *Options) SetEnableBlobFiles(enabled bool) {
	if o == nil || o.ptr == nil {
		return
	}
	C.rocksdb_options_set_enable_blob_files(o.ptr, boolToChar(enabled))
}

// SetMinBlobSize 用于设置最小多大的 value 会进入 blob 文件。
// 示例：调用 SetMinBlobSize(4 << 10) 表示大小 >= 4KiB 的 value 会被当作 blob。
func (o *Options) SetMinBlobSize(size uint64) {
	if o == nil || o.ptr == nil || size == 0 {
		return
	}
	C.rocksdb_options_set_min_blob_size(o.ptr, C.uint64_t(size))
}

// SetBlobFileSize 用于控制 blob 文件的目标大小。
// 示例：调用 SetBlobFileSize(256 << 20) 表示目标 blob 文件大小为 256MiB。
func (o *Options) SetBlobFileSize(size uint64) {
	if o == nil || o.ptr == nil || size == 0 {
		return
	}
	C.rocksdb_options_set_blob_file_size(o.ptr, C.uint64_t(size))
}

// SetEnableBlobGC 用于控制是否开启过期 blob value 的后台回收。
// 示例：调用 SetEnableBlobGC(true) 表示允许 RocksDB 回收过期 blob。
func (o *Options) SetEnableBlobGC(enabled bool) {
	if o == nil || o.ptr == nil {
		return
	}
	C.rocksdb_options_set_enable_blob_gc(o.ptr, boolToChar(enabled))
}

// SetBlobCompressionType 用于设置 blob 文件压缩算法。
// 示例：调用 SetBlobCompressionType(4) 表示 blob 文件使用 LZ4 压缩。
func (o *Options) SetBlobCompressionType(code int) {
	if o == nil || o.ptr == nil {
		return
	}
	C.rocksdb_options_set_blob_compression_type(o.ptr, C.int(code))
}

// SetBlobGCAgeCutoff 用于控制 blob 文件多“旧”后才会被 GC 考虑。
// 示例：调用 SetBlobGCAgeCutoff(0.25) 表示优先考虑较老的 blob 文件。
func (o *Options) SetBlobGCAgeCutoff(value float64) {
	if o == nil || o.ptr == nil || value <= 0 {
		return
	}
	C.rocksdb_options_set_blob_gc_age_cutoff(o.ptr, C.double(value))
}

// SetCreateIfMissing 用于控制数据库目录不存在时是否自动创建。
// 示例：Create() 会使用 SetCreateIfMissing(true)。
func (o *Options) SetCreateIfMissing(enabled bool) {
	if o == nil || o.ptr == nil {
		return
	}
	C.rocksdb_options_set_create_if_missing(o.ptr, boolToChar(enabled))
}

// SetErrorIfExists 用于控制数据库已存在时是否直接报错。
// 示例：调用 SetErrorIfExists(true) 会让 Create() 在路径已存在时报错。
func (o *Options) SetErrorIfExists(enabled bool) {
	if o == nil || o.ptr == nil {
		return
	}
	C.rocksdb_options_set_error_if_exists(o.ptr, boolToChar(enabled))
}

// SetBlockBasedTableFactory 用于将 block-based table 配置绑定到主 options 对象。
// 示例：调用 SetBlockBasedTableFactory(blockOpts) 会把 cache/filter/index 配置应用到 SST 表。
func (o *Options) SetBlockBasedTableFactory(blockOpts *BlockBasedOptions) {
	if o == nil || o.ptr == nil || blockOpts == nil || blockOpts.ptr == nil {
		return
	}
	C.rocksdb_options_set_block_based_table_factory(o.ptr, blockOpts.ptr)
}

// SetRateLimiter 用于给后台 IO 挂上 rate limiter。
// 示例：调用 SetRateLimiter(limiter) 可以限制共享磁盘上的 compaction 吞吐。
func (o *Options) SetRateLimiter(rateLimiter *RateLimiter) {
	if o == nil || o.ptr == nil || rateLimiter == nil || rateLimiter.ptr == nil {
		return
	}
	C.rocksdb_options_set_ratelimiter(o.ptr, rateLimiter.ptr)
}

// SetCompression 用于设置普通 SST 文件的压缩算法。
// 示例：调用 SetCompression(4) 表示普通 SST 使用 LZ4。
func (o *Options) SetCompression(code int) {
	if o == nil || o.ptr == nil {
		return
	}
	C.rocksdb_options_set_compression(o.ptr, C.int(code))
}

// SetBottommostCompression 用于设置最底层冷数据的压缩算法。
// 示例：调用 SetBottommostCompression(7) 表示最底层使用 ZSTD。
func (o *Options) SetBottommostCompression(code int) {
	if o == nil || o.ptr == nil {
		return
	}
	C.rocksdb_options_set_bottommost_compression(o.ptr, C.int(code))
}

// SetWriteBufferSizeToMaintain 用于保留历史 memtable 数据，以支持冲突检查。
// 示例：调用 SetWriteBufferSizeToMaintain(64 << 20) 表示额外保留 64MiB 历史 memtable 数据。
func (o *Options) SetWriteBufferSizeToMaintain(size uint64) {
	if o == nil || o.ptr == nil {
		return
	}
	C.rocksdb_options_set_max_write_buffer_size_to_maintain(o.ptr, C.int64_t(size))
}

// AddCompactOnDeletionCollectorFactory 用于配置基于删除数量的 compaction 触发策略。
// 示例： AddCompactOnDeletionCollectorFactory(100000, 5000) may trigger compaction
// 当 100000 个 key 的滑动窗口中删除数量超过 5000 时会生效。
func (o *Options) AddCompactOnDeletionCollectorFactory(windowSize, numDelsTrigger uint) {
	if o == nil || o.ptr == nil || windowSize == 0 || numDelsTrigger == 0 {
		return
	}
	C.rocksdb_options_add_compact_on_deletion_collector_factory(o.ptr, C.size_t(windowSize), C.size_t(numDelsTrigger))
}

// Close 用于释放当前持有的 rocksdb_options_t 句柄。
// 示例：NewOptions() 之后使用 defer opts.Close() 可以避免原生 options 内存泄漏。
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
