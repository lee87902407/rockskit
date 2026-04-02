#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ROCKSDB_DIR="$ROOT_DIR/third_party/rocksdb"
NATIVE_DIR="$ROOT_DIR/native"

parse_args() {
  TARGETS=()
  while [[ $# -gt 0 ]]; do
    case "$1" in
      --linux-builder-image)
        LINUX_BUILDER_IMAGE="$2"
        shift 2
        ;;
      --replace-apt-source)
        REPLACE_APT_SOURCE="$2"
        shift 2
        ;;
      --docker-cache-dir)
        DOCKER_CACHE_DIR="$2"
        shift 2
        ;;
      *)
        TARGETS+=("$1")
        shift
        ;;
    esac
  done
}

LINUX_BUILDER_IMAGE="${LINUX_BUILDER_IMAGE:-harbor.clever8790.top/yanjie/debian:trixie-slim}"
REPLACE_APT_SOURCE="${REPLACE_APT_SOURCE:-1}"
DOCKER_CACHE_DIR="${DOCKER_CACHE_DIR:-$ROOT_DIR/.cache/rocksdb-docker}"

parse_args "$@"

if [[ ${#TARGETS[@]} -eq 0 ]]; then
  TARGETS=(darwin_arm64 linux_amd64 linux_arm64)
fi

PORTABLE="${PORTABLE:-1}"
USE_COROUTINES="${USE_COROUTINES:-ON}"
USE_RTTI="${USE_RTTI:-ON}"
WITH_SNAPPY="${WITH_SNAPPY:-ON}"
WITH_LZ4="${WITH_LZ4:-ON}"
WITH_ZLIB="${WITH_ZLIB:-ON}"
WITH_JEMALLOC="${WITH_JEMALLOC:-ON}"
WITH_ZSTD="${WITH_ZSTD:-ON}"
WITH_BZ2="${WITH_BZ2:-ON}"
WITH_NUMA="${WITH_NUMA:-ON}"
WITH_TBB="${WITH_TBB:-ON}"
WITH_LIBURING="${WITH_LIBURING:-ON}"
WITH_GFLAGS="${WITH_GFLAGS:-ON}"
FORCE_AVX="${FORCE_AVX:-ON}"
FORCE_SSE42="${FORCE_SSE42:-ON}"

WITH_JEMALLOC_DARWIN="OFF"
WITH_JEMALLOC_LINUX="ON"

if [[ "$WITH_JEMALLOC" != "ON" && "$WITH_JEMALLOC" != "OFF" ]]; then
  echo "WITH_JEMALLOC must be ON or OFF" >&2
  exit 1
fi

if [[ "$WITH_JEMALLOC" == "OFF" ]]; then
  WITH_JEMALLOC_LINUX="OFF"
fi

if [[ ! -f "$ROCKSDB_DIR/include/rocksdb/c.h" ]]; then
  echo "rocksdb submodule is missing: $ROCKSDB_DIR" >&2
  exit 1
fi

copy_headers() {
  local target_dir="$1"
  mkdir -p "$target_dir/include"
  rm -rf "$target_dir/include/rocksdb"
  cp -R "$ROCKSDB_DIR/include/rocksdb" "$target_dir/include/rocksdb"
}

prepare_target_dir() {
  local target="$1"
  mkdir -p "$NATIVE_DIR/$target"
  rm -f "$NATIVE_DIR/$target/librocksdb.a"
  rm -rf "$NATIVE_DIR/$target/deps"
  copy_headers "$NATIVE_DIR/$target"
}

build_darwin() {
  local arch="$1"
  local target="darwin_arm64"
  local build_dir="$ROCKSDB_DIR/build_prebuilt_${target}"

  if ! command -v brew >/dev/null 2>&1; then
    echo "brew is required to install macOS build dependencies" >&2
    exit 1
  fi

  brew install snappy lz4 zlib zstd bzip2 gflags || true

  prepare_target_dir "$target"
  rm -rf "$build_dir"
  mkdir -p "$build_dir"

  cmake -S "$ROCKSDB_DIR" -B "$build_dir" \
    -DCMAKE_BUILD_TYPE=Release \
    -DCMAKE_POSITION_INDEPENDENT_CODE=ON \
    -DCMAKE_OSX_ARCHITECTURES="$arch" \
    -DPORTABLE=1 \
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
    -DWITH_LIBURING=OFF \
    -DWITH_JEMALLOC="$WITH_JEMALLOC_DARWIN" \
    -DWITH_NUMA=OFF \
    -DWITH_TBB=OFF \
    -DUSE_RTTI="$USE_RTTI" \
    -DWITH_TRACE_TOOLS=OFF \
    -DWITH_CORE_TOOLS=OFF \
    -DUSE_FOLLY=OFF \
    -DUSE_COROUTINES=OFF \
    -DWITH_ZSTD="$WITH_ZSTD" \
    -DWITH_BZ2="$WITH_BZ2"

  cmake --build "$build_dir" --target rocksdb -j"$(sysctl -n hw.ncpu)"
  cp "$build_dir/librocksdb.a" "$NATIVE_DIR/$target/librocksdb.a"
}

build_linux() {
  local arch="$1"
  local target="linux_${arch}"
  local platform
  local image_ref
  case "$arch" in
    amd64) platform="linux/amd64" ;;
    arm64) platform="linux/arm64" ;;
    *) echo "unsupported linux arch: $arch" >&2; exit 1 ;;
  esac

  prepare_target_dir "$target"

  image_ref="$(resolve_linux_builder_image "$arch")"
  mkdir -p "$DOCKER_CACHE_DIR/apt-$arch" "$DOCKER_CACHE_DIR/build-$arch"

  docker run --rm --platform "$platform" \
    -v "$ROOT_DIR:/src" \
    -v "$NATIVE_DIR/$target:/out" \
    -v "$DOCKER_CACHE_DIR/apt-$arch:/var/cache/apt" \
    -v "$DOCKER_CACHE_DIR/build-$arch:/tmp/rocksdb-build" \
    -e ROOT_DIR=/src \
    -e OUT_DIR=/out \
    -e ROCKSDB_DIR=/src/third_party/rocksdb \
    -e TARGET_ARCH="$arch" \
    -e REPLACE_APT_SOURCE="$REPLACE_APT_SOURCE" \
    -e PORTABLE="$PORTABLE" \
    -e USE_COROUTINES="$USE_COROUTINES" \
    -e USE_RTTI="$USE_RTTI" \
    -e WITH_SNAPPY="$WITH_SNAPPY" \
    -e WITH_TBB="$WITH_TBB" \
    -e WITH_NUMA="$WITH_NUMA" \
    -e WITH_LZ4="$WITH_LZ4" \
    -e WITH_ZLIB="$WITH_ZLIB" \
    -e WITH_LIBURING="$WITH_LIBURING" \
    -e WITH_JEMALLOC="$WITH_JEMALLOC_LINUX" \
    -e WITH_ZSTD="$WITH_ZSTD" \
    -e WITH_BZ2="$WITH_BZ2" \
    -e WITH_GFLAGS="$WITH_GFLAGS" \
    -e FORCE_AVX="$FORCE_AVX" \
    -e FORCE_SSE42="$FORCE_SSE42" \
    "$image_ref" bash /src/scripts/build-rocksdb-linux.sh
}

