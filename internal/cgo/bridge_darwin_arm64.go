//go:build darwin && arm64

package cgo

/*
#cgo CFLAGS: -I${SRCDIR}/../../native/darwin_arm64/include
#cgo LDFLAGS: -L/opt/homebrew/opt/snappy/lib -L/opt/homebrew/opt/lz4/lib -L/opt/homebrew/opt/zstd/lib -L/opt/homebrew/opt/zlib/lib -L/opt/homebrew/opt/bzip2/lib ${SRCDIR}/../../native/darwin_arm64/librocksdb.a -lsnappy -llz4 -lzstd -lz -lbz2 -lc++ -lm
*/
import "C"
