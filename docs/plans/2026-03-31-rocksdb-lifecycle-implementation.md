# RocksDB 生命周期实现计划

> **给 Claude：** 必须使用 `superpowers:executing-plans` 子技能按任务逐步执行本计划。

**目标：** 将 `facebook/rocksdb` 作为子模块接入，并把占位版 `rocksdb.Create`/`Open`/`Close` 实现替换成真实的 RocksDB 生命周期实现，同时基于迁移后的配置构建选项。

**架构：** 对外生命周期 API 保持在 `rocksdb/`；所有原生调用和所有权处理都收敛到 `internal/cgo/`；ANSI C shim 保留在 `c/`。第一阶段只迁移生命周期创建所需的最小配置面，不一次性照搬全部参考实现。

**技术栈：** Go 1.25.3、CGO、RocksDB C API、本地 `facebook/rocksdb` git 子模块、用于编译原生库的 shell 脚本。

### 任务 1：添加 RocksDB 子模块和原生构建入口

**文件：**
- 修改：`.gitmodules`
- 新建：`third_party/rocksdb/`（git submodule）
- 新建：`scripts/build-rocksdb.sh`

**步骤 1：先写失败验证命令**

```bash
test -f third_party/rocksdb/include/rocksdb/c.h
```

预期：FAIL，因为当时子模块尚不存在。

**步骤 2：添加子模块**

```bash
GIT_MASTER=1 git submodule add ssh://org-69631@github.com:facebook/rocksdb.git third_party/rocksdb
```

**步骤 3：添加构建脚本**

创建 `scripts/build-rocksdb.sh`，要求：
- 校验 `third_party/rocksdb` 存在
- 在子模块中完成 RocksDB 静态库构建
- 产出一个稳定路径，供 CGO 侧链接使用

**步骤 4：验证构建脚本可用**

```bash
./scripts/build-rocksdb.sh
```

预期：退出码为 0，并且生成可被 CGO 链接的 RocksDB 库文件。

### 任务 2：迁移最小生命周期配置面

**文件：**
- 新建：`rocksdb/options.go`
- 新建：`rocksdb/config.go`
- 测试：`rocksdb/options_test.go`

**步骤 1：先写失败测试**

增加最小迁移配置面的测试：
- env 配置：cache/block size
- db 配置：write buffer size、max bytes for level base、target file size base、level0 compaction trigger、compression type
- 空值或非法值的默认值/校验行为

**步骤 2：运行测试并确认失败**

```bash
go test ./rocksdb -run 'Test(BuildOptions|DefaultConfig|ConfigValidation)' -v
```

预期：FAIL，因为配置类型和选项构建器当时还不存在。

**步骤 3：补最小实现**

按参考实现思路引入 Go 侧配置结构：
- `EnvConfig`：`LRUSize`、`BlockSize`
- `Config`：只保留生命周期必需字段
- 默认值以本地开发可用为准

**步骤 4：重新运行测试**

```bash
go test ./rocksdb -run 'Test(BuildOptions|DefaultConfig|ConfigValidation)' -v
```

预期：PASS。

### 任务 3：增加生命周期所需的 C shim 和选项接线

**文件：**
- 新建：`c/include/rockskit.h`
- 新建：`c/src/db.c`
- 新建：`c/src/options.c`
- 新建：`c/src/error.c`

**步骤 1：先写面向 CGO 的失败测试**

增加包级生命周期测试，要求原生 create/open/close 在 bridge 存在后能够工作。

**步骤 2：运行测试并确认失败**

```bash
CGO_ENABLED=1 go test ./rocksdb -run 'TestNative(Create|Open|Close)' -v
```

预期：FAIL，因为原生 bridge 当时尚不存在。

**步骤 3：实现最小 C 层功能**

只实现：
- 生命周期所需的 option 分配/释放辅助函数
- 基于 RocksDB C API 的 db create/open 包装
- db close 包装
- 集中的 errptr 辅助处理

**步骤 4：重新运行测试**

```bash
CGO_ENABLED=1 go test ./rocksdb -run 'TestNative(Create|Open|Close)' -v
```

预期：bridge 能编译，或者只在下一层未实现逻辑上失败，而不是因为缺少 C 包装函数失败。

### 任务 4：增加 CGO bridge 并替换占位生命周期实现

**文件：**
- 新建：`internal/cgo/bridge.go`
- 新建：`internal/cgo/db.go`
- 新建：`internal/cgo/options.go`
- 新建：`internal/cgo/errors.go`
- 修改：`rocksdb/db.go`
- 修改：`rocksdb/db_test.go`

**步骤 1：扩展失败中的生命周期测试**

覆盖以下行为：
- `Create` 能基于配置构建选项并创建新 RocksDB 数据库
- 目标已存在时 `Create` 应失败
- `Open` 只能打开已存在且已初始化的数据库
- `Close` 释放原生句柄并保持幂等

**步骤 2：运行定向测试并观察失败**

```bash
CGO_ENABLED=1 go test ./rocksdb -run 'Test(Create|Open|Close)' -v
```

预期：FAIL，因为当时 `rocksdb/db.go` 仍然使用 marker-file 占位实现。

**步骤 3：用最小原生实现替换占位逻辑**

实现：
- 持有原生句柄所有权的 `DB`
- `Create(path string, cfg *Config)` 或语义清晰且能保持对外 API 的等价构造函数
- `Open(...)`
- `Close()`

所有选项组装都必须经过新的配置面和原生 bridge。

**步骤 4：重新运行定向测试**

```bash
CGO_ENABLED=1 go test ./rocksdb -run 'Test(Create|Open|Close)' -v
```

预期：PASS。

### 任务 5：完整验证和手工 QA

**文件：**
- 如有必要，修改：`README.md`

**步骤 1：运行诊断**

对所有改动过的 Go 文件运行 LSP 诊断，确保 0 error。

**步骤 2：运行包级与全项目测试**

```bash
CGO_ENABLED=1 go test ./rocksdb -v
CGO_ENABLED=1 go test ./...
go build ./...
```

预期：PASS。

**步骤 3：手工 QA**

针对临时目录执行一次真实 create/open/close 场景：

```bash
CGO_ENABLED=1 go test ./rocksdb -run 'TestCreateCloseThenOpen' -v
```

预期：PASS，并且体现真实的原生 RocksDB 生命周期行为。

**步骤 4：记录残余风险**

如果 macOS 链接仍需要额外本地依赖，则在 `README.md` 中准确记录所需命令或环境变量。
