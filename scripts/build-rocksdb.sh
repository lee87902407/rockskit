#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ROCKSDB_DIR="$ROOT_DIR/third_party/rocksdb"

if [[ ! -f "$ROCKSDB_DIR/include/rocksdb/c.h" ]]; then
  echo "rocksdb submodule is missing: $ROCKSDB_DIR" >&2
  exit 1
fi

make_cmd="make"
if command -v gmake >/dev/null 2>&1; then
  make_cmd="gmake"
fi

export PORTABLE=1
export DISABLE_WARNING_AS_ERROR=1
export ROCKSDB_DISABLE_GFLAGS=1
export ROCKSDB_DISABLE_SNAPPY=1
export ROCKSDB_DISABLE_ZLIB=1
export ROCKSDB_DISABLE_BZIP=1
export ROCKSDB_DISABLE_LZ4=1
export ROCKSDB_DISABLE_ZSTD=1
export ROCKSDB_DISABLE_JEMALLOC=1

"$make_cmd" -C "$ROCKSDB_DIR" static_lib

if [[ ! -f "$ROCKSDB_DIR/librocksdb.a" ]]; then
  echo "rocksdb static library was not produced" >&2
  exit 1
fi

echo "$ROCKSDB_DIR/librocksdb.a"
