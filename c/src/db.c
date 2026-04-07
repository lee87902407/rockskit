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

rocksdb_pinnableslice_t* rockskit_get_pinned_cf(
    rocksdb_t* db,
    rocksdb_column_family_handle_t* column_family,
    const rocksdb_readoptions_t* options,
    const char* key,
    size_t keylen,
    const char** data,
    size_t* datalen,
    char** errptr) {
  rocksdb_pinnableslice_t* slice = rocksdb_get_pinned_cf(db, options, column_family, key, keylen, errptr);
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

void rockskit_put_list(
    rocksdb_t* db,
    const rocksdb_writeoptions_t* options,
    size_t num_pairs,
    const char* const* keys_list,
    const size_t* keys_list_sizes,
    const char* const* values_list,
    const size_t* values_list_sizes,
    char** errptr) {
  rocksdb_writebatch_t* batch = rocksdb_writebatch_create();
  size_t i;
  for (i = 0; i < num_pairs; i++) {
    rocksdb_writebatch_put(batch, keys_list[i], keys_list_sizes[i],
                           values_list[i], values_list_sizes[i]);
  }
  rocksdb_write(db, options, batch, errptr);
  rocksdb_writebatch_destroy(batch);
}

void rockskit_put_list_cf(
    rocksdb_t* db,
    const rocksdb_writeoptions_t* options,
    rocksdb_column_family_handle_t* column_family,
    size_t num_pairs,
    const char* const* keys_list,
    const size_t* keys_list_sizes,
    const char* const* values_list,
    const size_t* values_list_sizes,
    char** errptr) {
  rocksdb_writebatch_t* batch = rocksdb_writebatch_create();
  size_t i;
  for (i = 0; i < num_pairs; i++) {
    rocksdb_writebatch_put_cf(batch, column_family,
                              keys_list[i], keys_list_sizes[i],
                              values_list[i], values_list_sizes[i]);
  }
  rocksdb_write(db, options, batch, errptr);
  rocksdb_writebatch_destroy(batch);
}

void rockskit_delete_list(
    rocksdb_t* db,
    const rocksdb_writeoptions_t* options,
    size_t num_keys,
    const char* const* keys_list,
    const size_t* keys_list_sizes,
    char** errptr) {
  rocksdb_writebatch_t* batch = rocksdb_writebatch_create();
  size_t i;
  for (i = 0; i < num_keys; i++) {
    rocksdb_writebatch_delete(batch, keys_list[i], keys_list_sizes[i]);
  }
  rocksdb_write(db, options, batch, errptr);
  rocksdb_writebatch_destroy(batch);
}

void rockskit_delete_list_cf(
    rocksdb_t* db,
    const rocksdb_writeoptions_t* options,
    rocksdb_column_family_handle_t* column_family,
    size_t num_keys,
    const char* const* keys_list,
    const size_t* keys_list_sizes,
    char** errptr) {
  rocksdb_writebatch_t* batch = rocksdb_writebatch_create();
  size_t i;
  for (i = 0; i < num_keys; i++) {
    rocksdb_writebatch_delete_cf(batch, column_family,
                                 keys_list[i], keys_list_sizes[i]);
  }
  rocksdb_write(db, options, batch, errptr);
  rocksdb_writebatch_destroy(batch);
}
