import { forwardRef } from 'react';
import { cn } from '@/lib/utils';
import { Button } from '@/components/ui/button';
import { HStack } from '@/components/layout';
import { ChevronRight } from 'lucide-react';
import type { LucideIcon } from 'lucide-react';

export interface SectionHeaderProps extends React.HTMLAttributes<HTMLDivElement> {
  /** 标题 */
  title: string;
  /** 副标题 */
  subtitle?: string;
  /** 图标 */
  icon?: LucideIcon;
  /** "查看更多"链接 */
  moreLink?: string;
  /** "查看更多"文本 */
  moreText?: string;
  /** "查看更多"点击回调 */
  onMore?: () => void;
  /** 右侧自定义内容 */
  extra?: React.ReactNode;
  /** 尺寸 */
  size?: 'sm' | 'md' | 'lg';
  /** 是否显示下划线 */
  underline?: boolean;
  /** 是否显示左侧装饰条 */
  decorated?: boolean;
}

const sizeConfig = {
  sm: {
    title: 'text-sm font-medium',
    subtitle: 'text-xs',
    icon: 'w-4 h-4',
    gap: 2,
  },
  md: {
    title: 'text-base font-semibold',
    subtitle: 'text-sm',
    icon: 'w-5 h-5',
    gap: 2,
  },
  lg: {
    title: 'text-lg font-bold',
    subtitle: 'text-base',
    icon: 'w-6 h-6',
    gap: 3,
  },
};

export const SectionHeader = forwardRef<HTMLDivElement, SectionHeaderProps>(({
  title,
  subtitle,
  icon: Icon,
  moreLink,
  moreText = '查看更多',
  onMore,
  extra,
  size = 'md',
  underline = false,
  decorated = false,
  className,
  children,
  ...props
}, ref) => {
  const config = sizeConfig[size];
  const showMore = moreLink || onMore;

  const handleMore = () => {
    if (onMore) {
      onMore();
    } else if (moreLink) {
      window.location.href = moreLink;
    }
  };

  return (
    <div
      ref={ref}
      className={cn(
        'flex items-center justify-between',
        underline && 'pb-3 border-b border-border/50',
        className
      )}
      {...props}
    >
      <HStack spacing={config.gap} align="center">
        {/* 左侧装饰条 */}
        {decorated && (
          <div className="w-1 h-5 rounded-full bg-primary" />
        )}
        
        {/* 图标 */}
        {Icon && (
          <Icon className={cn(config.icon, 'text-primary shrink-0')} />
        )}
        
        {/* 标题和副标题 */}
        <div>
          <h3 className={cn(config.title, 'text-foreground')}>
            {title}
          </h3>
          {subtitle && (
            <p className={cn(config.subtitle, 'text-muted-foreground mt-0.5')}>
              {subtitle}
            </p>
          )}
        </div>
      </HStack>
      
      {/* 右侧内容 */}
      <HStack spacing={2} align="center">
        {extra}
        {children}
        {showMore && (
          <Button
            variant="ghost"
            size="sm"
            onClick={handleMore}
            className="text-muted-foreground hover:text-primary group"
          >
            <span className="text-sm">{moreText}</span>
            <ChevronRight className="w-4 h-4 ml-1 transition-transform group-hover:translate-x-0.5" />
          </Button>
        )}
      </HStack>
    </div>
  );
});

SectionHeader.displayName = 'SectionHeader';
