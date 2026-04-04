# PROJECT KNOWLEDGE BASE

**Generated:** 2026-03-31 Asia/Shanghai
**Commit:** N/A (not a git repo)
**Branch:** N/A (not a git repo)

## OVERVIEW

`rockskit` should evolve into a layered Go wrapper around RocksDB: ANSI C shim for common RocksDB calls, CGO bridge for type/error/memory handling, and a Go-facing package for the public API.
Current workspace still contains only a root-level pure Go placeholder implementation backed by `store.json`; treat that as a temporary prototype, not the final storage engine design.

## CURRENT STATE

```text
rockskit/
├── storage.go        # Placeholder pure-Go storage implementation
├── types.go          # Current public DTO-style structs and enums
├── storage_test.go   # Single CRUD/scan test against the placeholder backend
├── go.mod            # Module path + Go version
├── go.sum            # Empty; no external deps yet
└── .idea/            # IDE metadata; not product code
```

- All shipping code still lives at repo root.
- There are no `*.c`, `*.h`, `#cgo`, `import "C"`, RocksDB link flags, or build scripts yet.
- `storage.go` is the only behavior implementation today; it uses an in-memory map plus JSON persistence.
- `types.go` and the exported method surface are the closest thing to a stable API contract right now.

## TARGET STRUCTURE

```text
rockskit/
├── AGENTS.md
├── go.mod
├── README.md
├── c/
│   ├── include/
│   │   └── rockskit.h          # ANSI C wrapper exposed to CGO layer
│   └── src/
│       ├── db.c                # open/close/get/put/delete wrappers
│       ├── batch.c             # write batch wrappers
│       ├── iterator.c          # iterator/prefix/range helpers
│       ├── options.c           # rocksdb options wrappers
│       ├── error.c             # errptr helpers and message ownership
│       └── callbacks.c         # comparator/merge/filter callback glue only
├── internal/
│   └── cgo/
│       ├── bridge.go           # import "C", build tags, CFLAGS/LDFLAGS
│       ├── db.go               # Go <-> C DB bridge
│       ├── batch.go            # Go <-> C batch bridge
│       ├── iterator.go         # iterator bridge
│       ├── options.go          # options bridge and lifetime tracking
│       ├── errors.go           # char** errptr -> Go error conversion
│       └── bytes.go            # []byte/string <-> C buffer helpers
├── rocksdb/
│   ├── db.go                   # Public Go API: Open/Get/Put/Delete/Close
│   ├── batch.go                # Public WriteBatch API
│   ├── iterator.go             # Public iterator/scan API
│   ├── options.go              # Public Go-style options
│   ├── types.go                # Public Go structs/enums/results
│   └── errors.go               # Public error values and wrappers
├── examples/
│   └── basic/
│       └── main.go             # Minimal runnable example
├── scripts/
│   ├── build-rocksdb.sh        # Build or verify native RocksDB deps
│   └── test.sh                 # Unified local test entrypoint
└── tests/
    ├── db_test.go              # Public API integration tests
    ├── batch_test.go           # Batch behavior tests
    ├── iterator_test.go        # Scan/iterator tests
    └── cgo_smoke_test.go       # CGO smoke tests and leak guards
```

## LAYER RESPONSIBILITIES

| Layer | Path | Owns | Must not own |
|------|------|------|------|
| ANSI C shim | `c/` | Minimal wrapper functions, callback glue, errptr normalization | Go-visible API semantics, business rules, Go memory management |
| CGO bridge | `internal/cgo/` | `import "C"`, pointer conversions, lifetime tracking, `char** errptr` handling, link flags | Public package API, domain logic, ad-hoc unsafe scattered across repo |
| Public Go API | `rocksdb/` | Go-style types, exported methods, validation, ergonomic wrappers, contextual errors | Raw C calls, direct `unsafe`, link configuration |

## WHERE TO LOOK

| Task | Current location | Future location | Notes |
|------|------------------|----------------|-------|
| Understand placeholder API | `storage.go:20` | `rocksdb/db.go` | `Open` should remain the construction entry point |
| Review current public structs | `types.go:3` | `rocksdb/types.go` | Preserve Go-friendly DTO surface while swapping backend |
| See current CRUD/scan behavior | `storage.go:44` | `rocksdb/db.go`, `rocksdb/iterator.go` | Existing semantics are the migration baseline |
| Find C wrapper boundary | n/a | `c/include/rockskit.h`, `c/src/*.c` | Only common RocksDB calls and callback glue belong here |
| Find CGO bridge boundary | n/a | `internal/cgo/*.go` | Every `import "C"` should be isolated here |
| Inspect end-to-end behavior | `storage_test.go:5` | `tests/*.go` | Public API tests should target Go layer, not raw C layer |

## CURRENT CODE MAP

