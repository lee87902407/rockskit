//go:build darwin && arm64

package cgo

/*
#cgo CFLAGS: -I${SRCDIR}/../../native/darwin_arm64/include
#cgo LDFLAGS: ${SRCDIR}/../../native/darwin_arm64/librocksdb.a -lc++ -lm
*/
import "C"
