# RocksDB 预编译产物实现计划

> **给 Claude：** 必须使用 `superpowers:executing-plans` 子技能按任务逐步执行本计划。

**目标：** 让 `rockskit` 以平台区分的预编译 RocksDB 静态库形式发布，使下游用户只需引入 Go module 即可直接使用，而不需要手动编译或安装 RocksDB。

**架构：** 对维护者保留 RocksDB 源码子模块；对仓库本身则在 `native/` 下按 `GOOS/GOARCH` 组织并提交预编译产物。`internal/cgo/` 通过平台专属构建文件自动选择正确的头文件与静态库，让下游构建时自动命中对应产物。

**技术栈：** Go 1.25.3、CGO、RocksDB C API、预编译静态库、shell 构建脚本，仅支持 darwin/linux 的 amd64/arm64。

### 任务 1：定义预编译产物布局并编写平台选择测试

**文件：**
- 新建：`internal/cgo/platform_test.go`
- 新建：`internal/cgo/platform.go`
- 新建：`native/.gitkeep`

**步骤 1：先写失败测试**

为确定性的产物布局增加测试：
- `darwin/amd64` -> `native/darwin_amd64/`
- `darwin/arm64` -> `native/darwin_arm64/`
- `linux/amd64` -> `native/linux_amd64/`
- `linux/arm64` -> `native/linux_arm64/`
- 不支持的平台返回清晰错误

**步骤 2：运行测试并确认失败**

```bash
go test ./internal/cgo -run 'TestNativeArtifactDir|TestUnsupportedPlatform' -v
```

预期：FAIL，因为当时平台选择辅助函数尚未实现。

**步骤 3：补最小实现**

在 `internal/cgo/platform.go` 中实现平台选择辅助函数：
- 将 `runtime.GOOS`/`runtime.GOARCH` 映射到预编译目录名
- 对不支持的平台显式报错

**步骤 4：重新运行测试**

```bash
go test ./internal/cgo -run 'TestNativeArtifactDir|TestUnsupportedPlatform' -v
```

预期：PASS。

### 任务 2：用平台专属预编译链接替换本地源码直接链接

**文件：**
- 修改：`internal/cgo/bridge.go`
- 新建：`internal/cgo/bridge_darwin_amd64.go`
- 新建：`internal/cgo/bridge_darwin_arm64.go`
- 新建：`internal/cgo/bridge_linux_amd64.go`
- 新建：`internal/cgo/bridge_linux_arm64.go`
- 新建：`internal/cgo/unsupported.go`

**步骤 1：先写失败测试**

增加一个小测试或 build 侧验证，要求当前包在编译时不再直接引用 `third_party/rocksdb/librocksdb.a`。

**步骤 2：运行测试/构建并确认失败**

```bash
go test ./internal/cgo -v
```

预期：FAIL，或者仍然能暴露出对子模块路径的硬编码链接。

**步骤 3：补最小实现**

实现按平台拆分的 bridge 文件，每个文件都提供精确的 `#cgo` 指令，直接指向对应 `native/<platform>/librocksdb.a` 和 include 目录。对不支持的平台，使用带 build tag 的 stub 返回编译期或运行期错误。

**步骤 4：重新运行验证**

```bash
go test ./internal/cgo -v
```

预期：在当前机器上 PASS。

### 任务 3：构建并落盘预编译 RocksDB 产物

**文件：**
- 修改：`scripts/build-rocksdb.sh`
- 新建：`scripts/build-rocksdb-prebuilt.sh`
- 新建：`native/darwin_amd64/include/rocksdb/c.h`
- 新建：`native/darwin_amd64/librocksdb.a`
- 新建：`native/darwin_arm64/include/rocksdb/c.h`
- 新建：`native/darwin_arm64/librocksdb.a`
- 新建：`native/linux_amd64/include/rocksdb/c.h`
- 新建：`native/linux_amd64/librocksdb.a`
- 新建：`native/linux_arm64/include/rocksdb/c.h`
- 新建：`native/linux_arm64/librocksdb.a`

**步骤 1：先写失败测试**

增加一个文件存在性测试，用于校验支持平台所需的 native 产物布局。

**步骤 2：运行测试并确认失败**

```bash
go test ./internal/cgo -run 'TestRequiredNativeArtifactsExist' -v
```

预期：FAIL，因为当时 `native/` 目录尚未填充产物。

**步骤 3：补最小实现**

实现一个面向维护者的脚本，要求：
- 接收目标平台元组
- 从子模块构建目标平台的 RocksDB
- 将裁剪后的 `librocksdb.a` 和必要公共头文件复制到 `native/<platform>/`
- 在合适的地方参考原始脚本的 CMake 优化策略

提交四个支持平台对应的产物。

**步骤 4：重新运行测试**

```bash
go test ./internal/cgo -run 'TestRequiredNativeArtifactsExist' -v
```

预期：PASS。

### 任务 4：验证下游消费行为

**文件：**
- 修改：`README.md`
- 新建：`examples/basic/main.go`

**步骤 1：先写失败的端到端验证**

增加集成测试或脚本验证，要求：
- 只使用发布后的 module 布局
- 不调用 `scripts/build-rocksdb.sh`
- 能构建一个最小消费程序，导入模块并完成 DB 的 open/close

**步骤 2：运行并确认当前会失败**

```bash
CGO_ENABLED=1 go build ./examples/basic
```

预期：在 native 产物选择或 example 接线完成前 FAIL。

**步骤 3：补最小实现**

添加示例程序，并更新文档说明：
- 支持的平台范围
- module 已内置 native 产物
- 用户不需要手动安装或编译 RocksDB
- Windows 不受支持

**步骤 4：重新运行验证**

```bash
CGO_ENABLED=1 go build ./examples/basic
CGO_ENABLED=1 go test ./...
```

预期：PASS。

### 任务 5：最终 QA 与发布前验证

**文件：**
- 如仍有注意事项，修改：`README.md`

**步骤 1：运行诊断**

对所有改动过的 Go 文件运行 LSP 诊断，要求 0 error。

**步骤 2：运行完整校验**

```bash
CGO_ENABLED=1 go test ./...
CGO_ENABLED=1 go build ./...
```

预期：PASS。

**步骤 3：手工 QA**

执行真实的 create/open/close 场景，并在不调用任何本地 RocksDB 构建脚本的前提下构建示例程序：

```bash
CGO_ENABLED=1 go test ./rocksdb -run 'TestCreateCloseThenOpen' -v
CGO_ENABLED=1 go build ./examples/basic
```

预期：PASS。

**步骤 4：记录发布注意事项**

在 `README.md` 中记录仓库体积影响和支持平台矩阵。
