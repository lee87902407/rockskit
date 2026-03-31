#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ROCKSDB_DIR="$ROOT_DIR/third_party/rocksdb"
NATIVE_DIR="$ROOT_DIR/native"

if [[ $# -gt 0 ]]; then
  TARGETS=("$@")
else
  TARGETS=(darwin_arm64 linux_amd64 linux_arm64)
fi

LINUX_BUILDER_IMAGE="harbor.clever8790.top/yanjie/debian:trixie-slim"

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
  copy_headers "$NATIVE_DIR/$target"
}

build_darwin() {
  local arch="$1"
  local target="darwin_arm64"
  local build_dir="$ROCKSDB_DIR/build_prebuilt_${target}"

  prepare_target_dir "$target"
  rm -rf "$build_dir"
  mkdir -p "$build_dir"

  cmake -S "$ROCKSDB_DIR" -B "$build_dir" \
    -DCMAKE_BUILD_TYPE=Release \
    -DCMAKE_POSITION_INDEPENDENT_CODE=ON \
    -DCMAKE_OSX_ARCHITECTURES="$arch" \
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

  cmake --build "$build_dir" --target rocksdb -j"$(sysctl -n hw.ncpu)"
  cp "$build_dir/librocksdb.a" "$NATIVE_DIR/$target/librocksdb.a"
}

linux_build_script() {
  cat <<'EOF'
set -euo pipefail
if [ -f /etc/apt/sources.list ]; then
  sed -i 's/deb.debian.org/mirrors.ustc.edu.cn/g' /etc/apt/sources.list
fi
if [ -f /etc/apt/sources.list.d/debian.sources ]; then
  sed -i 's/deb.debian.org/mirrors.ustc.edu.cn/g' /etc/apt/sources.list.d/debian.sources
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
EOF
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

  docker run --rm --platform "$platform" \
    -v "$ROOT_DIR:/src" \
    -v "$NATIVE_DIR/$target:/out" \
    "$image_ref" bash -lc "$(linux_build_script)"
}

resolve_linux_builder_image() {
  local arch="$1"
  local image_id
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
