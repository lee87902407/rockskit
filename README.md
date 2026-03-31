# rockskit

`rockskit` is a Go wrapper around RocksDB.

The current repository already includes the native RocksDB static libraries required by the Go CGO layer for these platforms only:

- `darwin/arm64`
- `linux/amd64`
- `linux/arm64`

Windows is not supported.

For normal users of this module, the goal is simple: after this repository is published, you only need to add the Go module dependency and build your Go program. You do **not** need to manually compile or install RocksDB first.

## Supported platform matrix

| OS | Architecture | Supported | Native artifact directory |
|---|---|---:|---|
| macOS | arm64 | yes | `native/darwin_arm64/` |
| Linux | amd64 | yes | `native/linux_amd64/` |
| Linux | arm64 | yes | `native/linux_arm64/` |
| macOS | amd64 | no | n/a |
| Windows | any | no | n/a |

## How downstream users use this module

If you are only consuming the published module, you do not need to run any of the maintainer build scripts below.

Typical usage:

```bash
go get github.com/yanjie/rockskit
```

Then in Go:

```go
package main

import (
	"log"

	"github.com/yanjie/rockskit/rocksdb"
)

func main() {
	envCfg := rocksdb.DefaultEnvConfig()
	dbCfg := rocksdb.DefaultConfig()

	db, err := rocksdb.Create("./example-db", envCfg, dbCfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}
```

The native static archive and public RocksDB headers are already expected to be present in this repository under `native/`, so downstream users do not need a separate RocksDB installation step.

## Repository layout for native artifacts

This repository stores prebuilt native artifacts in these directories:

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

The Go CGO bridge chooses the correct directory using platform-specific files in `internal/cgo/`.

## Maintainer build and refresh flow

This section is for maintainers who need to regenerate the prebuilt native artifacts before committing or publishing a new version.

### 1. Prerequisites

#### macOS side

The current host machine is expected to be macOS arm64.

Required local tools:

- Apple command line tools / Xcode command line tools
- `cmake`
- `docker` command available in shell (Docker or OrbStack are both fine)

Check them:

```bash
uname -s
uname -m
command -v cmake
command -v docker
xcode-select -p
```

Expected on the current supported maintainer host:

- `Darwin`
- `arm64`

#### Linux builder images

Linux builds are performed inside Docker-compatible containers.

This workflow depends on the image:

```text
harbor.clever8790.top/yanjie/debian:trixie-slim
```

The local machine should already have both Linux variants cached locally:

- one `linux/amd64` image
- one `linux/arm64` image

You can inspect local images with:

```bash
docker images --digests | grep 'harbor.clever8790.top/yanjie/debian'
```

The build script automatically scans the local image store and selects the matching builder image ID by architecture.

### 2. RocksDB source location

RocksDB source lives in the submodule:

```text
third_party/rocksdb/
```

Before rebuilding artifacts, make sure the submodule is initialized and on the intended version:

```bash
git submodule update --init --recursive
git -C third_party/rocksdb rev-parse HEAD
git -C third_party/rocksdb describe --tags --always
```

### 3. Debian apt mirror handling inside Linux builders

The Linux container build step may need to rewrite Debian sources to USTC mirrors.

The build script already applies these commands conditionally inside the container when the corresponding files exist:

```bash
sed -i 's/deb.debian.org/mirrors.ustc.edu.cn/g' /etc/apt/sources.list
sed -i 's/deb.debian.org/mirrors.ustc.edu.cn/g' /etc/apt/sources.list.d/debian.sources
```

That means maintainers usually do not need to run them manually.

### 4. Full native artifact rebuild

The main maintainer script is:

```bash
./scripts/build-rocksdb-prebuilt.sh
```

By default it builds exactly these targets:

- `darwin_arm64`
- `linux_amd64`
- `linux_arm64`

You can also rebuild only a subset:

```bash
./scripts/build-rocksdb-prebuilt.sh darwin_arm64
./scripts/build-rocksdb-prebuilt.sh linux_amd64 linux_arm64
```

### 5. What the build script actually does

#### For macOS arm64

The script performs a local CMake build from `third_party/rocksdb` and copies the output into `native/darwin_arm64/`.

The effective shape of the local build is:

```bash
cmake -S third_party/rocksdb -B third_party/rocksdb/build_prebuilt_darwin_arm64 \
  -DCMAKE_BUILD_TYPE=Release \
  -DCMAKE_POSITION_INDEPENDENT_CODE=ON \
  -DCMAKE_OSX_ARCHITECTURES=arm64 \
  -DPORTABLE=1 \
  -DWITH_TESTS=OFF \
  -DWITH_GFLAGS=OFF \
  -DWITH_BENCHMARK_TOOLS=OFF \
  -DWITH_TOOLS=OFF \
  -DWITH_MD_LIBRARY=OFF \
  -DWITH_RUNTIME_DEBUG=OFF \
  -DROCKSDB_BUILD_SHARED=OFF \
  -DWITH_SNAPPY=OFF \
  -DWITH_LZ4=OFF \
  -DWITH_ZLIB=OFF \
  -DWITH_LIBURING=OFF \
  -DWITH_JEMALLOC=OFF \
  -DWITH_NUMA=OFF \
  -DWITH_TBB=OFF \
  -DUSE_RTTI=ON \
  -DWITH_TRACE_TOOLS=OFF \
  -DWITH_CORE_TOOLS=OFF \
  -DUSE_FOLLY=OFF \
  -DUSE_COROUTINES=OFF \
  -DWITH_ZSTD=OFF \
  -DWITH_BZ2=OFF

cmake --build third_party/rocksdb/build_prebuilt_darwin_arm64 --target rocksdb -j"$(sysctl -n hw.ncpu)"
```

