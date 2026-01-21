/**
 * API Helper Functions for E2E Tests
 */

const API_URL = process.env.API_URL || 'http://localhost:5000/api/v1';

/**
 * Get auth token for test user
 */
export async function getTestUserToken(): Promise<string> {
  const response = await fetch(`${API_URL}/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      username: process.env.TEST_USERNAME || 'testuser',
      password: process.env.TEST_PASSWORD || 'Test123456!',
    }),
  });

  if (!response.ok) {
    throw new Error(`Failed to get auth token: ${response.statusText}`);
  }

  const data = await response.json();
  return data.data.token;
}

/**
 * Create test user via API
 */
export async function createTestUser(userData: {
  username: string;
  password: string;
  email: string;
  phone?: string;
}): Promise<{ id: number; token: string }> {
  const response = await fetch(`${API_URL}/auth/register`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(userData),
  });

  if (!response.ok) {
    throw new Error(`Failed to create user: ${response.statusText}`);
  }

  const data = await response.json();
  return {
    id: data.data.user?.id,
    token: data.data.token,
  };
}

/**
 * Clean up test user
 */
export async function cleanupTestUser(token: string): Promise<void> {
  // This would need a delete endpoint, for now just log
  console.log('Cleanup: User would be deleted with token:', token.substring(0, 10) + '...');
}

/**
 * Mock API responses for testing
 */
export function mockAuthEndpoints(page: any) {
  // Mock login endpoint
  page.route('**/api/v1/auth/login', async route => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({
        success: true,
        data: {
          token: 'mock-jwt-token',
          user: {
            id: 1,
            username: 'testuser',
            email: 'test@example.com',
            role: 'user',
          },
        },
      }),
    });
  });

  // Mock user info endpoint
  page.route('**/api/v1/user/info', async route => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({
        success: true,
        data: {
          id: 1,
          username: 'testuser',
          email: 'test@example.com',
          role: 'user',
          balance: 10000, // 100 yuan in cents
        },
      }),
    });
  });
}
