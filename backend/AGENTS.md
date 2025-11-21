# GameLink Backend - AI Agent Reference

## 🎯 Project Overview

GameLink is a modern gaming companion management platform (现代化游戏陪玩管理平台) built with Go backend + React frontend architecture. The platform provides efficient order distribution, user management, and companion management services for gaming escort businesses.

**Repository**: https://github.com/HXSLtim/GameLink.git
**Current Version**: 0.3.0
**Primary Language**: Chinese (中文)

### Core Features
- 🎮 Intelligent order distribution (智能订单分发)
- 👥 Multi-role management (多角色管理) 
- 💬 Real-time communication (实时通讯)
- 💰 Complete payment system (完整支付)
- 📊 Data monitoring (数据监控)
- 🔐 Secure authentication (安全认证)

## 🏗️ Architecture Overview

### Clean 4-Layer Architecture
```
🌐 HTTP Requests
   ↓
🎯 Handler Layer (API处理) - internal/handler/
   ↓
💼 Service Layer (业务逻辑) - internal/service/
   ↓
🗄️ Repository Layer (数据访问) - internal/repository/
   ↓
📊 Model Layer (数据模型) - internal/model/
```

### Technology Stack
- **Backend**: Go 1.25.3 + Gin Framework
- **Database**: GORM with PostgreSQL/SQLite/MySQL support
- **Cache**: Redis + Memory cache
- **Authentication**: JWT (golang-jwt/jwt/v5)
- **Documentation**: Swagger/OpenAPI 3.0
- **Testing**: Standard Go testing + Testify
- **Container**: Docker + Alpine Linux
- **Dependency Injection**: Google Wire

### Key Design Principles
1. **Unified Service Model**: One table manages all services (including gifts)
2. **Unified Order Processing**: One logic handles all order types
3. **Layer Separation**: Strict separation of concerns
4. **Comprehensive Testing**: Multiple test types (unit, integration, quick)
5. **Configuration-Driven**: Environment-based configuration

## 📁 Project Structure

```
backend/
├── cmd/                          # Application entry points
│   ├── main.go                  # Main application entry
│   └── main_test.go             # Entry point tests
├── internal/                     # Internal packages (Go convention)
│   ├── model/                   # Data models and entities
│   ├── repository/              # Data access layer
│   │   ├── interfaces/          # Repository interfaces
│   │   ├── implementations/     # Concrete implementations
│   │   └── [module]/           # Module-specific repositories
│   ├── service/                 # Business logic layer
│   │   ├── admin/              # Admin business logic
│   │   ├── auth/               # Authentication logic
│   │   ├── order/              # Order processing
│   │   ├── payment/            # Payment processing
│   │   └── [module]/           # Other business modules
│   ├── handler/                 # HTTP request handlers
│   │   ├── admin/              # Admin API handlers
│   │   ├── user/               # User API handlers
│   │   ├── player/             # Player API handlers
│   │   └── middleware/         # HTTP middlewares
│   ├── config/                  # Configuration management
│   ├── container/               # Dependency injection (Wire)
│   ├── auth/                    # JWT authentication
│   ├── cache/                   # Caching implementations
│   ├── db/                      # Database connection & migration
│   ├── router/                  # Route definitions
│   ├── scheduler/               # Background jobs
│   └── lifecycle/               # Service lifecycle management
├── configs/                      # Configuration files
│   ├── config.development.yaml  # Development config
│   └── config.production.yaml   # Production config
├── docs/                         # Documentation
├── scripts/                      # Utility scripts
├── archive/                      # Archived files and reports
├── var/                          # Runtime data (SQLite DB)
├── go.mod                        # Go module definition
├── Makefile                      # Build automation
└── Dockerfile                    # Container definition
```

## 🚀 Quick Start Commands

### Development Setup
```bash
# Install dependencies
go mod tidy

# Run in development mode
go run ./cmd/main.go

# Run with specific service
make run CMD=main

# Generate Swagger documentation
make swagger

# Run tests
go test ./...
make test

# Code linting
make lint

# Build application
make build
```

### Docker Deployment
```bash
# Build Docker image
docker build -t gamelink-backend .

# Run container
docker run -p 8080:8080 -e APP_ENV=development gamelink-backend
```

## 🔧 Configuration

### Environment Variables
- `APP_ENV`: Environment (development/test/staging/production)
- `GIN_MODE`: Gin framework mode (debug/release/test)
- `SERVICE_PORT`: Server port (default: 8080)

### Configuration Files
- **Development**: `configs/config.development.yaml`
- **Production**: `configs/config.production.yaml`

### Key Configuration Sections
```yaml
server:
  port: "8080"
  enable_swagger: true

database:
  type: "sqlite"  # sqlite/postgres/mysql
  dsn: "file:./var/dev.db?mode=rwc"

cache:
  type: "memory"  # memory/redis
  redis:
    addr: "127.0.0.1:6379"

auth:
  jwt_secret: "your-secret-key"
  token_ttl_hours: 24

super_admin:
  email: "admin@gamelink.com"
  password: "123456"
```

