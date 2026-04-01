# rockskit

`rockskit` 是一个基于 RocksDB 的 Go 封装库。

当前仓库已经内置了 Go CGO 层所需的 RocksDB 预编译静态库，支持的平台只有以下三种：

- `darwin/arm64`
- `linux/amd64`
- `linux/arm64`

不支持 Windows，也不再提供 `darwin/amd64` 产物。

对最终使用者来说，目标很明确：当这个仓库发布后，只需要在自己的项目里通过 `go mod` 引入依赖即可使用，不需要提前手工安装或编译 RocksDB。

## 一、支持的平台矩阵

| 系统 | 架构 | 是否支持 | 原生库目录 |
|---|---|---:|---|
| macOS | arm64 | 是 | `native/darwin_arm64/` |
| Linux | amd64 | 是 | `native/linux_amd64/` |
| Linux | arm64 | 是 | `native/linux_arm64/` |
| macOS | amd64 | 否 | 无 |
| Windows | 任意 | 否 | 无 |

## 二、普通使用者如何使用

如果你只是这个库的使用者，而不是维护者，那么你通常不需要执行本文后面的 RocksDB 编译脚本。

正常使用方式如下：

```bash
go get github.com/yanjie/rockskit
```

然后在 Go 代码里直接使用：

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

仓库里已经包含 `native/` 下的头文件和静态库目录，因此下游用户不需要单独安装 RocksDB。

## 三、仓库中的原生产物布局

仓库中预编译产物目录如下：

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

Go 的 CGO 桥接层通过 `internal/cgo/` 下的按平台拆分文件选择对应目录。

## 四、维护者编译与更新流程

这一节是给仓库维护者看的，用来重新生成和更新预编译 RocksDB 产物。

### 4.1 本机前置条件

当前推荐维护者主机为 macOS arm64。

本机需要：

- Apple Command Line Tools / Xcode Command Line Tools
- `cmake`
- `docker` 命令可用（Docker Desktop 或 OrbStack 都可以）
- Homebrew（用于安装 macOS 侧缺失依赖）

检查命令：

```bash
uname -s
uname -m
command -v cmake
command -v docker
command -v brew
xcode-select -p
```

期望输出至少满足：

- `Darwin`
- `arm64`

### 4.2 RocksDB 源码位置

RocksDB 源码通过 submodule 放在：

```text
third_party/rocksdb/
```

重新构建前，建议先确认 submodule 已初始化，并且在你想要的版本上：

```bash
git submodule update --init --recursive
git -C third_party/rocksdb rev-parse HEAD
git -C third_party/rocksdb describe --tags --always
```

### 4.3 Linux 构建镜像

Linux 产物构建依赖本地已有的 Debian builder 镜像。

当前默认镜像名是：

```text
harbor.clever8790.top/yanjie/debian:trixie-slim
```

你已经说明本地已有：

- `linux/amd64` 版本
- `linux/arm64` 版本

脚本会优先尝试使用你传入的 `LINUX_BUILDER_IMAGE`，如果这个 tag 本身不是多架构 manifest，就会自动在本地镜像列表中查找同仓库对应架构的 image id 来运行。

查看本地镜像：

```bash
docker images --digests | grep 'harbor.clever8790.top/yanjie/debian'
```

### 4.4 Linux 容器内可选换源逻辑

Linux 构建脚本支持通过参数控制是否替换 Debian 源。

对应命令如下：

```bash
sed -i 's/deb.debian.org/mirrors.ustc.edu.cn/g' /etc/apt/sources.list
sed -i 's/deb.debian.org/mirrors.ustc.edu.cn/g' /etc/apt/sources.list.d/debian.sources
```

当前实现里，默认会在 Linux 容器构建时执行该替换逻辑；如果你不希望替换，可以显式传参关闭。

## 五、维护者使用的脚本说明

### 5.1 `scripts/build-rocksdb.sh`

这是本地直接在源码目录上构建 RocksDB 静态库的脚本。它主要作为“用户可编辑参数入口”，方便将来维护者通过修改脚本顶部变量生成自己需要的编译产出。

当前这个脚本暴露的主要参数包括：

- `PORTABLE`
- `DISABLE_WARNING_AS_ERROR`
- `USE_COROUTINES`
- `USE_RTTI`
- `WITH_GFLAGS`
- `WITH_SNAPPY`
- `WITH_LZ4`
- `WITH_ZLIB`
- `WITH_JEMALLOC`
- `WITH_ZSTD`
- `WITH_BZ2`
- `WITH_NUMA`
- `WITH_TBB`
- `WITH_LIBURING`
- `FORCE_AVX`
- `FORCE_SSE42`
- `EXTRA_CFLAGS`
- `EXTRA_CXXFLAGS`

