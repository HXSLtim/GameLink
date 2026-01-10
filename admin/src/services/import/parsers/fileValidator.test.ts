/**
 * Property-Based Tests for File Validator
 *
 * **Feature: admin-phase3-improvements, Property 13: File Format Validation**
 * **Validates: Requirements 5.1**
 *
 * Tests that the file validator correctly accepts only .xlsx, .xls, and .csv files
 * up to 10MB in size, rejecting others with appropriate error messages.
 */

import { describe, it, expect } from 'vitest';
import * as fc from 'fast-check';
import {
  validateFile,
  getSupportedExtensions,
  getAcceptAttribute,
  isSupportedFileType,
  getFileType,
} from './fileValidator';
import { MAX_FILE_SIZE_BYTES, MAX_FILE_SIZE_MB } from './types';

/**
 * Create a mock File object for testing
 */
function createMockFile(
  name: string,
  size: number,
  type: string = 'application/octet-stream'
): File {
  // Create a blob with the specified size
  const content = new Uint8Array(Math.min(size, 1024)); // Limit actual content for performance
  const blob = new Blob([content], { type });

  // Create a File-like object with overridden size
  const file = new File([blob], name, { type });

  // Override size property for testing
  Object.defineProperty(file, 'size', {
    value: size,
    writable: false,
  });

  return file;
}

