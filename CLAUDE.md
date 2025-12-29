# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

GameLink is a modern game companion (陪玩) management platform - a full-stack system connecting users with gaming companions for paid gaming sessions. The platform features intelligent order distribution, multi-role management (users, players, admins), real-time chat, payment processing, and comprehensive marketing features (VIP, coupons, referrals, teams).

**Project Status**: Backend modules are 100% complete (36/36 modules). Frontend admin panel is at ~75% completion.

### Key Business Concepts

- **Commission Structure**: 15-25% platform commission via three-tier calculation (service item rate + player individual rate + monthly ranking adjustment)
- **Pricing Model**: System-controlled pricing ¥20-60+/hour based on player rank (players cannot customize rates)
- **User Roles**: user (customers), player (companions), admin (platform operators)
- **Order Types**: solo (single player), team (multiple players), gift (direct payment)
- **Dispute Handling**: Dual-CS mechanism (original + independent customer service), 30-minute SLA
- **Income Settlement**: T+7 holding period before players can withdraw earnings

> See [`.kiro/steering/01-product.md`](.kiro/steering/01-product.md) for complete product overview and business rules.

## Technology Stack

| Layer | Backend (Go) | Frontend (Admin) |
|-------|-------------|------------------|
| **Language** | Go 1.25+ | TypeScript 5.9+ |
| **Framework** | Gin (web), GORM (ORM) | React 19, Vite 7 |
| **UI Library** | - | Ant Design 6.0 |
| **Database** | PostgreSQL 16+, Redis 7+ | - |
| **Auth** | JWT (golang-jwt/jwt/v5) | crypto-js |
| **WebSocket** | gorilla/websocket | socket.io-client |
| **Testing** | testify, mockery | Vitest, Testing Library |
| **Docs** | Swagger (swaggo/swag) | - |
| **Security** | AES-256-CBC + SHA-256 | crypto-js |

## Common Commands

### Backend (api/)

```bash
cd api

# Run tests
make test              # Run all tests
make test-coverage     # With coverage report
make test-race         # With race detector
make test-integration  # Integration tests (requires PostgreSQL)

# Linting and formatting
make lint              # golangci-lint
make fmt               # go fmt
make check             # Run all checks (fmt, vet, lint, test)

# Code generation
make swagger          # Generate Swagger docs
make generate-mocks   # Generate mocks with mockery

# Dependencies
make deps             # Install dependencies
make test-tools       # Install test tools

# Run specific package test
make run-test PKG=service/user
make cover-pkg PKG=service/user

# Run application
go run cmd/main.go
```

### Frontend (admin/)

```bash
cd admin

npm install           # Install dependencies
npm run dev           # Dev server (localhost:5173)
npm run build         # Production build
npm run build:analyze # Build with bundle analysis
npm run lint          # ESLint
npm run test          # Vitest (watch mode)
npm run test:run      # Run tests once
```

### Docker

```bash
# Development environment
docker-compose up -d

# Test environment
docker-compose -f docker-compose.test.yml up -d

# Production deployment (encrypted, recommended)
.\scripts\deploy-production-encrypted.ps1

# With options
.\scripts\deploy-production-encrypted.ps1 -SkipBuild    # Skip Docker image build
.\scripts\deploy-production-encrypted.ps1 -RegenerateKeys  # Regenerate crypto keys
```

### Production Deployment Environment Variables

```bash
# Database
POSTGRES_USER=gamelink
POSTGRES_PASSWORD=<safe password, no special characters>
POSTGRES_DB=gamelink

# Redis
REDIS_PASSWORD=<safe password>

# JWT
JWT_SECRET_KEY=<32+ characters>

# Encryption (required for production)
CRYPTO_ENABLED=true
CRYPTO_SECRET_KEY=<32 characters>
CRYPTO_IV=<16 characters>

# Super Admin
SUPER_ADMIN_EMAIL=admin@gamelink.com
SUPER_ADMIN_PASSWORD=<8+ chars with upper/lower/number/special>
```

