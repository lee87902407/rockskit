package rocksdb

import "testing"

func TestParseByteSizeSupportsBinaryAndSuffixForms(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  uint64
	}{
		{name: "K", input: "1K", want: 1024},
		{name: "KB", input: "1KB", want: 1024},
		{name: "KiB", input: "1KiB", want: 1024},
		{name: "M", input: "2M", want: 2 * 1024 * 1024},
		{name: "MB", input: "2MB", want: 2 * 1024 * 1024},
		{name: "MiB", input: "2MiB", want: 2 * 1024 * 1024},
		{name: "G", input: "3G", want: 3 * 1024 * 1024 * 1024},
		{name: "GB", input: "3GB", want: 3 * 1024 * 1024 * 1024},
		{name: "GiB", input: "3GiB", want: 3 * 1024 * 1024 * 1024},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := nativeParseByteSize(tt.input)
			if err != nil {
				t.Fatalf("parseByteSize(%q) error = %v", tt.input, err)
			}
			if got != tt.want {
				t.Fatalf("parseByteSize(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestSmallDefaultLargeConfigPresets(t *testing.T) {
	small := SmallConfig()
	def := DefaultConfig()
	large := LargeConfig()

	if small == nil || def == nil || large == nil {
		t.Fatal("expected all config presets to be non-nil")
	}
	if small.WriteBufferSize == def.WriteBufferSize {
		t.Fatal("expected small preset to differ from default write buffer size")
	}
	if large.WriteBufferSize == def.WriteBufferSize {
		t.Fatal("expected large preset to differ from default write buffer size")
	}
	if small.LRUSize == def.LRUSize {
		t.Fatal("expected small preset to differ from default cache size")
	}
	if large.LRUSize == def.LRUSize {
		t.Fatal("expected large preset to differ from default cache size")
	}
	if small.MaxBackgroundJobs >= large.MaxBackgroundJobs {
		t.Fatal("expected large preset to allow more background jobs than small preset")
	}
}

func TestValidateConfigAcceptsDeletionCollectorFactory(t *testing.T) {
	cfg := DefaultConfig()
	cfg.DeletionCollectorFactory = "100000 5000"

	if err := ValidateConfig(cfg); err != nil {
		t.Fatalf("expected DeletionCollectorFactory to be accepted, got %v", err)
	}
}

func TestCloneConfigPreservesExplicitZeroValues(t *testing.T) {
	cfg := DefaultConfig()
	cfg.WriteBufferNumberToMerge = 0
	cfg.CompressionParallelThreads = 0
	cfg.RateBytesPerSec = ""

	cloned := cloneConfig(cfg)
	if cloned.WriteBufferNumberToMerge != 0 {
		t.Fatalf("expected explicit zero WriteBufferNumberToMerge to be preserved, got %d", cloned.WriteBufferNumberToMerge)
	}
	if cloned.CompressionParallelThreads != 0 {
		t.Fatalf("expected explicit zero CompressionParallelThreads to be preserved, got %d", cloned.CompressionParallelThreads)
	}
	if cloned.RateBytesPerSec != "" {
		t.Fatalf("expected empty RateBytesPerSec to be preserved, got %q", cloned.RateBytesPerSec)
	}
}

func TestDefaultEnvConfigProvidesCacheSettings(t *testing.T) {
	cfg := DefaultConfig()
	if cfg == nil {
		t.Fatal("DefaultConfig() returned nil")
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
	if cfg.WriteBufferNumber == 0 {
		t.Fatal("expected default WriteBufferNumber to be set")
	}
	if cfg.MaxBackgroundJobs == 0 {
		t.Fatal("expected default MaxBackgroundJobs to be set")
	}
}

func TestValidateConfigRejectsInvalidCompressionType(t *testing.T) {
	cfg := DefaultConfig()
	cfg.CompressionType = "invalid"

	if err := ValidateConfig(cfg); err == nil {
		t.Fatal("expected ValidateConfig() to reject invalid compression type")
	}
}

func TestValidateConfigAllowsEmptyWriteBufferSizeToSkipOption(t *testing.T) {
	cfg := DefaultConfig()
	cfg.WriteBufferSize = ""

	if err := ValidateConfig(cfg); err != nil {
		t.Fatalf("expected ValidateConfig() to allow empty write buffer size, got %v", err)
	}
}

func TestDefaultConfigProvidesBlockBasedTuning(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.FilterBitsPerKey == 0 {
		t.Fatal("expected default FilterBitsPerKey to be set")
	}
	if cfg.BlockBasedIndexType == 0 && cfg.BlockBasedDataBlockIndexType == 0 {
		t.Fatal("expected default block-based index tuning to be set")
	}
	if cfg.BlockBasedDataBlockHashRatio <= 0 {
		t.Fatal("expected default BlockBasedDataBlockHashRatio to be set")
	}
}
