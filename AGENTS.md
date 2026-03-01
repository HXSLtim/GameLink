# Agent 3 - Customer Service & Dashboard

## Task
实现在线客服后端接口和修复仪表盘接口

## Dashboard Endpoints
- GET /api/v1/admin/dashboard/stats
- GET /api/v1/admin/dashboard/user-behavior
- GET /api/v1/admin/dashboard/user-distribution

## Customer Service
- WebSocket: WS /api/v1/ws/customer-service
- REST: GET/POST /api/v1/user/customer-service/conversations

## Files
- Create: api/internal/model/conversation.go
- Create: api/internal/handler/user/customer_service_handler.go
- Create: api/internal/ws/customer_service.go
- Modify: api/internal/handler/admin/dashboard.go

## Success Criteria
- 仪表盘接口返回真实数据
- WebSocket实时消息功能正常
- 测试通过: cd api && make test
