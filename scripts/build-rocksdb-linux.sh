#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="${ROOT_DIR:-/src}"
OUT_DIR="${OUT_DIR:-/out}"
ROCKSDB_DIR="${ROCKSDB_DIR:-$ROOT_DIR/third_party/rocksdb}"
REPLACE_APT_SOURCE="${REPLACE_APT_SOURCE:-0}"
TARGET_ARCH="${TARGET_ARCH:-$(dpkg --print-architecture 2>/dev/null || uname -m)}"

PORTABLE="${PORTABLE:-1}"
USE_COROUTINES="${USE_COROUTINES:-ON}"
USE_RTTI="${USE_RTTI:-ON}"
WITH_SNAPPY="${WITH_SNAPPY:-ON}"
WITH_TBB="${WITH_TBB:-ON}"
WITH_NUMA="${WITH_NUMA:-ON}"
WITH_LZ4="${WITH_LZ4:-ON}"
WITH_ZLIB="${WITH_ZLIB:-ON}"
WITH_LIBURING="${WITH_LIBURING:-ON}"
WITH_JEMALLOC="${WITH_JEMALLOC:-ON}"
WITH_ZSTD="${WITH_ZSTD:-ON}"
WITH_BZ2="${WITH_BZ2:-ON}"
WITH_GFLAGS="${WITH_GFLAGS:-ON}"
FORCE_AVX="${FORCE_AVX:-ON}"
FORCE_SSE42="${FORCE_SSE42:-ON}"

if [[ "$TARGET_ARCH" != "amd64" && "$TARGET_ARCH" != "x86_64" ]]; then
  FORCE_AVX="OFF"
  FORCE_SSE42="OFF"
fi

EXTRA_C_FLAGS=""
if [[ "$FORCE_AVX" == "ON" ]]; then
  EXTRA_C_FLAGS="$EXTRA_C_FLAGS -mavx"
fi
if [[ "$FORCE_SSE42" == "ON" ]]; then
  EXTRA_C_FLAGS="$EXTRA_C_FLAGS -msse4.2"
fi

if [[ "$REPLACE_APT_SOURCE" == "1" ]]; then
  if [[ -f /etc/apt/sources.list ]]; then
    sed -i 's/deb.debian.org/mirrors.ustc.edu.cn/g' /etc/apt/sources.list
  fi
  if [[ -f /etc/apt/sources.list.d/debian.sources ]]; then
    sed -i 's/deb.debian.org/mirrors.ustc.edu.cn/g' /etc/apt/sources.list.d/debian.sources
  fi
fi

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

mkdir -p "$OUT_DIR"
cmake -S "$ROCKSDB_DIR" -B /tmp/rocksdb-build \
  -DCMAKE_BUILD_TYPE=Release \
  -DCMAKE_POSITION_INDEPENDENT_CODE=ON \
  -DPORTABLE="$PORTABLE" \
  -DCMAKE_C_FLAGS="$EXTRA_C_FLAGS" \
  -DCMAKE_CXX_FLAGS="$EXTRA_C_FLAGS" \
  -DWITH_TESTS=OFF \
  -DWITH_GFLAGS="$WITH_GFLAGS" \
  -DWITH_BENCHMARK_TOOLS=OFF \
  -DWITH_TOOLS=OFF \
  -DWITH_MD_LIBRARY=OFF \
  -DWITH_RUNTIME_DEBUG=OFF \
  -DROCKSDB_BUILD_SHARED=OFF \
  -DWITH_SNAPPY="$WITH_SNAPPY" \
  -DWITH_LZ4="$WITH_LZ4" \
  -DWITH_ZLIB="$WITH_ZLIB" \
  -DWITH_LIBURING="$WITH_LIBURING" \
  -DWITH_JEMALLOC="$WITH_JEMALLOC" \
  -DWITH_NUMA="$WITH_NUMA" \
  -DWITH_TBB="$WITH_TBB" \
  -DUSE_RTTI="$USE_RTTI" \
  -DWITH_TRACE_TOOLS=OFF \
  -DWITH_CORE_TOOLS=OFF \
  -DUSE_FOLLY=OFF \
  -DUSE_COROUTINES="$USE_COROUTINES" \
  -DWITH_ZSTD="$WITH_ZSTD" \
  -DWITH_BZ2="$WITH_BZ2"

cmake --build /tmp/rocksdb-build --target rocksdb -j"$(nproc)"
cp /tmp/rocksdb-build/librocksdb.a "$OUT_DIR/librocksdb.a"
