# Swagger 文档生成指南

## 支持泛型的 Swagger 配置

本项目已配置支持 Go 1.18+ 泛型的 Swagger 文档生成，使用以下参数：

### 核心参数说明

- `--parseDependency`: 解析依赖包中的注解
- `--parseInternal`: 解析 internal 包中的注解  
- `--parseDepth 10`: 解析深度为 10，支持多级依赖
- `--output`: 指定输出目录
- `--generalInfo`: 指定主入口文件
- `--dir`: 指定扫描目录

### 生成方法

#### 方法 1: 使用 PowerShell 脚本 (推荐 Windows)

```powershell
# 生成支持泛型的 Swagger 文档
.\scripts\generate-swagger.ps1

# 或使用 PowerShell 直接运行
cd backend\cmd
swag init --output "../docs" --generalInfo "main.go" --dir "." --parseDependency --parseInternal --parseDepth 10
```

#### 方法 2: 使用 Bash 脚本 (Linux/macOS)

```bash
# 生成并修复 Swagger 注解
./scripts/fix-swagger-annotations.sh

# 或直接运行 swag 命令
cd cmd
swag init -g main.go --output "../docs" --generalInfo "main.go" --dir "." --parseDependency --parseInternal --parseDepth 10
```

#### 方法 3: 使用 Makefile (Linux/macOS)

```bash
# 生成 Swagger 文档（支持泛型）
make swagger

# 生成 Swagger 文档（包含 vendor）
make swagger-vendor

# 清理 Swagger 文档
make clean-swagger
```

#### 方法 4: 手动运行 swag 命令

```bash
cd backend/cmd
swag init \
  -g main.go \
  --output "../docs" \
  --generalInfo "main.go" \
  --dir "." \
  --parseDependency \
  --parseInternal \
  --parseDepth 10
```

### 输出位置

生成的 Swagger 文档位于：`backend/docs/`

文件结构：
```
backend/docs/
├── docs.go          # Go 代码文件
├── swagger.json     # Swagger JSON 格式
└── swagger.yaml     # Swagger YAML 格式
```

### 访问 Swagger UI

启动服务后，访问：
```
http://localhost:8080/swagger/index.html
```

### 常见问题

#### 1. 泛型类型解析失败

如果遇到泛型类型解析错误，请确保：
- Go 版本 >= 1.18
- swag 版本 >= 1.8.4
- 使用 `--parseDependency` 和 `--parseInternal` 参数

#### 2. 内部包无法解析

使用 `--parseInternal` 参数来解析 internal 包中的注解。

#### 3. 依赖包中的注解无法解析

使用 `--parseDependency` 和 `--parseDepth` 参数来解析依赖包中的注解。

#### 4. vendor 目录支持

如果需要解析 vendor 目录中的注解，添加 `--parseVendor` 参数：

```bash
swag init -g main.go --parseVendor --parseDependency --parseInternal --parseDepth 10
```

### 验证 Swagger 文档

生成文档后，可以验证文档格式是否正确：

```bash
# 检查 Swagger 版本
grep -q "swagger.*2.0" docs/swagger.json && echo "✅ Swagger 版本正确"

# 统计 API 数量
API_COUNT=$(grep -o '"paths"' docs/swagger.json | wc -l)
echo "API 数量: $API_COUNT"
```

### 自动化脚本

本项目提供了自动化脚本：

- `scripts/generate-swagger.ps1` - Windows PowerShell 脚本
- `scripts/fix-swagger-annotations.sh` - Linux/macOS Bash 脚本

这些脚本已配置好所有必要参数，直接运行即可。

### 注意事项

1. **泛型支持**: 必须使用支持泛型的 swag 版本 (>= 1.8.4)
2. **解析深度**: `--parseDepth` 建议设置为 10，以确保足够的解析深度
3. **输出目录**: 所有文档统一输出到 `docs/` 目录
4. **入口文件**: 主入口文件为 `cmd/main.go`

### 版本要求

- Go: >= 1.18 (支持泛型)
- swag: >= 1.8.4
- swaggo/gin-swagger: >= 1.5.0
- swaggo/files: >= 1.0.0