当前脚本调用方式：

```bash
./scripts/build-rocksdb.sh
```

如果你想自定义，例如：

```bash
WITH_SNAPPY=ON WITH_ZSTD=ON USE_COROUTINES=1 ./scripts/build-rocksdb.sh
```

说明：

- 这个脚本走的是 RocksDB 自己的 `make static_lib`
- 对 `FORCE_AVX` 和 `FORCE_SSE42`，当前实现会把它们转换为额外编译参数（例如 `-mavx`、`-msse4.2`）
- 这类参数只适合 x86 系列平台，不适合 arm64

### 5.2 `scripts/build-rocksdb-prebuilt.sh`

这是维护者的主脚本，用于生成仓库里最终要提交的预编译产物。

默认构建目标为：

- `darwin_arm64`
- `linux_amd64`
- `linux_arm64`

默认调用：

```bash
./scripts/build-rocksdb-prebuilt.sh
```

只构建一部分目标：

```bash
./scripts/build-rocksdb-prebuilt.sh darwin_arm64
./scripts/build-rocksdb-prebuilt.sh linux_amd64 linux_arm64
```

显式传入 Linux 构建镜像：

```bash
./scripts/build-rocksdb-prebuilt.sh \
  --linux-builder-image harbor.clever8790.top/yanjie/debian:trixie-slim \
  linux_amd64 linux_arm64
```

显式关闭容器内换源：

```bash
./scripts/build-rocksdb-prebuilt.sh --replace-apt-source 0 linux_amd64
```

显式指定 Docker cache 目录：

```bash
./scripts/build-rocksdb-prebuilt.sh \
  --docker-cache-dir "$PWD/.cache/rocksdb-docker" \
  linux_amd64 linux_arm64
```

### 5.3 `scripts/build-rocksdb-linux.sh`

这是独立拆出来的 Linux 容器内构建脚本，只负责 Linux 下的 RocksDB 构建，不再内嵌在 `build-rocksdb-prebuilt.sh` 里。

它支持的核心输入参数来自环境变量，例如：

- `REPLACE_APT_SOURCE`
- `TARGET_ARCH`
- `PORTABLE`
- `USE_COROUTINES`
- `USE_RTTI`
- `WITH_SNAPPY`
- `WITH_TBB`
- `WITH_NUMA`
- `WITH_LZ4`
- `WITH_ZLIB`
- `WITH_LIBURING`
- `WITH_JEMALLOC`
- `WITH_ZSTD`
- `WITH_BZ2`
- `WITH_GFLAGS`
- `FORCE_AVX`
- `FORCE_SSE42`

当前实现中：

- Linux 下默认把 `USE_COROUTINES`、`WITH_SNAPPY`、`WITH_TBB`、`USE_RTTI`、`WITH_NUMA`、`WITH_LZ4`、`WITH_ZLIB`、`WITH_LIBURING`、`WITH_JEMALLOC`、`WITH_ZSTD`、`WITH_BZ2` 都设为 `ON`
- `FORCE_AVX`、`FORCE_SSE42` 默认也对 Linux 开启，但脚本会在 `arm64` 下自动收敛为 `OFF`，避免把 x86 指令集参数硬塞给 arm64

## 六、macOS arm64 构建说明

你要求在编译 mac 版本 RocksDB 时，将以下选项改为 `true`：

- `WITH_SNAPPY`
- `WITH_LZ4`
- `WITH_ZLIB`
- `WITH_JEMALLOC`
- `WITH_ZSTD`
- `WITH_BZ2`

当前实现已按此调整，并且在执行 mac 构建前会尝试通过 Homebrew 安装缺失组件：

```bash
brew install snappy lz4 zlib jemalloc zstd bzip2 gflags || true
```

mac 侧实际构建命令形态如下：

