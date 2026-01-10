/**
 * Import Template Types
 * Defines the structure for import templates
 *
 * @module services/import/templates/types
 */

/**
 * Supported import types
 */
export type ImportType = 'user' | 'player' | 'game';

/**
 * Column data types for import templates
 */
export type ColumnType = 'string' | 'number' | 'boolean' | 'date' | 'email' | 'phone';

/**
 * Column definition for import templates
 */
export interface ImportColumn {
  /** Column key (used as object property name) */
  key: string;
  /** Display label for the column */
  label: string;
  /** Chinese label for the column (used in templates) */
  labelZh: string;
  /** Whether this column is required */
  required: boolean;
  /** Data type of the column */
  type: ColumnType;
  /** Default value if not provided */
  defaultValue?: unknown;
  /** Description/help text for the column */
  description?: string;
  /** Example value for the template */
  example?: string;
  /** Custom validation function */
  validation?: (value: unknown) => boolean;
  /** Allowed values (for enum-like columns) */
  allowedValues?: string[];
}

/**
 * Import template definition
 */
export interface ImportTemplate {
  /** Type of import */
  type: ImportType;
  /** Display name for the template */
  name: string;
  /** Chinese name for the template */
  nameZh: string;
  /** Description of what this template imports */
  description: string;
  /** Column definitions */
  columns: ImportColumn[];
  /** Example rows for the template */
  exampleRows?: Record<string, unknown>[];
}

/**
 * Get required columns from a template
 */
export function getRequiredColumns(template: ImportTemplate): ImportColumn[] {
  return template.columns.filter((col) => col.required);
}

/**
 * Get optional columns from a template
 */
export function getOptionalColumns(template: ImportTemplate): ImportColumn[] {
  return template.columns.filter((col) => !col.required);
}

/**
 * Get column keys from a template
 */
export function getColumnKeys(template: ImportTemplate): string[] {
  return template.columns.map((col) => col.key);
}

/**
 * Get column labels from a template (Chinese)
 */
export function getColumnLabels(template: ImportTemplate): string[] {
  return template.columns.map((col) => col.labelZh);
}
