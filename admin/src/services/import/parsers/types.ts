/**
 * Import Parser Types
 * Common types used by file parsers
 *
 * @module services/import/parsers/types
 */

/**
 * Supported file types for import
 */
export type SupportedFileType = 'xlsx' | 'xls' | 'csv';

/**
 * File extension to type mapping
 */
export const FILE_EXTENSION_MAP: Record<string, SupportedFileType> = {
  '.xlsx': 'xlsx',
  '.xls': 'xls',
  '.csv': 'csv',
};

/**
 * Supported MIME types for import files
 */
export const SUPPORTED_MIME_TYPES: Record<SupportedFileType, string[]> = {
  xlsx: ['application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'],
  xls: ['application/vnd.ms-excel'],
  csv: ['text/csv', 'application/csv', 'text/plain'],
};

/**
 * Maximum file size in bytes (10MB)
 */
export const MAX_FILE_SIZE_BYTES = 10 * 1024 * 1024;

/**
 * Maximum file size in MB for display
 */
export const MAX_FILE_SIZE_MB = 10;

/**
 * Result of parsing a file
 */
export interface ParseResult {
  /** Whether parsing was successful */
  success: boolean;
  /** Parsed headers (column names) */
  headers: string[];
  /** Parsed data rows */
  rows: Record<string, unknown>[];
  /** Error message if parsing failed */
  error?: string;
  /** Total number of rows parsed */
  totalRows: number;
}

/**
 * File validation result
 */
export interface FileValidationResult {
  /** Whether the file is valid */
  valid: boolean;
  /** Error message if invalid */
  error?: string;
  /** Detected file type */
  fileType?: SupportedFileType;
}

/**
 * Parser interface that all file parsers must implement
 */
export interface FileParser {
  /**
   * Parse a file and return the data
   * @param file - The file to parse
   * @returns ParseResult with headers and rows
   */
  parse(file: File): Promise<ParseResult>;

  /**
   * Check if this parser can handle the given file
   * @param file - The file to check
   * @returns true if this parser can handle the file
   */
  canParse(file: File): boolean;
}