**Important Notes**:
- Production must enable encryption middleware (`CRYPTO_ENABLED=true`)
- Cache must be set to redis (`CACHE_TYPE=redis`)
- Super admin role slug is `superAdmin` (camelCase)
- Backend health check path: `/api/v1/healthz`
- Database password cannot contain URL special characters (like `%`)

### CI/CD Pipeline

| Pipeline | Purpose | Key Features |
|----------|---------|--------------|
| ci.yml | Continuous Integration | Change detection, race detection, 70% coverage check, Docker build |
| security.yml | Security Scanning | Gosec, govulncheck, Trivy, Gitleaks (runs weekly) |
| deploy.yml | Deployment | Tag/manual trigger, staging/production, auto health check, rollback |

### Integration Tests

```bash
# Start test database first
docker-compose -f docker-compose.test.yml up -d

# Run integration tests
go test ./api/internal/service/integration/... -v

# Run specific integration test
go test ./api/internal/service/integration/order_integration_test.go -v

# Run with coverage
go test ./api/internal/service/integration/... -cover -coverprofile=coverage.out
```

## Architecture

### Repository Structure

```
GameLink/
├── api/                  # Go backend (monolithic)
│   ├── cmd/main.go       # Application entry point
│   ├── internal/
│   │   ├── handler/      # HTTP handlers (admin/, user/, player/, middleware/)
│   │   ├── service/      # Business logic layer
│   │   ├── repository/   # Data access layer
│   │   ├── model/        # Data models
│   │   ├── router/       # Route definitions
│   │   └── ws/           # WebSocket handlers
│   ├── pkg/              # Public reusable packages
│   ├── configs/          # Config files (development.yaml, production.yaml)
│   └── Makefile
├── admin/                # React admin panel
│   ├── src/
│   │   ├── api/          # API clients (auth.ts, admin.ts, client.ts)
│   │   ├── components/   # Reusable components
│   │   ├── pages/        # Page components
│   │   └── ...
│   └── package.json
├── app/                  # Taro mini-program
├── client/               # User/Player frontend (to be developed)
├── scripts/              # Deployment scripts
└── docs/                 # Documentation
```

### Backend Layered Architecture

```
Handler → Service → Repository → Model
```

- **Handler** (`api/internal/handler/`): HTTP request handling, validation, response formatting
- **Service** (`api/internal/service/`): Business logic, transaction management, cross-module coordination
- **Repository** (`api/internal/repository/`): Data access abstraction, caching, query encapsulation
- **Model** (`api/internal/model/`): Data structures, DB mappings, validation rules

### Module Organization

The backend is organized into 36 modules across 4 categories:

| Category | Modules | Status |
|----------|---------|--------|
| **Core** | user, auth, order, payment, player, chat, dispute, etc. (19) | ✅ 100% |
| **New Business** | player-rank, order-timeout, user-block, vip (4) | ✅ 100% |
| **Marketing** | vip, coupon, recharge, activity, team, referral (6) | ✅ 100% |
| **Auxiliary** | commission, ranking, routing-rule, settlement-company, etc. (7) | ✅ 100% |

See `[.kiro/steering/06-project-management.md`](.kiro/steering/06-project-management.md) for detailed status.

## Naming Conventions

### Go

| Type | Convention | Example |
|------|-----------|---------|
| Files | camelCase | `userService.go`, NOT `user_service.go` |
| Packages | lowercase | `handler`, `service` |
| Types | PascalCase | `UserService`, `Order` |
| Exported functions | PascalCase | `CreateUser`, `GetOrder` |
| Private functions | camelCase | `validateInput`, `calculatePrice` |
| Variables | camelCase | `userID`, `orderService` |
| Test files | *_test.go | `userService_test.go` |

### TypeScript/React

| Type | Convention | Example |
|------|-----------|---------|
| Components | PascalCase | `UserProfile.tsx` |
| Utilities | camelCase | `formatDate.ts` |
| Types/Interfaces | PascalCase | `UserResponse` |
| Constants | UPPER_SNAKE_CASE | `API_BASE_URL` |
| CSS classes | kebab-case | `user-profile` |

## API Response Patterns

### Unified Response Format

