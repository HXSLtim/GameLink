import { forwardRef } from 'react';
import { cn } from '@/lib/utils';
import { Card, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Wallet, Plus, ArrowUpRight } from 'lucide-react';
import { HStack, VStack } from '@/components/layout';

export interface BalanceCardProps extends React.HTMLAttributes<HTMLDivElement> {
  /** 余额（分） */
  balanceCents: number;
  /** 冻结金额（分） */
  frozenCents?: number;
  /** 货币符号 */
  currency?: string;
  /** 是否显示充值按钮 */
  showRecharge?: boolean;
  /** 是否显示提现按钮 */
  showWithdraw?: boolean;
  /** 充值回调 */
  onRecharge?: () => void;
  /** 提现回调 */
  onWithdraw?: () => void;
  /** 变体 */
  variant?: 'default' | 'gradient' | 'minimal';
  /** 尺寸 */
  size?: 'sm' | 'md' | 'lg';
}

const sizeConfig = {
  sm: {
    padding: 'p-4',
    title: 'text-xs',
    balance: 'text-2xl',
    symbol: 'text-sm',
    height: 'h-[140px]',
    icon: 'h-24 w-24',
    button: 'sm' as const,
  },
  md: {
    padding: 'p-6',
    title: 'text-sm',
    balance: 'text-4xl',
    symbol: 'text-lg',
    height: 'h-[180px]',
    icon: 'h-36 w-36',
    button: 'default' as const,
  },
  lg: {
    padding: 'p-8',
    title: 'text-sm',
    balance: 'text-5xl',
    symbol: 'text-2xl',
    height: 'h-[240px]',
    icon: 'h-48 w-48',
    button: 'lg' as const,
  },
};

export const BalanceCard = forwardRef<HTMLDivElement, BalanceCardProps>(({
  balanceCents,
  frozenCents,
  currency = '¥',
  showRecharge = true,
  showWithdraw = true,
  onRecharge,
  onWithdraw,
  variant = 'gradient',
  size = 'md',
  className,
  ...props
}, ref) => {
  const config = sizeConfig[size];
  const balance = balanceCents / 100;
  const frozen = frozenCents ? frozenCents / 100 : 0;

  const variantClasses = {
    default: 'bg-card border',
    gradient: 'text-white border-0 shadow-2xl bg-gradient-to-br from-indigo-600 via-purple-600 to-primary',
    minimal: 'bg-muted/30 border-0',
  };

  const textClasses = {
    default: 'text-foreground',
    gradient: 'text-white',
    minimal: 'text-foreground',
  };

  const subtextClasses = {
    default: 'text-muted-foreground',
    gradient: 'text-indigo-100 opacity-80',
    minimal: 'text-muted-foreground',
  };

  const buttonClasses = {
    default: {
      primary: '',
      secondary: 'variant-outline',
    },
    gradient: {
      primary: 'bg-white text-primary hover:bg-white/90 font-bold shadow-lg border-0',
      secondary: 'bg-transparent border-white/20 text-white hover:bg-white/10 hover:border-white/40 backdrop-blur-md',
    },
    minimal: {
      primary: '',
      secondary: 'variant-outline',
    },
  };

  return (
    <Card
      ref={ref}
      className={cn(
        'relative overflow-hidden',
        variantClasses[variant],
        className
      )}
      {...props}
    >
      {/* 背景装饰 */}
      {variant === 'gradient' && (
        <div className="absolute top-0 right-0 p-8 opacity-10">
          <Wallet className={config.icon} />
        </div>
      )}

      <CardContent className={cn(
        config.padding,
        'flex flex-col justify-between relative z-10',
        config.height
      )}>
        <VStack spacing={2}>
          <p className={cn(
            'font-medium tracking-wide uppercase',
            config.title,
            subtextClasses[variant]
          )}>
            账户余额
          </p>
          <h2 className={cn(
            'font-bold tracking-tight flex items-baseline gap-2',
            config.balance,
            textClasses[variant]
          )}>
            <span className={cn('opacity-60', config.symbol)}>{currency}</span>
            {balance.toFixed(2)}
          </h2>
          {frozen > 0 && (
            <p className={cn('text-xs', subtextClasses[variant])}>
              冻结金额: {currency}{frozen.toFixed(2)}
            </p>
          )}
        </VStack>

        {(showRecharge || showWithdraw) && (
          <HStack spacing={3}>
            {showRecharge && (
              <Button
                size={config.button}
                onClick={onRecharge}
                className={buttonClasses[variant].primary}
              >
                <Plus className="h-4 w-4 mr-2" />
                充值
              </Button>
            )}
            {showWithdraw && (
              <Button
                size={config.button}
                variant="outline"
                onClick={onWithdraw}
                className={buttonClasses[variant].secondary}
              >
                <ArrowUpRight className="h-4 w-4 mr-2" />
                提现
              </Button>
            )}
          </HStack>
        )}
      </CardContent>
    </Card>
  );
});

BalanceCard.displayName = 'BalanceCard';
