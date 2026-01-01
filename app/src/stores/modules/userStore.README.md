# UserStore Documentation

## Overview

The `userStore` manages user profile, wallet, and VIP state for the GameLink Taro app. It provides actions for fetching user data, updating profiles, uploading avatars, and managing wallet information.

## Features

- **User Profile Management**: Fetch and update user information
- **Avatar Upload**: Upload user avatar images
- **Wallet Balance**: Track available and frozen balances
- **VIP Status**: Manage VIP level and expiration
- **Persistence**: Automatically persists user data to local storage
- **Computed Values**: Helper functions for VIP status and balance conversions

## Installation

The store is already set up in your project. Simply import it in your components:

```typescript
import { useUserStore } from '@/stores/modules/userStore';
```

## API Reference

### State

| Property | Type | Description |
|----------|------|-------------|
| `userInfo` | `UserInfo \| null` | Current user information |
| `loading` | `boolean` | Loading state for async operations |
| `wallet` | `WalletInfo` | Wallet and VIP information |

### WalletInfo Interface

```typescript
interface WalletInfo {
  balanceCents: number;       // Available balance in cents (分)
  frozenCents: number;        // Frozen balance in cents (分)
  vipLevel: number;           // VIP level (0 = non-VIP)
  vipExpireAt: string \| null; // VIP expiration time (ISO 8601)
}
```

### Actions

#### `fetchUserInfo()`

Fetches the current user's information from the server.

**API Endpoint**: `GET /api/v1/auth/me`

```typescript
const { fetchUserInfo, userInfo, loading } = useUserStore();

await fetchUserInfo();

if (loading) {
  console.log('Loading...');
}

if (userInfo) {
  console.log('User:', userInfo.name);
}
```

#### `updateProfile(data: Partial<UserInfo>)`

Updates the user's profile information.

**API Endpoint**: `PUT /api/v1/user/profile`

**Parameters**:
- `data`: Partial user info object (e.g., `{ name: 'New Name' }`)

```typescript
const { updateProfile } = useUserStore();

try {
  await updateProfile({ name: 'John Doe' });
  // Success toast shown automatically
} catch (error) {
  console.error('Update failed:', error);
  // Error toast shown automatically
}
```

#### `uploadAvatar(filePath: string)`

Uploads an avatar image and updates the user's avatar URL.

**API Endpoint**: `POST /api/v1/upload/avatar`

**Parameters**:
- `filePath`: Local file path to the image

**Returns**: Promise<string> - The uploaded avatar URL

```typescript
const { uploadAvatar } = useUserStore();

Taro.chooseImage({
  count: 1,
  success: async (res) => {
    try {
      const avatarUrl = await uploadAvatar(res.tempFilePaths[0]);
      console.log('Avatar uploaded:', avatarUrl);
    } catch (error) {
      console.error('Upload failed:', error);
    }
  },
});
```

#### `fetchWallet()`

Fetches the user's wallet balance from the server.

**API Endpoint**: `GET /api/v1/user/wallet/balance`

```typescript
const { fetchWallet, wallet } = useUserStore();

await fetchWallet();

console.log('Balance:', wallet.balanceCents);
console.log('Frozen:', wallet.frozenCents);
```

### Selectors (Computed Values)

#### `isVip()`

Checks if the user has an active VIP membership.

**Returns**: `boolean`

```typescript
const { isVip } = useUserStore();

if (isVip()) {
  console.log('User is a VIP member');
} else {
  console.log('User is not a VIP member');
}
```

#### `vipDaysLeft()`

Calculates the number of remaining VIP days.

**Returns**: `number` - Days remaining (0 if not VIP or expired)

```typescript
const { vipDaysLeft } = useUserStore();

const days = vipDaysLeft();
console.log(`${days} days of VIP remaining`);
```

#### `balanceYuan()`

Converts balance from cents to yuan (元).

**Returns**: `number` - Balance in yuan

```typescript
const { balanceYuan } = useUserStore();

const balance = balanceYuan();
console.log(`Balance: ¥${balance.toFixed(2)}`);
```

#### `frozenBalanceYuan()`

Converts frozen balance from cents to yuan (元).

**Returns**: `number` - Frozen balance in yuan

```typescript
const { frozenBalanceYuan } = useUserStore();

const frozen = frozenBalanceYuan();
console.log(`Frozen: ¥${frozen.toFixed(2)}`);
```

## Usage Examples

### Basic User Info Display

