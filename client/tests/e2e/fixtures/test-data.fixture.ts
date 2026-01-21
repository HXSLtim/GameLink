import { test as base } from '@playwright/test';

/**
 * Test data fixtures for E2E tests
 */
interface TestData {
  testUser: {
    username: string;
    password: string;
    email: string;
    phone?: string;
  };
  testPlayer: {
    username: string;
    nickname: string;
    gameName: string;
    price: number;
  };
}

export const test = base.extend<{ testData: TestData }>({
  testData: async (_params, use) => {
    // Test user credentials (can be overridden with env vars)
    const testData: TestData = {
      testUser: {
        username: process.env.TEST_USERNAME || 'testuser',
        password: process.env.TEST_PASSWORD || 'Test123456!',
        email: process.env.TEST_EMAIL || 'test@example.com',
        phone: process.env.TEST_PHONE || '13800138000',
      },
      testPlayer: {
        username: 'player1',
        nickname: 'ProGamer',
        gameName: 'Valorant',
        price: 50,
      },
    };
    // eslint-disable-next-line react-hooks/rules-of-hooks
    await use(testData);
  },
});

export { expect } from '@playwright/test';
