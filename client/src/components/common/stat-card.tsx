import { forwardRef } from 'react';
import { cn } from '@/lib/utils';
import { Card, CardContent } from '@/components/ui/card';
import { VStack, HStack } from '@/components/layout';
import type { LucideIcon } from 'lucide-react';
import { TrendingUp, TrendingDown, Minus } from 'lucide-react';

export interface StatCardProps extends React.HTMLAttributes<HTMLDivElement> {
  /** 标题 */
  title: string;
  /** 数值 */
  value: string | number;
  /** 图标 */
  icon?: LucideIcon;
  /** 变化百分比 */
  change?: number;
  /** 变化描述 */
  changeLabel?: string;
  /** 变化趋势 */
  trend?: 'up' | 'down' | 'neutral';
  /** 变体 */
  variant?: 'default' | 'outlined' | 'filled';
  /** 尺寸 */
  size?: 'sm' | 'md' | 'lg';
  /** 是否加载中 */
  loading?: boolean;
}

const sizeConfig = {
  sm: {
    padding: 'p-3',
    icon: 'w-4 h-4',
    iconBox: 'w-8 h-8',
    title: 'text-xs',
    value: 'text-lg font-bold',
    change: 'text-[10px]',
  },
  md: {
    padding: 'p-4',
    icon: 'w-5 h-5',
    iconBox: 'w-10 h-10',
    title: 'text-sm',
    value: 'text-2xl font-bold',
    change: 'text-xs',
  },
  lg: {
    padding: 'p-6',
    icon: 'w-6 h-6',
    iconBox: 'w-12 h-12',
    title: 'text-base',
    value: 'text-3xl font-bold',
    change: 'text-sm',
  },
};

const variantConfig = {
  default: 'bg-card',
  outlined: 'bg-transparent border-2',
  filled: 'bg-primary/5 border-primary/10',
};

export const StatCard = forwardRef<HTMLDivElement, StatCardProps>(({
  title,
  value,
  icon: Icon,
  change,
  changeLabel = '较上期',
  trend,
  variant = 'default',
  size = 'md',
  loading = false,
  className,
  ...props
}, ref) => {
  const config = sizeConfig[size];
  
  // 自动判断趋势
  const actualTrend = trend ?? (change !== undefined ? (change > 0 ? 'up' : change < 0 ? 'down' : 'neutral') : undefined);
  
  const TrendIcon = actualTrend === 'up' ? TrendingUp : actualTrend === 'down' ? TrendingDown : Minus;
  const trendColor = actualTrend === 'up' ? 'text-green-500' : actualTrend === 'down' ? 'text-red-500' : 'text-muted-foreground';

  if (loading) {
    return (
      <Card ref={ref} className={cn(variantConfig[variant], className)} {...props}>
        <CardContent className={config.padding}>
          <VStack spacing={2}>
            <div className="h-4 w-20 bg-muted animate-pulse rounded" />
            <div className="h-8 w-16 bg-muted animate-pulse rounded" />
            <div className="h-3 w-24 bg-muted animate-pulse rounded" />
          </VStack>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card ref={ref} className={cn(variantConfig[variant], className)} {...props}>
      <CardContent className={config.padding}>
        <HStack justify="between" align="start">
          <VStack spacing={1}>
            <span className={cn('text-muted-foreground', config.title)}>{title}</span>
            <span className={cn('text-foreground tabular-nums', config.value)}>{value}</span>
            
            {change !== undefined && (
              <HStack spacing={1} align="center" className={cn(config.change, trendColor)}>
                <TrendIcon className="w-3 h-3" />
                <span>{change > 0 ? '+' : ''}{change}%</span>
                <span className="text-muted-foreground">{changeLabel}</span>
              </HStack>
            )}
          </VStack>
          
          {Icon && (
            <div className={cn(
              'rounded-lg flex items-center justify-center',
              config.iconBox,
              'bg-primary/10 text-primary'
            )}>
              <Icon className={config.icon} />
            </div>
          )}
        </HStack>
      </CardContent>
    </Card>
  );
});

StatCard.displayName = 'StatCard';

// StatGrid 组件 - 统计卡片网格
export interface StatGridProps extends React.HTMLAttributes<HTMLDivElement> {
  /** 统计项 */
  stats: Array<Omit<StatCardProps, 'ref'>>;
  /** 列数 */
  columns?: 2 | 3 | 4;
  /** 尺寸 */
  size?: 'sm' | 'md' | 'lg';
}

const columnClasses = {
  2: 'grid-cols-1 sm:grid-cols-2',
  3: 'grid-cols-1 sm:grid-cols-2 lg:grid-cols-3',
  4: 'grid-cols-2 sm:grid-cols-2 lg:grid-cols-4',
};

export const StatGrid = forwardRef<HTMLDivElement, StatGridProps>(({
  stats,
  columns = 4,
  size = 'md',
  className,
  ...props
}, ref) => {
  return (
    <div
      ref={ref}
      className={cn('grid gap-4', columnClasses[columns], className)}
      {...props}
    >
      {stats.map((stat, index) => (
        <StatCard key={index} {...stat} size={stat.size ?? size} />
      ))}
    </div>
  );
});

StatGrid.displayName = 'StatGrid';