describe('File Validator - Property Tests', () => {
  /**
   * **Feature: admin-phase3-improvements, Property 13: File Format Validation**
   * **Validates: Requirements 5.1**
   *
   * For any import file upload, the service SHALL accept only .xlsx, .xls, and .csv
   * files up to 10MB, rejecting others with appropriate error messages.
   */
  describe('Property 13: File Format Validation', () => {
    it('should accept valid Excel files (.xlsx) within size limit', () => {
      // Generate valid file sizes (1 byte to 10MB)
      const validSizeArb = fc.integer({ min: 1, max: MAX_FILE_SIZE_BYTES });
      const filenameArb = fc.string({ minLength: 1, maxLength: 50 }).map((name) => {
        // Ensure valid filename characters
        const safeName = name.replace(/[<>:"/\\|?*]/g, '_');
        return `${safeName || 'file'}.xlsx`;
      });

      fc.assert(
        fc.property(filenameArb, validSizeArb, (filename, size) => {
          const file = createMockFile(
            filename,
            size,
            'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
          );
          const result = validateFile(file);

          expect(result.valid).toBe(true);
          expect(result.fileType).toBe('xlsx');
          expect(result.error).toBeUndefined();
        }),
        { numRuns: 100 }
      );
    });

    it('should accept valid Excel files (.xls) within size limit', () => {
      const validSizeArb = fc.integer({ min: 1, max: MAX_FILE_SIZE_BYTES });
      const filenameArb = fc.string({ minLength: 1, maxLength: 50 }).map((name) => {
        const safeName = name.replace(/[<>:"/\\|?*]/g, '_');
        return `${safeName || 'file'}.xls`;
      });

      fc.assert(
        fc.property(filenameArb, validSizeArb, (filename, size) => {
          const file = createMockFile(filename, size, 'application/vnd.ms-excel');
          const result = validateFile(file);

          expect(result.valid).toBe(true);
          expect(result.fileType).toBe('xls');
          expect(result.error).toBeUndefined();
        }),
        { numRuns: 100 }
      );
    });

    it('should accept valid CSV files within size limit', () => {
      const validSizeArb = fc.integer({ min: 1, max: MAX_FILE_SIZE_BYTES });
      const filenameArb = fc.string({ minLength: 1, maxLength: 50 }).map((name) => {
        const safeName = name.replace(/[<>:"/\\|?*]/g, '_');
        return `${safeName || 'file'}.csv`;
      });

      fc.assert(
        fc.property(filenameArb, validSizeArb, (filename, size) => {
          const file = createMockFile(filename, size, 'text/csv');
          const result = validateFile(file);

          expect(result.valid).toBe(true);
          expect(result.fileType).toBe('csv');
          expect(result.error).toBeUndefined();
        }),
        { numRuns: 100 }
      );
    });

    it('should reject files exceeding 10MB size limit', () => {
      // Generate sizes larger than 10MB
      const oversizedArb = fc.integer({
        min: MAX_FILE_SIZE_BYTES + 1,
        max: MAX_FILE_SIZE_BYTES * 2,
      });
      const extensionArb = fc.constantFrom('.xlsx', '.xls', '.csv');

      fc.assert(
        fc.property(oversizedArb, extensionArb, (size, extension) => {
          const file = createMockFile(`test${extension}`, size);
          const result = validateFile(file);

          expect(result.valid).toBe(false);
          expect(result.error).toBeDefined();
          expect(result.error).toContain(`${MAX_FILE_SIZE_MB}MB`);
        }),
        { numRuns: 100 }
      );
    });

    it('should reject unsupported file extensions', () => {
      const unsupportedExtensions = [
        '.txt',
        '.pdf',
        '.doc',
        '.docx',
        '.json',
        '.xml',
        '.html',
        '.zip',
        '.rar',
        '.exe',
        '.png',
        '.jpg',
      ];
      const extensionArb = fc.constantFrom(...unsupportedExtensions);
      const validSizeArb = fc.integer({ min: 1, max: MAX_FILE_SIZE_BYTES });

      fc.assert(
        fc.property(extensionArb, validSizeArb, (extension, size) => {
          const file = createMockFile(`test${extension}`, size);
          const result = validateFile(file);

          expect(result.valid).toBe(false);
          expect(result.error).toBeDefined();
          expect(result.error).toContain('Unsupported file type');
        }),
        { numRuns: 100 }
      );
    });

    it('should reject empty files (size = 0)', () => {
      const extensionArb = fc.constantFrom('.xlsx', '.xls', '.csv');

      fc.assert(
        fc.property(extensionArb, (extension) => {
          const file = createMockFile(`test${extension}`, 0);
          const result = validateFile(file);

          expect(result.valid).toBe(false);
          expect(result.error).toBeDefined();
          expect(result.error).toContain('empty');
        }),
        { numRuns: 10 }
      );
    });

    it('should handle files with no extension', () => {
      const validSizeArb = fc.integer({ min: 1, max: MAX_FILE_SIZE_BYTES });
      const filenameArb = fc.string({ minLength: 1, maxLength: 50 }).filter((name) => {
        // Ensure no dots in filename
        return !name.includes('.');
      });

      fc.assert(
        fc.property(filenameArb, validSizeArb, (filename, size) => {
          const safeName = filename.replace(/[<>:"/\\|?*]/g, '_') || 'noextension';
          const file = createMockFile(safeName, size);
          const result = validateFile(file);

          expect(result.valid).toBe(false);
          expect(result.error).toBeDefined();
          expect(result.error).toContain('Unsupported file type');
        }),
        { numRuns: 50 }
      );
    });
  });
});

describe('File Validator - Unit Tests', () => {
  describe('validateFile', () => {
    it('should accept xlsx file with correct MIME type', () => {
      const file = createMockFile(
        'test.xlsx',
        1024,
        'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
      );
      const result = validateFile(file);

      expect(result.valid).toBe(true);
      expect(result.fileType).toBe('xlsx');
    });

    it('should accept xlsx file with generic MIME type', () => {
      const file = createMockFile('test.xlsx', 1024, 'application/octet-stream');
      const result = validateFile(file);

      expect(result.valid).toBe(true);
      expect(result.fileType).toBe('xlsx');
    });

    it('should accept csv file with text/csv MIME type', () => {
      const file = createMockFile('data.csv', 1024, 'text/csv');
      const result = validateFile(file);

      expect(result.valid).toBe(true);
      expect(result.fileType).toBe('csv');
    });

    it('should accept csv file with text/plain MIME type', () => {
      const file = createMockFile('data.csv', 1024, 'text/plain');
      const result = validateFile(file);

      expect(result.valid).toBe(true);
      expect(result.fileType).toBe('csv');
    });

    it('should reject file exactly at size limit + 1 byte', () => {
      const file = createMockFile('test.xlsx', MAX_FILE_SIZE_BYTES + 1);
      const result = validateFile(file);

      expect(result.valid).toBe(false);
      expect(result.error).toContain('exceeds');
    });

    it('should accept file exactly at size limit', () => {
      const file = createMockFile(
        'test.xlsx',
        MAX_FILE_SIZE_BYTES,
        'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
      );
      const result = validateFile(file);

      expect(result.valid).toBe(true);
    });
  });

  describe('getSupportedExtensions', () => {
    it('should return all supported extensions', () => {
      const extensions = getSupportedExtensions();

      expect(extensions).toContain('.xlsx');
      expect(extensions).toContain('.xls');
      expect(extensions).toContain('.csv');
    });
  });

  describe('getAcceptAttribute', () => {
    it('should return valid accept attribute for file input', () => {
      const accept = getAcceptAttribute();

      expect(accept).toContain('.xlsx');
      expect(accept).toContain('.xls');
      expect(accept).toContain('.csv');
      expect(accept).toContain('application/vnd.openxmlformats-officedocument.spreadsheetml.sheet');
    });
  });

  describe('isSupportedFileType', () => {
    it('should return true for supported extensions', () => {
      expect(isSupportedFileType('test.xlsx')).toBe(true);
      expect(isSupportedFileType('test.xls')).toBe(true);
      expect(isSupportedFileType('test.csv')).toBe(true);
    });

    it('should return false for unsupported extensions', () => {
      expect(isSupportedFileType('test.txt')).toBe(false);
      expect(isSupportedFileType('test.pdf')).toBe(false);
      expect(isSupportedFileType('test')).toBe(false);
    });

    it('should be case insensitive', () => {
      expect(isSupportedFileType('test.XLSX')).toBe(true);
      expect(isSupportedFileType('test.CSV')).toBe(true);
    });
  });

  describe('getFileType', () => {
    it('should return correct file type for supported extensions', () => {
      expect(getFileType('test.xlsx')).toBe('xlsx');
      expect(getFileType('test.xls')).toBe('xls');
      expect(getFileType('test.csv')).toBe('csv');
    });

    it('should return undefined for unsupported extensions', () => {
      expect(getFileType('test.txt')).toBeUndefined();
      expect(getFileType('test')).toBeUndefined();
    });
  });
});
