# AGENTS.md - Backend Infrastructure Fix

## 🎯 Your Task
Fix backend infrastructure issues to prepare for frontend-backend integration.

## 📋 Task List

### Task 1: Fix Port Conflict (0.5 day)
**Problem**: Backend default port 8080 conflicts with admin frontend dev server.

**Solution**:
- Modify `api/cmd/main.go` to change default port from `:8080` to `:8081`
- Update any related configuration files
- Update documentation if needed (README.md, docker-compose files)

**Files to modify**:
- `api/cmd/main.go`
- `docker-compose.yml` (if port mapping exists)
- `docker-compose.dev.yml` (if port mapping exists)

### Task 2: Standardize DTO Models (1 day)
**Problem**: Some handlers use models directly instead of dedicated DTOs, causing:
- Data leakage risk
- API instability
- Unclear contract between frontend and backend

**Solution**:
Create dedicated DTO models in `api/internal/handler/admin/dto/`:

**Priority modules** (start with these):
1. Extend `api/internal/handler/admin/dto/player.go` if missing batch DTOs
2. Extend `api/internal/handler/admin/dto/order.go` if missing batch DTOs

**DTO Pattern to Follow** (根据项目规范):
```go
// api/internal/handler/admin/dto/batch.go
package dto

import "gamelink/internal/model"

// BatchRoleUpdateRequest 批量修改角色请求
type BatchRoleUpdateRequest struct {
    UserIDs []uint64     `json:"userIds" binding:"required,min=1,max=100"`
    Role    model.Role   `json:"role" binding:"required,oneof=user player admin"`
}

// BatchOperationResponse 批量操作响应
type BatchOperationResponse struct {
    Success bool              `json:"success"`
    Message string            `json:"message"`
    Data    BatchOperationResult `json:"data"`
}

// BatchOperationResult 批量操作结果
type BatchOperationResult struct {
    Updated int          `json:"updated"`
    Failed  int          `json:"failed"`
    Errors  []BatchError `json:"errors"`
}

// BatchError 批量操作错误
type BatchError struct {
    UserID uint64 `json:"userId"`
    Reason string `json:"reason"`
}
```

**Files to create/modify**:
- Create: `api/internal/handler/admin/dto/batch.go`
- Modify: `api/internal/handler/admin/user.go` (add batch methods)

## 🚫 Files to NOT Touch
- `api/internal/handler/user/**` (User-facing endpoints)
- `api/internal/handler/player/**` (Player endpoints)
- `api/internal/ws/**` (WebSocket handlers - Agent 3's domain)
- Frontend code in `admin/` or `app/`

## 📐 Patterns to Follow (遵循项目代码规范)

### Port Configuration
```go
// api/cmd/main.go
func main() {
    port := os.Getenv("APP_PORT")
    if port == "" {
        port = "8081" // Changed from 8080
    }
    r.Run(":" + port)
}
```

### Handler Pattern (使用项目现有模式)
```go
// BatchRoleUpdate
// @Summary      批量修改用户角色
// @Description  批量修改用户角色
// @Tags         Admin/Users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body      dto.BatchRoleUpdateRequest  true  "批量修改角色请求"
// @Success      200      {object}  dto.BatchOperationResponse
// @Router       /admin/users/batch/role [post]
func (h *UserHandler) BatchRoleUpdate(c *gin.Context) {
    var req dto.BatchRoleUpdateRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        respondError(c, apierr.BadRequest("参数错误").WithDetails(err.Error()))
        return
    }

    result, err := h.svc.BatchRoleUpdate(c.Request.Context(), &req)
    if err != nil {
        respondError(c, err)
        return
    }
    respondSuccess(c, result)
}
```

## ✅ Success Criteria
- [ ] Backend starts on port 8081 by default
- [ ] Batch DTOs created in `api/internal/handler/admin/dto/`
- [ ] All existing tests still pass: `cd api && make test`
- [ ] Code follows existing patterns in codebase
- [ ] Swagger docs updated: `cd api && make swagger`

## 🧪 Testing
```bash
cd api

# Run all tests
make test

# Generate swagger docs
make swagger

# Check code quality
make check
```

## 📚 Reference
- Existing DTO pattern: `api/internal/handler/admin/dto/user.go`
- Handler pattern: `api/internal/handler/admin/user.go`
- Error handling: `gamelink/pkg/apierr`
- CLAUDE.md for coding standards

## 💡 Notes
- 使用中文注释和错误消息
- 遵循项目现有的 DTO 模式
- 使用 `gamelink` 模块前缀
- 使用 `apierr` 包处理错误
