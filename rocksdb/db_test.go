package rocksdb

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateCreatesNewDB(t *testing.T) {
	path := filepath.Join(t.TempDir(), "db")
	envCfg := DefaultEnvConfig()
	dbCfg := DefaultConfig()

	db, err := Create(path, envCfg, dbCfg)
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
	envCfg := DefaultEnvConfig()
	dbCfg := DefaultConfig()

	db, err := Create(path, envCfg, dbCfg)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	defer db.Close()

	other, err := Create(path, envCfg, dbCfg)
	if err == nil {
		other.Close()
		t.Fatal("expected Create() to fail for existing db")
	}
}

func TestOpenFailsWhenDBIsMissing(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing")
	envCfg := DefaultEnvConfig()
	dbCfg := DefaultConfig()

	if _, err := Open(path, envCfg, dbCfg); err == nil {
		t.Fatal("expected Open() to fail for missing db")
	}
}

func TestOpenFailsForUninitializedDirectory(t *testing.T) {
	path := filepath.Join(t.TempDir(), "plain-dir")
	envCfg := DefaultEnvConfig()
	dbCfg := DefaultConfig()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	if _, err := Open(path, envCfg, dbCfg); err == nil {
		t.Fatal("expected Open() to fail for uninitialized db directory")
	}
}

func TestCreateCloseThenOpen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "db")
	envCfg := DefaultEnvConfig()
	dbCfg := DefaultConfig()

	db, err := Create(path, envCfg, dbCfg)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if err := db.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	reopened, err := Open(path, envCfg, dbCfg)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer reopened.Close()
}

func TestCloseIsIdempotent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "db")
	envCfg := DefaultEnvConfig()
	dbCfg := DefaultConfig()

	db, err := Create(path, envCfg, dbCfg)
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
