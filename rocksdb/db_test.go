package rocksdb

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestCreateCreatesNewDB(t *testing.T) {
	path := filepath.Join(t.TempDir(), "db")
	dbCfg := DefaultConfig()

	db, err := Create(path, dbCfg)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	defer db.Close()

	if _, err := os.Stat(filepath.Join(path, "CURRENT")); err != nil {
		t.Fatalf("expected CURRENT file from real rocksdb: %v", err)
	}
}

func TestCreateFailsWhenDBAlreadyExists(t *testing.T) {
	path := filepath.Join(t.TempDir(), "db")
	dbCfg := DefaultConfig()

	db, err := Create(path, dbCfg)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	defer db.Close()

	other, err := Create(path, dbCfg)
	if err == nil {
		other.Close()
		t.Fatal("expected Create() to fail for existing db")
	}
}

func TestOpenFailsWhenDBIsMissing(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing")
	dbCfg := DefaultConfig()

	if _, err := Open(path, dbCfg); err == nil {
		t.Fatal("expected Open() to fail for missing db")
	}
}

func TestOpenFailsForUninitializedDirectory(t *testing.T) {
	path := filepath.Join(t.TempDir(), "plain-dir")
	dbCfg := DefaultConfig()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	if _, err := Open(path, dbCfg); err == nil {
		t.Fatal("expected Open() to fail for uninitialized db directory")
	}
}

func TestCreateCloseThenOpen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "db")
	dbCfg := DefaultConfig()

	db, err := Create(path, dbCfg)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if err := db.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	reopened, err := Open(path, dbCfg)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer reopened.Close()
}

func TestCloseIsIdempotent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "db")
	dbCfg := DefaultConfig()

	db, err := Create(path, dbCfg)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if err := db.Close(); err != nil {
		t.Fatalf("first Close() error = %v", err)
	}
	if err := db.Close(); err != nil {
		t.Fatalf("second Close() error = %v", err)
	}
}

func TestCreateUsesExpandedConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "db")
	cfg := DefaultConfig()
	cfg.LRUSize = "256MB"
	cfg.BlockSize = "16KB"
	cfg.RateBytesPerSec = "200MB"
	cfg.WriteBufferNumber = 8

	db, err := Create(path, cfg)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	defer db.Close()
}

func TestCreateOrOpenCreatesThenOpens(t *testing.T) {
	path := filepath.Join(t.TempDir(), "db")
	cfg := DefaultConfig()

	db, err := createOrOpen(path, cfg)
	if err != nil {
		t.Fatalf("first createOrOpen() error = %v", err)
	}
	if err := db.Close(); err != nil {
		t.Fatalf("first Close() error = %v", err)
	}

	reopened, err := createOrOpen(path, cfg)
	if err != nil {
		t.Fatalf("second createOrOpen() error = %v", err)
	}
	defer reopened.Close()
}

func TestCreateWritesOptionsFileMatchingConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "db")
	cfg := DefaultConfig()
	cfg.WriteBufferSize = "32MiB"
	cfg.WriteBufferNumber = 5
	cfg.WriteBufferNumberToMerge = 2
	cfg.CompressionType = "lz4"
	cfg.BlockSize = "8KiB"
	cfg.LRUSize = "64MiB"

	db, err := Create(path, cfg)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if err := db.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	optionFiles, err := filepath.Glob(filepath.Join(path, "OPTIONS-*"))
	if err != nil {
		t.Fatalf("Glob() error = %v", err)
	}
	if len(optionFiles) == 0 {
		t.Fatal("expected OPTIONS file to exist after closing rocksdb")
	}
	contentBytes, err := os.ReadFile(optionFiles[len(optionFiles)-1])
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	content := string(contentBytes)

	writeBufferSize, err := nativeParseByteSize(cfg.WriteBufferSize)
	if err != nil {
		t.Fatalf("nativeParseByteSize(write_buffer_size) error = %v", err)
	}
	blockSize, err := nativeParseByteSize(cfg.BlockSize)
	if err != nil {
		t.Fatalf("nativeParseByteSize(block_size) error = %v", err)
	}

	wants := []string{
		"max_write_buffer_number=" + strconv.Itoa(cfg.WriteBufferNumber),
		"min_write_buffer_number_to_merge=" + strconv.Itoa(cfg.WriteBufferNumberToMerge),
		"write_buffer_size=" + strconv.FormatUint(writeBufferSize, 10),
		"compression=kLZ4Compression",
		"block_size=" + strconv.FormatUint(blockSize, 10),
		"no_block_cache=false",
	}
	for _, want := range wants {
		if !strings.Contains(content, want) {
			t.Fatalf("expected OPTIONS file %q to contain %q", optionFiles[len(optionFiles)-1], want)
		}
	}
}
