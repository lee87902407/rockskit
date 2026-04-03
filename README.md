# rockskit

`rockskit` is a Go wrapper for RocksDB with bundled prebuilt native libraries.

- [简体中文文档](./README.zh-CN.md)
- [English Documentation](./README.en.md)

## Quick Notes

- Supported platforms: `darwin/arm64`, `linux/amd64`, `linux/arm64`
- Unsupported platforms: Windows, `darwin/amd64`
- The repository already includes the native headers and static libraries needed by the CGO layer

## Quick Start

```bash
go get github.com/lee87902407/rockskit@v0.0.1
```

```go
package main

import (
    "log"

    "github.com/lee87902407/rockskit/rocksdb"
)

func main() {
    cfg := rocksdb.DefaultConfig()
    db, err := rocksdb.Create("./example-db", cfg)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
}
```

For full usage instructions, build details, and maintainer notes, use the language-specific documents above.
