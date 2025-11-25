# Go测试文件命名规范

## ✅ 正确命名规范

遵循Go官方规范，使用 `*_test.go` 命名格式。

### 支持的格式

```
{package_name}_test.go
```

**示例**:
- ✅ `jwt_test.go`
- ✅ `user_test.go`
- ✅ `payment_test.go`
- ✅ `order_service_test.go`
- ✅ `auth_handler_test.go`

### 命名规则说明

**必须遵循**:
- 文件名必须以 `_test.go` 结尾
- 前缀通常是包名或测试对象名
- 使用小写字母和数字
- 单词之间用下划线分隔

**错误示例**:
- ❌ `jwt.test.go` (Go工具链不支持)
- ❌ `test_jwt.go` (不符合约定)
- ❌ `JWT_test.go` (不应使用大写字母)

## 📖 为什么使用这个规范

### Go工具链要求

Go测试工具硬编码只识别 `*_test.go` 模式：

```
'Go test' recompiles each package along with any files with names matching
the file pattern "*_test.go".
```

### 优势

1. **标准兼容**: 完全符合Go语言规范
2. **工具支持**: IDE、编辑器、CI/CD完全支持
3. **社区实践**: 所有Go项目通用做法
4. **简单易用**: 无需额外配置

## 📋 实施检查清单

### 文件命名检查

```bash
# 检查是否有不符合规范的文件
find . -name "*.go" -type f | grep -v "_test.go$" | grep -v ".go$"

# 应该只返回非测试文件
```

### 测试运行验证

```bash
# 运行测试
go test ./... -v

# 检查覆盖率
go test ./... -cover
```

## 🎯 当前项目状态

**已确认**: 项目使用 `*_test.go` 命名规范

**已实施文件**:
- `internal/auth/jwt_test.go` ✅
- `internal/testutil/*.go` ✅ (测试工具包)

**待实施文件**:
- `internal/model/*_test.go` (计划中)
- `internal/repository/*_test.go` (计划中)
- `internal/service/*_test.go` (计划中)
- `internal/handler/*_test.go` (计划中)

## 🚀 下一步行动

继续按照100%测试覆盖率计划实施，所有测试文件将遵循 `*_test.go` 命名规范。

### 本周任务 (Week 1)

1. ✅ 测试工具包完成
2. 🔄 Auth模块测试 (`jwt_test.go` 已完成)
3. ⏳ Payment模块测试
4. ⏳ Order模块测试

### 预期成果

- 所有测试文件: `*_test.go`
- 覆盖率: 0% → 5%
- 测试通过率: 100%

---

**规范确认**: ✅ 已接受 `*_test.go` 命名规范
**实施状态**: 🔄 进行中
**预计完成**: 2026-02-15
