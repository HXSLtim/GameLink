/**
 * Import Parsers Module
 * Exports all file parsers and validation utilities
 *
 * @module services/import/parsers
 */

// Types
export * from './types';

// Parsers
export { ExcelParser, excelParser } from './excelParser';
export { CsvParser, csvParser } from './csvParser';

// Validation
export {
  validateFile,
  getSupportedExtensions,
  getAcceptAttribute,
  isSupportedFileType,
  getFileType,
} from './fileValidator';

import type { FileParser, ParseResult } from './types';
import { excelParser } from './excelParser';
import { csvParser } from './csvParser';
import { validateFile } from './fileValidator';

/**
 * Available parsers in order of preference
 */
const parsers: FileParser[] = [excelParser, csvParser];

/**
 * Parse a file using the appropriate parser
 * Automatically detects file type and uses the correct parser
 *
 * @param file - The file to parse
 * @returns ParseResult with headers and rows
 */
export async function parseFile(file: File): Promise<ParseResult> {
  // Validate file first
  const validation = validateFile(file);
  if (!validation.valid) {
    return {
      success: false,
      headers: [],
      rows: [],
      totalRows: 0,
      error: validation.error,
    };
  }

  // Find appropriate parser
  const parser = parsers.find((p) => p.canParse(file));
  if (!parser) {
    return {
      success: false,
      headers: [],
      rows: [],
      totalRows: 0,
      error: 'No parser available for this file type',
    };
  }

  // Parse the file
  return parser.parse(file);
}
