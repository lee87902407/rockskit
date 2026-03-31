//go:build linux && amd64

package cgo

/*
#cgo CFLAGS: -I${SRCDIR}/../../native/linux_amd64/include
#cgo LDFLAGS: ${SRCDIR}/../../native/linux_amd64/librocksdb.a -lstdc++ -lm -ldl -lpthread
*/
import "C"
