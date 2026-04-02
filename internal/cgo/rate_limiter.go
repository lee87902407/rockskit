package cgo

/*
#include <rocksdb/c.h>
*/
import "C"

type RateLimiter struct {
	ptr *C.rocksdb_ratelimiter_t
}

func NewRateLimiter(rateBytesPerSec int64, refillPeriodMicros int64, fairness int32, autoTuned bool) *RateLimiter {
	if autoTuned {
		return &RateLimiter{ptr: C.rocksdb_ratelimiter_create_auto_tuned(C.int64_t(rateBytesPerSec), C.int64_t(refillPeriodMicros), C.int32_t(fairness))}
	}
	return &RateLimiter{ptr: C.rocksdb_ratelimiter_create(C.int64_t(rateBytesPerSec), C.int64_t(refillPeriodMicros), C.int32_t(fairness))}
}

func (r *RateLimiter) Close() {
	if r == nil || r.ptr == nil {
		return
	}
	C.rocksdb_ratelimiter_destroy(r.ptr)
	r.ptr = nil
}
