#ifndef ROCKSKIT_H
#define ROCKSKIT_H

#include <stddef.h>

typedef struct rocksdb_t rocksdb_t;
typedef struct rocksdb_readoptions_t rocksdb_readoptions_t;
typedef struct rocksdb_pinnableslice_t rocksdb_pinnableslice_t;
typedef struct rocksdb_column_family_handle_t rocksdb_column_family_handle_t;
typedef struct rocksdb_writeoptions_t rocksdb_writeoptions_t;

/* PinnedSlice 读取：合并 get_pinned + value + 错误判断为单次 C 调用 */

rocksdb_pinnableslice_t* rockskit_get_pinned(
    rocksdb_t* db,
    const rocksdb_readoptions_t* options,
    const char* key,
    size_t keylen,
    const char** data,
    size_t* datalen,
    char** errptr);

rocksdb_pinnableslice_t* rockskit_get_pinned_cf(
    rocksdb_t* db,
    rocksdb_column_family_handle_t* column_family,
    const rocksdb_readoptions_t* options,
    const char* key,
    size_t keylen,
    const char** data,
    size_t* datalen,
    char** errptr);

/* 批量写入：将 writebatch 创建 + 批量 put/delete + write + destroy 合并为单次 C 调用 */

void rockskit_put_list(
    rocksdb_t* db,
    const rocksdb_writeoptions_t* options,
    size_t num_pairs,
    const char* const* keys_list,
    const size_t* keys_list_sizes,
    const char* const* values_list,
    const size_t* values_list_sizes,
    char** errptr);

void rockskit_put_list_cf(
    rocksdb_t* db,
    const rocksdb_writeoptions_t* options,
    rocksdb_column_family_handle_t* column_family,
    size_t num_pairs,
    const char* const* keys_list,
    const size_t* keys_list_sizes,
    const char* const* values_list,
    const size_t* values_list_sizes,
    char** errptr);

void rockskit_delete_list(
    rocksdb_t* db,
    const rocksdb_writeoptions_t* options,
    size_t num_keys,
    const char* const* keys_list,
    const size_t* keys_list_sizes,
    char** errptr);

void rockskit_delete_list_cf(
    rocksdb_t* db,
    const rocksdb_writeoptions_t* options,
    rocksdb_column_family_handle_t* column_family,
    size_t num_keys,
    const char* const* keys_list,
    const size_t* keys_list_sizes,
    char** errptr);

#endif
