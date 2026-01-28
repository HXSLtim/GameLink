import { forwardRef } from 'react';
import { cn } from '@/lib/utils';

export interface PriceTagProps extends React.HTMLAttributes<HTMLSpanElement> {
  /** 价格（分） */
  priceCents?: number;
  /** 价格（元），优先使用 */
  priceYuan?: number;
  /** 单位 */
  unit?: string;
  /** 原价（分） */
  originalPriceCents?: number;
  /** 原价（元） */
  originalPriceYuan?: number;
  /** 尺寸 */
  size?: 'xs' | 'sm' | 'md' | 'lg' | 'xl';
  /** 颜色变体 */
  variant?: 'default' | 'primary' | 'success' | 'warning';
  /** 是否显示货币符号 */
  showSymbol?: boolean;
  /** 是否显示免费文字 */
  showFreeText?: boolean;
}

const sizeClasses = {
  xs: 'text-xs',
  sm: 'text-sm',
  md: 'text-base',
  lg: 'text-lg',
  xl: 'text-2xl font-bold',
};

const symbolSizes = {
  xs: 'text-[10px]',
  sm: 'text-xs',
  md: 'text-sm',
  lg: 'text-base',
  xl: 'text-lg',
};

const variantClasses = {
  default: 'text-foreground',
  primary: 'text-primary',
  success: 'text-green-500',
  warning: 'text-amber-500',
};

export const PriceTag = forwardRef<HTMLSpanElement, PriceTagProps>(({
  priceCents,
  priceYuan,
  unit = '/小时',
  originalPriceCents,
  originalPriceYuan,
  size = 'md',
  variant = 'primary',
  showSymbol = true,
  showFreeText = true,
  className,
  ...props
}, ref) => {
  // 计算实际价格（元）
  const price = priceYuan ?? (priceCents !== undefined ? priceCents / 100 : undefined);
  const originalPrice = originalPriceYuan ?? (originalPriceCents !== undefined ? originalPriceCents / 100 : undefined);
  
  // 判断是否免费
  const isFree = price === 0 || price === undefined;
  
  if (isFree && showFreeText) {
    return (
      <span
        ref={ref}
        className={cn(
          'font-medium text-green-500',
          sizeClasses[size],
          className
        )}
        {...props}
      >
        免费
      </span>
    );
  }
  
  // 格式化价格
  const formatPrice = (value: number) => {
    return value % 1 === 0 ? value.toString() : value.toFixed(2);
  };

  return (
    <span
      ref={ref}
      className={cn('inline-flex items-baseline gap-0.5', className)}
      {...props}
    >
      {/* 原价（划线） */}
      {originalPrice !== undefined && originalPrice > (price || 0) && (
        <span className={cn(
          'line-through text-muted-foreground/60 mr-1',
          size === 'xs' ? 'text-[10px]' : 'text-xs'
        )}>
          ¥{formatPrice(originalPrice)}
        </span>
      )}
      
      {/* 货币符号 */}
      {showSymbol && (
        <span className={cn(
          'font-medium',
          symbolSizes[size],
          variantClasses[variant]
        )}>
          ¥
        </span>
      )}
      
      {/* 价格数值 */}
      <span className={cn(
        'font-semibold tabular-nums',
        sizeClasses[size],
        variantClasses[variant]
      )}>
        {price !== undefined ? formatPrice(price) : '--'}
      </span>
      
      {/* 单位 */}
      {unit && (
        <span className={cn(
          'text-muted-foreground/70 font-normal ml-0.5',
          size === 'xs' || size === 'sm' ? 'text-[10px]' : 'text-xs'
        )}>
          {unit}
        </span>
      )}
    </span>
  );
});

PriceTag.displayName = 'PriceTag';