Use the `resp` package for all HTTP responses:

```go
// Success responses
resp.OK(c, data)           // 200 with data
resp.Success(c, message, data)  // 200 with custom message
resp.Created(c, data)      // 201
resp.Updated(c, data)      // 200 for updates
resp.Deleted(c)            // 200 for delete

// List with pagination
resp.List(c, items, pagination)

// Error responses (using apierr package)
resp.Error(c, apierr.BadRequest("invalid input"))
resp.Error(c, apierr.NotFound("user not found"))
resp.Error(c, apierr.Unauthorized("authentication required"))
```

### Error Handling

Use the `apierr` package for standardized errors:

```go
// Common errors (predefined)
apierr.ErrNotFound
apierr.ErrUnauthorized
apierr.ErrForbidden
apierr.ErrInvalidInput
apierr.ErrInternal

// Create custom errors
apierr.BadRequest("invalid phone number")
apierr.NotFound("user not found")
apierr.Unauthorized("invalid token")

// Add details
apierr.BadRequest("validation failed").WithDetails("email is required").WithField("email")
```

### Frontend API Convention

Axios interceptor auto-parses `response.data.data`, so components use data directly:

```typescript
// No need to double-nest .data
const users = await getUsers(); // Returns User[] directly, not Response<User[]>
```

## Testing

### Test Structure

- **Unit tests**: `*_test.go` alongside source files
- **Integration tests**: `api/internal/service/integration/*_integration_test.go`
- **Mock files**: `api/internal/repository/mocks/*.go` (generated by mockery)

### Coverage Goals

Current average: ~80% (service layer)

| Module | Coverage | Status |
|--------|----------|--------|
| menu | 100.0% | ✅ |
| handler/resp | 96.0% | ✅ |
| item | 90.9% | ✅ |
| player | 86.8% | ✅ |
| user | 84.5% | ✅ |
| auth | 84.3% | ✅ |
| permission | 84.2% | ✅ |
| withdraw | 81.4% | ✅ |
| order | 79.2% | ⚠️ |
| payment | 75.8% | ⚠️ |
| routingrule | 50.7% | ⚠️ |
| admin | 1.2% | ❌ |

**Target**: 80%+ for all modules

### Integration Test Helpers

Located in `api/internal/service/integration/testdb.go`:

```go
// Database setup
SkipIfNoTestDB(t)      // Skip if TEST_DB_* not set
db := SetupTestDB(t)   // Initialize and auto-cleanup

// Create test data
CreateTestUser(t, db, "name")
CreateUniqueTestUser(t, db, "prefix")  // Ensures uniqueness
CreateTestPlayer(t, db, user)
CreateTestOrder(t, db, user, player, status)
CreateTestPayment(t, db, order, status)
CreateTestWallet(t, db, userID, balanceCents)
// ... and 30+ more helpers
```

### Test Naming

```go
// Service test
func Test{ServiceName}_{MethodName}(t *testing.T) {
    SkipIfNoTestDB(t)
    db := SetupTestDB(t)
    // ...
}

// Scenario test
func Test{ServiceName}_{MethodName}_{Scenario}(t *testing.T) {
    // Example: TestOrderService_CreateOrder_WithCoupon
}
```

## Configuration

### Backend Config Files

- `api/configs/config.development.yaml` - Dev (SQLite, memory cache, crypto disabled)
- `api/configs/config.production.yaml` - Prod (PostgreSQL, Redis, crypto enabled)

### Security Features

- **Production**: AES-256-CBC + SHA-256 signature enforced
- **Development**: Crypto can be disabled for debugging
- **Middleware**: `api/internal/handler/middleware/crypto.go` handles encryption

## Important Conventions

### Commit Messages

Follow Conventional Commits:

```
feat(user): add user registration feature
fix(order): resolve order status update issue
docs(api): update payment API documentation
refactor(payment): simplify payment logic
test(chat): add integration tests for chat service
```

### Import Organization (Go)

```go
import (
    // Standard library
    "context"
    "fmt"

    // Third-party packages
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"

    // Internal packages
    "gamelink/internal/model"
    "gamelink/pkg/auth"
)
```

