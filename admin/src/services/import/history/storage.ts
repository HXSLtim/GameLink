/**
 * Import History Storage Implementation
 * Provides localStorage-based storage for import history records
 *
 * @module services/import/history/storage
 */

import type {
  ImportHistoryRecord,
  ImportHistoryQueryParams,
  ImportHistoryPage,
  IImportHistoryStorage,
} from './types';

/**
 * Storage key for import history in localStorage
 */
const STORAGE_KEY = 'import_history';

/**
 * Maximum number of history records to keep
 */
const MAX_HISTORY_RECORDS = 100;

/**
 * LocalStorage-based import history storage implementation
 *
 * Features:
 * - Persists import history to localStorage
 * - Supports filtering by type, status, date range, and user
 * - Supports pagination
 * - Automatically prunes old records when limit is reached
 */
export class ImportHistoryStorage implements IImportHistoryStorage {
  private storage: Storage;

  constructor(storage?: Storage) {
    this.storage = storage ?? (typeof localStorage !== 'undefined' ? localStorage : createMemoryStorage());
  }

  /**
   * Save a new import history record
   */
  async save(record: ImportHistoryRecord): Promise<void> {
    const records = this.loadRecords();
    
    // Add new record at the beginning (most recent first)
    records.unshift(record);
    
    // Prune old records if limit exceeded
    if (records.length > MAX_HISTORY_RECORDS) {
      records.splice(MAX_HISTORY_RECORDS);
    }
    
    this.saveRecords(records);
  }

  /**
   * Update an existing import history record
   */
  async update(id: string, updates: Partial<ImportHistoryRecord>): Promise<void> {
    const records = this.loadRecords();
    const index = records.findIndex((r) => r.id === id);
    
    if (index !== -1) {
      records[index] = { ...records[index], ...updates };
      this.saveRecords(records);
    }
  }

  /**
   * Get import history record by ID
   */
  async getById(id: string): Promise<ImportHistoryRecord | null> {
    const records = this.loadRecords();
    return records.find((r) => r.id === id) ?? null;
  }

  /**
   * Query import history with filtering and pagination
   */
  async query(params: ImportHistoryQueryParams): Promise<ImportHistoryPage> {
    const records = this.loadRecords();
    
    // Apply filters
    let filtered = records;
    
    if (params.type) {
      filtered = filtered.filter((r) => r.type === params.type);
    }
    
    if (params.status) {
      filtered = filtered.filter((r) => r.status === params.status);
    }
    
    if (params.uploadedBy !== undefined) {
      filtered = filtered.filter((r) => r.uploadedBy === params.uploadedBy);
    }
    
    if (params.startDate) {
      const startDate = new Date(params.startDate);
      filtered = filtered.filter((r) => new Date(r.uploadedAt) >= startDate);
    }
    
    if (params.endDate) {
      const endDate = new Date(params.endDate);
      // Set end date to end of day
      endDate.setHours(23, 59, 59, 999);
      filtered = filtered.filter((r) => new Date(r.uploadedAt) <= endDate);
    }
    
    // Calculate pagination
    const page = params.page ?? 1;
    const pageSize = params.pageSize ?? 10;
    const total = filtered.length;
    const totalPages = Math.ceil(total / pageSize);
    
    // Apply pagination
    const startIndex = (page - 1) * pageSize;
    const paginatedRecords = filtered.slice(startIndex, startIndex + pageSize);
    
    return {
      records: paginatedRecords,
      total,
      page,
      pageSize,
      totalPages,
    };
  }

  /**
   * Delete import history record by ID
   */
  async delete(id: string): Promise<void> {
    const records = this.loadRecords();
    const filtered = records.filter((r) => r.id !== id);
    this.saveRecords(filtered);
  }

  /**
   * Clear all import history
   */
  async clear(): Promise<void> {
    this.saveRecords([]);
  }

  /**
   * Load records from storage
   */
  private loadRecords(): ImportHistoryRecord[] {
    try {
      const data = this.storage.getItem(STORAGE_KEY);
      if (data) {
        return JSON.parse(data) as ImportHistoryRecord[];
      }
    } catch {
      // Ignore parse errors, return empty array
    }
    return [];
  }

  /**
   * Save records to storage
   */
  private saveRecords(records: ImportHistoryRecord[]): void {
    this.storage.setItem(STORAGE_KEY, JSON.stringify(records));
  }
}

/**
 * Create an in-memory storage for testing or SSR environments
 */
function createMemoryStorage(): Storage {
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

/**
 * Generate a unique ID for import history records
 */
export function generateImportId(): string {
  // Use crypto.randomUUID if available, otherwise fallback
  if (typeof crypto !== 'undefined' && crypto.randomUUID) {
    return crypto.randomUUID();
  }
  
  // Fallback: timestamp + random string
  const timestamp = Date.now().toString(36);
  const random = Math.random().toString(36).substring(2, 10);
  return `${timestamp}-${random}`;
}

/**
 * Create a new import history record with initial values
 */
export function createImportHistoryRecord(
  type: ImportHistoryRecord['type'],
  fileName: string,
  fileSize: number,
  uploadedBy: number,
  uploadedByName?: string,
  totalRows: number = 0
): ImportHistoryRecord {
  return {
    id: generateImportId(),
    type,
    fileName,
    fileSize,
    uploadedBy,
    uploadedByName,
    uploadedAt: new Date().toISOString(),
    totalRows,
    importedCount: 0,
    skippedCount: 0,
    status: 'pending',
  };
}

// Export singleton instance
export const importHistoryStorage = new ImportHistoryStorage();
