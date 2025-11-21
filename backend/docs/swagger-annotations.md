# GameLink Swagger/OpenAPI 注解规范文档

## 📋 概述
本文档定义了 GameLink 项目中 Swagger/OpenAPI 注解的使用规范，旨在消除重复注解、统一注解格式，并提供最佳实践指导。

## 🎯 目标
- 消除重复的 Swagger 注解
- 统一注解格式和风格
- 提供标准化的注解模板
- 确保 Swagger 文档的一致性

## 📊 当前问题分析

### 1. 重复路由定义问题
**问题**: 在 `/mnt/c/Users/a2778/Desktop/code/GameLink/backend/internal/handler/admin/router.go` 和各个 handler 文件中存在重复的 `@Router` 注解定义。

**示例**:
```go
// 在 order.go 中定义:
// @Router       /admin/orders [get]
// @Router       /admin/orders [post]

// 在 router.go 中又重复定义:
// @Router       /admin/orders [get]  // ❌ 重复
// @Router       /admin/orders [post] // ❌ 重复
```

**影响文件**:
- `/mnt/c/Users/a2778/backend/internal/handler/admin/router.go` (33个重复路由)
- 所有 admin handler 文件

### 2. 重复的响应模型
**问题**: 大量使用相同的响应模型定义，缺乏标准化。

**统计**:
- `model.SuccessResponse` 使用超过 100 次
- `model.APIResponse[T]` 使用超过 70 次
- `model.ErrorResponse` 使用超过 50 次

**示例**:
```go
// 重复的定义模式:
// @Success      200  {object}  model.SuccessResponse
// @Failure      400  {object}  model.ErrorResponse
// @Failure      401  {object}  model.ErrorResponse
```

### 3. 注解格式不一致
**问题**: 同一种类型的注解在不同文件中格式不一致。

**示例**:
```go
// 格式1 (admin/commission.go):
// @Success      200            {object}  model.SuccessResponse

// 格式2 (admin/dashboard.go):
// @Success      200  {object}  model.SuccessResponse

// 格式3 (user/order.go):
// @Success      200       {object}  model.APIResponse[order.MyOrderListResponse]
```

### 4. 缺少描述信息
**问题**: 许多注解缺少必要的描述信息。

**示例**:
```go
// 问题示例:
// @Summary      列出用户
// @Description  API endpoint  // ❌ 描述过于简单

// 应该改进为:
// @Summary      获取用户列表
// @Description  获取分页的用户列表，支持按角色、状态、关键词等条件筛选
```

## ✅ 规范标准

### 1. 全局注解规范
**只允许在 main.go 中定义全局信息**:
```go
// @title           GameLink API
// @version         0.3.0
// @description     GameLink 平台 API，包含健康检查、认证与管理端能力
// @BasePath        /api/v1
// @schemes         http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
```

### 2. 路由注解规范
**每个 API 端点只在一个地方定义路由注解**:

**推荐做法**:
- 在具体的 handler 函数上定义路由注解
- 删除 router.go 中的重复路由注解
- 使用统一的标签格式

**标准格式**:
```go
// 统一的注解格式:
// @Summary      功能简述
// @Description  详细功能描述
// @Tags         模块/子模块
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        ...
// @Success      200  {object}  响应模型
// @Failure      400  {object}  model.ErrorResponse
// @Router       /路径 [method]
```

### 3. 响应模型标准化
**定义标准响应模型**:

#### 成功响应模型
```go
// 标准成功响应 (无数据返回)
// @Success 200 {object} model.SuccessResponse

// 标准成功响应 (有数据返回)
// @Success 200 {object} model.APIResponse[具体类型]

// 创建成功响应
// @Success 201 {object} model.APIResponse[具体类型]
```

#### 错误响应模型
```go
// 标准错误响应
// @Failure 400 {object} model.ErrorResponse  // 请求参数错误
// @Failure 401 {object} model.ErrorResponse  // 认证失败
// @Failure 403 {object} model.ErrorResponse  // 权限不足
// @Failure 404 {object} model.ErrorResponse  // 资源不存在
// @Failure 500 {object} model.ErrorResponse  // 服务器内部错误
```

### 4. 标签命名规范
**统一使用模块/子模块格式**:

```go
// 认证相关
// @Tags         Auth

// 用户端功能
// @Tags         User/Orders
// @Tags         User/Players
// @Tags         User/Payments

// 管理端功能
// @Tags         Admin/Users
// @Tags         Admin/Orders
// @Tags         Admin/Players

// 陪玩师端功能
// @Tags         Player/Orders
// @Tags         Player/Earnings
```

## 📋 注解模板

### 1. 基础 CRUD 操作模板

#### 列表查询
```go
// @Summary      获取列表
// @Description  获取分页的列表数据，支持筛选和排序
// @Tags         Admin/模块名
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        page       query     int      false  "页码"      default(1)
// @Param        pageSize   query     int      false  "每页数量"  default(20)
// @Param        keyword    query     string   false  "搜索关键词"
// @Success      200  {object} model.APIResponse[[]模型类型]
// @Failure      400  {object} model.ErrorResponse
// @Failure      401  {object} model.ErrorResponse
// @Router       /admin/模块名 [get]
```

