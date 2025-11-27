/**
 * 组件统一导出
 */

// 权限相关
export { PermissionGuard, PermissionButton, withPermission } from './PermissionGuard';
export type { PermissionGuardProps, PermissionButtonProps } from './PermissionGuard';

// 页面容器
export { default as PageContainer } from './PageContainer';
export type { PageContainerProps } from './PageContainer';

// 搜索表格
export { default as SearchTable } from './SearchTable';
export type { SearchTableProps, SearchField, ToolbarButton } from './SearchTable';

// 统计卡片
export { default as StatCard } from './StatCard';
export type { StatCardProps } from './StatCard';

// 主题切换
export { ThemeToggle } from './ThemeToggle';
