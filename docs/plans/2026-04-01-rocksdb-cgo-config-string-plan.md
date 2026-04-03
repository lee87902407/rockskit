# RocksDB CGO 配置字符串实现计划

> **给 Claude：** 必须使用 `superpowers:executing-plans` 子技能按任务逐步执行本计划。

**目标：** 重构底层 CGO 层，使 RocksDB C 句柄以仅持有指针的 Go 对象进行封装并统一提供 `Close()` 方法；扩展 `rocksdb.Config` 到用户要求的配置面；并让 `Create`/`Open` 能通过配置字符串结合 `rocksdb_get_options_from_string` 构建 RocksDB 选项。

**架构：** `internal/cgo` 仍然是唯一允许直接接触原始 `rocksdb/c.h` 指针和生命周期管理的位置。对外的 `rocksdb/` 代码只负责组装配置文本、调用底层封装，并返回具备显式 `Close()` 语义的 Go 风格 `DB` 包装对象。

**技术栈：** Go 1.25.3、CGO、仓库内 vendored 的 RocksDB C API、预编译静态库、Go 测试。

### 任务 1：先为扩展配置和基于配置字符串的 create/open 编写失败测试

**文件：**
- 修改：`rocksdb/options_test.go`
- 修改：`rocksdb/db_test.go`

**步骤 1：编写失败测试**

增加测试，要求：
- `Config` 暴露用户要求的字段
- `DefaultConfig()` 初始化关键新增字段
- 配置校验能拦截非法值
- `Create` 能通过配置字符串路径成功创建数据库
- `Open` 能沿用同一条路径，但使用不同的创建标志

**步骤 2：运行测试并确认失败**

执行：

```bash
CGO_ENABLED=1 go test ./rocksdb -run 'Test(Create|Open|Default|Validate)' -v
```

预期：FAIL，因为当时的配置面和底层配置字符串构造尚未实现。

### 任务 2：引入仅持有指针且带 Close 的底层包装对象

**文件：**
- 修改：`internal/cgo/options.go`
- 修改：`internal/cgo/db.go`
- 修改：`internal/cgo/errors.go`
- 新建：`internal/cgo/cache.go`
- 新建：`internal/cgo/block_based_options.go`
- 新建：`internal/cgo/rate_limiter.go`
- 新建：`internal/cgo/config_options.go`

**步骤 1：先写失败测试**

增加最小化测试或包内断言，要求所有指针包装对象都暴露 `Close()`，并且重复关闭时不会 panic。

**步骤 2：运行测试并确认失败**

执行：

```bash
CGO_ENABLED=1 go test ./internal/cgo -v
```

预期：FAIL，因为这些包装类型尚不存在。

**步骤 3：补最小实现**

先实现以下仅持有指针的包装对象：
- options 包装
- block-based options 包装
- cache 包装
- rate limiter 包装
- 如果 `rocksdb_get_options_from_string` 需要，则增加 config-options 包装
- db 包装

所有包装对象都必须提供 `Close()`。

**步骤 4：重新运行测试**

执行：

```bash
CGO_ENABLED=1 go test ./internal/cgo -v
```

预期：PASS，或者只在下一层未实现行为上失败，而不是缺少包装对象本身。

### 任务 3：扩展 Config 并生成 RocksDB 选项字符串

**文件：**
- 修改：`rocksdb/config.go`
- 修改：`rocksdb/options.go`

**步骤 1：实现用户要求的 `Config` 配置面**

加入用户明确要求的字段，并保留 YAML tag。

**步骤 2：补配置归一化与校验**

实现默认值和校验逻辑，确保非法字符串或数值在接触 RocksDB 之前就失败。

**步骤 3：构造 RocksDB 选项文本**

生成满足 RocksDB option-string 语法要求的配置字符串。第一版保持最小实现，只映射当前 `Create`/`Open` 实际会用到的字段。

**步骤 4：重新运行测试**

执行：

```bash
CGO_ENABLED=1 go test ./rocksdb -run 'Test(Default|Validate)' -v
```

预期：PASS。

### 任务 4：围绕 `rocksdb_get_options_from_string` 重做 Create/Open

**文件：**
- 修改：`rocksdb/db.go`

**步骤 1：替换原来的直接 setter 组装路径**

让 `Create(path, cfg)` 和 `Open(path, cfg)` 先组装配置字符串，再调用底层包装，通过 `rocksdb_get_options_from_string` 加载选项。

**步骤 2：保留 create/open 标志位差异**

确保 `Create` 和 `Open` 的唯一区别仍然是 create-if-missing / error-if-exists 等创建标志。

**步骤 3：确保所有临时底层对象都会被关闭**

在构造选项过程中创建的每个底层辅助对象都必须显式关闭。

**步骤 4：重新运行定向测试**

执行：

```bash
CGO_ENABLED=1 go test ./rocksdb -run 'Test(Create|Open)' -v
```

预期：PASS。

### 任务 5：最终验证与手工 QA

**文件：**
- 如有必要，修改：`README.md`

**步骤 1：运行诊断**

对所有改动过的 Go 文件运行 LSP 诊断，要求 0 error。

**步骤 2：运行完整验证**

执行：

```bash
CGO_ENABLED=1 go test ./internal/cgo -v
CGO_ENABLED=1 go test ./rocksdb -v
CGO_ENABLED=1 go test ./...
CGO_ENABLED=1 go build ./...
```

预期：PASS。

**步骤 3：手工 QA**

执行真实的 create/open/close 生命周期测试，确认结果成功。

```bash
CGO_ENABLED=1 go test ./rocksdb -run 'TestCreateCloseThenOpen' -v
```

预期：PASS。
