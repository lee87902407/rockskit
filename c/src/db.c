#if __has_include("../../native/darwin_arm64/include/rocksdb/c.h")
#include "../../native/darwin_arm64/include/rocksdb/c.h"
#elif __has_include("../../native/linux_amd64/include/rocksdb/c.h")
#include "../../native/linux_amd64/include/rocksdb/c.h"
#elif __has_include("../../native/linux_arm64/include/rocksdb/c.h")
#include "../../native/linux_arm64/include/rocksdb/c.h"
#else
#include <rocksdb/c.h>
#endif

#include "../include/rockskit.h"

rocksdb_pinnableslice_t* rockskit_get_pinned(
    rocksdb_t* db,
    const rocksdb_readoptions_t* options,
    const char* key,
    size_t keylen,
    const char** data,
    size_t* datalen,
    char** errptr) {
  rocksdb_pinnableslice_t* slice = rocksdb_get_pinned(db, options, key, keylen, errptr);
  if (data != NULL) {
    *data = NULL;
  }
  if (datalen != NULL) {
    *datalen = 0;
  }
  if (slice == NULL) {
    return NULL;
  }
  if (data != NULL) {
    *data = rocksdb_pinnableslice_value(slice, datalen);
  } else {
    rocksdb_pinnableslice_value(slice, datalen);
  }
  if ((data != NULL && *data == NULL) || (data == NULL && datalen != NULL && *datalen == 0)) {
    rocksdb_pinnableslice_destroy(slice);
    return NULL;
  }
  return slice;
}
