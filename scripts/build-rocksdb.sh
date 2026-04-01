#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ROCKSDB_DIR="$ROOT_DIR/third_party/rocksdb"

PORTABLE="${PORTABLE:-1}"
DISABLE_WARNING_AS_ERROR="${DISABLE_WARNING_AS_ERROR:-1}"
USE_COROUTINES="${USE_COROUTINES:-1}"
USE_RTTI="${USE_RTTI:-1}"

WITH_GFLAGS="${WITH_GFLAGS:-ON}"
WITH_SNAPPY="${WITH_SNAPPY:-ON}"
WITH_LZ4="${WITH_LZ4:-ON}"
WITH_ZLIB="${WITH_ZLIB:-ON}"
WITH_JEMALLOC="${WITH_JEMALLOC:-ON}"
WITH_ZSTD="${WITH_ZSTD:-ON}"
WITH_BZ2="${WITH_BZ2:-ON}"
WITH_NUMA="${WITH_NUMA:-ON}"
WITH_TBB="${WITH_TBB:-ON}"
WITH_LIBURING="${WITH_LIBURING:-ON}"

FORCE_AVX="${FORCE_AVX:-OFF}"
FORCE_SSE42="${FORCE_SSE42:-OFF}"

EXTRA_CFLAGS="${EXTRA_CFLAGS:-}"
EXTRA_CXXFLAGS="${EXTRA_CXXFLAGS:-}"

if [[ ! -f "$ROCKSDB_DIR/include/rocksdb/c.h" ]]; then
  echo "rocksdb submodule is missing: $ROCKSDB_DIR" >&2
  exit 1
fi

make_cmd="make"
if command -v gmake >/dev/null 2>&1; then
  make_cmd="gmake"
fi

export PORTABLE
export DISABLE_WARNING_AS_ERROR
export USE_COROUTINES
export USE_RTTI

if [[ "$WITH_GFLAGS" != "ON" ]]; then
  export ROCKSDB_DISABLE_GFLAGS=1
fi
if [[ "$WITH_SNAPPY" != "ON" ]]; then
  export ROCKSDB_DISABLE_SNAPPY=1
fi
if [[ "$WITH_ZLIB" != "ON" ]]; then
  export ROCKSDB_DISABLE_ZLIB=1
fi
if [[ "$WITH_BZ2" != "ON" ]]; then
  export ROCKSDB_DISABLE_BZIP=1
fi
if [[ "$WITH_LZ4" != "ON" ]]; then
  export ROCKSDB_DISABLE_LZ4=1
fi
if [[ "$WITH_ZSTD" != "ON" ]]; then
  export ROCKSDB_DISABLE_ZSTD=1
fi
if [[ "$WITH_JEMALLOC" != "ON" ]]; then
  export ROCKSDB_DISABLE_JEMALLOC=1
fi
if [[ "$FORCE_AVX" == "ON" ]]; then
  EXTRA_CFLAGS="$EXTRA_CFLAGS -mavx"
  EXTRA_CXXFLAGS="$EXTRA_CXXFLAGS -mavx"
fi
if [[ "$FORCE_SSE42" == "ON" ]]; then
  EXTRA_CFLAGS="$EXTRA_CFLAGS -msse4.2"
  EXTRA_CXXFLAGS="$EXTRA_CXXFLAGS -msse4.2"
fi

if [[ -n "$EXTRA_CFLAGS" ]]; then
  export EXTRA_CFLAGS
fi
if [[ -n "$EXTRA_CXXFLAGS" ]]; then
  export EXTRA_CXXFLAGS
fi

"$make_cmd" -C "$ROCKSDB_DIR" static_lib

if [[ ! -f "$ROCKSDB_DIR/librocksdb.a" ]]; then
  echo "rocksdb static library was not produced" >&2
  exit 1
fi

echo "$ROCKSDB_DIR/librocksdb.a"
