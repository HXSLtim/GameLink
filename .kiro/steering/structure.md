# Project Structure

## Repository Layout

```
GameLink/
├── backend/              # Go backend service
├── frontend/             # React frontend application
├── docs/                 # Project documentation
├── .kiro/                # Kiro configuration and steering rules
└── [root files]          # README, LICENSE, PRD documents
```

## Backend Structure (`backend/`)

```
backend/
├── cmd/
│   └── main.go                    # Application entry point
├── internal/                      # Private application code
│   ├── handler/                   # HTTP handlers (controllers)
│   │   ├── admin/                 # Admin panel endpoints
│   │   ├── middleware/            # HTTP middleware (auth, CORS, logging)
│   │   ├── notification/          # Notification endpoints
│   │   ├── player/                # Player/companion endpoints
│   │   ├── user/                  # User endpoints
│   │   ├── auth.go                # Authentication handlers
│   │   ├── error.go               # Error handling utilities
│   │   └── response.go            # Response formatting
│   ├── service/                   # Business logic layer
│   │   ├── admin/                 # Admin services
│   │   ├── auth/                  # Authentication services
│   │   ├── chat/                  # Chat services
│   │   ├── order/                 # Order management
│   │   ├── payment/               # Payment processing
│   │   ├── player/                # Player services
│   │   └── user/                  # User services
│   ├── repository/                # Data access layer
│   │   ├── implementations/       # Concrete implementations
│   │   ├── interfaces/            # Repository interfaces
│   │   ├── mocks/                 # Mock implementations for testing
│   │   ├── admin/                 # Admin data access
│   │   ├── order/                 # Order data access
│   │   ├── payment/               # Payment data access
│   │   └── user/                  # User data access
│   ├── model/                     # Data models and entities
│   │   ├── user.go                # User model
│   │   ├── order.go               # Order model
│   │   ├── player.go              # Player model
│   │   ├── payment.go             # Payment model
│   │   ├── chat.go                # Chat model
│   │   ├── role.go                # RBAC role model
│   │   └── permission.go          # RBAC permission model
│   ├── router/                    # Route definitions
│   ├── integration/               # Integration tests
│   └── ws/                        # WebSocket handlers
│       ├── hub.go                 # WebSocket hub
│       ├── client.go              # WebSocket client
│       └── message.go             # Message types
├── pkg/                           # Public/reusable packages
│   ├── auth/                      # JWT utilities
│   ├── cache/                     # Cache interfaces and implementations
│   ├── config/                    # Configuration management
│   ├── db/                        # Database utilities
│   ├── logging/                   # Logging utilities
│   ├── metrics/                   # Prometheus metrics
│   ├── safety/                    # Safety utilities (panic recovery)
│   ├── scheduler/                 # Scheduled jobs
│   └── testutil/                  # Testing utilities
├── configs/                       # Configuration files
│   ├── config.development.yaml    # Dev environment config
│   └── config.production.yaml     # Production config
├── docs/                          # Generated API documentation
│   ├── docs.go                    # Swagger docs
│   ├── swagger.json               # OpenAPI spec (JSON)
│   └── swagger.yaml               # OpenAPI spec (YAML)
├── scripts/                       # Utility scripts
├── var/                           # Runtime data (SQLite DB for dev)
├── go.mod                         # Go module definition
├── go.sum                         # Go dependencies checksum
├── Makefile                       # Build and test commands
└── Dockerfile                     # Container image definition
```

## Frontend Structure (`frontend/`)

```
frontend/
├── src/
│   ├── api/                       # API client modules
│   │   ├── auth.ts                # Authentication API
│   │   ├── order.ts               # Order API
│   │   ├── player.ts              # Player API
│   │   └── user.ts                # User API
│   ├── components/                # Reusable React components
│   │   ├── common/                # Common UI components
│   │   ├── layout/                # Layout components
│   │   └── [feature]/             # Feature-specific components
│   ├── pages/                     # Page components (routes)
│   │   ├── admin/                 # Admin panel pages
│   │   ├── player/                # Player dashboard pages
│   │   ├── user/                  # User-facing pages
│   │   ├── Login.tsx              # Login page
│   │   └── Register.tsx           # Registration page
│   ├── layouts/                   # Layout wrappers
│   │   ├── AdminLayout.tsx        # Admin panel layout
│   │   ├── PlayerLayout.tsx       # Player dashboard layout
│   │   └── UserLayout.tsx         # User-facing layout
│   ├── router/                    # Route configuration
│   │   └── index.tsx              # Route definitions
│   ├── services/                  # Business logic services
│   ├── context/                   # React Context providers
│   │   ├── AuthContext.tsx        # Authentication context
│   │   └── ThemeContext.tsx       # Theme context
│   ├── hooks/                     # Custom React hooks
│   │   ├── useAuth.ts             # Authentication hook
│   │   └── useWebSocket.ts        # WebSocket hook
│   ├── types/                     # TypeScript type definitions
│   │   ├── api.ts                 # API response types
│   │   ├── models.ts              # Data model types
│   │   └── index.ts               # Type exports
│   ├── utils/                     # Utility functions
│   │   ├── request.ts             # Axios configuration
│   │   ├── storage.ts             # LocalStorage utilities
│   │   └── validation.ts          # Form validation
│   ├── constants/                 # Constants and enums
│   ├── config/                    # App configuration
│   ├── assets/                    # Static assets (images, fonts)
│   ├── App.tsx                    # Root component
│   ├── main.tsx                   # Application entry point
│   └── index.css                  # Global styles
├── public/                        # Public static files
├── dist/                          # Build output (gitignored)
├── node_modules/                  # Dependencies (gitignored)
├── package.json                   # NPM dependencies and scripts
├── tsconfig.json                  # TypeScript configuration
├── vite.config.ts                 # Vite build configuration
├── eslint.config.js               # ESLint configuration
└── .prettierrc                    # Prettier configuration
```