```bash
cmake -S third_party/rocksdb -B third_party/rocksdb/build_prebuilt_darwin_arm64 \
  -DCMAKE_BUILD_TYPE=Release \
  -DCMAKE_POSITION_INDEPENDENT_CODE=ON \
  -DCMAKE_OSX_ARCHITECTURES=arm64 \
  -DPORTABLE=1 \
  -DWITH_TESTS=OFF \
  -DWITH_GFLAGS=ON \
  -DWITH_BENCHMARK_TOOLS=OFF \
  -DWITH_TOOLS=OFF \
  -DWITH_MD_LIBRARY=OFF \
  -DWITH_RUNTIME_DEBUG=OFF \
  -DROCKSDB_BUILD_SHARED=OFF \
  -DWITH_SNAPPY=ON \
  -DWITH_LZ4=ON \
  -DWITH_ZLIB=ON \
  -DWITH_LIBURING=OFF \
  -DWITH_JEMALLOC=ON \
  -DWITH_NUMA=OFF \
  -DWITH_TBB=OFF \
  -DUSE_RTTI=ON \
  -DWITH_TRACE_TOOLS=OFF \
  -DWITH_CORE_TOOLS=OFF \
  -DUSE_FOLLY=OFF \
  -DUSE_COROUTINES=OFF \
  -DWITH_ZSTD=ON \
  -DWITH_BZ2=ON
```

这里没有把 `USE_COROUTINES` 打开，是因为当前仓库的 macOS 预编译流程并不打算把 Folly/coroutines 依赖链也一起带进来；而你特别要求读取源码内关于 `USE_COROUTINES` 的编译内容，所以这里把结论写清楚：

- upstream `CMakeLists.txt` 里，`USE_COROUTINES` 会引出更特殊的编译要求
- 它会影响 gflags/Folly 相关逻辑
- 它并不是“无额外代价直接打开”的一个选项

## 七、Linux 构建说明

你要求在 Linux 下把以下能力改为开启：

- `USE_COROUTINES`
- `WITH_SNAPPY`
- `WITH_TBB`
- `USE_RTTI`
- `WITH_NUMA`
- `FORCE_AVX`
- `FORCE_SSE42`
- `WITH_LZ4`
- `WITH_ZLIB`
- `WITH_LIBURING`
- `WITH_JEMALLOC`
- `WITH_ZSTD`
- `WITH_BZ2`

当前 Linux 独立构建脚本已经按这个方向实现，并且会安装你列出的依赖。

容器内安装依赖命令如下：

```bash
apt-get update
apt-get install -y \
  zip \
  build-essential \
  libjemalloc-dev \
  libgflags-dev \
  zlib1g-dev \
  liblz4-dev \
  libzstd-dev \
  libnuma-dev \
  libtbb-dev \
  libgoogle-glog-dev \
  cmake \
  liburing-dev \
  libsnappy-dev \
  libbz2-dev \
  pkg-config
```

Linux 容器内的实际构建命令形态如下：

```bash
cmake -S /src/third_party/rocksdb -B /tmp/rocksdb-build \
  -DCMAKE_BUILD_TYPE=Release \
  -DCMAKE_POSITION_INDEPENDENT_CODE=ON \
  -DPORTABLE=1 \
  -DCMAKE_C_FLAGS="$EXTRA_C_FLAGS" \
  -DCMAKE_CXX_FLAGS="$EXTRA_C_FLAGS" \
  -DWITH_TESTS=OFF \
  -DWITH_GFLAGS=ON \
  -DWITH_BENCHMARK_TOOLS=OFF \
  -DWITH_TOOLS=OFF \
  -DWITH_MD_LIBRARY=OFF \
  -DWITH_RUNTIME_DEBUG=OFF \
  -DROCKSDB_BUILD_SHARED=OFF \
  -DWITH_SNAPPY=ON \
  -DWITH_LZ4=ON \
  -DWITH_ZLIB=ON \
  -DWITH_LIBURING=ON \
  -DWITH_JEMALLOC=ON \
  -DWITH_NUMA=ON \
  -DWITH_TBB=ON \
  -DUSE_RTTI=ON \
  -DWITH_TRACE_TOOLS=OFF \
  -DWITH_CORE_TOOLS=OFF \
  -DUSE_FOLLY=OFF \
  -DUSE_COROUTINES=ON \
  -DWITH_ZSTD=ON \
  -DWITH_BZ2=ON
```

额外说明：

- `FORCE_AVX` 和 `FORCE_SSE42` 对 `linux/amd64` 会转成 `-mavx`、`-msse4.2`
- 对 `linux/arm64` 会自动关闭，避免无效或错误的 x86 指令集参数

## 八、Docker / OrbStack 缓存复用

你要求 Docker 编译尽量放入 cache，后续可以复用。当前实现已经加上了主机缓存目录挂载：

- `apt` 缓存：`$DOCKER_CACHE_DIR/apt-$arch`
- RocksDB 构建目录缓存：`$DOCKER_CACHE_DIR/build-$arch`

