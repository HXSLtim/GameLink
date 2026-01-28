import { forwardRef } from 'react';
import { cn } from '@/lib/utils';
import { HStack, VStack } from '@/components/layout';
import { format } from 'date-fns';
import { zhCN } from 'date-fns/locale';
import {
  ArrowUpRight,
  ArrowDownLeft,
  Gift,
  ShoppingBag,
  Wallet,
  RefreshCcw,
  type LucideIcon,
} from 'lucide-react';

export type TransactionType =
  | 'recharge'
  | 'withdraw'
  | 'payment'
  | 'income'
  | 'refund'
  | 'gift'
  | 'bonus';

export interface TransactionItemProps extends React.HTMLAttributes<HTMLDivElement> {
  /** 交易 ID */
  transactionId: number;
  /** 交易类型 */
  type: TransactionType;
  /** 金额（分） */
  amountCents: number;
  /** 交易描述 */
  description?: string;
  /** 交易时间 */
  createdAt: string | Date;
  /** 状态 */
  status?: 'pending' | 'success' | 'failed';
  /** 货币符号 */
  currency?: string;
  /** 点击回调 */
  onClick?: (id: number) => void;
}

const typeConfig: Record<TransactionType, {
  label: string;
  icon: LucideIcon;
  color: string;
  isIncome: boolean;
}> = {
  recharge: { label: '充值', icon: ArrowDownLeft, color: 'text-green-500 bg-green-500/10', isIncome: true },
  withdraw: { label: '提现', icon: ArrowUpRight, color: 'text-orange-500 bg-orange-500/10', isIncome: false },
  payment: { label: '支付', icon: ShoppingBag, color: 'text-red-500 bg-red-500/10', isIncome: false },
  income: { label: '收入', icon: Wallet, color: 'text-green-500 bg-green-500/10', isIncome: true },
  refund: { label: '退款', icon: RefreshCcw, color: 'text-blue-500 bg-blue-500/10', isIncome: true },
  gift: { label: '礼物', icon: Gift, color: 'text-pink-500 bg-pink-500/10', isIncome: false },
  bonus: { label: '奖励', icon: Gift, color: 'text-yellow-500 bg-yellow-500/10', isIncome: true },
};

export const TransactionItem = forwardRef<HTMLDivElement, TransactionItemProps>(({
  transactionId,
  type,
  amountCents,
  description,
  createdAt,
  status = 'success',
  currency = '¥',
  onClick,
  className,
  ...props
}, ref) => {
  const config = typeConfig[type] || typeConfig.payment;
  const Icon = config.icon;
  const amount = amountCents / 100;
  const isIncome = config.isIncome;

  const handleClick = () => {
    onClick?.(transactionId);
  };

  return (
    <div
      ref={ref}
      className={cn(
        'flex items-center gap-4 p-4',
        'hover:bg-muted/50 transition-colors rounded-lg',
        onClick && 'cursor-pointer',
        className
      )}
      onClick={handleClick}
      {...props}
    >
      {/* 图标 */}
      <div className={cn(
        'w-10 h-10 rounded-full flex items-center justify-center shrink-0',
        config.color
      )}>
        <Icon className="w-5 h-5" />
      </div>

      {/* 信息 */}
      <VStack spacing={0.5} className="flex-1 min-w-0">
        <HStack justify="between" align="center" className="w-full">
          <span className="font-medium truncate">
            {description || config.label}
          </span>
          <span className={cn(
            'font-semibold tabular-nums shrink-0',
            isIncome ? 'text-green-500' : 'text-foreground'
          )}>
            {isIncome ? '+' : '-'}{currency}{amount.toFixed(2)}
          </span>
        </HStack>
        <HStack justify="between" align="center" className="w-full">
          <span className="text-xs text-muted-foreground">
            {format(new Date(createdAt), 'yyyy-MM-dd HH:mm', { locale: zhCN })}
          </span>
          {status !== 'success' && (
            <span className={cn(
              'text-xs',
              status === 'pending' ? 'text-yellow-500' : 'text-red-500'
            )}>
              {status === 'pending' ? '处理中' : '失败'}
            </span>
          )}
        </HStack>
      </VStack>
    </div>
  );
});

TransactionItem.displayName = 'TransactionItem';