## Documentation Structure (`docs/`)

```
docs/
├── api/                           # API design standards
│   ├── api-design-standards.md
│   └── go-coding-standards.md
├── backend/                       # Backend documentation
│   ├── README.md
│   ├── AGENTS.md
│   └── configs/
├── frontend/                      # Frontend documentation
│   ├── README.md
│   ├── DEVELOPER_GUIDE.md
│   └── features/
├── guides/                        # Development guides
│   ├── CONTRIBUTING.md
│   └── project-structure.md
├── reports/                       # Status and coverage reports
│   ├── coverage/
│   └── status/
├── archive/                       # Historical documentation
└── INDEX.md                       # Documentation index
```

## Naming Conventions

### Backend (Go)

- **Files**: `snake_case.go` (e.g., `user_service.go`, `order_repository.go`)
- **Packages**: lowercase, single word (e.g., `user`, `order`, `payment`)
- **Types**: PascalCase (e.g., `UserService`, `OrderRepository`)
- **Functions/Methods**: PascalCase for exported, camelCase for private
- **Variables**: camelCase (e.g., `userID`, `orderStatus`)
- **Constants**: PascalCase or UPPER_SNAKE_CASE for exported
- **Test files**: `*_test.go` (e.g., `user_service_test.go`)
- **Integration tests**: `*_integration_test.go`

### Frontend (TypeScript/React)

- **Files**: PascalCase for components (e.g., `UserProfile.tsx`), camelCase for utilities (e.g., `formatDate.ts`)
- **Components**: PascalCase (e.g., `UserProfile`, `OrderList`)
- **Functions**: camelCase (e.g., `fetchUserData`, `handleSubmit`)
- **Variables**: camelCase (e.g., `userId`, `orderStatus`)
- **Constants**: UPPER_SNAKE_CASE (e.g., `API_BASE_URL`, `MAX_RETRY_COUNT`)
- **Types/Interfaces**: PascalCase (e.g., `User`, `OrderResponse`)
- **CSS classes**: kebab-case (e.g., `user-profile`, `order-list`)

## Key Architectural Patterns

### Backend Patterns

1. **Dependency Injection**: Using Google Wire for compile-time DI
2. **Repository Pattern**: Abstraction over data access
3. **Service Layer**: Business logic separation
4. **Middleware Chain**: Request processing pipeline
5. **Error Wrapping**: Context-aware error propagation

### Frontend Patterns

1. **Component Composition**: Reusable, composable components
2. **Custom Hooks**: Shared stateful logic
3. **Context API**: Global state management
4. **Route-based Code Splitting**: Lazy loading pages
5. **API Client Abstraction**: Centralized HTTP requests

## File Organization Rules

1. **Group by feature**: Related files stay together (e.g., all user-related code in `user/`)
2. **Separate concerns**: Handler → Service → Repository → Model layers
3. **Test proximity**: Test files next to implementation files
4. **Shared code in pkg/**: Reusable utilities in `pkg/` directory
5. **Configuration separation**: Environment-specific configs in `configs/`
6. **Documentation co-location**: Feature docs near feature code when appropriate

## Import Organization

### Go Imports

```go
import (
    // Standard library
    "context"
    "fmt"
    "time"

    // Third-party packages
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"

    // Internal packages
    "gamelink/internal/model"
    "gamelink/internal/repository"
    "gamelink/pkg/auth"
)
```

### TypeScript Imports

```typescript
// React and third-party
import React, { useState, useEffect } from 'react';
import { Button, Form } from 'antd';
import axios from 'axios';

// Internal modules
import { User } from '@/types/models';
import { fetchUserData } from '@/api/user';
import { useAuth } from '@/hooks/useAuth';
```
