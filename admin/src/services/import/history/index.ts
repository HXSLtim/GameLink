/**
 * Import History Module
 * Provides import history storage and tracking functionality
 *
 * @module services/import/history
 */

// Types
export type {
  ImportStatus,
  ImportRowResult,
  ImportHistoryRecord,
  ImportHistoryQueryParams,
  ImportHistoryPage,
  IImportHistoryStorage,
} from './types';

// Storage
export {
  ImportHistoryStorage,
  importHistoryStorage,
  generateImportId,
  createImportHistoryRecord,
} from './storage';

// History Service
export type {
  ImportHistoryFilter,
  PaginationOptions,
  ImportHistoryDetails,
  IImportHistoryService,
} from './historyService';

export {
  ImportHistoryService,
  importHistoryService,
} from './historyService';

// Error Report
export type {
  ReportFormat,
  ErrorReportOptions,
  ErrorReportResult,
} from './errorReport';

export {
  generateErrorReport,
  downloadErrorReport,
  hasErrorDetails,
  getErrorSummary,
} from './errorReport';
