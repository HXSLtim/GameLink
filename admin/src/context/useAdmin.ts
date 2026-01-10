import { useContext } from 'react';
import AdminContext from './AdminContext';

/**
 * 获取管理员上下文Hook
 *
 * @example
 * ```tsx
 * const { permissions, hasPermission, isSuperAdmin } = useAdmin();
 *
 * if (hasPermission('admin.games.create')) {
 *     // 显示创建按钮
 * }
 * ```
 */
export const useAdmin = () => useContext(AdminContext);