```typescript
import { useUserStore } from '@/stores/modules/userStore';

const UserProfile = () => {
  const { userInfo, loading, fetchUserInfo } = useUserStore();

  useEffect(() => {
    fetchUserInfo();
  }, []);

  if (loading) return <Text>Loading...</Text>;
  if (!userInfo) return <Text>No user info</Text>;

  return (
    <View>
      <Text>Name: {userInfo.name}</Text>
      <Text>Phone: {userInfo.phone}</Text>
      <Image src={userInfo.avatar || '/default-avatar.png'} />
    </View>
  );
};
```

### Wallet Balance Display

```typescript
const WalletDisplay = () => {
  const { wallet, fetchWallet, balanceYuan, frozenBalanceYuan } = useUserStore();

  useEffect(() => {
    fetchWallet();
  }, []);

  return (
    <View>
      <Text>Available: ¥{balanceYuan().toFixed(2)}</Text>
      <Text>Frozen: ¥{frozenBalanceYuan().toFixed(2)}</Text>
      <Button onClick={fetchWallet}>Refresh</Button>
    </View>
  );
};
```

### VIP Status Check

```typescript
const VipBadge = () => {
  const { wallet, isVip, vipDaysLeft } = useUserStore();

  if (!isVip()) {
    return <Text>Regular User</Text>;
  }

  return (
    <View>
      <Text>VIP Level {wallet.vipLevel}</Text>
      <Text>{vipDaysLeft()} days left</Text>
    </View>
  );
};
```

### Update Profile Form

```typescript
const UpdateProfileForm = () => {
  const { userInfo, updateProfile } = useUserStore();
  const [name, setName] = useState(userInfo?.name || '');

  const handleSubmit = async () => {
    try {
      await updateProfile({ name });
    } catch (error) {
      console.error('Update failed:', error);
    }
  };

  return (
    <View>
      <Input value={name} onInput={(e) => setName(e.detail.value)} />
      <Button onClick={handleSubmit}>Update</Button>
    </View>
  );
};
```

## Persistence

The store uses Zustand's `persist` middleware to automatically save user data to local storage.

**What is persisted**:
- `userInfo`: User profile information
- `wallet.vipLevel`: VIP level
- `wallet.vipExpireAt`: VIP expiration time

**What is NOT persisted**:
- `loading`: Loading state
- `wallet.balanceCents`: Balance (always fetched from server)
- `wallet.frozenCents`: Frozen balance (always fetched from server)

This ensures that sensitive financial data is always fresh from the server, while basic user info is available offline.

## Error Handling

All actions include automatic error handling:

1. **Toast notifications** are shown for success and error states
2. **Errors are logged** to the console for debugging
3. **Loading state** is properly managed
4. **Silent failures** for non-critical operations (e.g., wallet fetch)

## Best Practices

1. **Fetch on mount**: Always call `fetchUserInfo()` when your component mounts
2. **Refresh after actions**: Call `fetchWallet()` after payments or balance changes
3. **Handle loading states**: Show loading indicators while data is being fetched
4. **Check VIP status**: Use `isVip()` selector instead of checking `wallet.vipLevel` directly
5. **Format currency**: Use `balanceYuan()` and `frozenBalanceYuan()` for display

## API Integration

The store integrates with the following backend endpoints:

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/api/v1/auth/me` | GET | Get current user info |
| `/api/v1/user/profile` | PUT | Update user profile |
| `/api/v1/upload/avatar` | POST | Upload avatar image |
| `/api/v1/user/wallet/balance` | GET | Get wallet balance |

All requests include the JWT token from local storage automatically.

## TypeScript Support

The store is fully typed with TypeScript:

```typescript
import type { UserState, WalletInfo } from '@/stores/modules/userStore';
```

All actions, state properties, and computed values have proper type definitions.

## Testing

When testing components that use `userStore`, you can mock the store:

```typescript
import { render } from '@testing-library/react';
import { useUserStore } from '@/stores/modules/userStore';

// Mock the store
jest.mock('@/stores/modules/userStore');

test('displays user name', () => {
  (useUserStore as jest.Mock).mockReturnValue({
    userInfo: { name: 'Test User' },
    loading: false,
  });

  // Test your component
});
```

## Related Stores

- **authStore**: Manages authentication tokens and login state
- **orderStore**: Manages user orders
- **chatStore**: Manages chat messages and rooms

## Contributing

When adding new features to the userStore:

1. Update the `UserState` interface
2. Implement the action with proper error handling
3. Add TypeScript types
4. Update this documentation
5. Add usage examples

## License

This store is part of the GameLink project.
