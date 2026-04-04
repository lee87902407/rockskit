#ifndef ROCKSKIT_H
#define ROCKSKIT_H

#include <stddef.h>

typedef struct rocksdb_t rocksdb_t;
typedef struct rocksdb_readoptions_t rocksdb_readoptions_t;
typedef struct rocksdb_pinnableslice_t rocksdb_pinnableslice_t;

rocksdb_pinnableslice_t* rockskit_get_pinned(
    rocksdb_t* db,
    const rocksdb_readoptions_t* options,
    const char* key,
    size_t keylen,
    const char** data,
    size_t* datalen,
    char** errptr);

#endif
