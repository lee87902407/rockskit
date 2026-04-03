# rockskit

`rockskit` is a Go wrapper for RocksDB with bundled prebuilt native libraries. It is designed so downstream users can import the Go module directly without manually installing or compiling RocksDB first.

[简体中文文档](./README.zh-CN.md)

## 1. Supported Platforms

| OS | Architecture | Supported | Native Directory |
|---|---|---:|---|
| macOS | arm64 | Yes | `native/darwin_arm64/` |
| Linux | amd64 | Yes | `native/linux_amd64/` |
| Linux | arm64 | Yes | `native/linux_arm64/` |
| macOS | amd64 | No | N/A |
| Windows | Any | No | N/A |

The repository already includes the native headers and static archives required by the CGO layer.

## 2. Installation and Basic Usage

If you are a consumer of this library, you usually do not need to run any RocksDB build script manually.

```bash
go get github.com/lee87902407/rockskit@v0.0.1
```

```go
package main

import (
    "log"

    "github.com/lee87902407/rockskit/rocksdb"
)

func main() {
    cfg := rocksdb.DefaultConfig()

    db, err := rocksdb.Create("./example-db", cfg)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
}
```

## 3. Native Artifact Layout

```text
native/
├── darwin_arm64/
│   ├── include/rocksdb/...
│   └── librocksdb.a
├── linux_amd64/
│   ├── include/rocksdb/...
│   └── librocksdb.a
└── linux_arm64/
    ├── include/rocksdb/...
    └── librocksdb.a
```

The CGO bridge under `internal/cgo/` selects the correct directory by platform-specific build files.

## 4. Maintainer Workflow

This section is for maintainers who need to regenerate or update the prebuilt RocksDB artifacts.

### 4.1 Host Requirements

The recommended maintainer host is macOS arm64. The local machine should provide:

- Apple Command Line Tools / Xcode Command Line Tools
- `cmake`
- a working `docker` command (Docker Desktop or OrbStack)
- Homebrew

### 4.2 RocksDB Source Location

The vendored RocksDB source lives in:

```text
third_party/rocksdb/
```

Before rebuilding, make sure the submodule is initialized and pinned to the expected revision:

```bash
git submodule update --init --recursive
git -C third_party/rocksdb rev-parse HEAD
git -C third_party/rocksdb describe --tags --always
```

### 4.3 Linux Builder Image

Linux artifact generation depends on the local Debian builder image:

```text
harbor.clever8790.top/yanjie/debian:trixie-slim
```

The scripts are written for the following Linux targets:

- `linux/amd64`
- `linux/arm64`

### 4.4 Build Scripts

- `scripts/build-rocksdb.sh`: builds RocksDB locally from source and acts as the editable parameter entry for maintainers
- `scripts/build-rocksdb-prebuilt.sh`: the main script used to generate the final committed prebuilt artifacts
- `scripts/build-rocksdb-linux.sh`: the standalone Linux-in-container build script

Typical commands:

```bash
./scripts/build-rocksdb-prebuilt.sh
./scripts/build-rocksdb-prebuilt.sh darwin_arm64
./scripts/build-rocksdb-prebuilt.sh linux_amd64 linux_arm64
```

Explicit Linux builder image:

```bash
./scripts/build-rocksdb-prebuilt.sh \
  --linux-builder-image harbor.clever8790.top/yanjie/debian:trixie-slim \
  linux_amd64 linux_arm64
```

### 4.5 Platform-Specific jemalloc Policy

Current build policy is:

- macOS arm64 prebuilt artifacts do **not** use jemalloc
- Linux builds **require** `WITH_JEMALLOC=ON`

### 4.6 macOS arm64 Build Notes

Before macOS builds, the script may install missing dependencies with Homebrew:

```bash
brew install snappy lz4 zlib zstd bzip2 gflags || true
```

The macOS prebuilt flow keeps `WITH_JEMALLOC=OFF` and `USE_COROUTINES=OFF`.

### 4.7 Linux Build Notes

The Linux build flow enables:

- `USE_COROUTINES`
- `WITH_SNAPPY`
- `WITH_TBB`
- `USE_RTTI`
- `WITH_NUMA`
- `WITH_LZ4`
- `WITH_ZLIB`
- `WITH_LIBURING`
- `WITH_JEMALLOC`
- `WITH_ZSTD`
- `WITH_BZ2`

The Linux script also requires `WITH_JEMALLOC=ON`; otherwise it exits with an error.

## 5. Notes

- Windows is not supported
- `darwin/amd64` artifacts are not provided
- The module path is aligned with GitHub and can be consumed directly from `github.com/lee87902407/rockskit`
