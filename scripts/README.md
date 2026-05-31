# Scripts 使用说明

本文档介绍 `scripts/` 目录下所有脚本的用途、使用方式和常见场景。

## 目录概览

- `fix_swagger.py`

## 使用约定

- 建议在仓库根目录执行所有命令（即 `backend/` 目录）。
- 若脚本依赖服务接口（注册、登录、WebSocket），请先启动后端服务。

---

## 1) fix_swagger.py

### 作用

修复 Swagger/OpenAPI 文档中参数对象的 `example` 字段，将其替换为 `x-example`，避免工具链兼容问题。

### 适用文件

- `docs.go`
- `swagger*.go`
- `*swagger*.json|yaml|yml`
- `*openapi*.json|yaml|yml`

### 常用命令

1. 处理单文件

```bash
python scripts/fix_swagger.py docs/docs.go
```

2. 只预览，不修改文件

```bash
python scripts/fix_swagger.py docs/docs.go --dry-run
```

3. 修改前先备份

```bash
python scripts/fix_swagger.py docs/docs.go --backup
```

4. 批量扫描目录并处理

```bash
python scripts/fix_swagger.py . --batch
```

5. 输出到指定目录

```bash
python scripts/fix_swagger.py docs/docs.go --output tmp/swagger_fixed
```

### 参数说明

- `--dry-run`：模拟运行，不落盘
- `--backup`：修改前生成备份文件（`.bak`）
- `--batch`：目录批量模式
- `--output` / `-o`：指定输出目录