/**
 * Import Error Report Generator
 * Generates downloadable error reports for failed imports
 *
 * @module services/import/history/errorReport
 */

import type { ImportHistoryRecord, ImportRowResult } from './types';
import type { ImportType } from '../templates/types';
import { getTemplate } from '../templates';

/**
 * Report format options
 */
export type ReportFormat = 'csv' | 'excel';

/**
 * Error report options
 */
export interface ErrorReportOptions {
  /** Report format (default: 'csv') */
  format?: ReportFormat;
  /** Include successful rows in report (default: false) */
  includeSuccessful?: boolean;
  /** Custom file name (default: generated from import info) */
  fileName?: string;
}

/**
 * Error report result
 */
export interface ErrorReportResult {
  /** Report blob for download */
  blob: Blob;
  /** Suggested file name */
  fileName: string;
  /** MIME type */
  mimeType: string;
  /** Number of rows in report */
  rowCount: number;
}

/**
 * Get column headers for a specific import type
 */
function getColumnHeaders(type: ImportType): string[] {
  const template = getTemplate(type);
  return template.columns.map((col) => col.labelZh);
}

/**
 * Get column keys for a specific import type
 */
function getColumnKeys(type: ImportType): string[] {
  const template = getTemplate(type);
  return template.columns.map((col) => col.key);
}

/**
 * Escape CSV cell value
 */
function escapeCsvCell(value: unknown): string {
  if (value === null || value === undefined) {
    return '';
  }
  
  const stringValue = String(value);
  
  // Escape quotes and wrap in quotes if contains comma, quote, or newline
  if (stringValue.includes(',') || stringValue.includes('"') || stringValue.includes('\n') || stringValue.includes('\r')) {
    return `"${stringValue.replace(/"/g, '""')}"`;
  }
  
  return stringValue;
}

/**
 * Convert row result to CSV row
 */
function rowResultToCsvRow(
  rowResult: ImportRowResult,
  columnKeys: string[]
): string[] {
  const row: string[] = [];
  
  // Add original data columns
  for (const key of columnKeys) {
    const value = rowResult.originalData[key];
    row.push(escapeCsvCell(value));
  }
  
  // Add status column
  row.push(rowResult.success ? '成功' : '失败');
  
  // Add error field column
  row.push(escapeCsvCell(rowResult.errorField));
  
  // Add error message column
  row.push(escapeCsvCell(rowResult.errorMessage));
  
  return row;
}

/**
 * Generate CSV content from import history record
 */
function generateCsvContent(
  record: ImportHistoryRecord,
  options: ErrorReportOptions
): string {
  const columnHeaders = getColumnHeaders(record.type);
  const columnKeys = getColumnKeys(record.type);
  
  // Build header row with additional status columns
  const headers = [...columnHeaders, '导入状态', '错误字段', '错误信息'];
  
  // Build data rows
  const rows: string[][] = [headers];
  
  if (record.rowResults) {
    for (const rowResult of record.rowResults) {
      // Skip successful rows if not included
      if (!options.includeSuccessful && rowResult.success) {
        continue;
      }
      
      rows.push(rowResultToCsvRow(rowResult, columnKeys));
    }
  }
  
  // Convert to CSV string
  const csvContent = rows
    .map((row) => row.join(','))
    .join('\n');
  
  return csvContent;
}

/**
 * Generate file name for error report
 */
function generateFileName(
  record: ImportHistoryRecord,
  format: ReportFormat
): string {
  const typeNames: Record<ImportType, string> = {
    user: '用户',
    player: '陪玩师',
    game: '游戏',
  };
  
  const typeName = typeNames[record.type] || record.type;
  const date = new Date(record.uploadedAt).toISOString().split('T')[0];
  const extension = format === 'excel' ? 'xlsx' : 'csv';
  
  return `${typeName}导入错误报告_${date}_${record.id.substring(0, 8)}.${extension}`;
}

/**
 * Generate error report for an import history record
 */
export function generateErrorReport(
  record: ImportHistoryRecord,
  options: ErrorReportOptions = {}
): ErrorReportResult {
  const format = options.format ?? 'csv';
  const includeSuccessful = options.includeSuccessful ?? false;
  
  // Generate CSV content
  const csvContent = generateCsvContent(record, { ...options, includeSuccessful });
  
  // Add BOM for Excel compatibility with Chinese characters
  const bom = '\uFEFF';
  const contentWithBom = bom + csvContent;
  
  // Create blob
  const mimeType = format === 'excel' 
    ? 'application/vnd.ms-excel' 
    : 'text/csv;charset=utf-8';
  
  const blob = new Blob([contentWithBom], { type: mimeType });
  
  // Calculate row count (excluding header)
  const rowCount = record.rowResults
    ? record.rowResults.filter((r) => includeSuccessful || !r.success).length
    : 0;
  
  // Generate file name
  const fileName = options.fileName ?? generateFileName(record, format);
  
  return {
    blob,
    fileName,
    mimeType,
    rowCount,
  };
}

/**
 * Download error report for an import history record
 * Creates a download link and triggers the download
 */
export function downloadErrorReport(
  record: ImportHistoryRecord,
  options: ErrorReportOptions = {}
): void {
  const report = generateErrorReport(record, options);
  
  // Create download link
  const url = URL.createObjectURL(report.blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = report.fileName;
  
  // Trigger download
  document.body.appendChild(link);
  link.click();
  
  // Cleanup
  document.body.removeChild(link);
  URL.revokeObjectURL(url);
}

/**
 * Check if an import record has error details available for report
 */
export function hasErrorDetails(record: ImportHistoryRecord): boolean {
  return (
    record.rowResults !== undefined &&
    record.rowResults.length > 0 &&
    record.rowResults.some((r) => !r.success)
  );
}

/**
 * Get error summary statistics from an import record
 */
export function getErrorSummary(record: ImportHistoryRecord): {
  totalRows: number;
  successCount: number;
  errorCount: number;
  errorsByField: Record<string, number>;
} {
  const errorsByField: Record<string, number> = {};
  let successCount = 0;
  let errorCount = 0;
  
  if (record.rowResults) {
    for (const row of record.rowResults) {
      if (row.success) {
        successCount++;
      } else {
        errorCount++;
        if (row.errorField) {
          errorsByField[row.errorField] = (errorsByField[row.errorField] || 0) + 1;
        }
      }
    }
  } else {
    successCount = record.importedCount;
    errorCount = record.skippedCount;
  }
  
  return {
    totalRows: record.totalRows,
    successCount,
    errorCount,
    errorsByField,
  };
}
