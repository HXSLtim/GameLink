/**
 * User Import Template
 * Defines the structure for importing user data
 *
 * @module services/import/templates/userTemplate
 */

import type { ImportTemplate } from './types';

/**
 * Email validation regex (RFC 5322 simplified)
 */
const EMAIL_REGEX = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

/**
 * Phone validation regex (Chinese mobile format)
 */
const PHONE_REGEX = /^1[3-9]\d{9}$/;

/**
 * User import template definition
 *
 * Columns:
 * - name (required): User's display name
 * - email (required): User's email address
 * - phone (required): User's phone number
 * - role (optional): User role (user/player/admin), defaults to 'user'
 * - status (optional): User status (active/banned/suspended), defaults to 'active'
 */
export const userImportTemplate: ImportTemplate = {
  type: 'user',
  name: 'User Import Template',
  nameZh: '用户导入模板',
  description: 'Template for importing user data including name, email, phone, role, and status',
  columns: [
    {
      key: 'name',
      label: 'Name',
      labelZh: '姓名',
      required: true,
      type: 'string',
      description: "User's display name",
      example: '张三',
      validation: (value) => typeof value === 'string' && value.trim().length > 0,
    },
    {
      key: 'email',
      label: 'Email',
      labelZh: '邮箱',
      required: true,
      type: 'email',
      description: "User's email address (must be unique)",
      example: 'zhangsan@example.com',
      validation: (value) => typeof value === 'string' && EMAIL_REGEX.test(value),
    },
    {
      key: 'phone',
      label: 'Phone',
      labelZh: '手机号',
      required: true,
      type: 'phone',
      description: "User's phone number (Chinese mobile format)",
      example: '13800138000',
      validation: (value) => typeof value === 'string' && PHONE_REGEX.test(value),
    },
    {
      key: 'role',
      label: 'Role',
      labelZh: '角色',
      required: false,
      type: 'string',
      defaultValue: 'user',
      description: 'User role: user, player, or admin',
      example: 'user',
      allowedValues: ['user', 'player', 'admin'],
      validation: (value) =>
        value === undefined ||
        value === '' ||
        ['user', 'player', 'admin'].includes(String(value).toLowerCase()),
    },
    {
      key: 'status',
      label: 'Status',
      labelZh: '状态',
      required: false,
      type: 'string',
      defaultValue: 'active',
      description: 'User status: active, banned, or suspended',
      example: 'active',
      allowedValues: ['active', 'banned', 'suspended'],
      validation: (value) =>
        value === undefined ||
        value === '' ||
        ['active', 'banned', 'suspended'].includes(String(value).toLowerCase()),
    },
  ],
  exampleRows: [
    {
      name: '张三',
      email: 'zhangsan@example.com',
      phone: '13800138000',
      role: 'user',
      status: 'active',
    },
    {
      name: '李四',
      email: 'lisi@example.com',
      phone: '13900139000',
      role: 'player',
      status: 'active',
    },
    {
      name: '王五',
      email: 'wangwu@example.com',
      phone: '13700137000',
      role: 'admin',
      status: 'active',
    },
  ],
};

/**
 * Validate user email format
 */
export function validateUserEmail(email: string): boolean {
  return EMAIL_REGEX.test(email);
}

/**
 * Validate user phone format
 */
export function validateUserPhone(phone: string): boolean {
  return PHONE_REGEX.test(phone);
}

/**
 * Validate user role
 */
export function validateUserRole(role: string): boolean {
  return ['user', 'player', 'admin'].includes(role.toLowerCase());
}

/**
 * Validate user status
 */
export function validateUserStatus(status: string): boolean {
  return ['active', 'banned', 'suspended'].includes(status.toLowerCase());
}
