//go:build linux && arm64

package cgo

/*
#cgo CFLAGS: -I${SRCDIR}/../../native/linux_arm64/include
#cgo LDFLAGS: -L${SRCDIR}/../../native/linux_arm64/deps ${SRCDIR}/../../native/linux_arm64/librocksdb.a -lsnappy -llz4 -lz -lzstd -lbz2 -ljemalloc -lnuma -ltbb -luring -lgflags -lglog -lstdc++ -lm -ldl -lpthread
*/
import "C"
