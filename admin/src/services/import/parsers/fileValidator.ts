/**
 * File Validator
 * Validates file type and size before parsing
 *
 * @module services/import/parsers/fileValidator
 */

import {
  type FileValidationResult,
  type SupportedFileType,
  FILE_EXTENSION_MAP,
  MAX_FILE_SIZE_BYTES,
  MAX_FILE_SIZE_MB,
  SUPPORTED_MIME_TYPES,
} from './types';

/**
 * Validate a file for import
 * Checks file type and size constraints
 *
 * @param file - The file to validate
 * @returns FileValidationResult with validation status
 */
export function validateFile(file: File): FileValidationResult {
  // Check file size
  if (file.size > MAX_FILE_SIZE_BYTES) {
    return {
      valid: false,
      error: `File size exceeds maximum allowed size of ${MAX_FILE_SIZE_MB}MB. Current size: ${(file.size / (1024 * 1024)).toFixed(2)}MB`,
    };
  }

  // Check file size is not zero
  if (file.size === 0) {
    return {
      valid: false,
      error: 'File is empty',
    };
  }

  // Get file extension
  const extension = getFileExtension(file.name);
  const fileType = FILE_EXTENSION_MAP[extension];

  // Check if extension is supported
  if (!fileType) {
    const supportedExtensions = Object.keys(FILE_EXTENSION_MAP).join(', ');
    return {
      valid: false,
      error: `Unsupported file type. Supported formats: ${supportedExtensions}`,
    };
  }

  // Verify MIME type matches extension (if MIME type is provided)
  if (file.type && !isValidMimeType(file.type, fileType)) {
    // Allow if extension is correct but MIME type is generic
    const genericMimeTypes = ['application/octet-stream', ''];
    if (!genericMimeTypes.includes(file.type)) {
      return {
        valid: false,
        error: `File MIME type (${file.type}) does not match expected type for ${extension} files`,
      };
    }
  }

  return {
    valid: true,
    fileType,
  };
}

/**
 * Get file extension from filename
 */
function getFileExtension(filename: string): string {
  const lastDot = filename.lastIndexOf('.');
  if (lastDot === -1) return '';
  return filename.slice(lastDot).toLowerCase();
}

/**
 * Check if MIME type is valid for the given file type
 */
function isValidMimeType(mimeType: string, fileType: SupportedFileType): boolean {
  const validMimeTypes = SUPPORTED_MIME_TYPES[fileType];
  return validMimeTypes.includes(mimeType);
}

/**
 * Get supported file extensions as a string for display
 */
export function getSupportedExtensions(): string {
  return Object.keys(FILE_EXTENSION_MAP).join(', ');
}

/**
 * Get accept attribute value for file input
 */
export function getAcceptAttribute(): string {
  const extensions = Object.keys(FILE_EXTENSION_MAP);
  const mimeTypes = Object.values(SUPPORTED_MIME_TYPES).flat();
  return [...extensions, ...mimeTypes].join(',');
}

/**
 * Check if a file type is supported
 */
export function isSupportedFileType(filename: string): boolean {
  const extension = getFileExtension(filename);
  return extension in FILE_EXTENSION_MAP;
}

/**
 * Get the file type from filename
 */
export function getFileType(filename: string): SupportedFileType | undefined {
  const extension = getFileExtension(filename);
  return FILE_EXTENSION_MAP[extension];
}
