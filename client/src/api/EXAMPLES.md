# GameLink Client API - Usage Examples

> **Complete guide for using the GameLink Client API modules**

## Table of Contents

- [Quick Start](#quick-start)
- [Authentication](#authentication)
- [Core APIs](#core-apis)
- [Business APIs](#business-apis)
- [Marketing APIs](#marketing-apis)
- [Error Handling](#error-handling)
- [Type Safety](#type-safety)

---

## Quick Start

### Import APIs

```typescript
// Import individual APIs
import { authApi, orderApi, playerApi } from '@/api';

// Or import all APIs
import * as api from '@/api';
```

### Basic Usage Pattern

All API modules follow a consistent pattern:

```typescript
// List with pagination
const { items, total, page, pageSize } = await orderApi.list({
    page: 1,
    pageSize: 20
});

// Get single item
const order = await orderApi.get(123);

// Create new item
const newOrder = await orderApi.create({
    serviceItemId: 1,
    playerId: 456,
    duration: 2
});

// Update item
const updated = await orderApi.update(123, { status: 'in_progress' });

// Delete item
await orderApi.delete(123);
```

---

## Authentication

### Login

```typescript
import { authApi } from '@/api';

try {
    const response = await authApi.login({
        username: 'user@example.com',
        password: 'password123'
    });

    console.log('Login successful:', response.user);
    // Token is automatically stored by the auth store
} catch (error) {
    console.error('Login failed:', getErrorMessage(error));
}
```

### Register

```typescript
const response = await authApi.register({
    phone: '+1234567890',
    password: 'SecurePass123!',
    name: 'John Doe'
});
```

### Get Current User

```typescript
const me = await authApi.me();
console.log('Current user:', me.user);
```

### Refresh Token

```typescript
// Token is automatically refreshed by the HTTP client
// You can manually trigger if needed:
const response = await authApi.refreshToken();
```

### Logout

```typescript
await authApi.logout();
// Clears token and redirects to login
```

---

## Core APIs

### User API

```typescript
import { userApi } from '@/api';

// Get user profile
const profile = await userApi.getProfile();

// Update profile
const updated = await userApi.updateProfile({
    name: 'New Name',
    avatar: 'https://example.com/avatar.jpg'
});

// Upload avatar
const formData = new FormData();
formData.append('file', avatarFile);
await userApi.uploadAvatar(formData);

// Get user preferences
const prefs = await userApi.getPreferences();

// Update preferences
await userApi.updatePreferences({
    language: 'zh-CN',
    theme: 'dark',
    notifications: {
        email: true,
        push: true,
        sms: false
    }
});
```

### Player API

```typescript
import { playerApi } from '@/api';

// Search players
const { items } = await playerApi.list({
    page: 1,
    pageSize: 20,
    gameId: 1,
    online: true,
    minRating: 4.0
});

// Get player profile
const profile = await playerApi.getProfile(123);

// Get player stats
const stats = await playerApi.getStats(123);

// Update own player profile
await playerApi.updateProfile({
    name: 'Professional Player',
    introduction: 'Experienced gamer ready to help!',
    games: [
        { gameId: 1, rankId: 5 }
    ]
});
```

### Order API

```typescript
import { orderApi } from '@/api';

// Create order
const order = await orderApi.create({
    playerId: 123,
    serviceItemId: 5,
    duration: 2,
    message: 'Looking for a ranked game partner'
});

// List my orders
const { items, total } = await orderApi.list({
    page: 1,
    pageSize: 20,
    status: 'pending'
});

// Get order details
const details = await orderApi.get(456);

// Accept order (for players)
await orderApi.accept(456);

// Start order
await orderApi.start(456);

// Complete order
await orderApi.complete(456);

// Cancel order
await orderApi.cancel(456, 'Changed my mind');
```

### Payment API

```typescript
import { paymentApi } from '@/api';

// Create payment
const payment = await paymentApi.create({
    orderId: 123,
    method: 'wechat',
    amount: 5000
});

// Get payment status
const status = await paymentApi.getStatus(payment.id);

// List payment history
const { items } = await paymentApi.list({
    page: 1,
    pageSize: 20
});
```

### Chat API

```typescript
import { chatApi } from '@/api';

// Get chat rooms
const { items } = await chatApi.getRooms({
    page: 1,
    pageSize: 20
});

// Get messages
const { items: messages } = await chatApi.getMessages(roomId, {
    page: 1,
    pageSize: 50
});

// Send text message
await chatApi.sendMessage(roomId, {
    content: 'Hello! When can we start?',
    type: 'text'
});

// Send image message
const formData = new FormData();
formData.append('image', imageFile);
await chatApi.sendImage(roomId, formData);

// Mark as read
await chatApi.markAsRead(roomId);
```

### Wallet API

```typescript
import { walletApi } from '@/api';

// Get wallet info
const wallet = await walletApi.getWallet();

// Get transactions
const { items } = await walletApi.getTransactions({
    page: 1,
    pageSize: 20,
    type: 'income'
});

// Request withdrawal
await walletApi.withdraw({
    amount: 10000,
    method: 'alipay',
    account: 'user@example.com'
});
```

---

## Business APIs

### Game API

```typescript
import { gameApi } from '@/api';

// Get all games
const games = await gameApi.list();

// Get game details
const game = await gameApi.get(1);

// Get game ranks
const ranks = await gameApi.getRanks(1);
```

### Review API

```typescript
import { reviewApi } from '@/api';

// Create review
const review = await reviewApi.create({
    orderId: 123,
    rating: 5,
    content: 'Excellent service! Very skilled player.',
    images: ['url1', 'url2']
});

// List reviews for a player
const { items } = await reviewApi.list({
    page: 1,
    pageSize: 20,
    playerId: 123
});

// Reply to review (for players)
await reviewApi.reply(reviewId, 'Thank you for your feedback!');
```

---

## Marketing APIs

### VIP API

```typescript
import { vipApi } from '@/api';

// Get my VIP info
const vipInfo = await vipApi.getMyInfo();

// Get all VIP levels
const levels = await vipApi.getLevels();

// Get VIP benefits
const benefits = await vipApi.getBenefits(levelId);
```

### Coupon API

```typescript
import { couponApi } from '@/api';

// Get available coupons
const { items } = await couponApi.list({
    status: 'active'
});

// Get my coupons
const myCoupons = await couponApi.getMyCoupons();
```

### Recharge API

```typescript
import { rechargeApi } from '@/api';

// Get recharge packages
const packages = await rechargeApi.getPackages();

// Create recharge
const recharge = await rechargeApi.create({
    packageId: 1,
    paymentMethod: 'wechat'
});

// Get recharge history
const { items } = await rechargeApi.getHistory({
    page: 1,
    pageSize: 20
});
```

---

## Error Handling

### Basic Error Handling

```typescript
import { getErrorMessage } from '@/lib/error';

try {
    const order = await orderApi.create(data);
    toast.success('Order created successfully');
} catch (error) {
    const message = getErrorMessage(error);
    toast.error(message);
    // Error messages are automatically localized to Chinese
    // Examples:
    // - "余额不足，请先充值"
    // - "陪玩师当前不可接单"
    // - "订单状态不允许此操作"
}
```

### Advanced Error Handling

```typescript
import { getErrorMessage, isAuthError, isValidationError } from '@/lib/error';

try {
    await orderApi.create(data);
} catch (error) {
    if (isAuthError(error)) {
        // Redirect to login
        router.push('/login');
    } else if (isValidationError(error)) {
        // Show validation errors
        const details = error.response?.data?.details;
        showFormErrors(details);
    } else {
        const message = getErrorMessage(error);
        toast.error(message);
    }
}
```

### Error Types

```typescript
// Auth errors (401, 403)
if (isAuthError(error)) {
    // Handle authentication/authorization errors
}

// Validation errors (400)
if (isValidationError(error)) {
    // Handle input validation errors
}

// Network errors
if (!error.response) {
    toast.error('Network error. Please check your connection');
}

// Server errors (500)
if (error.response?.status >= 500) {
    toast.error('Server error. Please try again later');
}
```

---

## Type Safety

All API modules are fully typed with TypeScript:

```typescript
import type {
    Order,
    CreateOrderRequest,
    Player,
    PaginatedResponse
} from '@/types/api';

// Type-safe API calls
const createOrder = async (data: CreateOrderRequest): Promise<Order> => {
    return await orderApi.create(data);
};

// Type-safe responses
const { items, total }: PaginatedResponse<Order> = await orderApi.list({
    page: 1,
    pageSize: 20
});

// Type definitions are auto-completed in IDEs
const order: Order = {
    id: 123,
    userId: 456,
    playerId: 789,
    serviceItemId: 1,
    type: 'solo',
    status: 'pending',
    // ... all required fields with type checking
};
```

---

## Best Practices

### 1. Always Handle Errors

```typescript
try {
    await orderApi.create(data);
} catch (error) {
    // Always handle errors gracefully
    toast.error(getErrorMessage(error));
}
```

### 2. Use Loading States

```typescript
const [loading, setLoading] = useState(false);

const handleCreate = async () => {
    setLoading(true);
    try {
        await orderApi.create(data);
        toast.success('Success');
    } catch (error) {
        toast.error(getErrorMessage(error));
    } finally {
        setLoading(false);
    }
};
```

### 3. Pagination Pattern

```typescript
const [orders, setOrders] = useState<Order[]>([]);
const [page, setPage] = useState(1);
const [total, setTotal] = useState(0);

useEffect(() => {
    const loadOrders = async () => {
        const response = await orderApi.list({ page, pageSize: 20 });
        setOrders(response.items);
        setTotal(response.total);
    };
    loadOrders();
}, [page]);
```

### 4. Optimistic Updates

```typescript
// Update local state immediately
setOrders(prev => [newOrder, ...prev]);

// Then sync with server
try {
    await orderApi.create(data);
} catch (error) {
    // Rollback on error
    setOrders(prev => prev.slice(1));
    toast.error(getErrorMessage(error));
}
```

---

## Complete Component Example

```typescript
import React, { useState, useEffect } from 'react';
import { orderApi, playerApi } from '@/api';
import type { Order, Player } from '@/types/api';
import { getErrorMessage } from '@/lib/error';

export function OrderList() {
    const [orders, setOrders] = useState<Order[]>([]);
    const [loading, setLoading] = useState(true);
    const [page, setPage] = useState(1);

    useEffect(() => {
        loadOrders();
    }, [page]);

    const loadOrders = async () => {
        try {
            setLoading(true);
            const response = await orderApi.list({
                page,
                pageSize: 20
            });
            setOrders(response.items);
        } catch (error) {
            toast.error(getErrorMessage(error));
        } finally {
            setLoading(false);
        }
    };

    const handleCreateOrder = async (playerId: number) => {
        try {
            const order = await orderApi.create({
                playerId,
                serviceItemId: 1,
                duration: 2
            });
            setOrders(prev => [order, ...prev]);
            toast.success('Order created successfully');
        } catch (error) {
            toast.error(getErrorMessage(error));
        }
    };

    if (loading) return <div>Loading...</div>;

    return (
        <div>
            {orders.map(order => (
                <div key={order.id}>
                    <h3>Order #{order.id}</h3>
                    <p>Status: {order.status}</p>
                    <p>Amount: ¥{order.amount}</p>
                </div>
            ))}
        </div>
    );
}
```

---

## More Examples

For more detailed examples, see:

- [Implementation Guide](../docs/CLIENT_IMPLEMENTATION_GUIDE.md)
- [API Completion Report](../docs/CLIENT_API_COMPLETION_REPORT.md)
- [Infrastructure Design](../docs/CLIENT_INFRASTRUCTURE_DESIGN.md)