默认缓存路径：

```text
.cache/rocksdb-docker/
```

例如：

```bash
./scripts/build-rocksdb-prebuilt.sh \
  --docker-cache-dir "$PWD/.cache/rocksdb-docker" \
  linux_amd64 linux_arm64
```

## 九、关于 INSTALL.md 和 USE_COROUTINES 的当前结论

基于当前仓库 vendored 的 `third_party/rocksdb/INSTALL.md` 与源码内容，可以确认：

1. `PORTABLE=1` 是 upstream 官方支持的兼容构建方式，适合发布可移植产物。
2. `WITH_SNAPPY`、`WITH_LZ4`、`WITH_ZLIB`、`WITH_ZSTD`、`WITH_BZ2`、`WITH_JEMALLOC`、`WITH_NUMA`、`WITH_TBB`、`WITH_LIBURING`、`USE_RTTI` 都在当前 RocksDB 的构建体系里有明确入口。
3. `USE_COROUTINES` 在 RocksDB 当前源码里确实被使用，但它不是一个“只打开一个宏就完事”的选项，会牵扯到更复杂的编译条件。
4. `FORCE_AVX`、`FORCE_SSE42` 并不是当前 upstream CMake 的正式标准开关，更适合作为脚本层对 x86 编译 flag 的增强控制，而不是直接作为通用 CMake 变量硬传给所有平台。

## 十、维护者推荐完整流程

### 10.1 重新生成全部预编译产物

```bash
git submodule update --init --recursive
git -C third_party/rocksdb describe --tags --always

./scripts/build-rocksdb-prebuilt.sh \
  --linux-builder-image harbor.clever8790.top/yanjie/debian:trixie-slim
```

### 10.2 只重建 Linux 产物

```bash
./scripts/build-rocksdb-prebuilt.sh \
  --linux-builder-image harbor.clever8790.top/yanjie/debian:trixie-slim \
  linux_amd64 linux_arm64
```

### 10.3 只重建 macOS arm64 产物

```bash
./scripts/build-rocksdb-prebuilt.sh darwin_arm64
```

## 十一、验证步骤

### 11.1 验证原生产物目录

```bash
ls native
```

期望：

```text
darwin_arm64
linux_amd64
linux_arm64
```

### 11.2 验证关键文件存在

```bash
test -f native/darwin_arm64/librocksdb.a
test -f native/linux_amd64/librocksdb.a
test -f native/linux_arm64/librocksdb.a

test -f native/darwin_arm64/include/rocksdb/c.h
test -f native/linux_amd64/include/rocksdb/c.h
test -f native/linux_arm64/include/rocksdb/c.h
```

### 11.3 验证平台选择逻辑

```bash
go test ./internal/cgo -run 'TestNativeArtifactDir|TestUnsupportedPlatform|TestRequiredNativeArtifactsExist' -v
```

### 11.4 验证当前 macOS arm64 CGO 链路

```bash
CGO_ENABLED=1 go test ./rocksdb -run 'Test(Create|Open|Close|Default|Validate)' -v
CGO_ENABLED=1 go test ./rocksdb -run 'TestCreateCloseThenOpen' -v
CGO_ENABLED=1 go test ./...
CGO_ENABLED=1 go build ./...
```

## 十二、发布前检查清单

发布新版本前，建议按下面顺序执行：

```bash
git submodule update --init --recursive
git -C third_party/rocksdb describe --tags --always

./scripts/build-rocksdb-prebuilt.sh \
  --linux-builder-image harbor.clever8790.top/yanjie/debian:trixie-slim

go test ./internal/cgo -run 'TestNativeArtifactDir|TestUnsupportedPlatform|TestRequiredNativeArtifactsExist' -v
CGO_ENABLED=1 go test ./...
CGO_ENABLED=1 go build ./...
```

最后再检查这些内容是否需要一并提交：

- `native/` 下的预编译产物
- `scripts/` 下的构建脚本
- `internal/cgo/` 下的平台选择逻辑
- `README.md`

## 十三、当前限制

1. 只支持 `darwin/arm64`、`linux/amd64`、`linux/arm64`。
2. 不支持 `darwin/amd64`。
3. 不支持 Windows。
4. Linux 产物构建依赖本地已有的 Debian builder 镜像。
5. `USE_COROUTINES` 已在 Linux 脚本中暴露并默认开启，但它本身属于更高要求的构建选项，后续如果碰到上游依赖链变化，需要以 RocksDB 当前版本源码为准重新调整。
