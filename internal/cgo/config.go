package cgo

// Config 用于描述 rockskit 如何将 Go 侧配置映射到 RocksDB 选项。
//
// 示例：
//
//	cfg := DefaultConfig()
//	cfg.WriteBufferSize = "128MiB"
//	cfg.LRUSize = "512MiB"
//	cfg.EnableBlobFiles = true
//
// 空字符串和数值零值会被视为“保持 RocksDB 默认行为不变”，
// 由选项构建流程负责跳过对应设置。
type Config struct {
	// WriteBufferSize 控制 memtable 刷盘前的目标大小。
	// 数值越大通常写入吞吐越高，但内存占用也越高。
	// 示例：普通场景可设为 "64MiB"，资源受限场景可设为 "16MiB"。
	WriteBufferSize string `yaml:"writeBufferSize"`
	// WriteBufferNumber 控制 RocksDB 在阻塞写入前最多保留多少个 memtable。
	// 数值越大越不容易触发写入阻塞，但会消耗更多内存。
	// 示例：小型设备可设为 4，写入密集型服务器可设为 8。
	WriteBufferNumber int `yaml:"writeBufferNumber"`
	// WriteBufferNumberToMerge 控制读路径上会参与合并的 immutable memtable 数量。
	// 数值更大可能降低写放大，但也可能提高读成本。
	// 示例：1 表示更保守的行为，2 常用于写入较重的场景。
	WriteBufferNumberToMerge int `yaml:"writeBufferNumberToMerge"`
	// WalSizeLimitMb 限制 WAL 日志在轮转前可累计的大小。
	// 示例：0 表示交给 RocksDB 使用内部默认行为。
	WalSizeLimitMb uint64 `yaml:"walSizeLimitMb"`
	// MaxTotalWalSize 限制整个数据库允许累计的 WAL 总大小。
	// 数值越大越不容易触发强制 flush，但会占用更多磁盘。
	// 示例：中等写入流量场景可设为 "512MiB"。
	MaxTotalWalSize string `yaml:"maxTotalWalSize"`
	// WriteBufferSizeToMaintain 控制为冲突检查额外保留多少历史 memtable 数据。
	// 示例："0" 表示关闭额外保留缓冲。
	WriteBufferSizeToMaintain string `yaml:"writeBufferSizeToMaintain"`
	// Level0FileNumCompactionTrigger 控制 Level-0 compaction 的启动阈值。
	// 数值越小越早开始 compaction；数值越大则会延后 compaction，但可能增加读放大。
	// 示例：4 是比较常见的默认值。
	Level0FileNumCompactionTrigger int `yaml:"level0FileNumCompactionTrigger"`
	// MaxBytesForLevelBase 定义第一个 leveled compaction 层的目标总大小。
	// 示例：默认可设为 "256MiB"，较小部署可设为 "64MiB"。
	MaxBytesForLevelBase string `yaml:"maxBytesForLevelBase"`
	// TargetFileSizeBase 定义第一层 SST 文件的目标大小。
	// 示例："64MiB" 可以让 SST 文件保持中等大小。
	TargetFileSizeBase string `yaml:"targetFileSizeBase"`
	// MaxBackgroundJobs 限制后台 flush 与 compaction 的总 worker 数量。
	// 示例：低核机器可设为 2，服务器可设为 8 或更高。
	MaxBackgroundJobs int `yaml:"maxBackgroundJobs"`
	// Level0StopWritesTrigger 控制 L0 文件过多时何时硬停止写入。
	// 示例：24 对通用场景来说相对保守。
	Level0StopWritesTrigger int `yaml:"level0StopWritesTrigger"`
	// Level0SlowDownWritesTrigger 控制在硬停止前何时开始减速写入。
	// 示例：20 表示在达到 24 前就开始减速。
	Level0SlowDownWritesTrigger int `yaml:"level0SlowDownWritesTrigger"`
	// MaxSubCompactions 允许单个 compaction 任务拆分为多个 worker 并行执行。
	// 示例：资源受限主机可设为 1，大型服务器可设为 4。
	MaxSubCompactions uint32 `yaml:"maxSubCompactions"`
	// MaxLogFileSize 限制每个 RocksDB INFO 日志文件的大小。
	// 示例："1GiB" 可以减少日志轮转频率。
	MaxLogFileSize string `yaml:"maxLogFileSize"`
	// KeepLogFileNum 控制保留多少个 INFO 日志文件。
	// 示例：10 可以保留近期历史，同时避免磁盘无限增长。
	KeepLogFileNum uint `yaml:"keepLogFileNum"`
	// MaxCompactionBytes 限制 compaction 的最大输入数据量。
	// 示例："0" 表示交给 RocksDB 默认策略决定。
	MaxCompactionBytes string `yaml:"maxCompactionBytes"`
	// PeriodicCompactionSeconds 控制老文件经过指定时间后触发周期性 compaction。
	// 示例：86400 表示每天触发一次周期性 compaction。
	PeriodicCompactionSeconds int `yaml:"periodicCompactionSeconds"`
	// CompressionParallelThreads 控制压缩阶段使用的额外 worker 线程数量。
	// 示例：默认可设为 1，CPU 充足的机器可设为 4。
	CompressionParallelThreads int `yaml:"compressionParallelThreads"`
	// RateBytesPerSec 通过 rate limiter 限制后台 IO 速率。
	// 示例："200MiB" 可以平滑 compaction 对共享磁盘的影响。
	RateBytesPerSec string `yaml:"rateBytesPerSec"`
	// RefillPeriodMicros 是 rate limiter 的令牌补充周期。
	// 示例：5000 表示每 5ms 补充一次令牌。
	RefillPeriodMicros int64 `yaml:"refillPeriodMicros"`
	// Fairness 控制 rate limiter 在不同请求之间分配 IO 令牌的均衡程度。
	// 示例：10 是比较均衡的设置。
	Fairness int32 `yaml:"fairness"`
	// AutoTuned 用于开启 RocksDB 的自适应 rate limiter 行为。
	// 示例：true 表示允许 RocksDB 根据负载变化自动调整后台 IO 限速。
	AutoTuned bool `yaml:"autoTuned"`
	// DeletionCollectorFactory 用于开启基于删除比例的 compaction 触发逻辑。
	// 格式："<windowSize> <numDelsTrigger>"。
	// 示例："100000 5000" 表示如果 100000 个 key 的滑动窗口内出现超过 5000 次删除，
	// RocksDB 可以更早触发 compaction 来清理 tombstone。
	DeletionCollectorFactory string `yaml:"deletionCollectorFactory"`
	// RecycleLogFileNum 控制可以复用而不是重新创建的 WAL 文件数量。
	// 示例：4 可以减少繁忙写路径上的文件创建抖动。
	RecycleLogFileNum uint `yaml:"recycleLogFileNum"`
	// EnableBlobFiles 控制是否开启 BlobDB 风格的大 value 分离。
	// 当该值为 false 时，其它所有 blob 相关配置都会被忽略。
	// 示例：大 value 较多的场景可设为 true，小 value KV 场景可设为 false。
	EnableBlobFiles bool `yaml:"enableBlobFiles"`
	// MinBlobSize 控制多大的 value 会被写入 blob 文件。
	// 示例："4KiB" 表示启用 blob 文件后，大小 >= 4KiB 的 value 会进入 blob 文件。
	MinBlobSize string `yaml:"minBlobSize"`
	// CompressionType 控制普通 SST 文件使用的压缩算法。
	// 当前仓库支持的取值有：none、snappy、lz4、zstd。
	// 示例："lz4" 适合在读写成本之间取平衡。
	CompressionType string `yaml:"compressionType"`
	// BottommostCompressionType 控制最底层 level 使用的压缩算法。
	// 示例：冷数据较多时可设为 "zstd" 以获得更高压缩率。
	BottommostCompressionType string `yaml:"bottommostCompressionType"`
	// BlobFileSize 控制启用 blob 文件时每个 blob 文件的目标大小。
	// 示例："256MiB" 可以生成中等大小的 blob 文件。
	BlobFileSize string `yaml:"blobFileSize"`
	// BlobCompressionType 单独控制 blob 文件的压缩算法，不受普通 SST 压缩设置影响。
	// 当前仓库支持的取值有：none、snappy、lz4、zstd。
	// 示例："lz4" 可以在保持 blob 读取速度的同时压缩数据。
	BlobCompressionType string `yaml:"blobCompressionType"`
	// EnableBlobCompression 控制是否实际启用 BlobCompressionType。
	// 示例：true 且 BlobCompressionType="lz4" 时，blob 文件会使用 LZ4 压缩。
	EnableBlobCompression bool `yaml:"enableBlobCompression"`
	// EnableBlobGC 控制后台 GC 是否回收已经过期的 blob value。
	// 示例：存在大量覆盖写和删除的长生命周期数据库可设为 true。
	EnableBlobGC bool `yaml:"enableBlobGC"`
	// BlobGCAgeCutoff 控制 blob 文件至少多“旧”后才会被 GC 考虑。
	// 示例：0.25 表示文件大约经过正常生命周期的 25% 后就可以被 GC 考虑。
	BlobGCAgeCutoff float64 `yaml:"blobGCAgeCutoff"`
	// LRUSize 控制 block cache 的容量。
	// 示例：小机器可设为 "64MiB"，大服务器可设为 "1GiB"。
	LRUSize string `yaml:"lruSize"`
	// LRUShardBits 控制 LRU block cache 内部分片的数量。
	// 示例：6 表示 64 个分片，可降低并发下的锁竞争。
	LRUShardBits int `yaml:"lruShardBits"`
	// FilterBitsPerKey 控制 Bloom filter 的位密度。
	// 数值越大误判率越低，但会消耗更多内存。
	// 示例：10 是比较常见的默认值。
	FilterBitsPerKey float64 `yaml:"filterBitsPerKey"`
	// BlockBasedPartitionFilters 控制是否为大索引启用分区 filter。
	// 示例：大型数据库可设为 true，以减少 filter 的内存占用。
	BlockBasedPartitionFilters bool `yaml:"blockBasedPartitionFilters"`
	// BlockBasedOptimizeFiltersForMemory 用更高的 CPU 成本换取更低的 filter 内存占用。
	// 示例：内存敏感型系统通常适合设为 true。
	BlockBasedOptimizeFiltersForMemory bool `yaml:"blockBasedOptimizeFiltersForMemory"`
	// BlockBasedWholeKeyFiltering 控制 Bloom filter 是否基于完整 key 进行过滤。
	// 示例：点查较多的场景通常设为 true。
	BlockBasedWholeKeyFiltering bool `yaml:"blockBasedWholeKeyFiltering"`
	// BlockBasedCacheIndexAndFilterBlocks 控制是否将 index 和 filter block 放入 block cache。
	// 示例：重复查询较多时设为 true 可以提升读局部性。
	BlockBasedCacheIndexAndFilterBlocks bool `yaml:"blockBasedCacheIndexAndFilterBlocks"`
	// BlockBasedCacheIndexAndFilterBlocksHighPriority 控制是否将 index/filter block 作为高优先级缓存项。
	// 示例：数据块频繁换入换出时，设为 true 有助于避免元数据被挤出缓存。
	BlockBasedCacheIndexAndFilterBlocksHighPriority bool `yaml:"blockBasedCacheIndexAndFilterBlocksHighPriority"`
	// BlockBasedPinL0FilterAndIndexBlocks 控制是否将 L0 的 filter 和 index 元数据钉在缓存中。
	// 示例：最近文件很多、写入较重的系统通常适合设为 true。
	BlockBasedPinL0FilterAndIndexBlocks bool `yaml:"blockBasedPinL0FilterAndIndexBlocks"`
	// BlockBasedPinTopLevelIndexAndFilter 控制是否将顶层分区元数据钉在缓存中。
	// 示例：分区索引较多时设为 true 有助于减少元数据重复 miss。
	BlockBasedPinTopLevelIndexAndFilter bool `yaml:"blockBasedPinTopLevelIndexAndFilter"`
	// BlockBasedIndexType 控制 SST index 的格式。
	// 示例：2 表示对大表使用 two-level index search。
	BlockBasedIndexType int `yaml:"blockBasedIndexType"`
	// BlockBasedDataBlockIndexType 控制 data block 内部索引的组织方式。
	// 示例：1 表示在 data block 内同时启用二分查找和 hash 查找。
	BlockBasedDataBlockIndexType int `yaml:"blockBasedDataBlockIndexType"`
	// BlockBasedDataBlockHashRatio 控制带有 hash index 的 data block 比例。
	// 示例：点读较重的场景可设为 0.75，以提升 hash 查找比例。
	BlockBasedDataBlockHashRatio float64 `yaml:"blockBasedDataBlockHashRatio"`
	// BlockBasedTopLevelIndexPinningTier 控制顶层分区索引的 pinning 行为。
	// 示例：0 表示回退到 RocksDB 默认行为。
	BlockBasedTopLevelIndexPinningTier int `yaml:"blockBasedTopLevelIndexPinningTier"`
	// BlockBasedPartitionPinningTier 控制分区 index/filter block 的 pinning 行为。
	// 示例：0 表示回退到 RocksDB 默认行为。
	BlockBasedPartitionPinningTier int `yaml:"blockBasedPartitionPinningTier"`
	// BlockBasedUnpartitionedPinningTier 控制未分区元数据的 pinning 行为。
	// 示例：0 表示回退到 RocksDB 默认行为。
	BlockBasedUnpartitionedPinningTier int `yaml:"blockBasedUnpartitionedPinningTier"`
	// BlockSize 控制 SST 文件中每个 data block 的目标大小。
	// 示例："16KiB" 更利于点查，"64KiB" 更利于扫描和减少元数据开销。
	BlockSize string `yaml:"blockSize"`
}