After build completes, the script copies:

- `third_party/rocksdb/build_prebuilt_darwin_arm64/librocksdb.a` -> `native/darwin_arm64/librocksdb.a`
- `third_party/rocksdb/include/rocksdb/...` -> `native/darwin_arm64/include/rocksdb/...`

#### For Linux amd64 and arm64

The script uses Docker/OrbStack to run a Debian builder container.

For each Linux target it:

1. picks the locally cached Debian image whose architecture matches the target
2. mounts the repository into `/src`
3. mounts the target native output directory into `/out`
4. rewrites apt source hosts to `mirrors.ustc.edu.cn` when needed
5. installs build tools with `apt-get`
6. runs a Release CMake build with RocksDB shared libs disabled
7. copies `librocksdb.a` into the target output directory
8. copies RocksDB public headers into the matching `native/<platform>/include/rocksdb/`

The effective Linux build shape is:

```bash
docker run --rm --platform linux/amd64 \
  -v "$PWD:/src" \
  -v "$PWD/native/linux_amd64:/out" \
  <local-amd64-image-id> \
  bash -lc '
    set -euo pipefail
    if [ -f /etc/apt/sources.list ]; then
      sed -i "s/deb.debian.org/mirrors.ustc.edu.cn/g" /etc/apt/sources.list
    fi
    if [ -f /etc/apt/sources.list.d/debian.sources ]; then
      sed -i "s/deb.debian.org/mirrors.ustc.edu.cn/g" /etc/apt/sources.list.d/debian.sources
    fi
    apt-get update
    apt-get install -y build-essential cmake pkg-config
    cmake -S /src/third_party/rocksdb -B /tmp/rocksdb-build \
      -DCMAKE_BUILD_TYPE=Release \
      -DCMAKE_POSITION_INDEPENDENT_CODE=ON \
      -DPORTABLE=1 \
      -DWITH_TESTS=OFF \
      -DWITH_GFLAGS=OFF \
      -DWITH_BENCHMARK_TOOLS=OFF \
      -DWITH_TOOLS=OFF \
      -DWITH_MD_LIBRARY=OFF \
      -DWITH_RUNTIME_DEBUG=OFF \
      -DROCKSDB_BUILD_SHARED=OFF \
      -DWITH_SNAPPY=OFF \
      -DWITH_LZ4=OFF \
      -DWITH_ZLIB=OFF \
      -DWITH_LIBURING=OFF \
      -DWITH_JEMALLOC=OFF \
      -DWITH_NUMA=OFF \
      -DWITH_TBB=OFF \
      -DUSE_RTTI=ON \
      -DWITH_TRACE_TOOLS=OFF \
      -DWITH_CORE_TOOLS=OFF \
      -DUSE_FOLLY=OFF \
      -DUSE_COROUTINES=OFF \
      -DWITH_ZSTD=OFF \
      -DWITH_BZ2=OFF
    cmake --build /tmp/rocksdb-build --target rocksdb -j"$(nproc)"
    cp /tmp/rocksdb-build/librocksdb.a /out/librocksdb.a
  '
```

`linux/arm64` uses the same flow with the local arm64 Debian builder image.

### 6. Verify native artifact output

After the script runs, verify all three supported directories exist:

```bash
ls native
```

Expected:

```text
darwin_arm64
linux_amd64
linux_arm64
```

Verify required files:

```bash
test -f native/darwin_arm64/librocksdb.a
test -f native/linux_amd64/librocksdb.a
test -f native/linux_arm64/librocksdb.a

test -f native/darwin_arm64/include/rocksdb/c.h
test -f native/linux_amd64/include/rocksdb/c.h
test -f native/linux_arm64/include/rocksdb/c.h
```

### 7. Verify Go platform selection logic

Run:

```bash
go test ./internal/cgo -run 'TestNativeArtifactDir|TestUnsupportedPlatform|TestRequiredNativeArtifactsExist' -v
```

This validates:

- supported platform to artifact-directory mapping
- unsupported platform rejection
- required artifact existence for the three supported targets

### 8. Verify the current macOS arm64 CGO link path

On the current machine, run:

```bash
CGO_ENABLED=1 go test ./rocksdb -run 'Test(Create|Open|Close|Default|Validate)' -v
CGO_ENABLED=1 go test ./rocksdb -run 'TestCreateCloseThenOpen' -v
CGO_ENABLED=1 go test ./...
CGO_ENABLED=1 go build ./...
```

These commands confirm that:

- `internal/cgo` selects `native/darwin_arm64/`
- `rocksdb.Create/Open/Close` still work against the shipped static archive
- the whole repository remains buildable on the current platform

### 9. Release checklist

Before publishing a new module version, maintainers should verify this sequence:

```bash
git submodule update --init --recursive
git -C third_party/rocksdb describe --tags --always

./scripts/build-rocksdb-prebuilt.sh

go test ./internal/cgo -run 'TestNativeArtifactDir|TestUnsupportedPlatform|TestRequiredNativeArtifactsExist' -v
CGO_ENABLED=1 go test ./...
CGO_ENABLED=1 go build ./...
```

Then review the native directories and commit:

- updated `native/` artifacts
- any `internal/cgo/` build-selection changes
- any script changes
- any docs updates

### 10. Important notes

1. macOS only ships `arm64` artifacts now. `darwin/amd64` is intentionally unsupported.
2. Linux ships both `amd64` and `arm64` artifacts.
3. Linux builds rely on the local Debian builder images already present in Docker/OrbStack.
4. If the local Debian builder image set changes, the maintainer must ensure both architectures are still available locally before regenerating artifacts.
5. Downstream users should not run the maintainer build script. They should only consume the published module.