## Key Documentation Files

### Steering Rules (.kiro/steering/)

These are the **authoritative source** for project conventions and must be consulted when making changes.

| File | Purpose |
|------|---------|
| [01-product.md](.kiro/steering/01-product.md) | **Product overview**, core features, business model, commission rules, user roles |
| [02-tech-stack.md](.kiro/steering/02-tech-stack.md) | **Technology stack**, CI/CD, deployment scripts, environment variables |
| [03-project-structure.md](.kiro/steering/03-project-structure.md) | **Repository structure**, naming conventions (Go files use camelCase, NOT snake_case) |
| [04-data-models.md](.kiro/steering/04-data-models.md) | **Data model hierarchy**, business rules for each module, commission calculation |
| [04c-enums-indexes.md](.kiro/steering/04c-enums-indexes.md) | **Enums and indexes** with changelog |
| [05-testing-standard.md](.kiro/steering/05-testing-standard.md) | **Testing standards** - 3-level testing approach (frontend/backend/database validation) |
| [06-project-management.md](.kiro/steering/06-project-management.md) | **Module completion status** - 100% backend complete (36/36 modules) |

### Other Documentation

| File | Purpose |
|------|---------|
| [docs/INTEGRATION_TEST_PLAN.md](docs/INTEGRATION_TEST_PLAN.md) | Integration test planning, test helpers, execution commands |
| [PROGRESS.md](PROGRESS.md) | Development progress tracking |
| [Program.md](Program.md) | Systematic project health report |

### Important Data Model Documents

- [04a-marketing-models.md](.kiro/steering/04a-marketing-models.md) - VIP, coupons, recharge, activity, referral
- [04b-team-models.md](.kiro/steering/04b-team-models.md) - Team system
- [04d-notification-models.md](.kiro/steering/04d-notification-models.md) - Notification system

## Business Context

> **IMPORTANT**: All business logic is documented in [`.kiro/steering/04-data-models.md`](.kiro/steering/04-data-models.md). Always reference this file before implementing business features.

### Commission Structure (Three-Tier)

| Priority | Source | Description |
|----------|--------|-------------|
| 1 | Player individual | `CommissionRule.PlayerID` - specific rate for a player |
| 2 | Service item | `ServiceItem.CommissionRate` - default 20% |
| 3 | Monthly ranking | `RankingCommissionConfig` - tiered reduction based on last month's rank |

**Formula**: Final commission = Base rate - Ranking discount

Example: ¥100 order, 20% base, player ranked #5 (5% discount) → 15% commission → ¥85 player income

### Order Types

| Type | Description | RequiredPlayers | Payment Flow |
|------|-------------|-----------------|--------------|
| solo | Single companion | 1 | Standard |
| team | Multiple companions | 2+ | Match all slots before starting |
| gift | Direct payment (no service) | 1 | Immediate completion, no T+7 |

### User Block System

**Effects**:
- No messaging between blocked users
- Hidden from each other's lists
- Order room isolation (blocked player invisible to user)

**Rules**:
- Orders in progress continue even after blocking
- Cannot place new orders to blocked players
- Blocking is directional - A blocking B ≠ B blocking A

### Dispute Handling

| Phase | Time Limit | Action |
|-------|-----------|--------|
| Filing | Order complete + 7 days | User or player can initiate |
| SLA | 30 minutes | CS must respond |
| Resolution | - | Full refund or reject |

**Dual-CS Mechanism**:
- Original CS (from order, if any)
- Independent CS (unrelated to order, for fairness)

### Income Settlement

```
Order complete → Income to FrozenCents (T+7 hold)
                ↓
              7 days
                ↓
            No issues → FrozenCents → BalanceCents (withdrawable)
            Dispute → SettlementStatus = disputed (continues frozen)
            Refund → Deduct from FrozenCents
```

### Testing Philosophy

> From [`.kiro/steering/05-testing-standard.md`](.kiro/steering/05-testing-standard.md)

**Three-Level Testing**:
1. **Level 1** - Frontend rendering (NOT sufficient alone)
2. **Level 2** - Request/response validation (⭐ focus)
3. **Level 3** - Database and business logic (⭐ focus)

