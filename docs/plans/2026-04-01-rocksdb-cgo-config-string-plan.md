# RocksDB CGO Config String Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Rework the low-level CGO layer so RocksDB C handles are wrapped in pointer-only Go objects with `Close()` methods, expand `rocksdb.Config` to the requested surface, and make `Create`/`Open` build RocksDB options from a config string via `rocksdb_get_options_from_string`.

**Architecture:** Keep `internal/cgo` as the only place that touches raw `rocksdb/c.h` pointers and lifetime management. Public `rocksdb/` code should only assemble config text, call the low-level wrappers, and return Go-facing `DB` wrappers with explicit `Close()` semantics.

**Tech Stack:** Go 1.25.3, CGO, vendored RocksDB C API, prebuilt static archives, Go tests.

### Task 1: Write failing tests for expanded config and config-string-based open/create

**Files:**
- Modify: `rocksdb/options_test.go`
- Modify: `rocksdb/db_test.go`

**Step 1: Write the failing tests**

Add tests that require:
- `Config` to expose the requested fields
- `DefaultConfig()` to initialize the critical new fields
- config validation to reject invalid required values
- `Create` to succeed using the config-string-based path
- `Open` to succeed using the same path with different create flags

**Step 2: Run the tests to verify they fail**

Run:

```bash
CGO_ENABLED=1 go test ./rocksdb -run 'Test(Create|Open|Default|Validate)' -v
```

Expected: FAIL because the config surface and low-level options-string assembly do not yet exist.

### Task 2: Introduce low-level pointer wrappers with Close methods

**Files:**
- Modify: `internal/cgo/options.go`
- Modify: `internal/cgo/db.go`
- Modify: `internal/cgo/errors.go`
- Create: `internal/cgo/cache.go`
- Create: `internal/cgo/block_based_options.go`
- Create: `internal/cgo/rate_limiter.go`
- Create: `internal/cgo/config_options.go`

**Step 1: Write the failing low-level tests**

Add minimal tests or package-local assertions that require all pointer wrappers to expose `Close()` and not panic when double-closed.

**Step 2: Run tests to verify failure**

Run:

```bash
CGO_ENABLED=1 go test ./internal/cgo -v
```

Expected: FAIL because the wrapper types do not exist yet.

**Step 3: Write the minimal implementation**

Implement pointer-only wrappers for the needed objects first:
- options wrapper
- block-based options wrapper
- cache wrapper
- rate limiter wrapper
- config-options wrapper if needed by `rocksdb_get_options_from_string`
- db wrapper

Every wrapper must have `Close()`.

**Step 4: Re-run tests**

Run:

```bash
CGO_ENABLED=1 go test ./internal/cgo -v
```

Expected: PASS or fail only at the next missing behavior layer.

### Task 3: Expand Config and build RocksDB option strings

**Files:**
- Modify: `rocksdb/config.go`
- Modify: `rocksdb/options.go`

**Step 1: Implement the requested `Config` surface**

Add the exact fields requested by the user, keeping the YAML tags.

**Step 2: Add config normalization and validation**

Implement defaulting and validation for required string/int fields so bad input fails before touching RocksDB.

**Step 3: Build RocksDB options text**

Generate an options string matching RocksDB option-string expectations for the supported fields. Keep the first version minimal and only map what is actually used by `Create`/`Open` today.

**Step 4: Re-run tests**

Run:

```bash
CGO_ENABLED=1 go test ./rocksdb -run 'Test(Default|Validate)' -v
```

Expected: PASS.

### Task 4: Rework Create/Open around `rocksdb_get_options_from_string`

**Files:**
- Modify: `rocksdb/db.go`

**Step 1: Replace direct setter-based option construction**

Make `Create(path, cfg)` and `Open(path, cfg)` assemble the config string and invoke the low-level wrapper that loads options through `rocksdb_get_options_from_string`.

**Step 2: Preserve create/open flag differences**

Ensure `Create` and `Open` differ only in create-if-missing / error-if-exists behavior.

**Step 3: Ensure all temporary low-level objects are closed**

Every helper object created during option assembly must be explicitly closed.

**Step 4: Re-run targeted tests**

Run:

```bash
CGO_ENABLED=1 go test ./rocksdb -run 'Test(Create|Open)' -v
```

Expected: PASS.

### Task 5: Final verification and manual QA

**Files:**
- No additional files required unless fixes are needed

**Step 1: Run diagnostics**

Run LSP diagnostics on all changed Go files and require zero errors.

**Step 2: Run full verification**

Run:

```bash
CGO_ENABLED=1 go test ./internal/cgo -v
CGO_ENABLED=1 go test ./rocksdb -v
CGO_ENABLED=1 go test ./...
CGO_ENABLED=1 go build ./...
```

Expected: PASS.

**Step 3: Manual QA**

Run a real create/open/close lifecycle test and confirm it succeeds.

```bash
CGO_ENABLED=1 go test ./rocksdb -run 'TestCreateCloseThenOpen' -v
```

Expected: PASS.
