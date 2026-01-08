/**
 * Import History Service
 * Provides high-level API for querying and managing import history
 *
 * @module services/import/history/historyService
 */

import type {
  ImportHistoryRecord,
  ImportHistoryQueryParams,
  ImportHistoryPage,
  IImportHistoryStorage,
  ImportStatus,
} from './types';
import type { ImportType } from '../templates/types';
import { importHistoryStorage } from './storage';

/**
 * Filter options for import history queries
 */
export interface ImportHistoryFilter {
  /** Filter by import type */
  type?: ImportType;
  /** Filter by status */
  status?: ImportStatus;
  /** Filter by date range start (ISO string or Date) */
  startDate?: string | Date;
  /** Filter by date range end (ISO string or Date) */
  endDate?: string | Date;
  /** Filter by user ID who performed the import */
  uploadedBy?: number;
}

/**
 * Pagination options
 */
export interface PaginationOptions {
  /** Page number (1-based, default: 1) */
  page?: number;
  /** Page size (default: 10) */
  pageSize?: number;
}

/**
 * Import history details with computed fields
 */
export interface ImportHistoryDetails extends ImportHistoryRecord {
  /** Success rate as percentage (0-100) */
  successRate: number;
  /** Duration in milliseconds (if completed) */
  durationMs?: number;
  /** Human-readable duration string */
  durationFormatted?: string;
  /** Whether the import has any errors */
  hasErrors: boolean;
  /** Number of errors */
  errorCount: number;
}

/**
 * Import History Service Interface
 */
export interface IImportHistoryService {
  /**
   * Get paginated import history with filtering
   * @param filter - Filter options
   * @param pagination - Pagination options
   * @returns Paginated import history records
   */
  getImportHistory(
    filter?: ImportHistoryFilter,
    pagination?: PaginationOptions
  ): Promise<ImportHistoryPage>;

  /**
   * Get detailed information about a specific import
   * @param id - Import history record ID
   * @returns Import history details or null if not found
   */
  getImportDetails(id: string): Promise<ImportHistoryDetails | null>;

  /**
   * Get recent imports (shortcut for common query)
   * @param limit - Maximum number of records to return (default: 10)
   * @returns Recent import history records
   */
  getRecentImports(limit?: number): Promise<ImportHistoryRecord[]>;

  /**
   * Get imports by status
   * @param status - Import status to filter by
   * @param pagination - Pagination options
   * @returns Paginated import history records
   */
  getImportsByStatus(
    status: ImportStatus,
    pagination?: PaginationOptions
  ): Promise<ImportHistoryPage>;

  /**
   * Get imports by type
   * @param type - Import type to filter by
   * @param pagination - Pagination options
   * @returns Paginated import history records
   */
  getImportsByType(
    type: ImportType,
    pagination?: PaginationOptions
  ): Promise<ImportHistoryPage>;

  /**
   * Get failed imports that may need attention
   * @param pagination - Pagination options
   * @returns Paginated failed import records
   */
  getFailedImports(pagination?: PaginationOptions): Promise<ImportHistoryPage>;

  /**
   * Check if an import exists
   * @param id - Import history record ID
   * @returns True if the import exists
   */
  exists(id: string): Promise<boolean>;
}

/**
 * Format duration in milliseconds to human-readable string
 */
function formatDuration(ms: number): string {
  if (ms < 1000) {
    return `${ms}ms`;
  }
  
  const seconds = Math.floor(ms / 1000);
  if (seconds < 60) {
    return `${seconds}s`;
  }
  
  const minutes = Math.floor(seconds / 60);
  const remainingSeconds = seconds % 60;
  if (minutes < 60) {
    return remainingSeconds > 0 ? `${minutes}m ${remainingSeconds}s` : `${minutes}m`;
  }
  
  const hours = Math.floor(minutes / 60);
  const remainingMinutes = minutes % 60;
  return remainingMinutes > 0 ? `${hours}h ${remainingMinutes}m` : `${hours}h`;
}

/**
 * Calculate success rate as percentage
 */
function calculateSuccessRate(record: ImportHistoryRecord): number {
  if (record.totalRows === 0) {
    return 0;
  }
  return Math.round((record.importedCount / record.totalRows) * 100);
}

/**
 * Calculate duration between upload and completion
 */
function calculateDuration(record: ImportHistoryRecord): number | undefined {
  if (!record.completedAt) {
    return undefined;
  }
  
  const startTime = new Date(record.uploadedAt).getTime();
  const endTime = new Date(record.completedAt).getTime();
  
  return endTime - startTime;
}

/**
 * Count errors in row results
 */
function countErrors(record: ImportHistoryRecord): number {
  if (!record.rowResults) {
    return record.skippedCount;
  }
  
  return record.rowResults.filter((r) => !r.success).length;
}

/**
 * Convert filter dates to ISO strings
 */
function normalizeDate(date: string | Date | undefined): string | undefined {
  if (!date) {
    return undefined;
  }
  
  if (date instanceof Date) {
    return date.toISOString();
  }
  
  return date;
}

/**
 * Import History Service Implementation
 */
export class ImportHistoryService implements IImportHistoryService {
  private storage: IImportHistoryStorage;

  constructor(storage?: IImportHistoryStorage) {
    this.storage = storage ?? importHistoryStorage;
  }

  /**
   * Get paginated import history with filtering
   */
  async getImportHistory(
    filter?: ImportHistoryFilter,
    pagination?: PaginationOptions
  ): Promise<ImportHistoryPage> {
    const queryParams: ImportHistoryQueryParams = {
      type: filter?.type,
      status: filter?.status,
      startDate: normalizeDate(filter?.startDate),
      endDate: normalizeDate(filter?.endDate),
      uploadedBy: filter?.uploadedBy,
      page: pagination?.page ?? 1,
      pageSize: pagination?.pageSize ?? 10,
    };

    return this.storage.query(queryParams);
  }

  /**
   * Get detailed information about a specific import
   */
  async getImportDetails(id: string): Promise<ImportHistoryDetails | null> {
    const record = await this.storage.getById(id);
    
    if (!record) {
      return null;
    }

    const durationMs = calculateDuration(record);
    
    return {
      ...record,
      successRate: calculateSuccessRate(record),
      durationMs,
      durationFormatted: durationMs !== undefined ? formatDuration(durationMs) : undefined,
      hasErrors: record.skippedCount > 0 || (record.rowResults?.some((r) => !r.success) ?? false),
      errorCount: countErrors(record),
    };
  }

  /**
   * Get recent imports
   */
  async getRecentImports(limit: number = 10): Promise<ImportHistoryRecord[]> {
    const result = await this.storage.query({
      page: 1,
      pageSize: limit,
    });
    
    return result.records;
  }

  /**
   * Get imports by status
   */
  async getImportsByStatus(
    status: ImportStatus,
    pagination?: PaginationOptions
  ): Promise<ImportHistoryPage> {
    return this.getImportHistory({ status }, pagination);
  }

  /**
   * Get imports by type
   */
  async getImportsByType(
    type: ImportType,
    pagination?: PaginationOptions
  ): Promise<ImportHistoryPage> {
    return this.getImportHistory({ type }, pagination);
  }

  /**
   * Get failed imports
   */
  async getFailedImports(pagination?: PaginationOptions): Promise<ImportHistoryPage> {
    return this.getImportsByStatus('failed', pagination);
  }

  /**
   * Check if an import exists
   */
  async exists(id: string): Promise<boolean> {
    const record = await this.storage.getById(id);
    return record !== null;
  }
}

// Export singleton instance
export const importHistoryService = new ImportHistoryService();