**Test Checklist**:
- Request sent with correct parameters
- Response status and code correct
- Database state changed correctly
- Page feedback displayed correctly
- Exception scenarios tested

## Development Workflow

When working on this codebase:

1. **Before implementing**: Check [`.kiro/steering/04-data-models.md`](.kiro/steering/04-data-models.md) for business rules
2. **For new features**: Check [`.kiro/steering/06-project-management.md`](.kiro/steering/06-project-management.md) - most backend modules are complete
3. **For bug fixes**: Follow Handler→Service→Repository pattern; update tests
4. **For frontend work**: The admin panel uses React 19 + Ant Design; check existing API clients
5. **After changes**: Run `make test` and `make lint` before committing
6. **Model changes**: Update [`.kiro/steering/04-data-models.md`](.kiro/steering/04-data-models.md) to match

## ⚠️ IMPORTANT: Document Sync Requirements

> **CRITICAL**: After completing ANY work, you MUST update the relevant steering documents to maintain project consistency.

### Required Document Updates After Work Completion

| Work Type | Documents to Update | Description |
|-----------|-------------------|-------------|
| **Model Changes** | `.kiro/steering/04-data-models.md` | Add/update model definitions, fields, enums |
| **New Features** | `.kiro/steering/06-project-management.md` | Mark module status, update progress |
| **API Changes** | `.kiro/steering/02-tech-stack.md` | Update API documentation status |
| **Business Rules** | `.kiro/steering/01-product.md` | Document new business logic or rules |
| **Test Coverage** | `.kiro/steering/05-testing-standard.md` | Update test coverage statistics |
| **Enum Changes** | `.kiro/steering/04c-enums-indexes.md` | Add enum values with changelog entry |
| **Module Completion** | `PROGRESS.md` | Update overall project progress |
| **Test Plans** | `docs/INTEGRATION_TEST_PLAN.md` | Add new test scenarios, update status |

### Update Workflow

After completing a task, always follow this sequence:

```bash
# 1. Code changes completed
# Example: Added new fields to Order model

# 2. Run tests and ensure they pass
make test

# 3. Update steering documents
# - Edit .kiro/steering/04-data-models.md to document new fields
# - Edit .kiro/steering/06-project-management.md if module status changed
# - Edit PROGRESS.md to reflect progress

# 4. Commit with proper message
git add .
git commit -m "feat(order): add new fields for XYZ feature"
```

### Steering Document Ownership

| Document | Maintainer | Update Frequency |
|----------|-----------|------------------|
| `01-product.md` | Product/Business | Per feature change |
| `02-tech-stack.md` | Tech Lead | Per tech stack change |
| `03-project-structure.md` | Tech Lead | Per structural change |
| `04-data-models.md` | Backend | **Per model change** ⚠️ |
| `04c-enums-indexes.md` | Backend | **Per enum change** ⚠️ |
| `05-testing-standard.md` | QA/Backend | Per test change |
| `06-project-management.md` | Project Lead | **Per module completion** ⚠️ |

### Why This Matters

The `.kiro/steering/` directory is the **single source of truth** for:
- Business rules and logic
- Data model definitions
- Module completion status
- Testing standards
- Project structure conventions

**Failing to update these documents causes**:
- Confusion about actual project status
- Outdated business rules documentation
- Misalignment between code and documentation
- wasted time for future developers

### Example: Proper Completion Workflow

```markdown
# Task: Add new field `DeliveryMethod` to Order model

## Step 1: Implement code
- Add field to Order model
- Update migration
- Add tests

## Step 2: Update documentation
- Edit .kiro/steering/04-data-models.md:
  - Add DeliveryMethod to Order table
  - Document enum values (standard/express)
  - Update business rules section

## Step 3: Update progress
- Edit PROGRESS.md if this completes a feature
- Edit .kiro/steering/06-project-management.md if module status changed

## Step 4: Commit
git add .kiro/steering/04-data-models.md PROGRESS.md
git commit -m "feat(order): add delivery method field"
```

