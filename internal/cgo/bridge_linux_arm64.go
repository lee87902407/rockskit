//go:build linux && arm64

package cgo

/*
#cgo CFLAGS: -I${SRCDIR}/../../native/linux_arm64/include
#cgo LDFLAGS: ${SRCDIR}/../../native/linux_arm64/librocksdb.a -lstdc++ -lm -ldl -lpthread
*/
import "C"
