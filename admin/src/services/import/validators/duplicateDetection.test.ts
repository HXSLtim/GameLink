/**
 * Property-Based Tests for Import Duplicate Detection
 *
 * **Feature: admin-phase3-improvements, Property 17: Import Duplicate Detection**
 * **Validates: Requirements 6.2, 7.2, 8.2**
 *
 * Tests that the validators correctly detect and report all duplicates
 * (emails for users, user references for players, keys for games).
 */

import { describe, it, expect } from 'vitest';
import * as fc from 'fast-check';
import { findDuplicates, validateUserData } from './userDataValidator';
import { findDuplicateUserEmails, validatePlayerData } from './playerDataValidator';
import { findDuplicateGameKeys, validateGameData } from './gameDataValidator';

describe('Import Duplicate Detection - Property Tests', () => {
  /**
   * **Feature: admin-phase3-improvements, Property 17: Import Duplicate Detection**
   * **Validates: Requirements 6.2, 7.2, 8.2**
   *
   * For any import data containing duplicates (emails for users, user references
   * for players, keys for games), the validation SHALL detect and report all duplicates.
   */
  describe('Property 17: Import Duplicate Detection', () => {
    // Simple email generator
    const emailArb = fc.integer({ min: 1, max: 9999 }).map((n) => `user${n}@example.com`);

    // Simple phone generator (valid Chinese format)
    const phoneArb = fc.integer({ min: 0, max: 99999999 }).map((n) => `138${String(n).padStart(8, '0')}`);

    // Simple game key generator
    const gameKeyArb = fc.integer({ min: 1, max: 9999 }).map((n) => `game_${n}`);

    describe('User Email Duplicate Detection', () => {
      it('should detect all duplicate emails in import data', () => {
        fc.assert(
          fc.property(
            fc.array(emailArb, { minLength: 2, maxLength: 6 }),
            fc.integer({ min: 0, max: 2 }),
            (emails, duplicateCount) => {
              const uniqueEmails = [...new Set(emails)];
              if (uniqueEmails.length < 2) return true;

              const rows: Array<{ rowNumber: number; data: Record<string, unknown> }> = [];
              let rowNum = 2;

              for (const email of uniqueEmails) {
                rows.push({
                  rowNumber: rowNum++,
                  data: { name: 'Test User', email, phone: '13800138000' },
                });
              }

              const numDuplicates = Math.min(duplicateCount, uniqueEmails.length);
              for (let i = 0; i < numDuplicates; i++) {
                rows.push({
                  rowNumber: rowNum++,
                  data: { name: 'Duplicate User', email: uniqueEmails[i], phone: '13900139000' },
                });
              }

              const duplicates = findDuplicates(rows, 'email');

              for (let i = 0; i < numDuplicates; i++) {
                const email = uniqueEmails[i].toLowerCase();
                expect(duplicates.has(email)).toBe(true);
                expect(duplicates.get(email)!.length).toBeGreaterThanOrEqual(2);
              }

              return true;
            }
          ),
          { numRuns: 20 }
        );
      });

      it('should not report false positives for unique emails', () => {
        fc.assert(
          fc.property(
            fc.array(emailArb, { minLength: 1, maxLength: 8 }),
            (emails) => {
              const uniqueEmails = [...new Set(emails)];
              const rows = uniqueEmails.map((email, i) => ({
                rowNumber: i + 2,
                data: { name: `User ${i}`, email, phone: `138${String(i).padStart(8, '0')}` },
              }));

              const duplicates = findDuplicates(rows, 'email');
              expect(duplicates.size).toBe(0);
            }
          ),
          { numRuns: 20 }
        );
      });

      it('should be case-insensitive for email duplicate detection', () => {
        fc.assert(
          fc.property(emailArb, (email) => {
            const rows = [
              { rowNumber: 2, data: { name: 'User 1', email: email.toLowerCase(), phone: '13800138000' } },
              { rowNumber: 3, data: { name: 'User 2', email: email.toUpperCase(), phone: '13900139000' } },
            ];

            const duplicates = findDuplicates(rows, 'email');
            expect(duplicates.size).toBe(1);
            expect(duplicates.get(email.toLowerCase())).toContain(2);
            expect(duplicates.get(email.toLowerCase())).toContain(3);
          }),
          { numRuns: 10 }
        );
      });
    });

    describe('User Phone Duplicate Detection', () => {
      it('should detect all duplicate phones in import data', () => {
        fc.assert(
          fc.property(
            fc.array(phoneArb, { minLength: 2, maxLength: 6 }),
            fc.integer({ min: 0, max: 2 }),
            (phones, duplicateCount) => {
              const uniquePhones = [...new Set(phones)];
              if (uniquePhones.length < 2) return true;

              const rows: Array<{ rowNumber: number; data: Record<string, unknown> }> = [];
              let rowNum = 2;

              for (let i = 0; i < uniquePhones.length; i++) {
                rows.push({
                  rowNumber: rowNum++,
                  data: { name: `User ${i}`, email: `user${i}@example.com`, phone: uniquePhones[i] },
                });
              }

              const numDuplicates = Math.min(duplicateCount, uniquePhones.length);
              for (let i = 0; i < numDuplicates; i++) {
                rows.push({
                  rowNumber: rowNum++,
                  data: { name: `Duplicate ${i}`, email: `dup${i}@example.com`, phone: uniquePhones[i] },
                });
              }

              const duplicates = findDuplicates(rows, 'phone');

              for (let i = 0; i < numDuplicates; i++) {
                expect(duplicates.has(uniquePhones[i].toLowerCase())).toBe(true);
              }

              return true;
            }
          ),
          { numRuns: 20 }
        );
      });
    });

    describe('Player User Email Duplicate Detection', () => {
      it('should detect all duplicate user emails in player import data', () => {
        fc.assert(
          fc.property(
            fc.array(emailArb, { minLength: 2, maxLength: 6 }),
            fc.integer({ min: 0, max: 2 }),
            (emails, duplicateCount) => {
              const uniqueEmails = [...new Set(emails)];
              if (uniqueEmails.length < 2) return true;

              const rows: Array<{ rowNumber: number; data: Record<string, unknown> }> = [];
              let rowNum = 2;

              for (const email of uniqueEmails) {
                rows.push({
                  rowNumber: rowNum++,
                  data: { userEmail: email, nickname: 'Player' },
                });
              }

              const numDuplicates = Math.min(duplicateCount, uniqueEmails.length);
              for (let i = 0; i < numDuplicates; i++) {
                rows.push({
                  rowNumber: rowNum++,
                  data: { userEmail: uniqueEmails[i], nickname: 'Duplicate Player' },
                });
              }

              const duplicates = findDuplicateUserEmails(rows);

              for (let i = 0; i < numDuplicates; i++) {
                expect(duplicates.has(uniqueEmails[i].toLowerCase())).toBe(true);
              }

              return true;
            }
          ),
          { numRuns: 20 }
        );
      });
    });

    describe('Game Key Duplicate Detection', () => {
      it('should detect all duplicate game keys in import data', () => {
        fc.assert(
          fc.property(
            fc.array(gameKeyArb, { minLength: 2, maxLength: 6 }),
            fc.integer({ min: 0, max: 2 }),
            (keys, duplicateCount) => {
              const uniqueKeys = [...new Set(keys)];
              if (uniqueKeys.length < 2) return true;

              const rows: Array<{ rowNumber: number; data: Record<string, unknown> }> = [];
              let rowNum = 2;

              for (const key of uniqueKeys) {
                rows.push({
                  rowNumber: rowNum++,
                  data: { key, name: `Game ${key}` },
                });
              }

              const numDuplicates = Math.min(duplicateCount, uniqueKeys.length);
              for (let i = 0; i < numDuplicates; i++) {
                rows.push({
                  rowNumber: rowNum++,
                  data: { key: uniqueKeys[i], name: 'Duplicate Game' },
                });
              }

              const duplicates = findDuplicateGameKeys(rows);

              for (let i = 0; i < numDuplicates; i++) {
                expect(duplicates.has(uniqueKeys[i].toLowerCase())).toBe(true);
              }

              return true;
            }
          ),
          { numRuns: 20 }
        );
      });

      it('should be case-insensitive for game key duplicate detection', () => {
        fc.assert(
          fc.property(gameKeyArb, (key) => {
            const rows = [
              { rowNumber: 2, data: { key: key.toLowerCase(), name: 'Game 1' } },
              { rowNumber: 3, data: { key: key.toUpperCase(), name: 'Game 2' } },
            ];

            const duplicates = findDuplicateGameKeys(rows);
            expect(duplicates.size).toBe(1);
            expect(duplicates.get(key.toLowerCase())).toContain(2);
            expect(duplicates.get(key.toLowerCase())).toContain(3);
          }),
          { numRuns: 10 }
        );
      });
    });

    describe('Full Validation with Duplicate Detection', () => {
      it('should mark rows with duplicate emails as invalid in user validation', () => {
        fc.assert(
          fc.property(emailArb, phoneArb, (email, phone) => {
            const rows = [
              { rowNumber: 2, data: { name: 'User 1', email, phone } },
              { rowNumber: 3, data: { name: 'User 2', email, phone: `139${phone.slice(3)}` } },
            ];

            const result = validateUserData(rows, { checkInternalDuplicates: true });

            expect(result.valid).toBe(false);
            expect(result.invalidRows.length).toBe(2);
            expect(result.duplicateEmails.has(email.toLowerCase())).toBe(true);
          }),
          { numRuns: 10 }
        );
      });

      it('should mark rows with duplicate user emails as invalid in player validation', () => {
        fc.assert(
          fc.property(emailArb, (email) => {
            const rows = [
              { rowNumber: 2, data: { userEmail: email, nickname: 'Player 1' } },
              { rowNumber: 3, data: { userEmail: email, nickname: 'Player 2' } },
            ];

            const result = validatePlayerData(rows, {
              existingUserEmails: new Set([email.toLowerCase()]),
              checkInternalDuplicates: true,
            });

            expect(result.valid).toBe(false);
            expect(result.invalidRows.length).toBe(2);
            expect(result.duplicateUserEmails.has(email.toLowerCase())).toBe(true);
          }),
          { numRuns: 10 }
        );
      });

      it('should mark rows with duplicate keys as invalid in game validation', () => {
        fc.assert(
          fc.property(gameKeyArb, (key) => {
            const rows = [
              { rowNumber: 2, data: { key, name: 'Game 1' } },
              { rowNumber: 3, data: { key, name: 'Game 2' } },
            ];

            const result = validateGameData(rows, { checkInternalDuplicates: true });

            expect(result.valid).toBe(false);
            expect(result.invalidRows.length).toBe(2);
            expect(result.duplicateKeys.has(key.toLowerCase())).toBe(true);
          }),
          { numRuns: 10 }
        );
      });
    });
  });
});

