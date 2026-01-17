# GameLink Client Infrastructure - Implementation Guide

> **Version**: 1.0
> **Date**: 2026-01-17
> **Author**: Super Dev Team
> **Status**: Ready for Implementation

---

## 📋 Table of Contents

1. [Quick Start](#quick-start)
2. [Environment Setup](#environment-setup)
3. [Core Components](#core-components)
4. [Usage Examples](#usage-examples)
5. [Testing Strategy](#testing-strategy)
6. [Deployment Checklist](#deployment-checklist)
7. [Troubleshooting](#troubleshooting)

---

## 🚀 Quick Start

### Prerequisites

```bash
# Required dependencies (already in package.json)
- axios: ^1.6.0
- crypto-js: ^4.2.0
- zustand: ^4.4.0
- react-router-dom: ^6.20.0
```

### Installation Steps

```bash
# 1. Install dependencies (if not already installed)
cd client
npm install

# 2. Configure environment variables
cp .env.example .env.local

# 3. Update .env.local with required values
VITE_API_BASE_URL=http://localhost:8080/api/v1
VITE_CRYPTO_ENABLED=true
VITE_CRYPTO_SECRET_KEY=your-32-character-secret-key-here
VITE_CRYPTO_IV=your-16-char-iv
VITE_CRYPTO_USE_SIGNATURE=true

# 4. Start development server
npm run dev
```

---

## ⚙️ Environment Setup

### Development Environment (.env.development)

```env
# API Configuration
VITE_API_BASE_URL=http://localhost:8080/api/v1

# Encryption (optional in development)
VITE_CRYPTO_ENABLED=false
VITE_CRYPTO_SECRET_KEY=dev-secret-key-32-characters!!
VITE_CRYPTO_IV=dev-iv-16-chars
VITE_CRYPTO_USE_SIGNATURE=true
```

### Production Environment (.env.production)

```env
# API Configuration
VITE_API_BASE_URL=https://api.gamelink.com/api/v1

# Encryption (MANDATORY in production)
VITE_CRYPTO_ENABLED=true
VITE_CRYPTO_SECRET_KEY=<generate-secure-32-char-key>
VITE_CRYPTO_IV=<generate-secure-16-char-iv>
VITE_CRYPTO_USE_SIGNATURE=true
```

**Security Note**: Generate production keys using:
```bash
# Secret Key (32 characters)
openssl rand -base64 32 | cut -c1-32

# IV (16 characters)
openssl rand -base64 16 | cut -c1-16
```

---

## 🏗️ Core Components

### 1. Enhanced HTTP Client

**Location**: `client/src/lib/http.ts`

**Features**:
- ✅ Proactive JWT token refresh (5-minute buffer)
- ✅ Request encryption (AES-256-CBC)
- ✅ Request queue during token refresh
- ✅ Auto-unwrap API responses
- ✅ Unified error handling

**Usage**:

```typescript
import { http } from '@/lib/http';

// GET request
const users = await http.get<User[]>('/users');

// POST request with encryption
const order = await http.post<Order>('/order/create', {
    serviceItemId: 1,
    duration: 2
});

// Error handling
try {
    await http.post('/order/create', data);
} catch (error) {
    console.error(getErrorMessage(error));
}
```

### 2. Encryption Utilities

**Location**: `client/src/lib/crypto.ts`

**Features**:
- ✅ AES-256-CBC encryption
- ✅ SHA-256 signature generation
- ✅ Conditional encryption (POST/PUT/PATCH only)
- ✅ Production safety checks

**Usage**:

```typescript
import { encryptRequest, shouldEncrypt } from '@/lib/crypto';

// Automatic encryption (handled by HTTP client)
// Manual encryption (if needed)
const encrypted = encryptRequest({ sensitive: 'data' });
```

### 3. Error Handling

**Location**: `client/src/lib/error.ts`

**Features**:
- ✅ Error code to message mapping (40+ codes)
- ✅ User-friendly Chinese messages
- ✅ Error type checking utilities
- ✅ Axios error parsing

**Usage**:

```typescript
import { getErrorMessage, isAuthError, ERROR_CODE_MESSAGES } from '@/lib/error';

try {
    await orderApi.create(data);
} catch (error) {
    // Get user-friendly message
    const message = getErrorMessage(error);
    toast.error(message);

    // Check error type
    if (isAuthError(error)) {
        navigate('/login');
    }
}
```

### 4. Route Guards

**Location**: `client/src/router/Guard.tsx`

**Features**:
- ✅ Authentication check
- ✅ Role-based access control
- ✅ View mode enforcement (user/player)
- ✅ Zustand hydration handling

**Usage**:

```typescript
import { RouteGuard, PublicOnlyGuard } from '@/router/Guard';

// Protected route (requires authentication)
<RouteGuard requiresAuth>
    <ProfilePage />
</RouteGuard>

// Player-only route
<RouteGuard requiresAuth roles={['player']} viewMode="player">
    <PlayerDashboard />
</RouteGuard>

// Public-only route (login page)
<PublicOnlyGuard>
    <LoginPage />
</PublicOnlyGuard>
```

### 5. API Layer

**Location**: `client/src/api/*.ts`

**Structure**:
```
api/
├── auth.ts          # Authentication
├── order.ts         # Order management
├── payment.ts       # Payment processing
├── player.ts        # Player profiles
├── user.ts          # User profiles
├── wallet.ts        # Wallet operations
├── chat.ts          # Chat messaging
├── dispute.ts       # Dispute handling
└── ... (28+ more modules)
```

**Usage**:

```typescript
import { authApi } from '@/api/auth';
import { orderApi } from '@/api/order';

// Login
const { token, user } = await authApi.login({
    username: 'user@example.com',
    password: 'password123'
});

// Create order
const order = await orderApi.create({
    serviceItemId: 1,
    playerId: 123,
    duration: 2
});

// Get order list
const { items, total } = await orderApi.list({
    page: 1,
    pageSize: 20,
    status: 'pending'
});
```

---

## 💡 Usage Examples

### Example 1: User Login Flow

```typescript
// pages/auth/LoginPage.tsx
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuthStore } from '@/stores/modules/auth-store';
import { getErrorMessage } from '@/lib/error';

export function LoginPage() {
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');
    const login = useAuthStore((state) => state.login);
    const navigate = useNavigate();

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError('');

        try {
            await login({
                username: email,
                password: password
            });
            navigate('/');
        } catch (err) {
            setError(getErrorMessage(err));
        } finally {
            setLoading(false);
        }
    };

    return (
        <form onSubmit={handleSubmit}>
            {/* Form fields */}
            {error && <div className="error">{error}</div>}
            <button type="submit" disabled={loading}>
                {loading ? '登录中...' : '登录'}
            </button>
        </form>
    );
}
```

### Example 2: Create Order with Error Handling

```typescript
// pages/order/CreateOrderPage.tsx
import { useState } from 'react';
import { orderApi } from '@/api/order';
import { getErrorMessage, isValidationError, parseError } from '@/lib/error';

export function CreateOrderPage() {
    const [loading, setLoading] = useState(false);

    const handleCreateOrder = async (data: CreateOrderRequest) => {
        setLoading(true);

        try {
            const order = await orderApi.create(data);
            toast.success('订单创建成功');
            navigate(`/order/${order.id}`);
        } catch (error) {
            // Parse error for detailed information
            const parsed = parseError(error);

            if (isValidationError(error)) {
                // Handle validation error
                toast.error(`验证失败: ${parsed.message}`);
                if (parsed.field) {
                    // Highlight specific field
                    setFieldError(parsed.field, parsed.message);
                }
            } else {
                // Generic error
                toast.error(getErrorMessage(error));
            }
        } finally {
            setLoading(false);
        }
    };

    return (
        <div>
            {/* Order form */}
        </div>
    );
}
```

### Example 3: Protected Route Setup

```typescript
// App.tsx or router configuration
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { RouteGuard, PublicOnlyGuard } from '@/router/Guard';

function App() {
    return (
        <BrowserRouter>
            <Routes>
                {/* Public routes */}
                <Route path="/" element={<HomePage />} />

                {/* Auth routes (public only) */}
                <Route path="/login" element={
                    <PublicOnlyGuard>
                        <LoginPage />
                    </PublicOnlyGuard>
                } />

                {/* Protected user routes */}
                <Route path="/profile" element={
                    <RouteGuard requiresAuth>
                        <ProfilePage />
                    </RouteGuard>
                } />

                {/* Protected player routes */}
                <Route path="/player/dashboard" element={
                    <RouteGuard requiresAuth roles={['player']} viewMode="player">
                        <PlayerDashboard />
                    </RouteGuard>
                } />

                {/* 403 Forbidden */}
                <Route path="/403" element={<ForbiddenPage />} />
            </Routes>
        </BrowserRouter>
    );
}
```

### Example 4: View Mode Switching

```typescript
// components/ViewModeSwitcher.tsx
import { useAuthStore } from '@/stores/modules/auth-store';

export function ViewModeSwitcher() {
    const viewMode = useAuthStore((state) => state.viewMode);
    const isPlayer = useAuthStore((state) => state.isPlayer);
    const switchToPlayerMode = useAuthStore((state) => state.switchToPlayerMode);
    const switchToUserMode = useAuthStore((state) => state.switchToUserMode);

    if (!isPlayer) return null;

    return (
        <div className="view-mode-switcher">
            <button
                onClick={switchToUserMode}
                className={viewMode === 'user' ? 'active' : ''}
            >
                用户模式
            </button>
            <button
                onClick={switchToPlayerMode}
                className={viewMode === 'player' ? 'active' : ''}
            >
                陪玩模式
            </button>
        </div>
    );
}
```

---

## 🧪 Testing Strategy

### Unit Tests

```typescript
// lib/__tests__/http.test.ts
import { describe, it, expect, vi } from 'vitest';
import { HttpClient, parseJWT, isTokenExpiringSoon } from '../http';

describe('HTTP Client', () => {
    it('should parse JWT token correctly', () => {
        const token = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...';
        const payload = parseJWT(token);
        expect(payload).toBeDefined();
        expect(payload?.sub).toBe(123);
    });

    it('should detect expiring token', () => {
        const expiringToken = createTokenExpiringIn(4 * 60); // 4 minutes
        expect(isTokenExpiringSoon(expiringToken)).toBe(true);
    });
});
```

### Integration Tests

```typescript
// api/__tests__/auth.integration.test.ts
import { describe, it, expect } from 'vitest';
import { authApi } from '../auth';

describe('Auth API Integration', () => {
    it('should login successfully', async () => {
        const response = await authApi.login({
            username: 'test@example.com',
            password: 'password123'
        });

        expect(response.token).toBeDefined();
        expect(response.user).toBeDefined();
    });

    it('should handle invalid credentials', async () => {
        await expect(
            authApi.login({
                username: 'invalid@example.com',
                password: 'wrong'
            })
        ).rejects.toThrow();
    });
});
```

---

## ✅ Deployment Checklist

### Pre-Production

- [ ] Environment variables configured (`.env.production`)
- [ ] Encryption enabled (`VITE_CRYPTO_ENABLED=true`)
- [ ] Secure keys generated (32-char secret, 16-char IV)
- [ ] API base URL points to production backend
- [ ] All unit tests passing (`npm run test`)
- [ ] No console errors in production build
- [ ] Bundle size optimized (<500KB gzipped)

### Security Checklist

- [ ] JWT tokens stored in localStorage (acceptable for SPA)
- [ ] No sensitive data logged in production
- [ ] HTTPS enforced (backend configuration)
- [ ] CSRF tokens enabled (backend)
- [ ] Content Security Policy headers configured (backend)
- [ ] Rate limiting enabled on sensitive endpoints (backend)

### Performance Checklist

- [ ] Code splitting implemented for routes
- [ ] Images optimized and lazy-loaded
- [ ] API responses cached where appropriate
- [ ] Bundle analyzed (`npm run build:analyze`)
- [ ] Lighthouse score >90 for performance

---

## 🔧 Troubleshooting

### Issue: Token Refresh Loop

**Symptoms**: Infinite token refresh requests

**Solution**:
```typescript
// Check if token is actually expired
const token = useAuthStore.getState().token;
if (token && !isTokenExpired(token)) {
    // Token is still valid, don't refresh
    return;
}
```

### Issue: Encryption Errors in Production

**Symptoms**: "Encryption configuration error" in console

**Solution**:
1. Verify `.env.production` has correct keys
2. Ensure keys are exactly 32 and 16 characters
3. Check that `VITE_CRYPTO_ENABLED=true`

```bash
# Regenerate keys
openssl rand -base64 32 | cut -c1-32
openssl rand -base64 16 | cut -c1-16
```

### Issue: Route Guard Not Working

**Symptoms**: Unauthorized users accessing protected routes

**Solution**:
1. Check Zustand hydration: `useAuthStore.persist?.hasHydrated()`
2. Verify auth state: `console.log(useAuthStore.getState())`
3. Ensure `requiresAuth` prop is set on `RouteGuard`

### Issue: CORS Errors

**Symptoms**: "CORS policy" errors in console

**Solution** (Backend):
```go
// api/internal/handler/middleware/cors.go
router.Use(cors.New(cors.Config{
    AllowOrigins: []string{"http://localhost:5173", "https://gamelink.com"},
    AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
    AllowHeaders: []string{"Authorization", "Content-Type"},
    AllowCredentials: true,
}))
```

---

## 📚 Additional Resources

- [Architecture Design](./CLIENT_INFRASTRUCTURE_DESIGN.md)
- [Security Threat Model](./CLIENT_SECURITY_THREAT_MODEL.md)
- [Project Quick Reference](../.kiro/steering/QUICKSTART.md)
- [Data Models](../.kiro/steering/04-data-models.md)

---

**Document Status**: ✅ Complete
**Implementation Status**: Ready for Development
**Next Steps**: Begin implementation following this guide
