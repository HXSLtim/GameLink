/**
 * Property-Based Tests for ImportService
 */

import { describe, it, vi, beforeEach } from 'vitest';
import * as fc from 'fast-check';
import {
  ImportService,
  generateSecurePassword,
} from './importService';

vi.mock('../../api/admin', () => ({
  adminApi: {
    getUsers: vi.fn(),
    createUser: vi.fn(),
    getGames: vi.fn(),
    createGame: vi.fn(),
    updateGame: vi.fn(),
    createPlayer: vi.fn(),
  },
}));

describe('ImportService - Property Tests', () => {
  let _importService: ImportService;

  beforeEach(() => {
    _importService = new ImportService();
    vi.clearAllMocks();
  });

  describe('Property 18: Password Generation Security', () => {
    it('should always generate passwords with length >= 8', () => {
      fc.assert(
        fc.property(fc.integer({ min: 1, max: 100 }), () => {
          const password = generateSecurePassword();
          return password.length >= 8;
        }),
        { numRuns: 100 }
      );
    });

    it('should always contain uppercase letter', () => {
      fc.assert(
        fc.property(fc.integer({ min: 1, max: 100 }), () => {
          return /[A-Z]/.test(generateSecurePassword());
        }),
        { numRuns: 100 }
      );
    });
  });
});