## 🧪 Testing Strategy

### Test Types
1. **Unit Tests**: `*_test.go` files alongside source code
2. **Integration Tests**: Full service integration testing
3. **Quick Tests**: Fast validation tests
4. **Coverage Tests**: Comprehensive code coverage

### Test Naming Conventions
- Unit tests: `TestFunctionName`
- Integration tests: `TestIntegration_Scenario`
- Quick tests: `TestQuick_Scenario`

### Running Tests
```bash
# All tests
go test ./...

# Specific package
go test ./internal/service/order/...

# With coverage
go test -cover ./...

# Verbose output
go test -v ./...
```

## 📊 Business Domain

### User Roles
1. **Users (客户)**: Browse services, place orders, send gifts
2. **Players (陪玩师)**: Provide services, receive orders, earn income
3. **Admins (管理员)**: Manage platform, create services, handle disputes

### Service Categories
- **护航服务 (Escort Services)**: Solo/Team gaming companionship
- **礼物服务 (Gift Services)**: Virtual gifts for players

### Order Flow
1. Admin creates service items in `service_items` table
2. User browses and selects services
3. User places order in `orders` table
4. System processes payment
5. Player provides service
6. Order completed → Automatic commission calculation
7. Monthly settlement on 1st of each month

### Commission Model
- Default commission rate: 20%
- Configurable per service item
- Automatic calculation: `Commission = TotalPrice × CommissionRate`
- Player income: `TotalPrice - Commission`

## 🔒 Security Considerations

### Authentication
- JWT-based authentication with configurable TTL
- Role-based access control (RBAC)
- Admin authentication modes: JWT or traditional admin

### Data Protection
- Optional request/response encryption middleware
- CSRF protection
- Rate limiting
- Security headers

### Production Security
- Strong password requirements (min 8 chars, mixed case, digits, symbols)
- JWT secret must be changed from defaults
- Crypto keys must be secure (16/24/32 bytes)
- Database credentials must be properly configured

## 📚 API Documentation

### Swagger/OpenAPI
- Available at: `http://localhost:8080/swagger/index.html`
- Auto-generated from code annotations
- Supports generics and complex types

### API Structure
- Base path: `/api/v1`
- Authentication: Bearer token in Authorization header
- Response format: Standardized JSON with code/message/data

### Key Endpoints
```
GET  /api/v1/health              # Health check
POST /api/v1/auth/login          # User login
POST /api/v1/auth/admin/login    # Admin login
GET  /api/v1/user/services       # Browse services
POST /api/v1/user/orders         # Place order
GET  /api/v1/player/earnings     # View earnings
```

## 🛠️ Development Guidelines

### Code Style
- Follow standard Go conventions
- Use meaningful variable and function names
- Keep functions focused and small
- Add comments for complex logic
- Use constants for magic values

### File Naming
- Keep names concise and descriptive
- Avoid redundant suffixes (e.g., use `user.go` not `user_repository.go`)
- Test files: `source_test.go`

### Git Commits
Follow conventional commit format:
```
<type>(<scope>): <description>

feat(service): add user creation functionality
fix(handler): resolve order status update issue
docs(readme): update installation instructions
```

### Error Handling
- Return errors explicitly
- Wrap errors with context
- Use structured error responses
- Log errors appropriately

## 📈 Monitoring & Metrics

### Built-in Metrics
- Prometheus metrics endpoint
- Request/response metrics
- Database query metrics
- Cache hit/miss ratios

### Health Checks
- `/api/v1/health`: Basic health status
- `/api/v1/health/detailed`: Detailed system status
- Database connectivity checks
- Cache availability checks

## 🚨 Common Issues & Solutions

### Database Issues
- **SQLite locking**: Use proper connection pooling
- **Migration failures**: Check database permissions
- **Connection timeouts**: Verify network connectivity

### Configuration Issues
- **JWT validation failures**: Check secret key configuration
- **Port conflicts**: Verify port availability
- **Missing environment variables**: Check config validation

### Build Issues
- **Module dependencies**: Run `go mod tidy`
- **Swagger generation**: Ensure swag tool is installed
- **Docker builds**: Check Dockerfile syntax

## 📞 Support & Documentation

### Internal Documentation
- Complete docs: `docs/` directory
- Architecture summary: `docs/ARCHITECTURE_SUMMARY.md`
- Quick start: `docs/QUICK_START_UNIFIED.md`
- API response format: `docs/API_RESPONSE_FORMAT_*.md`

### Development Tools
- **golangci-lint**: Code quality enforcement
- **swag**: Swagger documentation generation
- **air**: Hot reload for development
- **testify**: Enhanced testing capabilities

### Getting Help
1. Check existing documentation in `docs/`
2. Review test files for usage examples
3. Check archived reports in `archive/`
4. Follow established patterns in existing code

---

**Remember**: This is a Chinese-language project with specific gaming industry domain requirements. Always consider the business context when making changes and maintain the existing code style and architectural patterns.