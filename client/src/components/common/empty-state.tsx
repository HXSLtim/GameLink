import { forwardRef } from 'react';
import { cn } from '@/lib/utils';
import { Button } from '@/components/ui/button';
import { VStack } from '@/components/layout';
import type { LucideIcon } from 'lucide-react';
import { PackageOpen } from 'lucide-react';

export interface EmptyStateProps extends React.HTMLAttributes<HTMLDivElement> {
  /** 图标 */
  icon?: LucideIcon;
  /** 标题 */
  title?: string;
  /** 描述 */
  description?: string;
  /** 操作按钮文本 */
  actionLabel?: string;
  /** 操作按钮回调 */
  onAction?: () => void;
  /** 尺寸 */
  size?: 'sm' | 'md' | 'lg';
}

const sizeConfig = {
  sm: {
    icon: 'w-10 h-10',
    title: 'text-sm',
    description: 'text-xs',
    spacing: 2,
    padding: 'py-6',
  },
  md: {
    icon: 'w-14 h-14',
    title: 'text-base',
    description: 'text-sm',
    spacing: 3,
    padding: 'py-10',
  },
  lg: {
    icon: 'w-20 h-20',
    title: 'text-lg',
    description: 'text-base',
    spacing: 4,
    padding: 'py-16',
  },
};

export const EmptyState = forwardRef<HTMLDivElement, EmptyStateProps>(({
  icon: Icon = PackageOpen,
  title = '暂无数据',
  description,
  actionLabel,
  onAction,
  size = 'md',
  className,
  children,
  ...props
}, ref) => {
  const config = sizeConfig[size];

  return (
    <div
      ref={ref}
      className={cn(
        'flex items-center justify-center w-full',
        config.padding,
        className
      )}
      {...props}
    >
      <VStack spacing={config.spacing} align="center" className="text-center max-w-xs">
        <div className={cn(
          'rounded-2xl bg-muted/50 p-4 text-muted-foreground/60',
          'transition-colors duration-200'
        )}>
          <Icon className={config.icon} strokeWidth={1.5} />
        </div>
        
        <VStack spacing={1} align="center">
          <p className={cn('font-medium text-foreground', config.title)}>
            {title}
          </p>
          {description && (
            <p className={cn('text-muted-foreground', config.description)}>
              {description}
            </p>
          )}
        </VStack>
        
        {actionLabel && onAction && (
          <Button
            variant="outline"
            size={size === 'sm' ? 'sm' : 'default'}
            onClick={onAction}
            className="mt-2"
          >
            {actionLabel}
          </Button>
        )}
        
        {children}
      </VStack>
    </div>
  );
});

EmptyState.displayName = 'EmptyState';
