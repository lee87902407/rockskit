package cgo

import (
	"fmt"
	"strconv"
	"strings"
)

func buildArtifacts(cfg *Config, createIfMissing bool, errorIfExists bool) (*artifacts, error) {
	var rateLimiter *RateLimiter
	options := NewOptions()

	if cfg.WriteBufferSize != "" {
		writeBufferSize, err := parseConfigSize(cfg.WriteBufferSize)
		if err != nil {
			options.Close()
			return nil, err
		}
		options.SetWriteBufferSize(writeBufferSize)
	}
	if cfg.WriteBufferNumber > 0 {
		options.SetWriteBufferNumber(cfg.WriteBufferNumber)
	}
	if cfg.WriteBufferNumberToMerge > 0 {
		options.SetWriteBufferNumberToMerge(cfg.WriteBufferNumberToMerge)
	}
	if cfg.WalSizeLimitMb > 0 {
		options.SetWalSizeLimitMb(cfg.WalSizeLimitMb)
	}
	if cfg.MaxTotalWalSize != "" {
		maxTotalWalSize, err := parseConfigSize(cfg.MaxTotalWalSize)
		if err != nil {
			options.Close()
			return nil, err
		}
		options.SetMaxTotalWalSize(maxTotalWalSize)
	}
	if cfg.Level0FileNumCompactionTrigger > 0 {
		options.SetLevel0FileNumCompactionTrigger(cfg.Level0FileNumCompactionTrigger)
	}
	if cfg.MaxBytesForLevelBase != "" {
		maxBytesForLevelBase, err := parseConfigSize(cfg.MaxBytesForLevelBase)
		if err != nil {
			options.Close()
			return nil, err
		}
		options.SetMaxBytesForLevelBase(maxBytesForLevelBase)
	}
	if cfg.TargetFileSizeBase != "" {
		targetFileSizeBase, err := parseConfigSize(cfg.TargetFileSizeBase)
		if err != nil {
			options.Close()
			return nil, err
		}
		options.SetTargetFileSizeBase(targetFileSizeBase)
	}
	if cfg.MaxBackgroundJobs > 0 {
		options.SetMaxBackgroundJobs(cfg.MaxBackgroundJobs)
	}
	if cfg.Level0StopWritesTrigger > 0 {
		options.SetLevel0StopWritesTrigger(cfg.Level0StopWritesTrigger)
	}
	if cfg.Level0SlowDownWritesTrigger > 0 {
		options.SetLevel0SlowDownWritesTrigger(cfg.Level0SlowDownWritesTrigger)
	}
	if cfg.MaxSubCompactions > 0 {
		options.SetMaxSubCompactions(cfg.MaxSubCompactions)
	}
	if cfg.MaxLogFileSize != "" {
		maxLogFileSize, err := parseConfigSize(cfg.MaxLogFileSize)
		if err != nil {
			options.Close()
			return nil, err
		}
		options.SetMaxLogFileSize(maxLogFileSize)
	}
	if cfg.KeepLogFileNum > 0 {
		options.SetKeepLogFileNum(cfg.KeepLogFileNum)
	}
	if cfg.RecycleLogFileNum > 0 {
		options.SetRecycleLogFileNum(cfg.RecycleLogFileNum)
	}
	if cfg.MaxCompactionBytes != "" {
		maxCompactionBytes, err := parseConfigSize(cfg.MaxCompactionBytes)
		if err != nil {
			options.Close()
			return nil, err
		}
		options.SetMaxCompactionBytes(maxCompactionBytes)
	}
	if cfg.PeriodicCompactionSeconds > 0 {
		options.SetPeriodicCompactionSeconds(cfg.PeriodicCompactionSeconds)
	}
	if cfg.CompressionParallelThreads > 0 {
		options.SetCompressionParallelThreads(cfg.CompressionParallelThreads)
	}
	if cfg.WriteBufferSizeToMaintain != "" {
		writeBufferSizeToMaintain, err := parseConfigSize(cfg.WriteBufferSizeToMaintain)
		if err != nil {
			options.Close()
			return nil, err
		}
		options.SetWriteBufferSizeToMaintain(writeBufferSizeToMaintain)
	}
	if cfg.CompressionType != "" {
		compressionKind, err := compressionCode(cfg.CompressionType)
		if err != nil {
			options.Close()
			return nil, err
		}
		options.SetCompression(compressionKind)
	}
	if cfg.BottommostCompressionType != "" {
		bottommostCompressionKind, err := compressionCode(cfg.BottommostCompressionType)
		if err != nil {
			options.Close()
			return nil, err
		}
		options.SetBottommostCompression(bottommostCompressionKind)
	}
	if strings.TrimSpace(cfg.DeletionCollectorFactory) != "" {
		windowSize, numDelsTrigger, err := parseDeletionCollectorFactory(cfg.DeletionCollectorFactory)
		if err != nil {
			options.Close()
			return nil, err
		}
		options.AddCompactOnDeletionCollectorFactory(windowSize, numDelsTrigger)
	}

	options.SetCreateIfMissing(createIfMissing)
	options.SetErrorIfExists(errorIfExists)

	if cfg.EnableBlobFiles {
		options.SetEnableBlobFiles(true)
		if cfg.MinBlobSize != "" {
			minBlobSize, err := parseConfigSize(cfg.MinBlobSize)
			if err != nil {
				options.Close()
				return nil, err
			}
			options.SetMinBlobSize(minBlobSize)
		}
		if cfg.BlobFileSize != "" {
			blobFileSize, err := parseConfigSize(cfg.BlobFileSize)
			if err != nil {
				options.Close()
				return nil, err
			}
			options.SetBlobFileSize(blobFileSize)
		}
		options.SetEnableBlobGC(cfg.EnableBlobGC)
		if cfg.BlobGCAgeCutoff > 0 {
			options.SetBlobGCAgeCutoff(cfg.BlobGCAgeCutoff)
		}
		if cfg.EnableBlobCompression && cfg.BlobCompressionType != "" {
			blobCompressionKind, err := compressionCode(cfg.BlobCompressionType)
			if err != nil {
				options.Close()
				return nil, err
			}
			options.SetBlobCompressionType(blobCompressionKind)
		}
	}

	cacheSize := uint64(0)
	if cfg.LRUSize != "" {
		var err error
		cacheSize, err = parseConfigSize(cfg.LRUSize)
		if err != nil {
			options.Close()
			return nil, err
		}
	}

	blockOptions := NewBlockBasedOptions()
	if cacheSize == 0 {
		blockOptions.SetNoBlockCache(true)
		options.SetBlockBasedTableFactory(blockOptions)
		if strings.TrimSpace(cfg.RateBytesPerSec) != "" {
			rate, err := parseByteSizeToInt64(cfg.RateBytesPerSec)
			if err != nil {
				options.Close()
				blockOptions.Close()
				return nil, err
			}
			rateLimiter = NewRateLimiter(rate, cfg.RefillPeriodMicros, cfg.Fairness, cfg.AutoTuned)
			options.SetRateLimiter(rateLimiter)
		}
		return &artifacts{options: options, blockOptions: blockOptions, rateLimiter: rateLimiter}, nil
	}

	filterPolicy := NewBloomFilterPolicy(cfg.FilterBitsPerKey)
	blockOptions.SetFilterPolicy(filterPolicy)
	if cfg.BlockSize != "" {
		blockSize, err := parseConfigSize(cfg.BlockSize)
		if err != nil {
			options.Close()
			blockOptions.Close()
			filterPolicy.Close()
			return nil, err
		}
		blockOptions.SetBlockSize(blockSize)
	}
	blockOptions.SetPartitionFilters(cfg.BlockBasedPartitionFilters)
	blockOptions.SetOptimizeFiltersForMemory(cfg.BlockBasedOptimizeFiltersForMemory)
	blockOptions.SetWholeKeyFiltering(cfg.BlockBasedWholeKeyFiltering)
	blockOptions.SetCacheIndexAndFilterBlocks(cfg.BlockBasedCacheIndexAndFilterBlocks)
	blockOptions.SetCacheIndexAndFilterBlocksWithHighPriority(cfg.BlockBasedCacheIndexAndFilterBlocksHighPriority)
	blockOptions.SetPinL0FilterAndIndexBlocksInCache(cfg.BlockBasedPinL0FilterAndIndexBlocks)
	blockOptions.SetPinTopLevelIndexAndFilter(cfg.BlockBasedPinTopLevelIndexAndFilter)
	if cfg.BlockBasedIndexType >= 0 {
		blockOptions.SetIndexType(cfg.BlockBasedIndexType)
	}
	if cfg.BlockBasedDataBlockIndexType >= 0 {
		blockOptions.SetDataBlockIndexType(cfg.BlockBasedDataBlockIndexType)
	}
	if cfg.BlockBasedDataBlockHashRatio > 0 {
		blockOptions.SetDataBlockHashRatio(cfg.BlockBasedDataBlockHashRatio)
	}
	if cfg.BlockBasedTopLevelIndexPinningTier > 0 {
		blockOptions.SetTopLevelIndexPinningTier(cfg.BlockBasedTopLevelIndexPinningTier)
	}
	if cfg.BlockBasedPartitionPinningTier > 0 {
		blockOptions.SetPartitionPinningTier(cfg.BlockBasedPartitionPinningTier)
	}
	if cfg.BlockBasedUnpartitionedPinningTier > 0 {
		blockOptions.SetUnpartitionedPinningTier(cfg.BlockBasedUnpartitionedPinningTier)
	}

	if cacheSize > 0 {
		cache, err := NewLRUCache(cacheSize, cfg.LRUShardBits)
		if err != nil {
			options.Close()
			blockOptions.Close()
			filterPolicy.Close()
			return nil, err
		}
		blockOptions.SetBlockCache(cache)
		options.SetBlockBasedTableFactory(blockOptions)
		if strings.TrimSpace(cfg.RateBytesPerSec) != "" {
			rate, err := parseByteSizeToInt64(cfg.RateBytesPerSec)
			if err != nil {
				options.Close()
				blockOptions.Close()
				filterPolicy.Close()
				cache.Close()
				return nil, err
			}
			rateLimiter = NewRateLimiter(rate, cfg.RefillPeriodMicros, cfg.Fairness, cfg.AutoTuned)
			options.SetRateLimiter(rateLimiter)
		}
		return &artifacts{options: options, blockOptions: blockOptions, cache: cache, rateLimiter: rateLimiter, filterPolicy: filterPolicy}, nil
	}

	options.SetBlockBasedTableFactory(blockOptions)
	if strings.TrimSpace(cfg.RateBytesPerSec) != "" {
		rate, err := parseByteSizeToInt64(cfg.RateBytesPerSec)
		if err != nil {
			options.Close()
			blockOptions.Close()
			filterPolicy.Close()
			return nil, err
		}
		rateLimiter = NewRateLimiter(rate, cfg.RefillPeriodMicros, cfg.Fairness, cfg.AutoTuned)
		options.SetRateLimiter(rateLimiter)
	}

	return &artifacts{options: options, blockOptions: blockOptions, rateLimiter: rateLimiter, filterPolicy: filterPolicy}, nil
}

