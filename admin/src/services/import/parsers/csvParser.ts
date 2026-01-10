/**
 * CSV File Parser
 * Parses .csv files with proper handling of quoted fields and various delimiters
 *
 * @module services/import/parsers/csvParser
 */

import type { FileParser, ParseResult } from './types';

/**
 * Parser for CSV files (.csv)
 */
export class CsvParser implements FileParser {
  private readonly supportedExtensions = ['.csv'];
  private readonly supportedMimeTypes = ['text/csv', 'application/csv', 'text/plain'];

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
   * Parse a CSV file and return the data
   */
  async parse(file: File): Promise<ParseResult> {
    try {
      const text = await file.text();

      if (!text.trim()) {
        return {
          success: false,
          headers: [],
          rows: [],
          totalRows: 0,
          error: 'CSV file is empty',
        };
      }

      // Detect delimiter (comma, semicolon, or tab)
      const delimiter = this.detectDelimiter(text);

      // Parse CSV content
      const lines = this.parseCSVLines(text, delimiter);

      if (lines.length === 0) {
        return {
          success: false,
          headers: [],
          rows: [],
          totalRows: 0,
          error: 'CSV file contains no data',
        };
      }

      // First line is headers
      const headers = lines[0].map((h) => h.trim());

      // Rest are data rows
      const rows: Record<string, unknown>[] = [];
      for (let i = 1; i < lines.length; i++) {
        const rowData = lines[i];

        // Skip completely empty rows
        const hasData = rowData.some((cell) => cell.trim() !== '');
        if (!hasData) continue;

        const row: Record<string, unknown> = {};
        headers.forEach((header, index) => {
          const value = rowData[index] ?? '';
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
        error: `Failed to parse CSV file: ${error instanceof Error ? error.message : 'Unknown error'}`,
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
   * Detect the delimiter used in the CSV file
   */
  private detectDelimiter(text: string): string {
    const firstLine = text.split(/\r?\n/)[0] || '';

    // Count occurrences of common delimiters
    const commaCount = (firstLine.match(/,/g) || []).length;
    const semicolonCount = (firstLine.match(/;/g) || []).length;
    const tabCount = (firstLine.match(/\t/g) || []).length;

    // Return the most common delimiter
    if (tabCount > commaCount && tabCount > semicolonCount) {
      return '\t';
    }
    if (semicolonCount > commaCount) {
      return ';';
    }
    return ',';
  }

  /**
   * Parse CSV content into lines and fields, handling quoted fields
   */
  private parseCSVLines(text: string, delimiter: string): string[][] {
    const lines: string[][] = [];
    let currentLine: string[] = [];
    let currentField = '';
    let inQuotes = false;

    for (let i = 0; i < text.length; i++) {
      const char = text[i];
      const nextChar = text[i + 1];

      if (inQuotes) {
        if (char === '"') {
          if (nextChar === '"') {
            // Escaped quote
            currentField += '"';
            i++; // Skip next quote
          } else {
            // End of quoted field
            inQuotes = false;
          }
        } else {
          currentField += char;
        }
      } else {
        if (char === '"') {
          // Start of quoted field
          inQuotes = true;
        } else if (char === delimiter) {
          // End of field
          currentLine.push(currentField);
          currentField = '';
        } else if (char === '\r') {
          // Carriage return - skip if followed by newline
          if (nextChar === '\n') {
            continue;
          }
          // Otherwise treat as newline
          currentLine.push(currentField);
          lines.push(currentLine);
          currentLine = [];
          currentField = '';
        } else if (char === '\n') {
          // End of line
          currentLine.push(currentField);
          lines.push(currentLine);
          currentLine = [];
          currentField = '';
        } else {
          currentField += char;
        }
      }
    }

    // Don't forget the last field and line
    if (currentField || currentLine.length > 0) {
      currentLine.push(currentField);
      lines.push(currentLine);
    }

    return lines;
  }

  /**
   * Normalize cell values to appropriate types
   */
  private normalizeValue(value: string): unknown {
    const trimmed = value.trim();

    if (trimmed === '') {
      return '';
    }

    // Try to parse boolean strings
    if (trimmed.toLowerCase() === 'true' || trimmed === '是') {
      return true;
    }
    if (trimmed.toLowerCase() === 'false' || trimmed === '否') {
      return false;
    }

    // Try to parse numbers
    if (/^-?\d+(\.\d+)?$/.test(trimmed)) {
      const num = parseFloat(trimmed);
      if (!isNaN(num)) {
        return num;
      }
    }

    return trimmed;
  }
}

/**
 * Singleton instance of CsvParser
 */
export const csvParser = new CsvParser();
