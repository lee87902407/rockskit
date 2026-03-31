package rocksdb

import (
	"fmt"
	"strconv"
	"strings"
)

type EnvConfig struct {
	LRUSize   string
	BlockSize string
}

type Config struct {
	WriteBufferSize                string
	Level0FileNumCompactionTrigger int
	MaxBytesForLevelBase           string
	TargetFileSizeBase             string
	CompressionType                string
	CreateIfMissing                bool
	ErrorIfExists                  bool
}

func DefaultEnvConfig() *EnvConfig {
	return &EnvConfig{
		LRUSize:   "64MB",
		BlockSize: "64KB",
	}
}

func DefaultConfig() *Config {
	return &Config{
		WriteBufferSize:                "64MB",
		Level0FileNumCompactionTrigger: 4,
		MaxBytesForLevelBase:           "256MB",
		TargetFileSizeBase:             "64MB",
		CompressionType:                "none",
	}
}

func ValidateConfig(cfg *Config) error {
	if cfg == nil {
		return fmt.Errorf("rocksdb config is nil")
	}
	if strings.TrimSpace(cfg.WriteBufferSize) == "" {
		return fmt.Errorf("rocksdb config write buffer size is empty")
	}
	if strings.TrimSpace(cfg.TargetFileSizeBase) == "" {
		return fmt.Errorf("rocksdb config target file size base is empty")
	}
	if strings.TrimSpace(cfg.MaxBytesForLevelBase) == "" {
		return fmt.Errorf("rocksdb config max bytes for level base is empty")
	}
	if cfg.Level0FileNumCompactionTrigger <= 0 {
		return fmt.Errorf("rocksdb config level0 compaction trigger must be positive")
	}

	switch normalizeCompressionType(cfg.CompressionType) {
	case "none", "lz4", "zstd":
		return nil
	default:
		return fmt.Errorf("rocksdb config compression type %q is invalid", cfg.CompressionType)
	}
}

func normalizeCompressionType(value string) string {
	trimmed := strings.TrimSpace(strings.ToLower(value))
	if trimmed == "" {
		return "none"
	}
	return trimmed
}

func parseByteSize(value string) (uint64, error) {
	trimmed := strings.TrimSpace(strings.ToUpper(value))
	trimmed = strings.TrimSuffix(trimmed, "B")

	multiplier := uint64(1)
	switch {
	case strings.HasSuffix(trimmed, "K"):
		multiplier = 1024
		trimmed = strings.TrimSuffix(trimmed, "K")
	case strings.HasSuffix(trimmed, "M"):
		multiplier = 1024 * 1024
		trimmed = strings.TrimSuffix(trimmed, "M")
	case strings.HasSuffix(trimmed, "G"):
		multiplier = 1024 * 1024 * 1024
		trimmed = strings.TrimSuffix(trimmed, "G")
	case strings.HasSuffix(trimmed, "T"):
		multiplier = 1024 * 1024 * 1024 * 1024
		trimmed = strings.TrimSuffix(trimmed, "T")
	}

	n, err := strconv.ParseUint(strings.TrimSpace(trimmed), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse byte size %q: %w", value, err)
	}
	return n * multiplier, nil
}
