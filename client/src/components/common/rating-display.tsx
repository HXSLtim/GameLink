import { forwardRef } from 'react';
import { cn } from '@/lib/utils';
import { Star, StarHalf } from 'lucide-react';
import { HStack } from '@/components/layout';

export interface RatingDisplayProps extends React.HTMLAttributes<HTMLDivElement> {
  /** 评分（0-5） */
  rating: number;
  /** 评价数量 */
  count?: number;
  /** 是否显示数值 */
  showValue?: boolean;
  /** 是否显示评价数量 */
  showCount?: boolean;
  /** 尺寸 */
  size?: 'xs' | 'sm' | 'md' | 'lg';
  /** 只显示一颗星和数值 */
  compact?: boolean;
  /** 星星颜色 */
  starColor?: string;
}

const sizeConfig = {
  xs: { star: 'w-3 h-3', text: 'text-xs', gap: 0.5 },
  sm: { star: 'w-3.5 h-3.5', text: 'text-sm', gap: 1 },
  md: { star: 'w-4 h-4', text: 'text-base', gap: 1 },
  lg: { star: 'w-5 h-5', text: 'text-lg', gap: 1.5 },
};

export const RatingDisplay = forwardRef<HTMLDivElement, RatingDisplayProps>(({
  rating,
  count,
  showValue = true,
  showCount = true,
  size = 'sm',
  compact = false,
  starColor = 'text-amber-400',
  className,
  ...props
}, ref) => {
  const config = sizeConfig[size];
  const clampedRating = Math.min(5, Math.max(0, rating));
  const fullStars = Math.floor(clampedRating);
  const hasHalfStar = clampedRating % 1 >= 0.5;
  const emptyStars = 5 - fullStars - (hasHalfStar ? 1 : 0);

  // 紧凑模式：只显示一颗星和数值
  if (compact) {
    return (
      <div
        ref={ref}
        className={cn('inline-flex items-center', className)}
        {...props}
      >
        <HStack spacing={config.gap} align="center">
          <Star className={cn(config.star, starColor, 'fill-current')} />
          {showValue && (
            <span className={cn('font-medium', config.text)}>
              {clampedRating.toFixed(1)}
            </span>
          )}
          {showCount && count !== undefined && (
            <span className={cn('text-muted-foreground', config.text)}>
              ({count})
            </span>
          )}
        </HStack>
      </div>
    );
  }

  return (
    <div
      ref={ref}
      className={cn('inline-flex items-center', className)}
      {...props}
    >
      <HStack spacing={config.gap} align="center">
        {/* 星星 */}
        <HStack spacing={0.5} align="center">
          {/* 满星 */}
          {Array.from({ length: fullStars }).map((_, i) => (
            <Star
              key={`full-${i}`}
              className={cn(config.star, starColor, 'fill-current')}
            />
          ))}
          
          {/* 半星 */}
          {hasHalfStar && (
            <div className="relative">
              <Star className={cn(config.star, 'text-muted-foreground/30')} />
              <StarHalf
                className={cn(config.star, starColor, 'fill-current absolute inset-0')}
              />
            </div>
          )}
          
          {/* 空星 */}
          {Array.from({ length: emptyStars }).map((_, i) => (
            <Star
              key={`empty-${i}`}
              className={cn(config.star, 'text-muted-foreground/30')}
            />
          ))}
        </HStack>
        
        {/* 数值 */}
        {showValue && (
          <span className={cn('font-medium', config.text)}>
            {clampedRating.toFixed(1)}
          </span>
        )}
        
        {/* 评价数量 */}
        {showCount && count !== undefined && (
          <span className={cn('text-muted-foreground', config.text)}>
            ({count >= 1000 ? `${(count / 1000).toFixed(1)}k` : count})
          </span>
        )}
      </HStack>
    </div>
  );
});

RatingDisplay.displayName = 'RatingDisplay';
