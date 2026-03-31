# RocksDB Prebuilt Artifacts Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Publish `rockskit` with platform-specific prebuilt RocksDB static libraries so downstream users can import the Go module and use it directly without manually compiling or installing RocksDB.

**Architecture:** Keep RocksDB source as a submodule for maintainers, but generate and commit prebuilt native artifacts under a repository-owned `native/` layout keyed by `GOOS/GOARCH`. Select the matching headers and static archive in `internal/cgo/` using platform-specific build files so consumers get the correct artifact automatically at build time.

**Tech Stack:** Go 1.25.3, CGO, RocksDB C API, prebuilt static archives, shell build script, darwin/linux amd64/arm64 only.

### Task 1: Define prebuilt artifact layout and platform selection tests

**Files:**
- Create: `internal/cgo/platform_test.go`
- Create: `internal/cgo/platform.go`
- Create: `native/.gitkeep`

**Step 1: Write the failing test**

Add tests for a deterministic artifact layout:
- `darwin/amd64` -> `native/darwin_amd64/`
- `darwin/arm64` -> `native/darwin_arm64/`
- `linux/amd64` -> `native/linux_amd64/`
- `linux/arm64` -> `native/linux_arm64/`
- unsupported platforms return a clear error

**Step 2: Run the test to verify it fails**

Run:

```bash
go test ./internal/cgo -run 'TestNativeArtifactDir|TestUnsupportedPlatform' -v
```

Expected: FAIL because the selector helpers do not exist yet.

**Step 3: Write the minimal implementation**

Implement platform selection helpers in `internal/cgo/platform.go` that:
- map `runtime.GOOS`/`runtime.GOARCH` to the native artifact directory name
- reject unsupported platforms explicitly

**Step 4: Re-run the tests**

Run:

```bash
go test ./internal/cgo -run 'TestNativeArtifactDir|TestUnsupportedPlatform' -v
```

Expected: PASS.

### Task 2: Replace local-source linking with platform-specific prebuilt linkage

**Files:**
- Modify: `internal/cgo/bridge.go`
- Create: `internal/cgo/bridge_darwin_amd64.go`
- Create: `internal/cgo/bridge_darwin_arm64.go`
- Create: `internal/cgo/bridge_linux_amd64.go`
- Create: `internal/cgo/bridge_linux_arm64.go`
- Create: `internal/cgo/unsupported.go`

**Step 1: Write the failing test**

Add a small test or build-target verification that expects the current package to compile without referencing `third_party/rocksdb/librocksdb.a` directly.

**Step 2: Run the test/build to verify it fails**

Run:

```bash
go test ./internal/cgo -v
```

Expected: FAIL or still reveal hardcoded linkage to the submodule path.

**Step 3: Write the minimal implementation**

Implement platform-split bridge files with exact `#cgo` directives per supported platform, each pointing at the corresponding `native/<platform>/librocksdb.a` and include directory. Keep unsupported platforms behind a build-tagged stub that returns a compile-time or runtime error.

**Step 4: Re-run the verification**

Run:

```bash
go test ./internal/cgo -v
```

Expected: PASS on the current machine.

### Task 3: Build and stage prebuilt RocksDB artifacts

**Files:**
- Modify: `scripts/build-rocksdb.sh`
- Create: `scripts/build-rocksdb-prebuilt.sh`
- Create: `native/darwin_amd64/include/rocksdb/c.h`
- Create: `native/darwin_amd64/librocksdb.a`
- Create: `native/darwin_arm64/include/rocksdb/c.h`
- Create: `native/darwin_arm64/librocksdb.a`
- Create: `native/linux_amd64/include/rocksdb/c.h`
- Create: `native/linux_amd64/librocksdb.a`
- Create: `native/linux_arm64/include/rocksdb/c.h`
- Create: `native/linux_arm64/librocksdb.a`

**Step 1: Write the failing test**

Add a file-existence test for required native artifact layout on supported platforms.

**Step 2: Run the test to verify it fails**

Run:

```bash
go test ./internal/cgo -run 'TestRequiredNativeArtifactsExist' -v
```

Expected: FAIL because `native/` artifacts are not populated yet.

**Step 3: Write the minimal implementation**

Implement a maintainer-facing script that:
- accepts target platform tuple(s)
- builds RocksDB from the submodule for the target
- copies stripped `librocksdb.a` and required public headers into `native/<platform>/`
- follows the reference script’s optimization-oriented CMake approach where it makes sense

Commit the produced artifacts for the four supported targets.

**Step 4: Re-run the test**

Run:

```bash
go test ./internal/cgo -run 'TestRequiredNativeArtifactsExist' -v
```

Expected: PASS.

### Task 4: Verify downstream-consumption behavior

**Files:**
- Modify: `README.md`
- Create: `examples/basic/main.go`

**Step 1: Write the failing end-to-end test**

Add an integration test or scripted verification that:
- uses the published module layout only
- does not invoke `scripts/build-rocksdb.sh`
- builds a tiny consumer program that imports the module and opens/closes a DB

**Step 2: Run it to verify the intended failure mode**

Run:

```bash
CGO_ENABLED=1 go build ./examples/basic
```

Expected: FAIL before native artifact selection or example wiring is complete.

**Step 3: Write the minimal implementation**

Add the example and update docs to explain:
- supported platforms
- that native artifacts are already shipped in the module
- that users do not need to install or compile RocksDB manually
- that Windows is unsupported

**Step 4: Re-run the verification**

Run:

```bash
CGO_ENABLED=1 go build ./examples/basic
CGO_ENABLED=1 go test ./...
```

Expected: PASS.

### Task 5: Final QA and release-readiness verification

**Files:**
- Modify: `README.md` if any caveats remain

**Step 1: Run diagnostics**

Run LSP diagnostics on all changed Go files and require zero errors.

**Step 2: Run full validation**

Run:

```bash
CGO_ENABLED=1 go test ./...
CGO_ENABLED=1 go build ./...
```

Expected: PASS.

**Step 3: Manual QA**

Run a real create/open/close scenario and an example build without calling any local RocksDB build script.

```bash
CGO_ENABLED=1 go test ./rocksdb -run 'TestCreateCloseThenOpen' -v
CGO_ENABLED=1 go build ./examples/basic
```

Expected: PASS.

**Step 4: Record release caveats**

Document repository size impact and supported platform matrix in `README.md`.
