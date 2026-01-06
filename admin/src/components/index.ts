/**
 * 组件统一导出
 */

// 权限相关
export { PermissionGuard, PermissionButton } from './PermissionGuard';
export { withPermission } from './withPermission';
export type { PermissionGuardProps, PermissionButtonProps } from './PermissionGuard';

// 页面容器
export { default as PageContainer } from './PageContainer';
export type { PageContainerProps } from './PageContainer';

// 搜索表格
export { SearchTable } from './SearchTable';
export type { SearchTableProps, SearchField, ToolbarButton } from './SearchTable';

// 统计卡片
export { default as StatCard } from './StatCard';
export type { StatCardProps } from './StatCard';

// 主题切换
export { ThemeToggle } from './ThemeToggle';

// 权限树组件
export { PermissionTree } from './PermissionTree';
export type { PermissionTreeProps } from './PermissionTree';

// 错误边界
export { ErrorBoundary, PageErrorBoundary } from './ErrorBoundary';

// 统一按钮
export { Button, IconButton } from './Button';
export type { ButtonProps, IconButtonProps, ButtonSize, BtnVariant } from './Button';

// 统一卡片
export { Card, StatisticCard, ContentCard } from './Card';
export type { CardProps, CardVariant, CardPadding } from './Card';

// 可折叠区块
export { CollapsibleSection } from './CollapsibleSection';
export type { CollapsibleSectionProps } from './CollapsibleSection';

// 动画容器
export {
  AnimatedContainer,
  AnimatedListItem,
  PageTransition,
  animationVariants,
  transitions,
} from './AnimatedContainer';
export type {
  AnimatedContainerProps,
  AnimatedListItemProps,
  AnimationType,
  TransitionType,
} from './AnimatedContainer';
