import { forwardRef } from 'react';
import { cn } from '@/lib/utils';
import { Button } from '@/components/ui/button';
import { HStack, VStack } from '@/components/layout';
import { ChevronLeft, MoreHorizontal } from 'lucide-react';
import type { LucideIcon } from 'lucide-react';
import { ActionMenu, type ActionMenuItem } from './action-menu';

export interface PageHeaderProps extends React.HTMLAttributes<HTMLDivElement> {
  /** 标题 */
  title: string;
  /** 副标题 */
  subtitle?: string;
  /** 图标 */
  icon?: LucideIcon;
  /** 是否显示返回按钮 */
  showBack?: boolean;
  /** 返回按钮回调 */
  onBack?: () => void;
  /** 返回按钮文本 */
  backText?: string;
  /** 主操作按钮 */
  primaryAction?: {
    label: string;
    icon?: LucideIcon;
    onClick: () => void;
    loading?: boolean;
    disabled?: boolean;
  };
  /** 次要操作按钮列表 */
  secondaryActions?: ActionMenuItem[];
  /** 右侧自定义内容 */
  extra?: React.ReactNode;
  /** 是否固定在顶部 */
  sticky?: boolean;
  /** 是否显示底部边框 */
  bordered?: boolean;
  /** 尺寸 */
  size?: 'sm' | 'md' | 'lg';
}

const sizeConfig = {
  sm: {
    padding: 'py-3',
    title: 'text-lg font-semibold',
    subtitle: 'text-xs',
    icon: 'w-5 h-5',
  },
  md: {
    padding: 'py-4',
    title: 'text-xl font-bold',
    subtitle: 'text-sm',
    icon: 'w-6 h-6',
  },
  lg: {
    padding: 'py-6',
    title: 'text-2xl font-bold',
    subtitle: 'text-base',
    icon: 'w-7 h-7',
  },
};

export const PageHeader = forwardRef<HTMLDivElement, PageHeaderProps>(({
  title,
  subtitle,
  icon: Icon,
  showBack = false,
  onBack,
  backText,
  primaryAction,
  secondaryActions,
  extra,
  sticky = false,
  bordered = true,
  size = 'md',
  className,
  children,
  ...props
}, ref) => {
  const config = sizeConfig[size];

  const handleBack = () => {
    if (onBack) {
      onBack();
    } else {
      window.history.back();
    }
  };

  return (
    <header
      ref={ref}
      className={cn(
        'bg-background/95 backdrop-blur-sm',
        config.padding,
        bordered && 'border-b border-border/50',
        sticky && 'sticky top-0 z-40',
        className
      )}
      {...props}
    >
      <HStack justify="between" align="center">
        {/* 左侧：返回按钮 + 标题 */}
        <HStack spacing={3} align="center">
          {showBack && (
            <Button
              variant="ghost"
              size="sm"
              onClick={handleBack}
              className="-ml-2"
            >
              <ChevronLeft className="w-4 h-4 mr-1" />
              {backText || '返回'}
            </Button>
          )}
          
          {Icon && (
            <div className="p-2 rounded-lg bg-primary/10 text-primary">
              <Icon className={config.icon} />
            </div>
          )}
          
          <VStack spacing={0.5}>
            <h1 className={cn(config.title, 'text-foreground')}>
              {title}
            </h1>
            {subtitle && (
              <p className={cn(config.subtitle, 'text-muted-foreground')}>
                {subtitle}
              </p>
            )}
          </VStack>
        </HStack>
        
        {/* 右侧：操作按钮 */}
        <HStack spacing={2} align="center">
          {extra}
          {children}
          
          {primaryAction && (
            <Button
              onClick={primaryAction.onClick}
              disabled={primaryAction.disabled || primaryAction.loading}
            >
              {primaryAction.loading ? (
                <span className="animate-spin mr-2">⟳</span>
              ) : primaryAction.icon ? (
                <primaryAction.icon className="w-4 h-4 mr-2" />
              ) : null}
              {primaryAction.label}
            </Button>
          )}
          
          {secondaryActions && secondaryActions.length > 0 && (
            <ActionMenu
              items={secondaryActions}
              trigger={
                <Button variant="outline" size="icon">
                  <MoreHorizontal className="w-4 h-4" />
                </Button>
              }
            />
          )}
        </HStack>
      </HStack>
    </header>
  );
});

PageHeader.displayName = 'PageHeader';

// PageHeaderWithStats 组件 - 带统计的页面头部
export interface PageHeaderWithStatsProps extends Omit<PageHeaderProps, 'children'> {
  /** 统计数据 */
  stats?: Array<{
    label: string;
    value: string | number;
    change?: number;
  }>;
}

export const PageHeaderWithStats = forwardRef<HTMLDivElement, PageHeaderWithStatsProps>(({
  stats,
  className,
  ...props
}, ref) => {
  return (
    <VStack spacing={4} className={className}>
      <PageHeader ref={ref} bordered={false} {...props} />
      
      {stats && stats.length > 0 && (
        <HStack spacing={6} className="px-4 py-3 bg-muted/30 rounded-lg">
          {stats.map((stat, index) => (
            <VStack key={index} spacing={0.5} className="text-center min-w-[80px]">
              <span className="text-2xl font-bold tabular-nums">{stat.value}</span>
              <span className="text-xs text-muted-foreground">{stat.label}</span>
              {stat.change !== undefined && (
                <span className={cn(
                  'text-[10px]',
                  stat.change > 0 ? 'text-green-500' : stat.change < 0 ? 'text-red-500' : 'text-muted-foreground'
                )}>
                  {stat.change > 0 ? '+' : ''}{stat.change}%
                </span>
              )}
            </VStack>
          ))}
        </HStack>
      )}
    </VStack>
  );
});

PageHeaderWithStats.displayName = 'PageHeaderWithStats';
