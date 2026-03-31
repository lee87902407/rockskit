# RocksDB Lifecycle Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add `facebook/rocksdb` as a submodule and replace the placeholder `rocksdb.Create`/`Open`/`Close` implementation with a real RocksDB-backed lifecycle that builds options from migrated config.

**Architecture:** Keep the public lifecycle API in `rocksdb/`, isolate all native calls and ownership handling in `internal/cgo/`, and keep the ANSI C shim in `c/`. Migrate only the config fields needed to build options for lifecycle creation first, not the entire reference surface.

**Tech Stack:** Go 1.25.3, CGO, RocksDB C API, local `facebook/rocksdb` git submodule, shell build script for native library compilation.

### Task 1: Add RocksDB submodule and native build entrypoint

**Files:**
- Modify: `.gitmodules`
- Create: `third_party/rocksdb/` (git submodule)
- Create: `scripts/build-rocksdb.sh`

**Step 1: Add a failing verification command**

Run:

```bash
test -f third_party/rocksdb/include/rocksdb/c.h
```

Expected: FAIL because the submodule does not exist yet.

**Step 2: Add the submodule**

Run:

```bash
GIT_MASTER=1 git submodule add ssh://org-69631@github.com:facebook/rocksdb.git third_party/rocksdb
```

**Step 3: Add the build script**

Create `scripts/build-rocksdb.sh` that:
- validates `third_party/rocksdb` exists
- runs RocksDB's static library build in the submodule
- emits a stable output path for CGO to link against

**Step 4: Verify the build script works**

Run:

```bash
./scripts/build-rocksdb.sh
```

Expected: exits 0 and produces a buildable RocksDB library artifact.

### Task 2: Migrate the minimal lifecycle config surface

**Files:**
- Create: `rocksdb/options.go`
- Create: `rocksdb/config.go`
- Test: `rocksdb/options_test.go`

**Step 1: Write the failing test**

Add tests for the minimum migrated config surface:
- env config for cache/block size
- db config for write buffer size, max bytes for level base, target file size base, level0 compaction trigger, compression type
- defaulting/validation behavior for empty or invalid values

**Step 2: Run the test to verify it fails**

Run:

```bash
go test ./rocksdb -run 'Test(BuildOptions|DefaultConfig|ConfigValidation)' -v
```

Expected: FAIL because config types and option builders do not exist yet.

**Step 3: Write the minimal implementation**

Implement Go-facing config structs adapted from the reference:
- `EnvConfig` with `LRUSize`, `BlockSize`
- `Config` with only lifecycle-critical fields
- defaults that keep local development usable

**Step 4: Re-run the tests**

Run:

```bash
go test ./rocksdb -run 'Test(BuildOptions|DefaultConfig|ConfigValidation)' -v
```

Expected: PASS.

### Task 3: Add the C shim for lifecycle and option wiring

**Files:**
- Create: `c/include/rockskit.h`
- Create: `c/src/db.c`
- Create: `c/src/options.c`
- Create: `c/src/error.c`

**Step 1: Write the failing CGO-focused test**

Add a package-level lifecycle test that expects native create/open/close to work once the bridge exists.

**Step 2: Run the test to verify it fails**

Run:

```bash
CGO_ENABLED=1 go test ./rocksdb -run 'TestNative(Create|Open|Close)' -v
```

Expected: FAIL because the native bridge does not exist yet.

**Step 3: Write the minimal C implementation**

Implement only:
- option allocation/free helpers needed by lifecycle
- db create/open wrappers using RocksDB C API
- db close wrapper
- central errptr helpers

**Step 4: Re-run the test**

Run:

```bash
CGO_ENABLED=1 go test ./rocksdb -run 'TestNative(Create|Open|Close)' -v
```

Expected: bridge compiles or fails at the next missing layer, not at missing C wrappers.

### Task 4: Add the CGO bridge and replace placeholder lifecycle implementation

**Files:**
- Create: `internal/cgo/bridge.go`
- Create: `internal/cgo/db.go`
- Create: `internal/cgo/options.go`
- Create: `internal/cgo/errors.go`
- Modify: `rocksdb/db.go`
- Modify: `rocksdb/db_test.go`

**Step 1: Extend the failing lifecycle tests**

Cover these behaviors:
- `Create` builds options from config and creates a new RocksDB database
- `Create` fails when the target already exists
- `Open` opens only an existing initialized database
- `Close` releases the native handle and stays idempotent

**Step 2: Run the targeted tests and watch them fail**

Run:

```bash
CGO_ENABLED=1 go test ./rocksdb -run 'Test(Create|Open|Close)' -v
```

Expected: FAIL because `rocksdb/db.go` still uses the marker-file placeholder.

**Step 3: Replace the placeholder with the minimal native implementation**

Implement:
- `DB` holding native handle ownership
- `Create(path string, cfg *Config)` or equivalent config-aware constructor chosen to preserve the public API clearly
- `Open(...)`
- `Close()`

Apply option building only through the new config surface and native bridge.

**Step 4: Re-run targeted tests**

Run:

```bash
CGO_ENABLED=1 go test ./rocksdb -run 'Test(Create|Open|Close)' -v
```

Expected: PASS.

### Task 5: Full verification and manual QA

**Files:**
- Modify: `README.md` (only if build/run instructions changed materially)

**Step 1: Run diagnostics**

Run LSP diagnostics on every changed Go file and ensure zero errors.

**Step 2: Run package and full-project tests**

Run:

```bash
CGO_ENABLED=1 go test ./rocksdb -v
CGO_ENABLED=1 go test ./... 
go build ./...
```

Expected: PASS.

**Step 3: Manual QA**

Run a real create/open/close scenario against a temp dir.

```bash
CGO_ENABLED=1 go test ./rocksdb -run 'TestCreateCloseThenOpen' -v
```

Expected: PASS with actual native RocksDB lifecycle behavior.

**Step 4: Record residual risks**

If macOS linker requirements force extra local dependencies, document the exact command or environment needed in `README.md`.
