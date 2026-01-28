import { forwardRef } from 'react';
import { cn } from '@/lib/utils';
import { ImageBox } from '@/components/layout';

export interface GameCardProps extends React.HTMLAttributes<HTMLDivElement> {
  /** 游戏 ID */
  gameId: number;
  /** 游戏名称 */
  name: string;
  /** 游戏图标 URL */
  iconUrl?: string;
  /** 是否选中 */
  selected?: boolean;
  /** 尺寸 */
  size?: 'sm' | 'md' | 'lg';
  /** 变体 */
  variant?: 'default' | 'compact' | 'detailed';
  /** 描述（detailed 变体显示） */
  description?: string;
  /** 玩家数量（detailed 变体显示） */
  playerCount?: number;
  /** 点击回调 */
  onSelect?: (gameId: number) => void;
}

const sizeConfig = {
  sm: {
    container: 'w-14',
    icon: 'w-12 h-12',
    text: 'text-[10px]',
    radius: 'rounded-xl',
  },
  md: {
    container: 'w-18',
    icon: 'w-14 h-14',
    text: 'text-xs',
    radius: 'rounded-2xl',
  },
  lg: {
    container: 'w-22',
    icon: 'w-18 h-18',
    text: 'text-sm',
    radius: 'rounded-2xl',
  },
};

export const GameCard = forwardRef<HTMLDivElement, GameCardProps>(({
  gameId,
  name,
  iconUrl,
  selected = false,
  size = 'md',
  variant = 'default',
  description,
  playerCount,
  onSelect,
  className,
  ...props
}, ref) => {
  const config = sizeConfig[size];

  const handleClick = () => {
    onSelect?.(gameId);
  };

  // 详细变体
  if (variant === 'detailed') {
    return (
      <div
        ref={ref}
        className={cn(
          'flex items-center gap-3 p-3 rounded-xl cursor-pointer',
          'border border-border/50 bg-card',
          'hover:bg-muted/50 hover:border-primary/30',
          'transition-all duration-200',
          selected && 'ring-2 ring-primary border-primary bg-primary/5',
          className
        )}
        onClick={handleClick}
        {...props}
      >
        <ImageBox
          src={iconUrl}
          alt={name}
          className={cn('w-12 h-12 shrink-0', config.radius)}
          fallback={name.charAt(0)}
        />
        <div className="flex-1 min-w-0">
          <h4 className="font-medium text-sm truncate">{name}</h4>
          {description && (
            <p className="text-xs text-muted-foreground truncate mt-0.5">{description}</p>
          )}
          {playerCount !== undefined && (
            <p className="text-xs text-muted-foreground mt-1">
              {playerCount} 位陪玩
            </p>
          )}
        </div>
      </div>
    );
  }

  // 紧凑变体
  if (variant === 'compact') {
    return (
      <div
        ref={ref}
        className={cn(
          'flex items-center gap-2 px-3 py-2 rounded-lg cursor-pointer',
          'hover:bg-muted/50',
          'transition-all duration-200',
          selected && 'bg-primary/10 text-primary',
          className
        )}
        onClick={handleClick}
        {...props}
      >
        <ImageBox
          src={iconUrl}
          alt={name}
          className="w-6 h-6 rounded-md shrink-0"
          fallback={name.charAt(0)}
        />
        <span className="text-sm truncate">{name}</span>
      </div>
    );
  }

  // 默认变体 - 图标 + 名称垂直排列
  return (
    <div
      ref={ref}
      className={cn(
        'flex flex-col items-center gap-2 cursor-pointer group',
        config.container,
        className
      )}
      onClick={handleClick}
      {...props}
    >
      <div className={cn(
        'relative overflow-hidden',
        config.icon,
        config.radius,
        'transition-all duration-300',
        'group-hover:scale-105 group-hover:shadow-lg',
        selected && 'ring-2 ring-primary ring-offset-2 ring-offset-background'
      )}>
        <ImageBox
          src={iconUrl}
          alt={name}
          className="w-full h-full"
          fallback={name.charAt(0)}
        />
      </div>
      <span className={cn(
        'text-center truncate w-full',
        config.text,
        'text-muted-foreground group-hover:text-foreground',
        'transition-colors duration-200',
        selected && 'text-primary font-medium'
      )}>
        {name}
      </span>
    </div>
  );
});

GameCard.displayName = 'GameCard';
