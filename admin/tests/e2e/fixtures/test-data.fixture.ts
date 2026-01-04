import { test as base } from '@playwright/test';

/**
 * Test data fixture for E2E tests
 * Provides reusable test data for admin operations
 */

export interface AdminTestData {
  adminUser: {
    username: string;
    password: string;
    email: string;
  };
  testUser: {
    name: string;
    email: string;
    phone: string;
    password: string;
    role: string;
  };
  testPlayer: {
    nickname: string;
    bio: string;
    rank: string;
    hourlyRateCents: number;
    mainGameId: number;
  };
  testOrder: {
    title: string;
    description: string;
    scheduledStart: string;
    scheduledEnd: string;
  };
}

// Generate unique test data with timestamp to avoid conflicts
const generateTestData = (): AdminTestData => ({
  adminUser: {
    username: process.env.TEST_ADMIN_USERNAME || 'admin@gamelink.com',
    password: process.env.TEST_ADMIN_PASSWORD || 'admin123456',
    email: process.env.TEST_ADMIN_EMAIL || 'admin@gamelink.com',
  },
  testUser: {
    name: `Test User ${Date.now()}`,
    email: `testuser${Date.now()}@example.com`,
    phone: `1${Math.floor(Math.random() * 10000000000)}`,
    password: 'Test123!@#',
    role: 'user',
  },
  testPlayer: {
    nickname: `Player ${Date.now()}`,
    bio: 'Professional gamer with 5+ years experience',
    rank: 'diamond',
    hourlyRateCents: 5000, // $50/hour
    mainGameId: 1,
  },
  testOrder: {
    title: `Test Order ${Date.now()}`,
    description: 'Test order for E2E testing',
    scheduledStart: new Date(Date.now() + 86400000).toISOString().slice(0, 16),
    scheduledEnd: new Date(Date.now() + 90000000).toISOString().slice(0, 16),
  },
});

// Extend base test with test data fixture
export const test = base.extend<{
  testData: AdminTestData;
}>({
  testData: async (
    // eslint-disable-next-line no-empty-pattern
    { },
    use
  ) => {
    // eslint-disable-next-line react-hooks/rules-of-hooks
    await use(generateTestData());
  },
});

export { expect } from '@playwright/test';
