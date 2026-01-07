/**
 * Structure Validator
 * Validates the structure of imported data against templates
 *
 * @module services/import/validators/structureValidator
 */

import type { ImportTemplate, ImportColumn } from '../templates/types';

/**
 * Result of structure validation
 */
export interface StructureValidationResult {
  /** Whether the structure is valid */
  valid: boolean;
  /** List of missing required columns */
  missingColumns: string[];
  /** List of extra columns not in template */
  extraColumns: string[];
  /** List of matched columns */
  matchedColumns: string[];
  /** Error messages for display */
  errors: string[];
}

/**
 * Column mapping result
 */
export interface ColumnMapping {
  /** Template column key */
  templateKey: string;
  /** Actual header from file */
  fileHeader: string;
  /** Whether this is a required column */
  required: boolean;
}

/**
 * Validate the structure of imported headers against a template
 *
 * @param headers - Headers from the imported file
 * @param template - Import template to validate against
 * @returns StructureValidationResult with validation details
 */
export function validateStructure(
  headers: string[],
  template: ImportTemplate
): StructureValidationResult {
  const result: StructureValidationResult = {
    valid: true,
    missingColumns: [],
    extraColumns: [],
    matchedColumns: [],
    errors: [],
  };

  // Normalize headers for comparison
  const normalizedHeaders = headers.map((h) => normalizeHeader(h));
  const headerSet = new Set(normalizedHeaders);

  // Check for required columns
  for (const column of template.columns) {
    const normalizedKey = normalizeHeader(column.key);
    const normalizedLabel = normalizeHeader(column.label);
    const normalizedLabelZh = normalizeHeader(column.labelZh);

    // Check if column exists (by key, label, or Chinese label)
    const found =
      headerSet.has(normalizedKey) ||
      headerSet.has(normalizedLabel) ||
      headerSet.has(normalizedLabelZh);

    if (found) {
      result.matchedColumns.push(column.key);
    } else if (column.required) {
      result.missingColumns.push(column.labelZh);
      result.valid = false;
    }
  }

  // Check for extra columns
  const templateKeys = new Set(
    template.columns.flatMap((col) => [
      normalizeHeader(col.key),
      normalizeHeader(col.label),
      normalizeHeader(col.labelZh),
    ])
  );

  for (const header of headers) {
    const normalizedHeader = normalizeHeader(header);
    if (!templateKeys.has(normalizedHeader)) {
      result.extraColumns.push(header);
    }
  }

  // Generate error messages
  if (result.missingColumns.length > 0) {
    result.errors.push(`缺少必填列: ${result.missingColumns.join(', ')}`);
  }

  if (result.extraColumns.length > 0) {
    // Extra columns are warnings, not errors
    result.errors.push(`发现未知列 (将被忽略): ${result.extraColumns.join(', ')}`);
  }

  return result;
}

/**
 * Create a mapping from file headers to template columns
 *
 * @param headers - Headers from the imported file
 * @param template - Import template
 * @returns Map of file header index to template column key
 */
export function createColumnMapping(
  headers: string[],
  template: ImportTemplate
): Map<number, ColumnMapping> {
  const mapping = new Map<number, ColumnMapping>();

  for (let i = 0; i < headers.length; i++) {
    const header = headers[i];
    const normalizedHeader = normalizeHeader(header);

    // Find matching template column
    const column = findMatchingColumn(normalizedHeader, template.columns);
    if (column) {
      mapping.set(i, {
        templateKey: column.key,
        fileHeader: header,
        required: column.required,
      });
    }
  }

  return mapping;
}

/**
 * Find a template column that matches the given header
 */
function findMatchingColumn(
  normalizedHeader: string,
  columns: ImportColumn[]
): ImportColumn | undefined {
  return columns.find((col) => {
    const normalizedKey = normalizeHeader(col.key);
    const normalizedLabel = normalizeHeader(col.label);
    const normalizedLabelZh = normalizeHeader(col.labelZh);

    return (
      normalizedHeader === normalizedKey ||
      normalizedHeader === normalizedLabel ||
      normalizedHeader === normalizedLabelZh
    );
  });
}

/**
 * Normalize a header string for comparison
 * - Converts to lowercase
 * - Removes whitespace
 * - Removes special characters
 */
function normalizeHeader(header: string): string {
  return header
    .toLowerCase()
    .trim()
    .replace(/[\s_-]+/g, '')
    .replace(/[()（）]/g, '');
}

/**
 * Map a data row using the column mapping
 *
 * @param row - Raw data row from file
 * @param headers - File headers
 * @param mapping - Column mapping
 * @param template - Import template for default values
 * @returns Mapped data object with template keys
 */
export function mapRowToTemplate(
  row: Record<string, unknown>,
  headers: string[],
  mapping: Map<number, ColumnMapping>,
  template: ImportTemplate
): Record<string, unknown> {
  const mappedRow: Record<string, unknown> = {};

  // Apply default values first
  for (const column of template.columns) {
    if (column.defaultValue !== undefined) {
      mappedRow[column.key] = column.defaultValue;
    }
  }

  // Map values from row
  for (let i = 0; i < headers.length; i++) {
    const header = headers[i];
    const columnMapping = mapping.get(i);

    if (columnMapping) {
      const value = row[header];
      if (value !== undefined && value !== '') {
        mappedRow[columnMapping.templateKey] = value;
      }
    }
  }

  return mappedRow;
}

/**
 * Get required columns that are missing from headers
 */
export function getMissingRequiredColumns(
  headers: string[],
  template: ImportTemplate
): ImportColumn[] {
  const normalizedHeaders = new Set(headers.map((h) => normalizeHeader(h)));

  return template.columns.filter((col) => {
    if (!col.required) return false;

    const normalizedKey = normalizeHeader(col.key);
    const normalizedLabel = normalizeHeader(col.label);
    const normalizedLabelZh = normalizeHeader(col.labelZh);

    return (
      !normalizedHeaders.has(normalizedKey) &&
      !normalizedHeaders.has(normalizedLabel) &&
      !normalizedHeaders.has(normalizedLabelZh)
    );
  });
}

/**
 * Check if all required columns are present
 */
export function hasAllRequiredColumns(headers: string[], template: ImportTemplate): boolean {
  return getMissingRequiredColumns(headers, template).length === 0;
}
