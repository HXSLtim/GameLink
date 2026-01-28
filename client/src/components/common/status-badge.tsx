import { forwardRef } from 'react';
import { cn } from '@/lib/utils';
import { Badge } from '@/components/ui/badge';
import type { LucideIcon } from 'lucide-react';
import {
  CheckCircle,
  XCircle,
  Clock,
  AlertCircle,
  Loader2,
  Circle,
} from 'lucide-react';

export type StatusType =
  | 'success'
  | 'error'
  | 'warning'
  | 'info'
  | 'pending'
  | 'processing'
  | 'default';

export interface StatusBadgeProps extends React.HTMLAttributes<HTMLDivElement> {
  /** 状态类型 */
  status: StatusType;
  /** 显示文本 */
  label?: string;
  /** 是否显示图标 */
  showIcon?: boolean;
  /** 自定义图标 */
  icon?: LucideIcon;
  /** 尺寸 */
  size?: 'sm' | 'md' | 'lg';
  /** 是否带有脉冲动画（用于在线状态等） */
  pulse?: boolean;
  /** 是否只显示圆点 */
  dot?: boolean;
}

const statusConfig: Record<StatusType, {
  icon: LucideIcon;
  classes: string;
  dotColor: string;
}> = {
  success: {
    icon: CheckCircle,
    classes: 'bg-green-500/10 text-green-600 border-green-500/20',
    dotColor: 'bg-green-500',
  },
  error: {
    icon: XCircle,
    classes: 'bg-red-500/10 text-red-600 border-red-500/20',
    dotColor: 'bg-red-500',
  },
  warning: {
    icon: AlertCircle,
    classes: 'bg-amber-500/10 text-amber-600 border-amber-500/20',
    dotColor: 'bg-amber-500',
  },
  info: {
    icon: Circle,
    classes: 'bg-blue-500/10 text-blue-600 border-blue-500/20',
    dotColor: 'bg-blue-500',
  },
  pending: {
    icon: Clock,
    classes: 'bg-gray-500/10 text-gray-600 border-gray-500/20',
    dotColor: 'bg-gray-500',
  },
  processing: {
    icon: Loader2,
    classes: 'bg-primary/10 text-primary border-primary/20',
    dotColor: 'bg-primary',
  },
  default: {
    icon: Circle,
    classes: 'bg-muted text-muted-foreground border-border',
    dotColor: 'bg-muted-foreground',
  },
};

const sizeConfig = {
  sm: { badge: 'text-xs px-2 py-0.5', icon: 'w-3 h-3', dot: 'w-1.5 h-1.5' },
  md: { badge: 'text-sm px-2.5 py-1', icon: 'w-3.5 h-3.5', dot: 'w-2 h-2' },
  lg: { badge: 'text-base px-3 py-1.5', icon: 'w-4 h-4', dot: 'w-2.5 h-2.5' },
};

export const StatusBadge = forwardRef<HTMLDivElement, StatusBadgeProps>(({
  status,
  label,
  showIcon = true,
  icon,
  size = 'md',
  pulse = false,
  dot = false,
  className,
  ...props
}, ref) => {
  const config = statusConfig[status];
  const sizeConf = sizeConfig[size];
  const Icon = icon || config.icon;
  const isProcessing = status === 'processing';

  // 只显示圆点模式
  if (dot) {
    return (
      <div
        ref={ref}
        className={cn('inline-flex items-center gap-1.5', className)}
        {...props}
      >
        <span className="relative flex">
          <span className={cn(
            'rounded-full',
            sizeConf.dot,
            config.dotColor
          )} />
          {pulse && (
            <span className={cn(
              'absolute inset-0 rounded-full animate-ping opacity-75',
              config.dotColor
            )} />
          )}
        </span>
        {label && (
          <span className={cn(
            'font-medium',
            size === 'sm' ? 'text-xs' : size === 'md' ? 'text-sm' : 'text-base'
          )}>
            {label}
          </span>
        )}
      </div>
    );
  }

  return (
    <Badge
      ref={ref as React.Ref<HTMLDivElement>}
      variant="outline"
      className={cn(
        'inline-flex items-center gap-1.5 font-medium border',
        config.classes,
        sizeConf.badge,
        className
      )}
      {...props}
    >
      {showIcon && (
        <Icon className={cn(
          sizeConf.icon,
          isProcessing && 'animate-spin'
        )} />
      )}
      {label}
    </Badge>
  );
});

StatusBadge.displayName = 'StatusBadge';

// 预定义的状态映射（方便业务使用）
export const orderStatusMap: Record<string, { status: StatusType; label: string }> = {
  pending: { status: 'pending', label: '待支付' },
  paid: { status: 'info', label: '已支付' },
  accepted: { status: 'processing', label: '进行中' },
  completed: { status: 'success', label: '已完成' },
  cancelled: { status: 'default', label: '已取消' },
  refunded: { status: 'warning', label: '已退款' },
  disputed: { status: 'error', label: '争议中' },
};

export const onlineStatusMap: Record<string, { status: StatusType; label: string }> = {
  online: { status: 'success', label: '在线' },
  busy: { status: 'warning', label: '忙碌' },
  away: { status: 'pending', label: '离开' },
  offline: { status: 'default', label: '离线' },
  in_game: { status: 'info', label: '游戏中' },
};