func DefaultConfig() *Config {
	return &Config{
		WriteBufferSize:                     "512MiB",
		WriteBufferNumber:                   8,
		WriteBufferNumberToMerge:            2,
		MaxTotalWalSize:                     "256MiB",
		WriteBufferSizeToMaintain:           "0",
		Level0FileNumCompactionTrigger:      2,
		MaxBytesForLevelBase:                "2GiB",
		TargetFileSizeBase:                  "512MiB",
		MaxBackgroundJobs:                   16,
		Level0StopWritesTrigger:             32,
		Level0SlowDownWritesTrigger:         24,
		MaxSubCompactions:                   8,
		MaxLogFileSize:                      "1GiB",
		KeepLogFileNum:                      4,
		MaxCompactionBytes:                  "8GiB",
		CompressionParallelThreads:          8,
		RateBytesPerSec:                     "4GiB",
		RefillPeriodMicros:                  5000,
		Fairness:                            10,
		EnableBlobCompression:               false,
		CompressionType:                     "lz4",
		BottommostCompressionType:           "zstd",
		BlobCompressionType:                 "zstd",
		MinBlobSize:                         "16KiB",
		BlobFileSize:                        "256MiB",
		BlobGCAgeCutoff:                     0.25,
		LRUSize:                             "4GiB",
		FilterBitsPerKey:                    10,
		BlockBasedPartitionFilters:          true,
		BlockBasedOptimizeFiltersForMemory:  true,
		BlockBasedWholeKeyFiltering:         true,
		BlockBasedCacheIndexAndFilterBlocks: true,
		BlockBasedCacheIndexAndFilterBlocksHighPriority: true,
		BlockBasedPinL0FilterAndIndexBlocks:             true,
		BlockBasedPinTopLevelIndexAndFilter:             true,
		BlockBasedIndexType:                             2,
		BlockBasedDataBlockIndexType:                    1,
		BlockBasedDataBlockHashRatio:                    0.75,
		BlockSize:                                       "64KiB",
	}
}

