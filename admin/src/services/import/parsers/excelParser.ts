/**
 * Excel File Parser
 * Parses .xlsx and .xls files using the xlsx library
 *
 * @module services/import/parsers/excelParser
 */

import * as XLSX from 'xlsx';
import type { FileParser, ParseResult } from './types';

/**
 * Parser for Excel files (.xlsx, .xls)
 */
export class ExcelParser implements FileParser {
  private readonly supportedExtensions = ['.xlsx', '.xls'];
  private readonly supportedMimeTypes = [
    'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
    'application/vnd.ms-excel',
  ];

  /**
   * Check if this parser can handle the given file
   */
  canParse(file: File): boolean {
    const extension = this.getFileExtension(file.name);
    const mimeTypeMatch = this.supportedMimeTypes.includes(file.type);
    const extensionMatch = this.supportedExtensions.includes(extension);

    return mimeTypeMatch || extensionMatch;
  }

  /**
   * Parse an Excel file and return the data
   */
  async parse(file: File): Promise<ParseResult> {
    try {
      const arrayBuffer = await file.arrayBuffer();
      const workbook = XLSX.read(arrayBuffer, { type: 'array' });

      // Get the first sheet
      const firstSheetName = workbook.SheetNames[0];
      if (!firstSheetName) {
        return {
          success: false,
          headers: [],
          rows: [],
          totalRows: 0,
          error: 'Excel file contains no sheets',
        };
      }

      const worksheet = workbook.Sheets[firstSheetName];
      if (!worksheet) {
        return {
          success: false,
          headers: [],
          rows: [],
          totalRows: 0,
          error: 'Failed to read worksheet',
        };
      }

      // Convert to JSON with header row
      const jsonData = XLSX.utils.sheet_to_json(worksheet, {
        header: 1,
        defval: '',
        blankrows: false,
      }) as unknown[][];

      if (jsonData.length === 0) {
        return {
          success: false,
          headers: [],
          rows: [],
          totalRows: 0,
          error: 'Excel file is empty',
        };
      }

      // First row is headers
      const headers = (jsonData[0] as unknown[]).map((h) => String(h).trim());

      // Rest are data rows
      const rows: Record<string, unknown>[] = [];
      for (let i = 1; i < jsonData.length; i++) {
        const rowData = jsonData[i] as unknown[];
        const row: Record<string, unknown> = {};

        // Skip completely empty rows
        const hasData = rowData.some((cell) => cell !== '' && cell !== null && cell !== undefined);
        if (!hasData) continue;

        headers.forEach((header, index) => {
          const value = rowData[index];
          row[header] = this.normalizeValue(value);
        });

        rows.push(row);
      }

      return {
        success: true,
        headers,
        rows,
        totalRows: rows.length,
      };
    } catch (error) {
      return {
        success: false,
        headers: [],
        rows: [],
        totalRows: 0,
        error: `Failed to parse Excel file: ${error instanceof Error ? error.message : 'Unknown error'}`,
      };
    }
  }

  /**
   * Get file extension from filename
   */
  private getFileExtension(filename: string): string {
    const lastDot = filename.lastIndexOf('.');
    if (lastDot === -1) return '';
    return filename.slice(lastDot).toLowerCase();
  }

  /**
   * Normalize cell values to appropriate types
   */
  private normalizeValue(value: unknown): unknown {
    if (value === null || value === undefined || value === '') {
      return '';
    }

    // Handle numbers
    if (typeof value === 'number') {
      return value;
    }

    // Handle booleans
    if (typeof value === 'boolean') {
      return value;
    }

    // Handle strings
    const strValue = String(value).trim();

    // Try to parse boolean strings
    if (strValue.toLowerCase() === 'true' || strValue === '是') {
      return true;
    }
    if (strValue.toLowerCase() === 'false' || strValue === '否') {
      return false;
    }

    return strValue;
  }
}

/**
 * Singleton instance of ExcelParser
 */
export const excelParser = new ExcelParser();
