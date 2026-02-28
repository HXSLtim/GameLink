/// <reference types="node" />
/**
 * API Helper Functions for E2E Tests
 * Provides direct API access for test setup and teardown
 */

interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
}

const API_BASE = process.env.API_URL || 'http://localhost:8080/api/v1';

// Token cache to avoid repeated login requests (rate limit: 10/min)
let cachedToken: string | null = null;
let tokenExpiry: number = 0;
const TOKEN_VALIDITY_MS = 20 * 60 * 1000; // 20 minutes (conservative, actual is 24h)

/**
 * Sleep for specified milliseconds
 */
function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

/**
 * Get authentication token for admin user with caching and retry logic
 * Uses token caching to avoid hitting rate limits (10 requests/minute for login)
 */
export async function getAdminToken(): Promise<string> {
  // Return cached token if still valid
  if (cachedToken && Date.now() < tokenExpiry) {
    return cachedToken;
  }

  const adminPassword = process.env.TEST_ADMIN_PASSWORD || process.env.SUPER_ADMIN_PASSWORD;
  if (!adminPassword) {
    console.warn(
      '⚠️ TEST_ADMIN_PASSWORD environment variable not set. ' +
      'Using default password. Set it via: $env:TEST_ADMIN_PASSWORD="YourPassword"'
    );
  }

  const maxRetries = 5;
  let lastError: Error | null = null;

  for (let attempt = 1; attempt <= maxRetries; attempt++) {
    try {
      const response = await fetch(`${API_BASE}/auth/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          username: process.env.TEST_ADMIN_USERNAME || 'admin@gamelink.com',
          password: adminPassword || 'Admin123456',
        }),
      });

      if (response.status === 429) {
        // Rate limited - wait with exponential backoff
        const waitTime = Math.min(1000 * Math.pow(2, attempt), 30000); // Max 30s
        console.warn(`Rate limited on login attempt ${attempt}, waiting ${waitTime}ms...`);
        await sleep(waitTime);
        continue;
      }

      if (!response.ok) {
        throw new Error(`Failed to get admin token: ${response.status} ${response.statusText}`);
      }

      const result = (await response.json()) as ApiResponse<{ token: string }>;

      // Cache the token
      cachedToken = result.data.token;
      tokenExpiry = Date.now() + TOKEN_VALIDITY_MS;

      return cachedToken;
    } catch (error) {
      lastError = error as Error;
      if (attempt < maxRetries) {
        const waitTime = Math.min(1000 * Math.pow(2, attempt), 30000);
        console.warn(`Login attempt ${attempt} failed, retrying in ${waitTime}ms...`);
        await sleep(waitTime);
      }
    }
  }

  throw lastError || new Error('Failed to get admin token after max retries');
}

/**
 * Clear the cached token (useful for testing token refresh)
 */
export function clearTokenCache(): void {
  cachedToken = null;
  tokenExpiry = 0;
}

/**
 * Make an authenticated API request with retry logic
 */
async function authenticatedFetch(
  url: string,
  options: RequestInit,
  token: string,
  maxRetries = 3
): Promise<Response> {
  for (let attempt = 1; attempt <= maxRetries; attempt++) {
    const response = await fetch(url, {
      ...options,
      headers: {
        ...options.headers,
        Authorization: `Bearer ${token}`,
      },
    });

    if (response.status === 429) {
      // Rate limited - wait with exponential backoff
      const waitTime = Math.min(500 * Math.pow(2, attempt), 10000);
      console.warn(`Rate limited on API call, waiting ${waitTime}ms...`);
      await sleep(waitTime);
      continue;
    }

    return response;
  }

  throw new Error(`API request failed after ${maxRetries} retries due to rate limiting`);
}

/**
 * Create a test user via API
 */
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export async function createTestUser(token: string, userData: any): Promise<number> {
  const response = await authenticatedFetch(
    `${API_BASE}/admin/users`,
    {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(userData),
    },
    token
  );

  if (!response.ok) {
    const text = await response.text();
    throw new Error(`Failed to create test user: ${response.status} ${response.statusText} - ${text}`);
  }

  const result = (await response.json()) as ApiResponse<{ id: number }>;
  return result.data.id;
}

/**
 * Delete a test user via API
 */
export async function deleteTestUser(token: string, userId: number): Promise<void> {
  const response = await authenticatedFetch(
    `${API_BASE}/admin/users/${userId}`,
    { method: 'DELETE' },
    token
  );

  if (!response.ok) {
    console.warn(`Failed to delete test user ${userId}: ${response.statusText}`);
  }
}

/**
 * Create a test player via API
 */
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export async function createTestPlayer(token: string, playerData: any): Promise<number> {
  const response = await authenticatedFetch(
    `${API_BASE}/admin/players`,
    {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(playerData),
    },
    token
  );

  if (!response.ok) {
    const text = await response.text();
    throw new Error(`Failed to create test player: ${response.status} ${response.statusText} - ${text}`);
  }

  const result = (await response.json()) as ApiResponse<{ id: number }>;
  return result.data.id;
}

/**
 * Delete a test player via API
 */
export async function deleteTestPlayer(token: string, playerId: number): Promise<void> {
  const response = await authenticatedFetch(
    `${API_BASE}/admin/players/${playerId}`,
    { method: 'DELETE' },
    token
  );

  if (!response.ok) {
    console.warn(`Failed to delete test player ${playerId}: ${response.statusText}`);
  }
}

/**
 * Get a list of users for testing
 */
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export async function getUsers(token: string, params?: any): Promise<any[]> {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const queryParams = new URLSearchParams(params as any).toString();
  const response = await authenticatedFetch(
    `${API_BASE}/admin/users?${queryParams}`,
    { method: 'GET' },
    token
  );

  if (!response.ok) {
    throw new Error(`Failed to get users: ${response.statusText}`);
  }

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const result = (await response.json()) as ApiResponse<any>;
  // Handle both array and paginated response formats
  if (Array.isArray(result.data)) {
    return result.data;
  }
  // If data has items property (paginated), return items
  if (result.data?.items && Array.isArray(result.data.items)) {
    return result.data.items;
  }
  return [];
}

/**
 * Health check for the backend API
 */
export async function healthCheck(): Promise<boolean> {
  try {
    const response = await fetch(`${API_BASE.replace('/api/v1', '')}/api/v1/healthz`);
    return response.ok;
  } catch {
    return false;
  }
}

/**
 * Cleanup utility to delete multiple test users
 */
export async function cleanupTestUsers(token: string, userIds: number[]): Promise<void> {
  await Promise.allSettled(userIds.map((id) => deleteTestUser(token, id)));
}

/**
 * Cleanup utility to delete multiple test players
 */
export async function cleanupTestPlayers(token: string, playerIds: number[]): Promise<void> {
  await Promise.allSettled(playerIds.map((id) => deleteTestPlayer(token, id)));
}

// ==================== Order API Helpers ====================

export interface CreateOrderData {
  userId: number;
  playerId?: number;
  itemId: number;
  gameId: number;
  title?: string;
  description?: string;
  totalPriceCents: number;
  currency?: string;
  scheduledStart?: string;
  scheduledEnd?: string;
}

export interface OrderResponse {
  id: number;
  orderNo: string;
  status: string;
  totalPriceCents: number;
}

/**
 * Create a test order via API
 */
export async function createTestOrder(
  token: string,
  orderData: CreateOrderData
): Promise<OrderResponse> {
  const response = await authenticatedFetch(
    `${API_BASE}/admin/orders`,
    {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        user_id: orderData.userId,
        player_id: orderData.playerId,
        item_id: orderData.itemId,
        game_id: orderData.gameId,
        title: orderData.title || 'E2E测试订单',
        description: orderData.description || '',
        total_price_cents: orderData.totalPriceCents,
        currency: orderData.currency || 'CNY',
        scheduled_start: orderData.scheduledStart,
        scheduled_end: orderData.scheduledEnd,
      }),
    },
    token
  );

  if (!response.ok) {
    const text = await response.text();
    throw new Error(`Failed to create test order: ${response.status} ${response.statusText} - ${text}`);
  }

  const result = (await response.json()) as ApiResponse<OrderResponse>;
  return result.data;
}

/**
 * Get orders list via API
 */
export async function getOrders(
  token: string,
  params?: { status?: string; page?: number; page_size?: number }
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
): Promise<any[]> {
  const queryParams = new URLSearchParams(params as Record<string, string>).toString();
  const response = await authenticatedFetch(
    `${API_BASE}/admin/orders?${queryParams}`,
    { method: 'GET' },
    token
  );

  if (!response.ok) {
    throw new Error(`Failed to get orders: ${response.statusText}`);
  }

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const result = (await response.json()) as ApiResponse<any[]>;
  return result.data || [];
}

/**
 * Update order status via API
 */
export async function updateOrderStatus(
  token: string,
  orderId: number,
  status: string
): Promise<void> {
  const response = await authenticatedFetch(
    `${API_BASE}/admin/orders/${orderId}`,
    {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ status }),
    },
    token
  );

  if (!response.ok) {
    const text = await response.text();
    throw new Error(`Failed to update order status: ${response.status} - ${text}`);
  }
}

/**
 * Delete a test order via API
 */
export async function deleteTestOrder(token: string, orderId: number): Promise<void> {
  const response = await authenticatedFetch(
    `${API_BASE}/admin/orders/${orderId}`,
    { method: 'DELETE' },
    token
  );

  if (!response.ok) {
    console.warn(`Failed to delete test order ${orderId}: ${response.statusText}`);
  }
}

/**
 * Get service items for creating orders
 */
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export async function getServiceItems(token: string): Promise<any[]> {
  const response = await authenticatedFetch(
    `${API_BASE}/admin/service-items?page_size=100`,
    { method: 'GET' },
    token
  );

  if (!response.ok) {
    throw new Error(`Failed to get service items: ${response.statusText}`);
  }

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const result = (await response.json()) as ApiResponse<any>;
  // Handle both array and paginated response formats
  if (Array.isArray(result.data)) {
    return result.data;
  }
  if (result.data?.items && Array.isArray(result.data.items)) {
    return result.data.items;
  }
  return [];
}

/**
 * Get games for creating orders
 */
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export async function getGames(token: string): Promise<any[]> {
  const response = await authenticatedFetch(
    `${API_BASE}/admin/games?page_size=100`,
    { method: 'GET' },
    token
  );

  if (!response.ok) {
    throw new Error(`Failed to get games: ${response.statusText}`);
  }

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const result = (await response.json()) as ApiResponse<any>;
  // Handle both array and paginated response formats
  if (Array.isArray(result.data)) {
    return result.data;
  }
  if (result.data?.items && Array.isArray(result.data.items)) {
    return result.data.items;
  }
  return [];
}

/**
 * Get players for creating orders
 */
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export async function getPlayers(token: string): Promise<any[]> {
  const response = await authenticatedFetch(
    `${API_BASE}/admin/players?page_size=100`,
    { method: 'GET' },
    token
  );

  if (!response.ok) {
    throw new Error(`Failed to get players: ${response.statusText}`);
  }

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const result = (await response.json()) as ApiResponse<any>;
  // Handle both array and paginated response formats
  if (Array.isArray(result.data)) {
    return result.data;
  }
  if (result.data?.items && Array.isArray(result.data.items)) {
    return result.data.items;
  }
  return [];
}

/**
 * Cleanup utility to delete multiple test orders
 */
export async function cleanupTestOrders(token: string, orderIds: number[]): Promise<void> {
  await Promise.allSettled(orderIds.map((id) => deleteTestOrder(token, id)));
}

/**
 * Create E2E test orders for order management tests
 * Returns order IDs for cleanup
 */
export async function createE2ETestOrders(token: string): Promise<{
  pendingOrderIds: number[];
  completedOrderIds: number[];
}> {
  // Get required data with better error handling
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  let users: any[], players: any[], serviceItems: any[];

  try {
    [users, players, serviceItems] = await Promise.all([
      getUsers(token, { page_size: 10 }),
      getPlayers(token),
      getServiceItems(token),
    ]);
  } catch (error) {
    console.error('Failed to fetch required data for test orders:', error);
    return { pendingOrderIds: [], completedOrderIds: [] };
  }

  // Debug logging


  // Validate data exists
  if (!users || !Array.isArray(users) || users.length === 0) {
    console.warn('No users found for creating test orders, skipping order creation');
    return { pendingOrderIds: [], completedOrderIds: [] };
  }
  if (!players || !Array.isArray(players) || players.length === 0) {
    console.warn('No players found for creating test orders, skipping order creation');
    return { pendingOrderIds: [], completedOrderIds: [] };
  }
  if (!serviceItems || !Array.isArray(serviceItems) || serviceItems.length === 0) {
    console.warn('No service items found for creating test orders, skipping order creation');
    return { pendingOrderIds: [], completedOrderIds: [] };
  }

  // Find a service item that has a gameId (required for order creation)
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const serviceItem = serviceItems.find((item: any) => item.gameId || item.game_id) || serviceItems[0];
  // Use the gameId from serviceItem to ensure they match (backend validates this)
  const gameId = serviceItem.gameId || serviceItem.game_id;

  if (!gameId) {
    console.warn('Service item has no gameId, cannot create orders:', serviceItem);
    return { pendingOrderIds: [], completedOrderIds: [] };
  }

  // Find a user with role 'user' (not admin or player)
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const user = users.find((u: any) => u.role === 'user') || users[0];
  const player = players[0];

  // Validate objects have required properties
  if (!user?.id || !player?.id || !serviceItem?.id) {
    console.warn('Invalid data structure - missing id fields:', {
      userId: user?.id,
      playerId: player?.id,
      serviceItemId: serviceItem?.id,
    });
    return { pendingOrderIds: [], completedOrderIds: [] };
  }

  const timestamp = Date.now();
  const pendingOrderIds: number[] = [];
  const completedOrderIds: number[] = [];



  // Create pending orders for cancel tests
  for (let i = 0; i < 3; i++) {
    try {
      const order = await createTestOrder(token, {
        userId: user.id,
        playerId: player.id,
        itemId: serviceItem.id,
        gameId: gameId,
        title: `E2E测试-待取消订单-${timestamp}-${i + 1}`,
        description: 'E2E测试专用：用于测试订单取消功能',
        totalPriceCents: 8800 + i * 1000,
      });
      pendingOrderIds.push(order.id);

    } catch (error) {
      console.warn(`Failed to create pending order ${i + 1}:`, error);
    }
  }

  // Create completed orders for refund tests (note: API creates orders as pending, we'll use existing completed orders)
  // The admin API doesn't support creating orders with status directly
  // So we'll just create more pending orders and rely on existing seed data for completed orders
  for (let i = 0; i < 3; i++) {
    try {
      const order = await createTestOrder(token, {
        userId: user.id,
        playerId: player.id,
        itemId: serviceItem.id,
        gameId: gameId,
        title: `E2E测试-待退款订单-${timestamp}-${i + 1}`,
        description: 'E2E测试专用：用于测试订单退款功能',
        totalPriceCents: 9800 + i * 1000,
      });
      // Note: These will be created as pending, not completed
      // For refund tests, we'll rely on existing seed data
      completedOrderIds.push(order.id);

    } catch (error) {
      console.warn(`Failed to create order for refund test ${i + 1}:`, error);
    }
  }

  return { pendingOrderIds, completedOrderIds };
}
