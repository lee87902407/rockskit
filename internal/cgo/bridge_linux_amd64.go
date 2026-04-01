//go:build linux && amd64

package cgo

/*
#cgo CFLAGS: -I${SRCDIR}/../../native/linux_amd64/include
#cgo LDFLAGS: ${SRCDIR}/../../native/linux_amd64/librocksdb.a -lsnappy -llz4 -lz -lzstd -lbz2 -ljemalloc -lnuma -ltbb -luring -lgflags -lglog -lstdc++ -lm -ldl -lpthread
*/
import "C"