| Symbol | Type | Location | Role |
|--------|------|----------|------|
| `Storage` | struct | `storage.go:13` | Placeholder in-memory map + file path + RWMutex |
| `Open` | func | `storage.go:20` | Current factory; good public entrypoint to preserve |
| `(*Storage).Get` | method | `storage.go:44` | Current read contract with defensive copy |
| `(*Storage).Put` | method | `storage.go:54` | Current write contract; sync persistence placeholder |
| `(*Storage).Delete` | method | `storage.go:61` | Current delete contract |
| `(*Storage).WriteBatch` | method | `storage.go:68` | Current batch mutation contract |
| `(*Storage).PrefixScan` | method | `storage.go:84` | Current prefix query baseline |
| `(*Storage).RangeScan` | method | `storage.go:100` | Current range query baseline |
| `OpType` | type | `types.go:3` | Current batch op enum |
| `KVOp` | struct | `types.go:10` | Current batch payload |
| `KeyValue` | struct | `types.go:16` | Current scan result payload |

## CONVENTIONS

- Prefer one public Go package (`rocksdb/` or equivalent main package) and hide all CGO details under `internal/cgo/`.
- Keep the Go-visible API surface small and stable; do not mirror the full RocksDB C API unless the project actually uses it.
- Isolate all `import "C"`, `#cgo`, `unsafe`, and pointer conversion helpers inside `internal/cgo/`.
- 代码注释统一使用中文，避免在同一项目中混用中英文注释。
- 计划、设计、实现说明等相关文档统一使用中文编写，避免在同一项目中混用中英文文档。
- Use the C layer only for repeated RocksDB call normalization and callback glue; simple CRUD wrappers do not need extra abstraction beyond what reduces duplication.
- C shim 只封装包含实质性逻辑（如多步调用、错误归一化、资源所有权转移）的 RocksDB 操作；仅做单行透传的函数不放入 C shim，CGO 层直接调用 RocksDB 的 `c.h` 中的原始函数。
- 如果某个 Go 函数对 CGO 的调用次数超过 2 次，应将多次 C 调用合并为一个 C shim 函数，减少 CGO 訨开开销。
- Preserve defensive copy semantics for user-visible `[]byte` values when crossing from native buffers into Go-owned memory.
- Centralize `char** errptr` extraction and freeing in one CGO helper path; every failing RocksDB call should become a normal Go `error`.
- Callback-capable features such as comparator, merge operator, compaction filter, and slice transform should use C glue plus Go-side registries; never pass raw Go pointers into C callbacks.

## ANTI-PATTERNS (THIS PROJECT)

- Do not claim `c/`, `internal/cgo/`, or `rocksdb/` already exist until the files are actually created; this document describes a target architecture as well as current reality.
- Do not place `import "C"` in public Go API files; CGO belongs in `internal/cgo/` only.
- Do not spread `unsafe.Pointer`, `C.CString`, `C.CBytes`, or `C.GoBytes` logic across multiple packages; keep conversions centralized.
- Do not forget to free native allocations, RocksDB-returned buffers, or `char** errptr` strings; encode ownership at the bridge boundary.
- Do not wrap every RocksDB symbol “just because it exists”; keep the boundary minimal and driven by actual product needs.
- Do not pass Go pointers directly into long-lived C state or callback registrations; use registries/index handles for callback-capable features.
- Do not expose RocksDB iterator or slice lifetimes directly to end users without a Go-owned cleanup contract.
- Do not add `AGENTS.md` under `.idea/`, vendor-like directories, or generated native build output.

## UNIQUE STYLES

- Public API should feel like idiomatic Go even if the implementation is native and CGO-backed.
- The C layer should be thin and boring; most readability should live in the Go layer, not in C macros or complex native helper stacks.
- Native resource ownership must be visible in type and method design: `Close`, `Destroy`, `Free`, and iterator cleanup should be explicit and consistent.
- The current root-level Go prototype is useful as a semantics reference during migration, not as the intended final file layout.

## COMMANDS

```bash
go test ./...
go vet ./...
go build ./...

# future native/CGO workflow
./scripts/build-rocksdb.sh
./scripts/test.sh
CGO_ENABLED=1 go test ./...
```

## NOTES

- Current workspace is not a git repo, so git-derived metadata and workflows do not apply yet.
- `go.mod` currently declares `go 1.25.3`; verify local toolchain support before adding CGO-specific build assumptions.
- Official RocksDB C API uses opaque pointers plus `char** errptr`; model the bridge around that ownership pattern.
- Mature wrappers such as `gorocksdb` and `grocksdb` keep direct C calls and conversion helpers localized instead of leaking them into public API files.
- Current hierarchy decision remains root `AGENTS.md` only; create child `AGENTS.md` files later only when `c/`, `internal/cgo/`, or `rocksdb/` become real directories with distinct local rules.