describe('Import Duplicate Detection - Unit Tests', () => {
  describe('findDuplicates', () => {
    it('should find duplicate emails', () => {
      const rows = [
        { rowNumber: 2, data: { email: 'test@example.com' } },
        { rowNumber: 3, data: { email: 'other@example.com' } },
        { rowNumber: 4, data: { email: 'test@example.com' } },
      ];

      const duplicates = findDuplicates(rows, 'email');

      expect(duplicates.size).toBe(1);
      expect(duplicates.get('test@example.com')).toEqual([2, 4]);
    });

    it('should return empty map when no duplicates', () => {
      const rows = [
        { rowNumber: 2, data: { email: 'a@example.com' } },
        { rowNumber: 3, data: { email: 'b@example.com' } },
        { rowNumber: 4, data: { email: 'c@example.com' } },
      ];

      const duplicates = findDuplicates(rows, 'email');
      expect(duplicates.size).toBe(0);
    });

    it('should handle empty rows', () => {
      const duplicates = findDuplicates([], 'email');
      expect(duplicates.size).toBe(0);
    });

    it('should ignore empty values', () => {
      const rows = [
        { rowNumber: 2, data: { email: '' } },
        { rowNumber: 3, data: { email: '' } },
        { rowNumber: 4, data: { email: 'test@example.com' } },
      ];

      const duplicates = findDuplicates(rows, 'email');
      expect(duplicates.size).toBe(0);
    });
  });

  describe('findDuplicateUserEmails', () => {
    it('should find duplicate user emails in player data', () => {
      const rows = [
        { rowNumber: 2, data: { userEmail: 'player@example.com' } },
        { rowNumber: 3, data: { userEmail: 'other@example.com' } },
        { rowNumber: 4, data: { userEmail: 'player@example.com' } },
      ];

      const duplicates = findDuplicateUserEmails(rows);

      expect(duplicates.size).toBe(1);
      expect(duplicates.get('player@example.com')).toEqual([2, 4]);
    });
  });

  describe('findDuplicateGameKeys', () => {
    it('should find duplicate game keys', () => {
      const rows = [
        { rowNumber: 2, data: { key: 'game_one' } },
        { rowNumber: 3, data: { key: 'game_two' } },
        { rowNumber: 4, data: { key: 'game_one' } },
      ];

      const duplicates = findDuplicateGameKeys(rows);

      expect(duplicates.size).toBe(1);
      expect(duplicates.get('game_one')).toEqual([2, 4]);
    });

    it('should be case-insensitive', () => {
      const rows = [
        { rowNumber: 2, data: { key: 'Game_One' } },
        { rowNumber: 3, data: { key: 'GAME_ONE' } },
      ];

      const duplicates = findDuplicateGameKeys(rows);

      expect(duplicates.size).toBe(1);
      expect(duplicates.get('game_one')).toEqual([2, 3]);
    });
  });

  describe('validateUserData with duplicates', () => {
    it('should detect and report duplicate emails', () => {
      const rows = [
        { rowNumber: 2, data: { name: 'User 1', email: 'test@example.com', phone: '13800138000' } },
        { rowNumber: 3, data: { name: 'User 2', email: 'test@example.com', phone: '13900139000' } },
      ];

      const result = validateUserData(rows);

      expect(result.valid).toBe(false);
      expect(result.duplicateEmails.has('test@example.com')).toBe(true);
      expect(result.invalidRows.length).toBe(2);
    });

    it('should detect and report duplicate phones', () => {
      const rows = [
        { rowNumber: 2, data: { name: 'User 1', email: 'a@example.com', phone: '13800138000' } },
        { rowNumber: 3, data: { name: 'User 2', email: 'b@example.com', phone: '13800138000' } },
      ];

      const result = validateUserData(rows);

      expect(result.valid).toBe(false);
      expect(result.duplicatePhones.has('13800138000')).toBe(true);
    });
  });

  describe('validatePlayerData with duplicates', () => {
    it('should detect and report duplicate user emails', () => {
      const rows = [
        { rowNumber: 2, data: { userEmail: 'player@example.com', nickname: 'Player 1' } },
        { rowNumber: 3, data: { userEmail: 'player@example.com', nickname: 'Player 2' } },
      ];

      const result = validatePlayerData(rows, {
        existingUserEmails: new Set(['player@example.com']),
      });

      expect(result.valid).toBe(false);
      expect(result.duplicateUserEmails.has('player@example.com')).toBe(true);
    });
  });

  describe('validateGameData with duplicates', () => {
    it('should detect and report duplicate keys', () => {
      const rows = [
        { rowNumber: 2, data: { key: 'game_one', name: 'Game 1' } },
        { rowNumber: 3, data: { key: 'game_one', name: 'Game 2' } },
      ];

      const result = validateGameData(rows);

      expect(result.valid).toBe(false);
      expect(result.duplicateKeys.has('game_one')).toBe(true);
    });
  });
});
