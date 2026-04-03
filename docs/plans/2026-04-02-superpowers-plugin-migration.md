# Superpowers 插件迁移实施计划

> **给 Claude：** 必须使用 `superpowers:executing-plans` 子技能按任务逐步执行本计划。

**目标：** 将当前 OpenCode 中 superpowers 的安装方式从旧的 symlink/clone 方案迁移为官方 plugin 方案，并验证技能发现机制仍然可用。

**架构：** 更新全局 OpenCode 配置 `~/.config/opencode/opencode.json`，让 `plugin` 字段包含官方上游 superpowers git 插件引用；然后移除旧的符号链接插件入口和本地 clone 的 superpowers 目录，使 OpenCode 重启后只保留一条权威安装路径。

**技术栈：** OpenCode 配置文件（`opencode.json`）、superpowers 插件安装流程、本地文件系统清理、OpenCode 技能发现。

### 任务 1：确认当前生效配置和旧安装残留

**文件：**
- 修改：`docs/plans/2026-04-02-superpowers-plugin-migration.md`
- 检查：`/Users/yanjie/.config/opencode/opencode.json`
- 检查：`/Users/yanjie/.config/opencode/plugin/superpowers.js`
- 检查：`/Users/yanjie/.config/opencode/superpowers`

**步骤 1：验证当前全局配置中的 plugin 数组**

执行：

```bash
python3 - <<'PY'
import json, pathlib
p = pathlib.Path('/Users/yanjie/.config/opencode/opencode.json')
obj = json.loads(p.read_text())
print(obj.get('plugin', []))
PY
```

预期：输出 plugin 列表，且其中尚不包含官方的 `superpowers@git+https://github.com/obra/superpowers.git` 条目。

**步骤 2：验证旧的符号链接安装存在**

执行：

```bash
ls -la /Users/yanjie/.config/opencode/plugin
```

预期：能看到 `superpowers.js`，并且它指向 `~/.config/opencode/superpowers/...` 下的 clone 仓库。

### 任务 2：把官方插件安装项加入当前配置

**文件：**
- 修改：`/Users/yanjie/.config/opencode/opencode.json`

**步骤 1：写入最小配置改动**

向已有 `plugin` 数组增加这一项：

```json
"superpowers@git+https://github.com/obra/superpowers.git"
```

**步骤 2：重新读取配置，确认该条目只出现一次**

执行：

```bash
python3 - <<'PY'
import json, pathlib
p = pathlib.Path('/Users/yanjie/.config/opencode/opencode.json')
obj = json.loads(p.read_text())
plugins = obj.get('plugin', [])
print(plugins)
print(plugins.count('superpowers@git+https://github.com/obra/superpowers.git'))
PY
```

预期：plugin 列表包含新条目，且计数为 `1`。

### 任务 3：删除旧的 symlink/clone 安装

**文件：**
- 删除：`/Users/yanjie/.config/opencode/plugin/superpowers.js`
- 删除：`/Users/yanjie/.config/opencode/superpowers`

**步骤 1：删除旧的符号链接插件文件**

执行：

```bash
rm -f /Users/yanjie/.config/opencode/plugin/superpowers.js
```

预期：再次执行 `ls -la /Users/yanjie/.config/opencode/plugin` 时，不再看到 `superpowers.js`。

**步骤 2：删除旧的 superpowers clone 目录**

执行：

```bash
rm -rf /Users/yanjie/.config/opencode/superpowers
```

预期：再次执行 `ls -la /Users/yanjie/.config/opencode` 时，不再看到 `superpowers/` 目录。

### 任务 4：验证迁移后的状态

**文件：**
- 检查：`/Users/yanjie/.config/opencode/opencode.json`

**步骤 1：验证文件系统清理结果**

执行：

```bash
ls -la /Users/yanjie/.config/opencode && ls -la /Users/yanjie/.config/opencode/plugin
```

预期：配置文件仍然存在，旧符号链接消失，旧 clone 目录也已移除。

**步骤 2：验证 OpenCode 重启后能识别插件**

重启后执行：

```bash
opencode run --print-logs "hello" 2>&1 | rg -i superpowers
```

预期：日志中能看到插件正常加载且没有报错。

**步骤 3：验证技能发现**

在 OpenCode 内执行：

```bash
find_skills
```

预期：superpowers 提供的技能仍然能够被发现。