#### 单条记录查询
```go
// @Summary      获取详情
// @Description  根据ID获取单条记录的详细信息
// @Tags         Admin/模块名
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path   int   true   "记录ID"
// @Success      200  {object} model.APIResponse[模型类型]
// @Failure      400  {object} model.ErrorResponse
// @Failure      401  {object} model.ErrorResponse
// @Failure      404  {object} model.ErrorResponse
// @Router       /admin/模块名/{id} [get]
```

#### 创建记录
```go
// @Summary      创建记录
// @Description  创建新的记录
// @Tags         Admin/模块名
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request   body   创建请求类型   true   "创建信息"
// @Success      201  {object} model.APIResponse[模型类型]
// @Failure      400  {object} model.ErrorResponse
// @Failure      401  {object} model.ErrorResponse
// @Router       /admin/模块名 [post]
```

#### 更新记录
```go
// @Summary      更新记录
// @Description  根据ID更新记录信息
// @Tags         Admin/模块名
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id        path   int           true   "记录ID"
// @Param        request   body   更新请求类型  true   "更新信息"
// @Success      200  {object} model.SuccessResponse
// @Failure      400  {object} model.ErrorResponse
// @Failure      401  {object} model.ErrorResponse
// @Failure      404  {object} model.ErrorResponse
// @Router       /admin/模块名/{id} [put]
```

#### 删除记录
```go
// @Summary      删除记录
// @Description  根据ID删除记录（软删除）
// @Tags         Admin/模块名
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path   int   true   "记录ID"
// @Success      200  {object} model.SuccessResponse
// @Failure      400  {object} model.ErrorResponse
// @Failure      401  {object} model.ErrorResponse
// @Failure      404  {object} model.ErrorResponse
// @Router       /admin/模块名/{id} [delete]
```

### 2. 认证相关模板

#### 登录
```go
// @Summary      用户登录
// @Description  使用用户名/邮箱/手机号 + 密码进行登录，返回 JWT Token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request   body   LoginRequest   true   "登录凭据"
// @Success      200  {object} model.APIResponse[LoginResponse]
// @Failure      400  {object} model.ErrorResponse
// @Failure      401  {object} model.ErrorResponse
// @Router       /auth/login [post]
```

#### 注册
```go
// @Summary      用户注册
// @Description  创建新用户账号，默认角色为普通用户
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request   body   RegisterRequest   true   "注册信息"
// @Success      201  {object} model.APIResponse[RegisterResponse]
// @Failure      400  {object} model.ErrorResponse
// @Router       /auth/register [post]
```

## 🔧 实施步骤

### 第一步：清理重复路由注解
1. 备份现有代码
2. 删除 `/mnt/c/Users/a2778/Desktop/code/GameLink/backend/internal/handler/admin/router.go` 中的所有 Swagger 注解
3. 确保每个 handler 函数都有完整的路由注解

### 第二步：标准化响应模型
1. 统一使用 `model.APIResponse[T]` 作为成功响应
2. 统一使用 `model.ErrorResponse` 作为错误响应
3. 移除不规范的响应定义

### 第三步：格式化注解
1. 统一注解缩进和格式
2. 补充缺失的描述信息
3. 标准化标签命名

### 第四步：验证和测试
1. 运行 `swag init` 生成 Swagger 文档
2. 验证文档正确性
3. 测试所有 API 端点

## 📋 检查清单

### 代码审查检查项
- [ ] 每个 API 端点都有完整的路由注解
- [ ] 响应模型使用统一的格式
- [ ] 标签命名符合模块/子模块规范
- [ ] 描述信息详细且准确
- [ ] 参数注解完整且正确
- [ ] 错误响应状态码定义正确

### Swagger 文档检查项
- [ ] 文档可以正常生成
- [ ] 所有路由都能正确显示
- [ ] 请求参数和响应模型正确
- [ ] 标签分类清晰合理
- [ ] 安全定义正确配置

## 🚨 注意事项

### 1. 避免常见错误
```go
// ❌ 错误的注解格式
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/users [get]

// ✅ 正确的注解格式
// @Success 200 {object} model.APIResponse[[]User]
// @Router /users [get]
```

### 2. 参数注解规范
```go
// ❌ 缺少参数描述
// @Param id path int true ""

// ✅ 完整的参数描述
// @Param id path int true "用户ID" minimum(1)
```

### 3. 响应状态码规范
```go
// 使用标准的 HTTP 状态码
// 200 - 成功（有数据返回）
// 201 - 创建成功
// 204 - 成功（无数据返回）
// 400 - 请求参数错误
// 401 - 认证失败
// 403 - 权限不足
// 404 - 资源不存在
// 500 - 服务器内部错误
```

## 📚 参考资源

- [Swagger 官方文档](https://swagger.io/docs/)
- [Swaggo 文档](https://github.com/swaggo/swag)
- [OpenAPI 规范](https://spec.openapis.org/oas/v3.1.0)

## 📞 支持

如有疑问，请联系项目维护者或参考本文档的最佳实践部分。

---

**文档版本**: 1.0.0
**最后更新**: 2025-01-16
**维护者**: Claude Code