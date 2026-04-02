package rocksdb

import (
	"fmt"
	"strings"

	native "github.com/yanjie/rockskit/internal/cgo"
)

type Config = native.Config

func DefaultConfig() *Config {
	return native.DefaultConfig()
}

func SmallConfig() *Config {
	return native.SmallConfig()
}

func LargeConfig() *Config {
	return native.LargeConfig()
}

func ValidateConfig(cfg *Config) error {
	if cfg == nil {
		return fmt.Errorf("rocksdb config is nil")
	}
	if cfg.WriteBufferSize != "" && !looksLikeByteSize(cfg.WriteBufferSize) {
		return fmt.Errorf("rocksdb config writeBufferSize is invalid")
	}
	if cfg.MaxTotalWalSize != "" && !looksLikeByteSize(cfg.MaxTotalWalSize) {
		return fmt.Errorf("rocksdb config maxTotalWalSize is invalid")
	}
	if cfg.WriteBufferSizeToMaintain != "" && !looksLikeByteSize(cfg.WriteBufferSizeToMaintain) {
		return fmt.Errorf("rocksdb config writeBufferSizeToMaintain is invalid")
	}
	if cfg.MaxBytesForLevelBase != "" && !looksLikeByteSize(cfg.MaxBytesForLevelBase) {
		return fmt.Errorf("rocksdb config maxBytesForLevelBase is invalid")
	}
	if cfg.TargetFileSizeBase != "" && !looksLikeByteSize(cfg.TargetFileSizeBase) {
		return fmt.Errorf("rocksdb config targetFileSizeBase is invalid")
	}
	if cfg.MaxLogFileSize != "" && !looksLikeByteSize(cfg.MaxLogFileSize) {
		return fmt.Errorf("rocksdb config maxLogFileSize is invalid")
	}
	if cfg.MaxCompactionBytes != "" && !looksLikeByteSize(cfg.MaxCompactionBytes) {
		return fmt.Errorf("rocksdb config maxCompactionBytes is invalid")
	}
	if cfg.RateBytesPerSec != "" && !looksLikeByteSize(cfg.RateBytesPerSec) {
		return fmt.Errorf("rocksdb config rateBytesPerSec is invalid")
	}
	if cfg.LRUSize != "" && !looksLikeByteSize(cfg.LRUSize) {
		return fmt.Errorf("rocksdb config lruSize is invalid")
	}
	if cfg.BlockSize != "" && !looksLikeByteSize(cfg.BlockSize) {
		return fmt.Errorf("rocksdb config blockSize is invalid")
	}
	if cfg.WriteBufferNumber < 0 {
		return fmt.Errorf("rocksdb config writeBufferNumber must be positive")
	}
	if cfg.WriteBufferNumberToMerge < 0 {
		return fmt.Errorf("rocksdb config writeBufferNumberToMerge must be positive")
	}
	if cfg.Level0FileNumCompactionTrigger < 0 {
		return fmt.Errorf("rocksdb config level0FileNumCompactionTrigger must be positive")
	}
	if cfg.MaxBackgroundJobs < 0 {
		return fmt.Errorf("rocksdb config maxBackgroundJobs must be positive")
	}
	if cfg.BlobGCAgeCutoff < 0 || cfg.BlobGCAgeCutoff > 1 {
		return fmt.Errorf("rocksdb config blobGCAgeCutoff must be between 0 and 1")
	}
	if cfg.FilterBitsPerKey < 0 {
		return fmt.Errorf("rocksdb config filterBitsPerKey must be positive")
	}
	if cfg.BlockBasedIndexType < 0 || cfg.BlockBasedIndexType > 3 {
		return fmt.Errorf("rocksdb config blockBasedIndexType must be between 0 and 2")
	}
	if cfg.BlockBasedDataBlockIndexType < 0 || cfg.BlockBasedDataBlockIndexType > 1 {
		return fmt.Errorf("rocksdb config blockBasedDataBlockIndexType must be between 0 and 1")
	}
	if cfg.BlockBasedDataBlockHashRatio < 0 || cfg.BlockBasedDataBlockHashRatio > 1 {
		return fmt.Errorf("rocksdb config blockBasedDataBlockHashRatio must be between 0 and 1")
	}
	if cfg.BlockBasedTopLevelIndexPinningTier < 0 || cfg.BlockBasedTopLevelIndexPinningTier > 3 {
		return fmt.Errorf("rocksdb config blockBasedTopLevelIndexPinningTier must be between 0 and 3")
	}
	if cfg.BlockBasedPartitionPinningTier < 0 || cfg.BlockBasedPartitionPinningTier > 3 {
		return fmt.Errorf("rocksdb config blockBasedPartitionPinningTier must be between 0 and 3")
	}
	if cfg.BlockBasedUnpartitionedPinningTier < 0 || cfg.BlockBasedUnpartitionedPinningTier > 3 {
		return fmt.Errorf("rocksdb config blockBasedUnpartitionedPinningTier must be between 0 and 3")
	}
	if cfg.CompressionType != "" {
		if _, err := compressionOptionName(cfg.CompressionType); err != nil {
			return err
		}
	}
	if cfg.BottommostCompressionType != "" {
		if _, err := compressionOptionName(cfg.BottommostCompressionType); err != nil {
			return err
		}
	}
	if cfg.BlobCompressionType != "" {
		if _, err := compressionOptionName(cfg.BlobCompressionType); err != nil {
			return err
		}
	}
	if strings.TrimSpace(cfg.DeletionCollectorFactory) != "" {
		parts := strings.Fields(cfg.DeletionCollectorFactory)
		if len(parts) != 2 {
			return fmt.Errorf("rocksdb config deletionCollectorFactory is invalid")
		}
	}
	return nil
}

func compressionOptionName(value string) (string, error) {
	switch normalizeCompressionType(value) {
	case "none":
		return "kNoCompression", nil
	case "snappy":
		return "kSnappyCompression", nil
	case "lz4":
		return "kLZ4Compression", nil
	case "zstd":
		return "kZSTD", nil
	default:
		return "", fmt.Errorf("rocksdb config compression type %q is invalid", value)
	}
}

func looksLikeByteSize(value string) bool {
	_, err := nativeParseByteSize(value)
	return err == nil
}

func nativeParseByteSize(value string) (uint64, error) {
	return native.ParseByteSizeForTest(value)
}

func normalizeCompressionType(value string) string {
	trimmed := strings.TrimSpace(strings.ToLower(value))
	if trimmed == "" {
		return "none"
	}
	return trimmed
}
