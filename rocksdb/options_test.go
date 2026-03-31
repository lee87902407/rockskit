package rocksdb

import "testing"

func TestDefaultEnvConfigProvidesCacheSettings(t *testing.T) {
	cfg := DefaultEnvConfig()
	if cfg == nil {
		t.Fatal("DefaultEnvConfig() returned nil")
	}
	if cfg.LRUSize == "" {
		t.Fatal("expected default LRUSize to be set")
	}
	if cfg.BlockSize == "" {
		t.Fatal("expected default BlockSize to be set")
	}
}

func TestDefaultConfigProvidesLifecycleSettings(t *testing.T) {
	cfg := DefaultConfig()
	if cfg == nil {
		t.Fatal("DefaultConfig() returned nil")
	}
	if cfg.WriteBufferSize == "" {
		t.Fatal("expected default WriteBufferSize to be set")
	}
	if cfg.TargetFileSizeBase == "" {
		t.Fatal("expected default TargetFileSizeBase to be set")
	}
	if cfg.MaxBytesForLevelBase == "" {
		t.Fatal("expected default MaxBytesForLevelBase to be set")
	}
}

func TestValidateConfigRejectsInvalidCompressionType(t *testing.T) {
	cfg := DefaultConfig()
	cfg.CompressionType = "invalid"

	if err := ValidateConfig(cfg); err == nil {
		t.Fatal("expected ValidateConfig() to reject invalid compression type")
	}
}

func TestValidateConfigRejectsEmptyWriteBufferSize(t *testing.T) {
	cfg := DefaultConfig()
	cfg.WriteBufferSize = ""

	if err := ValidateConfig(cfg); err == nil {
		t.Fatal("expected ValidateConfig() to reject empty write buffer size")
	}
}