func parseConfigSize(value string) (uint64, error) {
	if strings.TrimSpace(value) == "" {
		return 0, nil
	}
	parsed, err := parseByteSize(value)
	if err != nil {
		return 0, err
	}
	return parsed, nil
}

func compressionCode(value string) (int, error) {
	switch normalizeCompressionType(value) {
	case "none":
		return 0, nil
	case "snappy":
		return 1, nil
	case "lz4":
		return 4, nil
	case "zstd":
		return 7, nil
	default:
		return 0, fmt.Errorf("rocksdb config compression type %q is invalid", value)
	}
}

func normalizeCompressionType(value string) string {
	trimmed := strings.TrimSpace(strings.ToLower(value))
	if trimmed == "" {
		return "none"
	}
	return trimmed
}

func parseByteSizeToInt64(value string) (int64, error) {
	v, err := parseByteSize(value)
	if err != nil {
		return 0, err
	}
	return int64(v), nil
}

func ParseByteSizeForTest(value string) (uint64, error) {
	return parseByteSize(value)
}

func parseByteSize(value string) (uint64, error) {
	trimmed := strings.TrimSpace(strings.ToUpper(value))
	trimmed = strings.ReplaceAll(trimmed, "IB", "B")

	multipliers := map[string]uint64{
		"K":  1024,
		"KB": 1024,
		"M":  1024 * 1024,
		"MB": 1024 * 1024,
		"G":  1024 * 1024 * 1024,
		"GB": 1024 * 1024 * 1024,
		"T":  1024 * 1024 * 1024 * 1024,
		"TB": 1024 * 1024 * 1024 * 1024,
	}

	for suffix, multiplier := range multipliers {
		if strings.HasSuffix(trimmed, suffix) {
			n, err := strconv.ParseUint(strings.TrimSpace(strings.TrimSuffix(trimmed, suffix)), 10, 64)
			if err != nil {
				return 0, fmt.Errorf("parse byte size %q: %w", value, err)
			}
			return n * multiplier, nil
		}
	}

	n, err := strconv.ParseUint(strings.TrimSpace(trimmed), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse byte size %q: %w", value, err)
	}
	return n, nil
}

func parseDeletionCollectorFactory(value string) (uint, uint, error) {
	parts := strings.Fields(value)
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("rocksdb config DeletionCollectorFactory error")
	}
	windowSize, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("rocksdb config DeletionCollectorFactory error")
	}
	numDelsTrigger, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("rocksdb config DeletionCollectorFactory error")
	}
	return uint(windowSize), uint(numDelsTrigger), nil
}