resolve_linux_builder_image() {
  local arch="$1"
  local image_id
  if docker image inspect "$LINUX_BUILDER_IMAGE" >/dev/null 2>&1; then
    if [[ "$(docker image inspect "$LINUX_BUILDER_IMAGE" --format '{{.Architecture}}' 2>/dev/null)" == "$arch" ]]; then
      printf '%s\n' "$LINUX_BUILDER_IMAGE"
      return 0
    fi
  fi
  while IFS= read -r image_id; do
    if [[ -n "$image_id" ]] && [[ "$(docker image inspect "$image_id" --format '{{.Architecture}}' 2>/dev/null)" == "$arch" ]]; then
      printf '%s\n' "$image_id"
      return 0
    fi
  done < <(docker images --format '{{.Repository}}:{{.Tag}} {{.ID}}' | awk '$1 ~ /^harbor.clever8790.top\/yanjie\/debian:/ {print $2}')

  echo "missing local Linux builder image for $arch from $LINUX_BUILDER_IMAGE" >&2
  exit 1
}

for target in "${TARGETS[@]}"; do
  case "$target" in
    darwin_arm64) build_darwin arm64 ;;
    linux_amd64) build_linux amd64 ;;
    linux_arm64) build_linux arm64 ;;
    *)
      echo "unsupported target: $target" >&2
      exit 1
      ;;
  esac
done

echo "prebuilt artifacts generated under $NATIVE_DIR"
