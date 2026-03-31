package cgo

/*
#include <stdlib.h>
#include <rocksdb/c.h>
*/
import "C"

import (
	"fmt"
	"unsafe"
)

func errFromC(errptr *C.char, action string) error {
	if errptr == nil {
		return nil
	}
	defer C.rocksdb_free(unsafe.Pointer(errptr))
	return fmt.Errorf("%s: %s", action, C.GoString(errptr))
}
