package rockskit

import "testing"

func TestStorageCRUDAndScans(t *testing.T) {
	s, err := Open(t.TempDir())
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer s.Close()
	if err := s.Put([]byte("doc\x00users\x001"), []byte(`{"name":"alice"}`)); err != nil {
		t.Fatalf("Put() error = %v", err)
	}
	if err := s.Put([]byte("doc\x00users\x002"), []byte(`{"name":"bob"}`)); err != nil {
		t.Fatalf("Put() error = %v", err)
	}
	rows, err := s.PrefixScan([]byte("doc\x00users\x00"))
	if err != nil {
		t.Fatalf("PrefixScan() error = %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(rows))
	}
	if err := s.Delete([]byte("doc\x00users\x001")); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	value, err := s.Get([]byte("doc\x00users\x001"))
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if value != nil {
		t.Fatalf("expected nil value after delete, got %q", value)
	}
}
