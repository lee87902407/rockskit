package cgo

import (
	"bytes"
	"path/filepath"
	"testing"
)

func TestDBGetPinnedReturnsValueView(t *testing.T) {
	path := filepath.Join(t.TempDir(), "db")
	db, err := Create(path, SmallConfig())
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	defer db.Close()

	key := []byte("alpha")
	value := []byte("value-from-pinnable-slice")
	if err := db.Put(nil, key, value); err != nil {
		t.Fatalf("Put() error = %v", err)
	}

	pinned, err := db.GetPinned(nil, key)
	if err != nil {
		t.Fatalf("GetPinned() error = %v", err)
	}
	if pinned == nil {
		t.Fatal("GetPinned() returned nil for existing key")
	}

	got := pinned.value

	if !bytes.Equal(got, value) {
		t.Fatalf("value = %q, want %q", got, value)
	}

	if len(got) == 0 {
		t.Fatal("value returned empty data for existing key")
	}

	pinnedAgain, err := db.GetPinned(nil, key)
	if err != nil {
		t.Fatalf("second GetPinned() error = %v", err)
	}
	if pinnedAgain == nil {
		t.Fatal("second GetPinned() returned nil for existing key")
	}
	defer pinnedAgain.Close()

	gotAgain := pinnedAgain.value
	if !bytes.Equal(gotAgain, value) {
		t.Fatalf("second value = %q, want %q", gotAgain, value)
	}

	closedView := pinnedAgain.value
	if len(closedView) != len(value) {
		t.Fatalf("value len = %d, want %d", len(closedView), len(value))
	}
	pinnedAgain.Close()
	if afterClose := pinnedAgain.value; afterClose != nil {
		t.Fatalf("value after Close() = %v, want nil", afterClose)
	}
	pinned.Close()
	db.Close()
	if err := db.Put(nil, []byte("beta"), []byte("later")); err == nil {
		t.Fatal("expected Put() on closed db to fail")
	}
	if _, err := db.GetPinned(nil, key); err == nil {
		t.Fatal("expected GetPinned() on closed db to fail")
	}
}

func TestDBGetPinnedReturnsNilForMissingKey(t *testing.T) {
	path := filepath.Join(t.TempDir(), "db")
	db, err := Create(path, SmallConfig())
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	defer db.Close()

	pinned, err := db.GetPinned(nil, []byte("missing"))
	if err != nil {
		t.Fatalf("GetPinned() error = %v", err)
	}
	if pinned != nil {
		pinned.Close()
		t.Fatal("expected missing key to return nil pinnable slice")
	}
}

func TestDBGetPinnedPreservesEmptyValue(t *testing.T) {
	path := filepath.Join(t.TempDir(), "db")
	db, err := Create(path, SmallConfig())
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	defer db.Close()

	key := []byte("empty")
	if err := db.Put(nil, key, []byte{}); err != nil {
		t.Fatalf("Put() error = %v", err)
	}

	pinned, err := db.GetPinned(nil, key)
	if err != nil {
		t.Fatalf("GetPinned() error = %v", err)
	}
	if pinned == nil {
		t.Fatal("GetPinned() returned nil for empty value")
	}
	defer pinned.Close()

	got := pinned.value
	if got == nil {
		t.Fatal("value returned nil for empty value")
	}
	if len(got) != 0 {
		t.Fatalf("len(value) = %d, want 0", len(got))
	}
}
