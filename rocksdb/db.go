package rocksdb

import (
	"fmt"
	"path/filepath"
	"sync"

	native "github.com/yanjie/rockskit/internal/cgo"
)

type DB struct {
	path string
	hdl  *native.DB

	mu     sync.Mutex
	closed bool
}

func Create(path string, envCfg *EnvConfig, cfg *Config) (*DB, error) {
	cleanPath := filepath.Clean(path)
	_ = cloneEnvConfig(envCfg)
	cfg = cloneConfig(cfg)
	cfg.CreateIfMissing = true
	cfg.ErrorIfExists = true
	if err := ValidateConfig(cfg); err != nil {
		return nil, err
	}

	options, err := buildNativeOptions(cfg)
	if err != nil {
		return nil, err
	}
	defer options.Close()

	hdl, err := native.Open(cleanPath, options)
	if err != nil {
		return nil, fmt.Errorf("create rocksdb at %q: %w", cleanPath, err)
	}

	return &DB{path: cleanPath, hdl: hdl}, nil
}

func Open(path string, envCfg *EnvConfig, cfg *Config) (*DB, error) {
	cleanPath := filepath.Clean(path)
	_ = cloneEnvConfig(envCfg)
	cfg = cloneConfig(cfg)
	cfg.CreateIfMissing = false
	cfg.ErrorIfExists = false
	if err := ValidateConfig(cfg); err != nil {
		return nil, err
	}

	options, err := buildNativeOptions(cfg)
	if err != nil {
		return nil, err
	}
	defer options.Close()

	hdl, err := native.Open(cleanPath, options)
	if err != nil {
		return nil, fmt.Errorf("open rocksdb at %q: %w", cleanPath, err)
	}

	return &DB{path: cleanPath, hdl: hdl}, nil
}

func (db *DB) Close() error {
	if db == nil {
		return nil
	}

	db.mu.Lock()
	defer db.mu.Unlock()
	if db.closed {
		return nil
	}
	if db.hdl != nil {
		db.hdl.Close()
		db.hdl = nil
	}
	db.closed = true
	return nil
}

func buildNativeOptions(cfg *Config) (*native.Options, error) {
	writeBufferSize, err := parseByteSize(cfg.WriteBufferSize)
	if err != nil {
		return nil, err
	}
	targetFileSizeBase, err := parseByteSize(cfg.TargetFileSizeBase)
	if err != nil {
		return nil, err
	}
	maxBytesForLevelBase, err := parseByteSize(cfg.MaxBytesForLevelBase)
	if err != nil {
		return nil, err
	}

	return native.NewOptions(native.RawOptions{
		WriteBufferSize:                writeBufferSize,
		Level0FileNumCompactionTrigger: cfg.Level0FileNumCompactionTrigger,
		TargetFileSizeBase:             targetFileSizeBase,
		MaxBytesForLevelBase:           maxBytesForLevelBase,
		CreateIfMissing:                cfg.CreateIfMissing,
		ErrorIfExists:                  cfg.ErrorIfExists,
		Compression:                    compressionCode(cfg.CompressionType),
	}), nil
}

func compressionCode(value string) int {
	switch normalizeCompressionType(value) {
	case "lz4":
		return 4
	case "zstd":
		return 7
	default:
		return 0
	}
}
