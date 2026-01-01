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

/**
 * Get authentication token for admin user
 */
export async function getAdminToken(): Promise<string> {
  const response = await fetch(`${API_BASE}/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      username: process.env.TEST_ADMIN_USERNAME || 'admin',
      password: process.env.TEST_ADMIN_PASSWORD || 'admin123',
    }),
  });

  if (!response.ok) {
    throw new Error(`Failed to get admin token: ${response.statusText}`);
  }

  const result = (await response.json()) as ApiResponse<{ token: string }>;
  return result.data.token;
}

/**
 * Create a test user via API
 */
export async function createTestUser(token: string, userData: any): Promise<number> {
  const response = await fetch(`${API_BASE}/admin/users`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(userData),
  });

  if (!response.ok) {
    throw new Error(`Failed to create test user: ${response.statusText}`);
  }

  const result = (await response.json()) as ApiResponse<{ id: number }>;
  return result.data.id;
}

/**
 * Delete a test user via API
 */
export async function deleteTestUser(token: string, userId: number): Promise<void> {
  const response = await fetch(`${API_BASE}/admin/users/${userId}`, {
    method: 'DELETE',
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });

  if (!response.ok) {
    console.warn(`Failed to delete test user ${userId}: ${response.statusText}`);
  }
}

/**
 * Create a test player via API
 */
export async function createTestPlayer(token: string, playerData: any): Promise<number> {
  const response = await fetch(`${API_BASE}/admin/players`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(playerData),
  });

  if (!response.ok) {
    throw new Error(`Failed to create test player: ${response.statusText}`);
  }

  const result = (await response.json()) as ApiResponse<{ id: number }>;
  return result.data.id;
}

/**
 * Delete a test player via API
 */
export async function deleteTestPlayer(token: string, playerId: number): Promise<void> {
  const response = await fetch(`${API_BASE}/admin/players/${playerId}`, {
    method: 'DELETE',
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });

  if (!response.ok) {
    console.warn(`Failed to delete test player ${playerId}: ${response.statusText}`);
  }
}

/**
 * Get a list of users for testing
 */
export async function getUsers(token: string, params?: any): Promise<any[]> {
  const queryParams = new URLSearchParams(params as any).toString();
  const response = await fetch(`${API_BASE}/admin/users?${queryParams}`, {
    method: 'GET',
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });

  if (!response.ok) {
    throw new Error(`Failed to get users: ${response.statusText}`);
  }

  const result = (await response.json()) as ApiResponse<any[]>;
  return result.data;
}

/**
 * Health check for the backend API
 */
export async function healthCheck(): Promise<boolean> {
  try {
    const response = await fetch(`${API_BASE.replace('/api/v1', '')}/healthz`);
    return response.ok;
  } catch {
    return false;
  }
}

/**
 * Cleanup utility to delete multiple test users
 */
export async function cleanupTestUsers(token: string, userIds: number[]): Promise<void> {
  await Promise.allSettled(
    userIds.map((id) => deleteTestUser(token, id))
  );
}

/**
 * Cleanup utility to delete multiple test players
 */
export async function cleanupTestPlayers(token: string, playerIds: number[]): Promise<void> {
  await Promise.allSettled(
    playerIds.map((id) => deleteTestPlayer(token, id))
  );
}
