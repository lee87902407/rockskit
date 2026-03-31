package rockskit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

type Storage struct {
	path string
	file string
	mu   sync.RWMutex
	data map[string][]byte
}

func Open(path string) (*Storage, error) {
	cleanPath := filepath.Clean(path)
	if err := os.MkdirAll(cleanPath, 0o755); err != nil {
		return nil, fmt.Errorf("create data dir: %w", err)
	}
	file := filepath.Join(cleanPath, "store.json")
	s := &Storage{path: cleanPath, file: file, data: make(map[string][]byte)}
	if err := s.load(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Storage) Close() error {
	if s == nil {
		return nil
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.persistLocked()
}

func (s *Storage) Path() string { return s.path }

func (s *Storage) Get(key []byte) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.data[string(key)]
	if !ok {
		return nil, nil
	}
	return append([]byte(nil), value...), nil
}

func (s *Storage) Put(key, value []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[string(key)] = append([]byte(nil), value...)
	return s.persistLocked()
}

func (s *Storage) Delete(key []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, string(key))
	return s.persistLocked()
}

func (s *Storage) WriteBatch(ops []KVOp) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, op := range ops {
		switch op.Type {
		case OpPut:
			s.data[string(op.Key)] = append([]byte(nil), op.Value...)
		case OpDelete:
			delete(s.data, string(op.Key))
		default:
			return fmt.Errorf("unsupported batch op %q", op.Type)
		}
	}
	return s.persistLocked()
}

func (s *Storage) PrefixScan(prefix []byte) ([]KeyValue, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	results := make([]KeyValue, 0)
	keys := s.sortedKeysLocked()
	for _, rawKey := range keys {
		k := []byte(rawKey)
		if !bytes.HasPrefix(k, prefix) {
			continue
		}
		v := append([]byte(nil), s.data[rawKey]...)
		results = append(results, KeyValue{Key: k, Value: v})
	}
	return results, nil
}

func (s *Storage) RangeScan(start, end []byte) ([]KeyValue, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	results := make([]KeyValue, 0)
	keys := s.sortedKeysLocked()
	for _, rawKey := range keys {
		k := []byte(rawKey)
		if len(start) > 0 && bytes.Compare(k, start) < 0 {
			continue
		}
		if len(end) > 0 && bytes.Compare(k, end) >= 0 {
			break
		}
		v := append([]byte(nil), s.data[rawKey]...)
		results = append(results, KeyValue{Key: k, Value: v})
	}
	return results, nil
}

func (s *Storage) ScanAll() ([]KeyValue, error) {
	return s.RangeScan([]byte{}, nil)
}

func (s *Storage) load() error {
	content, err := os.ReadFile(s.file)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("load store file: %w", err)
	}
	if len(content) == 0 {
		return nil
	}
	if err := json.Unmarshal(content, &s.data); err != nil {
		return fmt.Errorf("decode store file: %w", err)
	}
	if s.data == nil {
		s.data = make(map[string][]byte)
	}
	return nil
}

func (s *Storage) persistLocked() error {
	content, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return fmt.Errorf("encode store file: %w", err)
	}
	if err := os.WriteFile(s.file, content, 0o644); err != nil {
		return fmt.Errorf("write store file: %w", err)
	}
	return nil
}

func (s *Storage) sortedKeysLocked() []string {
	keys := make([]string, 0, len(s.data))
	for key := range s.data {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
