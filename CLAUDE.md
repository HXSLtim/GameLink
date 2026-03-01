# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

GameLink is a gaming companion platform with three applications:
- **api/**: Go backend (Gin + GORM) - port 8080
- **admin/**: React admin panel (Ant Design) - port 5173
- **app/**: React user web app (shadcn/ui + Tailwind) - port 5175

## Commands

### Go Backend (api/)
```bash
make test              # Run all tests
make test-coverage     # Tests with coverage
make lint              # golangci-lint
make fmt               # Format code
make check             # fmt + vet + lint + test
make swagger           # Generate Swagger docs
make run-test PKG=pkg  # Test specific package
```

### Frontend (admin/ and app/)
```bash
pnpm dev      # Start dev server
pnpm build    # Production build
pnpm lint     # ESLint
pnpm test     # Vitest tests
```

### Docker
```bash
docker compose up -d                    # PostgreSQL + Redis
docker compose -f docker-compose.dev.yml up -d  # Full dev environment
```

## Architecture

### Backend (api/)

**Layered architecture:**
```
Handler (HTTP) → Service (Business Logic) → Repository (Data Access) → Model
```

**Key directories:**
- `internal/handler/` - HTTP handlers (admin/user/player/public groups)
- `internal/service/` - Business logic (57 modules)
- `internal/repository/` - Data access (56 modules)
- `internal/model/` - Data models (67 models)
- `internal/router/` - Route registration
- `internal/ws/` - WebSocket for real-time chat
- `pkg/` - Shared packages (auth, config, db, scheduler)

**API routes:**
| Prefix | Description | Auth |
|--------|-------------|------|
| `/api/v1/auth` | Login, register, refresh | Partial |
| `/api/v1/public` | Public endpoints | None |
| `/api/v1/user` | User endpoints | Required |
| `/api/v1/player` | Player endpoints | Required |
| `/api/v1/admin` | Admin endpoints | Required + RBAC |

### Frontend (admin/)

- **UI**: Ant Design 6
- **State**: Zustand
- **Testing**: Vitest + Playwright
- **Pages**: 40+ modules in `src/pages/`

### Frontend (app/)

- **UI**: shadcn/ui (Radix UI + Tailwind CSS 4)
- **State**: Zustand
- **Structure**:
  - `src/features/` - Business pages
  - `src/components/` - UI components
  - `src/services/` - API layer
  - `src/hooks/` - Custom hooks

## Code Style

### Go (api/)
- **Architecture**: Handler → Service → Repository pattern
- **Linter**: golangci-lint (govet, gofmt, goimports, gocyclo, typecheck, errcheck, dupl)
- **Complexity**: Max cyclomatic complexity 15
- **Imports**: Local prefix `gamelink` for goimports
- **Formatting**: gofmt with simplify enabled
- **Error handling**: Proper error returns, custom errors
- **Middleware**: For auth, logging, cross-cutting concerns
- **Testing**: testify package

### TypeScript (admin/ & app/)
- **Strict mode**: Enabled in tsconfig
- **Linter**: ESLint (typescript-eslint, react-hooks, react-refresh)
- **Formatter**: Prettier
  - semi: true, singleQuote: true, tabWidth: 2, printWidth: 100, trailingComma: es5
- **Types**:
  - `interface` for objects, `type` for unions/intersections
  - Avoid `any`, use `unknown` or generics
  - Use type guards over type assertions (`as`)
- **Package manager**: pnpm
- **Unused vars**: Prefix with `_` to ignore

## Database

- PostgreSQL 16+ with 80+ tables
- Redis 7+ for caching
- GORM for ORM
- Key features: RBAC, multi-tenancy, covering indexes

## Testing

- **Backend**: 159 test files, 70%+ coverage target
- **Admin**: 88 unit tests
- **CI**: GitHub Actions with quality gates

## Default Credentials

| Role | Email | Password |
|------|-------|----------|
| Super Admin | admin@gamelink.com | Admin123456 |

## Swagger

After starting backend: http://localhost:8080/swagger/index.html
