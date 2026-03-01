# Agent 2 - Batch Operations API

## Task
实现批量操作API接口

## Endpoints
- POST /api/v1/admin/users/batch/role
- POST /api/v1/admin/users/batch/status
- POST /api/v1/admin/users/batch/points
- POST /api/v1/admin/users/batch/notify

## Files
- Create: api/internal/handler/admin/dto/batch.go
- Modify: api/internal/handler/admin/user.go
- Modify: api/internal/service/admin/admin_service.go
- Modify: api/internal/router/admin.go

## Success Criteria
- 所有批量接口正常工作
- 使用事务保证数据一致性
- 测试通过: cd api && make test
