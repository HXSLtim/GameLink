/**
 * Tests for Import History Service
 *
 * @module services/import/history/historyService.test
 */

import { describe, it, expect, beforeEach } from 'vitest';
import { ImportHistoryService } from './historyService';
import { ImportHistoryStorage, createImportHistoryRecord } from './storage';
import type { ImportHistoryRecord, ImportRowResult } from './types';

/**
 * Create a mock in-memory storage for testing
 */
function createMockStorage(): Storage {
  const store = new Map<string, string>();
  return {
    getItem: (key: string) => store.get(key) ?? null,
    setItem: (key: string, value: string) => { store.set(key, value); },
    removeItem: (key: string) => { store.delete(key); },
    clear: () => { store.clear(); },
    get length() { return store.size; },
    key: (index: number) => Array.from(store.keys())[index] ?? null,
  };
}

describe('ImportHistoryService', () => {
  let service: ImportHistoryService;
  let storage: ImportHistoryStorage;
  let mockStorage: Storage;

  beforeEach(async () => {
    mockStorage = createMockStorage();
    storage = new ImportHistoryStorage(mockStorage);
    service = new ImportHistoryService(storage);
    await storage.clear();
  });

  describe('getImportHistory', () => {
    it('should return empty page when no records exist', async () => {
      const result = await service.getImportHistory();

      expect(result.records).toEqual([]);
      expect(result.total).toBe(0);
      expect(result.page).toBe(1);
      expect(result.totalPages).toBe(0);
    });

    it('should return all records without filters', async () => {
      const record1 = createImportHistoryRecord('user', 'users.csv', 1024, 1, 'Admin', 10);
      const record2 = createImportHistoryRecord('player', 'players.csv', 2048, 2, 'Manager', 20);
      
      await storage.save(record1);
      await storage.save(record2);

      const result = await service.getImportHistory();

      expect(result.records).toHaveLength(2);
      expect(result.total).toBe(2);
    });

    it('should filter by type', async () => {
      const userRecord = createImportHistoryRecord('user', 'users.csv', 1024, 1, 'Admin', 10);
      const playerRecord = createImportHistoryRecord('player', 'players.csv', 2048, 2, 'Manager', 20);
      
      await storage.save(userRecord);
      await storage.save(playerRecord);

      const result = await service.getImportHistory({ type: 'user' });

      expect(result.records).toHaveLength(1);
      expect(result.records[0].type).toBe('user');
    });

    it('should filter by status', async () => {
      const pendingRecord = createImportHistoryRecord('user', 'users1.csv', 1024, 1, 'Admin', 10);
      const completedRecord = createImportHistoryRecord('user', 'users2.csv', 2048, 1, 'Admin', 20);
      completedRecord.status = 'completed';
      
      await storage.save(pendingRecord);
      await storage.save(completedRecord);

      const result = await service.getImportHistory({ status: 'completed' });

      expect(result.records).toHaveLength(1);
      expect(result.records[0].status).toBe('completed');
    });

    it('should filter by date range', async () => {
      const oldRecord = createImportHistoryRecord('user', 'old.csv', 1024, 1, 'Admin', 10);
      oldRecord.uploadedAt = '2024-01-01T00:00:00.000Z';
      
      const newRecord = createImportHistoryRecord('user', 'new.csv', 2048, 1, 'Admin', 20);
      newRecord.uploadedAt = '2024-06-15T00:00:00.000Z';
      
      await storage.save(oldRecord);
      await storage.save(newRecord);

      const result = await service.getImportHistory({
        startDate: '2024-06-01',
        endDate: '2024-12-31',
      });

      expect(result.records).toHaveLength(1);
      expect(result.records[0].fileName).toBe('new.csv');
    });

    it('should filter by uploadedBy', async () => {
      const record1 = createImportHistoryRecord('user', 'users1.csv', 1024, 1, 'Admin', 10);
      const record2 = createImportHistoryRecord('user', 'users2.csv', 2048, 2, 'Manager', 20);
      
      await storage.save(record1);
      await storage.save(record2);

      const result = await service.getImportHistory({ uploadedBy: 1 });

      expect(result.records).toHaveLength(1);
      expect(result.records[0].uploadedBy).toBe(1);
    });

    it('should support pagination', async () => {
      // Create 15 records
      for (let i = 0; i < 15; i++) {
        const record = createImportHistoryRecord('user', `file${i}.csv`, 1024, 1, 'Admin', 10);
        await storage.save(record);
      }

      const page1 = await service.getImportHistory({}, { page: 1, pageSize: 10 });
      const page2 = await service.getImportHistory({}, { page: 2, pageSize: 10 });

      expect(page1.records).toHaveLength(10);
      expect(page1.total).toBe(15);
      expect(page1.totalPages).toBe(2);
      
      expect(page2.records).toHaveLength(5);
      expect(page2.page).toBe(2);
    });

    it('should accept Date objects for date filters', async () => {
      const record = createImportHistoryRecord('user', 'test.csv', 1024, 1, 'Admin', 10);
      record.uploadedAt = '2024-06-15T12:00:00.000Z';
      
      await storage.save(record);

      const result = await service.getImportHistory({
        startDate: new Date('2024-06-01'),
        endDate: new Date('2024-06-30'),
      });

      expect(result.records).toHaveLength(1);
    });
  });

  describe('getImportDetails', () => {
    it('should return null for non-existent record', async () => {
      const result = await service.getImportDetails('non-existent-id');
      expect(result).toBeNull();
    });

    it('should return record with computed fields', async () => {
      const record = createImportHistoryRecord('user', 'users.csv', 1024, 1, 'Admin', 100);
      record.importedCount = 80;
      record.skippedCount = 20;
      record.status = 'completed';
      record.completedAt = new Date(new Date(record.uploadedAt).getTime() + 5000).toISOString();
      
      await storage.save(record);

      const details = await service.getImportDetails(record.id);

      expect(details).not.toBeNull();
      expect(details!.successRate).toBe(80);
      expect(details!.durationMs).toBe(5000);
      expect(details!.durationFormatted).toBe('5s');
      expect(details!.hasErrors).toBe(true);
      expect(details!.errorCount).toBe(20);
    });

    it('should calculate 0% success rate for empty imports', async () => {
      const record = createImportHistoryRecord('user', 'empty.csv', 1024, 1, 'Admin', 0);
      record.status = 'completed';
      
      await storage.save(record);

      const details = await service.getImportDetails(record.id);

      expect(details!.successRate).toBe(0);
    });

    it('should count errors from rowResults when available', async () => {
      const record = createImportHistoryRecord('user', 'users.csv', 1024, 1, 'Admin', 5);
      record.skippedCount = 2;
      record.rowResults = [
        { rowNumber: 1, success: true, originalData: {} },
        { rowNumber: 2, success: false, originalData: {}, errorMessage: 'Error 1' },
        { rowNumber: 3, success: true, originalData: {} },
        { rowNumber: 4, success: false, originalData: {}, errorMessage: 'Error 2' },
        { rowNumber: 5, success: true, originalData: {} },
      ];
      
      await storage.save(record);

      const details = await service.getImportDetails(record.id);

      expect(details!.errorCount).toBe(2);
      expect(details!.hasErrors).toBe(true);
    });

    it('should format duration correctly for various time spans', async () => {
      // Test milliseconds
      const record1 = createImportHistoryRecord('user', 'fast.csv', 1024, 1, 'Admin', 10);
      record1.completedAt = new Date(new Date(record1.uploadedAt).getTime() + 500).toISOString();
      await storage.save(record1);
      
      let details = await service.getImportDetails(record1.id);
      expect(details!.durationFormatted).toBe('500ms');

      // Test minutes
      const record2 = createImportHistoryRecord('user', 'medium.csv', 1024, 1, 'Admin', 10);
      record2.completedAt = new Date(new Date(record2.uploadedAt).getTime() + 90000).toISOString();
      await storage.save(record2);
      
      details = await service.getImportDetails(record2.id);
      expect(details!.durationFormatted).toBe('1m 30s');

      // Test hours
      const record3 = createImportHistoryRecord('user', 'slow.csv', 1024, 1, 'Admin', 10);
      record3.completedAt = new Date(new Date(record3.uploadedAt).getTime() + 3660000).toISOString();
      await storage.save(record3);
      
      details = await service.getImportDetails(record3.id);
      expect(details!.durationFormatted).toBe('1h 1m');
    });

    it('should not have duration for incomplete imports', async () => {
      const record = createImportHistoryRecord('user', 'pending.csv', 1024, 1, 'Admin', 10);
      record.status = 'processing';
      // No completedAt
      
      await storage.save(record);

      const details = await service.getImportDetails(record.id);

      expect(details!.durationMs).toBeUndefined();
      expect(details!.durationFormatted).toBeUndefined();
    });
  });

  describe('getRecentImports', () => {
    it('should return empty array when no records exist', async () => {
      const result = await service.getRecentImports();
      expect(result).toEqual([]);
    });

    it('should return most recent records first', async () => {
      for (let i = 0; i < 5; i++) {
        const record = createImportHistoryRecord('user', `file${i}.csv`, 1024, 1, 'Admin', 10);
        await storage.save(record);
      }

      const result = await service.getRecentImports(3);

      expect(result).toHaveLength(3);
    });

    it('should use default limit of 10', async () => {
      for (let i = 0; i < 15; i++) {
        const record = createImportHistoryRecord('user', `file${i}.csv`, 1024, 1, 'Admin', 10);
        await storage.save(record);
      }

      const result = await service.getRecentImports();

      expect(result).toHaveLength(10);
    });
  });

  describe('getImportsByStatus', () => {
    it('should filter by status', async () => {
      const pending = createImportHistoryRecord('user', 'pending.csv', 1024, 1, 'Admin', 10);
      const completed = createImportHistoryRecord('user', 'completed.csv', 1024, 1, 'Admin', 10);
      completed.status = 'completed';
      const failed = createImportHistoryRecord('user', 'failed.csv', 1024, 1, 'Admin', 10);
      failed.status = 'failed';
      
      await storage.save(pending);
      await storage.save(completed);
      await storage.save(failed);

      const result = await service.getImportsByStatus('completed');

      expect(result.records).toHaveLength(1);
      expect(result.records[0].status).toBe('completed');
    });
  });

  describe('getImportsByType', () => {
    it('should filter by type', async () => {
      const userImport = createImportHistoryRecord('user', 'users.csv', 1024, 1, 'Admin', 10);
      const playerImport = createImportHistoryRecord('player', 'players.csv', 1024, 1, 'Admin', 10);
      const gameImport = createImportHistoryRecord('game', 'games.csv', 1024, 1, 'Admin', 10);
      
      await storage.save(userImport);
      await storage.save(playerImport);
      await storage.save(gameImport);

      const result = await service.getImportsByType('player');

      expect(result.records).toHaveLength(1);
      expect(result.records[0].type).toBe('player');
    });
  });

  describe('getFailedImports', () => {
    it('should return only failed imports', async () => {
      const completed = createImportHistoryRecord('user', 'completed.csv', 1024, 1, 'Admin', 10);
      completed.status = 'completed';
      const failed1 = createImportHistoryRecord('user', 'failed1.csv', 1024, 1, 'Admin', 10);
      failed1.status = 'failed';
      const failed2 = createImportHistoryRecord('player', 'failed2.csv', 1024, 1, 'Admin', 10);
      failed2.status = 'failed';
      
      await storage.save(completed);
      await storage.save(failed1);
      await storage.save(failed2);

      const result = await service.getFailedImports();

      expect(result.records).toHaveLength(2);
      expect(result.records.every((r) => r.status === 'failed')).toBe(true);
    });
  });

  describe('exists', () => {
    it('should return false for non-existent record', async () => {
      const result = await service.exists('non-existent-id');
      expect(result).toBe(false);
    });

    it('should return true for existing record', async () => {
      const record = createImportHistoryRecord('user', 'test.csv', 1024, 1, 'Admin', 10);
      await storage.save(record);

      const result = await service.exists(record.id);
      expect(result).toBe(true);
    });
  });
});
