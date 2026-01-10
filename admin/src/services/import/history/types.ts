/**
 * Import History Types
 * Defines types for import history storage and tracking
 *
 * @module services/import/history/types
 */

import type { ImportType } from '../templates/types';

/**
 * Import status enum
 */
export type ImportStatus = 'pending' | 'processing' | 'completed' | 'failed' | 'partial';

/**
 * Row-level import result
 */
export interface ImportRowResult {
  /** Row number in the original file (1-based, excluding header) */
  rowNumber: number;
  /** Whether this row was successfully imported */
  success: boolean;
  /** Original data from the row */
  originalData: Record<string, unknown>;
  /** Error message if failed */
  errorMessage?: string;
  /** Field that caused the error */
  errorField?: string;
  /** ID of the created record (if successful) */
  createdRecordId?: number;
}

/**
 * Import history record metadata
 */
export interface ImportHistoryRecord {
  /** Unique identifier for this import */
  id: string;
  /** Type of import (user, player, game) */
  type: ImportType;
  /** Original file name */
  fileName: string;
  /** File size in bytes */
  fileSize: number;
  /** User ID who performed the import */
  uploadedBy: number;
  /** User name who performed the import */
  uploadedByName?: string;
  /** Timestamp when import started */
  uploadedAt: string;
  /** Timestamp when import completed */
  completedAt?: string;
  /** Total number of rows in the file (excluding header) */
  totalRows: number;
  /** Number of successfully imported rows */
  importedCount: number;
  /** Number of skipped/failed rows */
  skippedCount: number;
  /** Import status */
  status: ImportStatus;
  /** Summary error message if failed */
  errorSummary?: string;
  /** Row-by-row results (stored for failed imports) */
  rowResults?: ImportRowResult[];
}

/**
 * Query parameters for import history
 */
export interface ImportHistoryQueryParams {
  /** Filter by import type */
  type?: ImportType;
  /** Filter by status */
  status?: ImportStatus;
  /** Filter by date range start */
  startDate?: string;
  /** Filter by date range end */
  endDate?: string;
  /** Filter by user ID */
  uploadedBy?: number;
  /** Page number (1-based) */
  page?: number;
  /** Page size */
  pageSize?: number;
}

/**
 * Paginated import history result
 */
export interface ImportHistoryPage {
  /** List of import history records */
  records: ImportHistoryRecord[];
  /** Total number of records matching the query */
  total: number;
  /** Current page number */
  page: number;
  /** Page size */
  pageSize: number;
  /** Total number of pages */
  totalPages: number;
}

/**
 * Import history storage interface
 */
export interface IImportHistoryStorage {
  /**
   * Save a new import history record
   */
  save(record: ImportHistoryRecord): Promise<void>;

  /**
   * Update an existing import history record
   */
  update(id: string, updates: Partial<ImportHistoryRecord>): Promise<void>;

  /**
   * Get import history record by ID
   */
  getById(id: string): Promise<ImportHistoryRecord | null>;

  /**
   * Query import history with filtering and pagination
   */
  query(params: ImportHistoryQueryParams): Promise<ImportHistoryPage>;

  /**
   * Delete import history record by ID
   */
  delete(id: string): Promise<void>;

  /**
   * Clear all import history (for testing/maintenance)
   */
  clear(): Promise<void>;
}