func SmallConfig() *Config {
	return &Config{
		WriteBufferSize:                     "256MiB",
		WriteBufferNumber:                   4,
		WriteBufferNumberToMerge:            2,
		MaxTotalWalSize:                     "256MiB",
		MaxBytesForLevelBase:                "1GiB",
		TargetFileSizeBase:                  "512MiB",
		MaxBackgroundJobs:                   8,
		Level0StopWritesTrigger:             16,
		Level0SlowDownWritesTrigger:         12,
		MaxSubCompactions:                   4,
		MaxLogFileSize:                      "256MiB",
		KeepLogFileNum:                      4,
		RateBytesPerSec:                     "1GiB",
		RefillPeriodMicros:                  5000,
		EnableBlobCompression:               false,
		Fairness:                            10,
		CompressionType:                     "lz4",
		BottommostCompressionType:           "zstd",
		BlobCompressionType:                 "zstd",
		MinBlobSize:                         "4KiB",
		BlobFileSize:                        "256MiB",
		BlobGCAgeCutoff:                     0.25,
		LRUSize:                             "2GiB",
		FilterBitsPerKey:                    8,
		BlockBasedPartitionFilters:          true,
		BlockBasedOptimizeFiltersForMemory:  true,
		BlockBasedWholeKeyFiltering:         true,
		BlockBasedCacheIndexAndFilterBlocks: true,
		BlockBasedDataBlockHashRatio:        0.5,
		BlockSize:                           "16KiB",
	}
}

