# rockskit

`rockskit` is currently a root-level pure Go prototype backed by `store.json`.

The long-term target is a layered RocksDB wrapper:

- `c/` for the ANSI C shim over common RocksDB C APIs
- `internal/cgo/` for Go-to-C bridging, pointer conversion, and error handling
- `rocksdb/` for the public Go-facing API

Those future directories are scaffolded in this repository, but no RocksDB dependency or CGO integration exists yet.
