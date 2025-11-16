# Swagger泛型语法修复指南

## 问题
Swagger (swaggo) 不完全支持 Go 泛型语法 `model.APIResponse[T]`，导致 `swag init` 失败。

## 解决方案

### 方案1：使用具体的响应类型（推荐）

为每个API端点定义具体的响应结构体，避免使用泛型：

```go
// 在handler文件中定义具体响应类型
type UserListResponse struct {
    Success bool          `json:"success"`
    Code    int           `json:"code"`
    Message string        `json:"message"`
    Data    []model.User  `json:"data"`
}

// 在Swagger注释中使用具体类型
// @Success 200 {object} UserListResponse
```

### 方案2：使用通用响应类型

对于不需要详细类型定义的API，使用 `model.SuccessResponse` 或 `model.ErrorResponse`：

```go
// @Success 200 {object} model.SuccessResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
```

### 方案3：批量替换（PowerShell）

运行以下命令批量替换所有Swagger注释中的泛型语法：

```powershell
# 在backend目录下执行

# 替换所有Failure注释
Get-ChildItem -Path "internal\handler" -Filter "*.go" -Recurse |
  Where-Object { $_.Name -notlike "*_test.go" } |
  ForEach-Object {
    $content = Get-Content $_.FullName -Raw
    $content = $content -replace '// @Failure\s+(\d+)\s+\{object\}\s+model\.APIResponse\[any\]', '// @Failure      $1            {object}  model.ErrorResponse'
    $content = $content -replace '// @Failure\s+(\d+)\s+\{object\}\s+model\.APIResponse\[interface\{\}\]', '// @Failure      $1            {object}  model.ErrorResponse'
    Set-Content -Path $_.FullName -Value $content -NoNewline
  }

# 替换通用Success注释
Get-ChildItem -Path "internal\handler" -Filter "*.go" -Recurse |
  Where-Object { $_.Name -notlike "*_test.go" } |
  ForEach-Object {
    $content = Get-Content $_.FullName -Raw
    $content = $content -replace '// @Success\s+(\d+)\s+\{object\}\s+model\.APIResponse\[any\]', '// @Success      $1            {object}  model.SuccessResponse'
    Set-Content -Path $_.FullName -Value $content -NoNewline
  }
```

### 方案4：升级swag版本

某些较新版本的swag可能支持泛型语法：

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag --version
```

## 已完成的修复

已经添加了以下辅助类型到 `internal/model/api.go`：

- `ErrorResponse` - 错误响应类型
- `SuccessResponse` - 通用成功响应类型

## 验证修复

运行以下命令验证Swagger文档生成：

```bash
cd backend
swag init -g cmd/main.go -o docs/swagger -q
```

如果成功，不会有错误输出，并且会在 `docs/swagger/` 目录下生成文档文件。

## 建议

对于新的API端点，建议采用以下最佳实践：

1. **有明确数据结构的Success响应**：定义具体的Response类型
2. **通用Success响应**：使用 `model.SuccessResponse`
3. **所有Failure响应**：统一使用 `model.ErrorResponse`

示例：

```go
// UserDetailResponse 用户详情响应
type UserDetailResponse struct {
    Success bool        `json:"success"`
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    model.User  `json:"data"`
}

// @Summary      获取用户详情
// @Success      200  {object}  UserDetailResponse
// @Failure      400  {object}  model.ErrorResponse
// @Failure      401  {object}  model.ErrorResponse
// @Router       /users/{id} [get]
func getUserHandler(c *gin.Context) {
    // ...
}
```