func LargeConfig() *Config {
	return &Config{
		WriteBufferSize:                     "512MiB",
		WriteBufferNumber:                   16,
		WriteBufferNumberToMerge:            2,
		MaxTotalWalSize:                     "2GiB",
		WriteBufferSizeToMaintain:           "256MiB",
		MaxBytesForLevelBase:                "4GiB",
		TargetFileSizeBase:                  "512MiB",
		MaxBackgroundJobs:                   32,
		Level0StopWritesTrigger:             36,
		Level0SlowDownWritesTrigger:         28,
		MaxSubCompactions:                   16,
		MaxLogFileSize:                      "512MiB",
		KeepLogFileNum:                      4,
		CompressionParallelThreads:          16,
		RateBytesPerSec:                     "16GiB",
		RefillPeriodMicros:                  5000,
		EnableBlobCompression:               false,
		Fairness:                            10,
		CompressionType:                     "lz4",
		BottommostCompressionType:           "zstd",
		BlobCompressionType:                 "zstd",
		MinBlobSize:                         "16KiB",
		BlobFileSize:                        "1GiB",
		BlobGCAgeCutoff:                     0.25,
		LRUSize:                             "8GiB",
		LRUShardBits:                        6,
		FilterBitsPerKey:                    10,
		BlockBasedPartitionFilters:          true,
		BlockBasedOptimizeFiltersForMemory:  true,
		BlockBasedWholeKeyFiltering:         true,
		BlockBasedCacheIndexAndFilterBlocks: true,
		BlockBasedCacheIndexAndFilterBlocksHighPriority: true,
		BlockBasedPinL0FilterAndIndexBlocks:             true,
		BlockBasedPinTopLevelIndexAndFilter:             true,
		BlockBasedIndexType:                             2,
		BlockBasedDataBlockIndexType:                    1,
		BlockBasedDataBlockHashRatio:                    0.75,
		BlockSize:                                       "64KiB",
	}
}
