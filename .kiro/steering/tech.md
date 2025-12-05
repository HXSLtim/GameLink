# Technology Stack

## Backend

- **Language**: Go 1.25.3+
- **Web Framework**: Gin (HTTP routing and middleware)
- **ORM**: GORM (database operations with auto-migration)
- **Authentication**: JWT (golang-jwt/jwt/v5)
- **Database**: PostgreSQL (production), SQLite with pure Go driver (testing - glebarez/sqlite)
- **Cache**: Redis 6.0+ (go-redis/v9)
- **WebSocket**: gorilla/websocket
- **API Documentation**: Swagger/OpenAPI (swaggo/swag)
- **Dependency Injection**: Google Wire
- **Testing**: Go testing + testify, mockery for mocks
- **Monitoring**: Prometheus client

## Frontend

- **Framework**: React 18.2+ with TypeScript 5.2+
- **Build Tool**: Vite 7.2+
- **UI Library**: Ant Design 6.0
- **Routing**: React Router 7.9+
- **HTTP Client**: Axios 1.13+
- **State Management**: React Context API
- **Styling**: Less 4.2
- **WebSocket**: socket.io-client 4.8+
- **Testing**: Vitest 4.0+, Testing Library, Playwright
- **Code Quality**: ESLint, Prettier, Husky, lint-staged

## Infrastructure

- **Containerization**: Docker + Docker Compose
- **Message Queue**: RabbitMQ (planned)
- **Search**: Elasticsearch (planned)
- **Reverse Proxy**: Nginx
- **CI/CD**: GitHub Actions (planned)

## Common Commands

### Backend

```bash
# Navigate to backend
cd backend

# Install dependencies
go mod download
go mod tidy

# Run application
go run cmd/main.go

# Run tests
make test                    # All tests
make test-coverage          # With coverage report
make test-coverage-html     # Generate HTML report
make test-race              # With race detector
go test ./internal/integration -v  # Integration tests only

# Code quality
make lint                   # Run golangci-lint
make fmt                    # Format code
make vet                    # Go vet
make check                  # All checks (fmt, vet, lint, test)

# Generate documentation
make swagger                # Generate Swagger docs
make swagger-vendor         # With vendor support

# Build
go build -o bin/gamelink cmd/main.go

# Database migrations (via GORM AutoMigrate in code)
```

### Frontend

```bash
# Navigate to frontend
cd frontend

# Install dependencies
npm install

# Development server
npm run dev                 # Starts on http://localhost:5173

# Build
npm run build               # Production build
npm run preview             # Preview production build

# Testing
npm run test                # Run tests
npm run test:ui             # Vitest UI
npm run test:run            # Single run (CI mode)

# Code quality
npm run lint                # ESLint
```

### Docker

```bash
# Build and start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down

# Rebuild specific service
docker-compose build backend
docker-compose up -d backend
```

## Project Structure Conventions

### Backend Layered Architecture

```
Handler → Service → Repository → Model
```

- **Handler**: HTTP request handling, parameter validation, response formatting
- **Service**: Business logic, transaction management, cross-module coordination
- **Repository**: Database operations, caching, query encapsulation
- **Model**: Data structures, database mappings, validation rules

### Error Handling

Three-tier error mechanism:
- Repository errors: Database-level errors
- Service errors: Business logic errors with context
- API errors: Standardized HTTP responses with error codes

Use `fmt.Errorf("context: %w", err)` for error wrapping.

## Testing Strategy

- **Unit Tests**: Table-driven tests, mock dependencies
- **Integration Tests**: Real database (SQLite), test fixtures
- **Concurrency Tests**: Race detector, stress testing
- **Coverage Target**: 76.4% (current), aiming for 80%+

## Configuration

- Development: `backend/configs/config.development.yaml`
- Production: `backend/configs/config.production.yaml`
- Environment variables override config files
- SQLite for development/testing, PostgreSQL for production
