import { forwardRef } from 'react';
import { cn } from '@/lib/utils';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip';

export interface AvatarGroupItem {
  id: string | number;
  name: string;
  imageUrl?: string;
}

export interface AvatarGroupProps extends React.HTMLAttributes<HTMLDivElement> {
  /** 头像列表 */
  items: AvatarGroupItem[];
  /** 最大显示数量 */
  max?: number;
  /** 尺寸 */
  size?: 'xs' | 'sm' | 'md' | 'lg';
  /** 是否显示 tooltip */
  showTooltip?: boolean;
  /** 剩余数量点击回调 */
  onRemainderClick?: () => void;
}

const sizeConfig = {
  xs: { avatar: 'w-6 h-6', text: 'text-[10px]', offset: '-ml-2' },
  sm: { avatar: 'w-8 h-8', text: 'text-xs', offset: '-ml-2.5' },
  md: { avatar: 'w-10 h-10', text: 'text-sm', offset: '-ml-3' },
  lg: { avatar: 'w-12 h-12', text: 'text-base', offset: '-ml-4' },
};

export const AvatarGroup = forwardRef<HTMLDivElement, AvatarGroupProps>(({
  items,
  max = 5,
  size = 'md',
  showTooltip = true,
  onRemainderClick,
  className,
  ...props
}, ref) => {
  const config = sizeConfig[size];
  const visibleItems = items.slice(0, max);
  const remainderCount = items.length - max;

  const renderAvatar = (item: AvatarGroupItem, index: number) => {
    const avatar = (
      <Avatar
        className={cn(
          config.avatar,
          'border-2 border-background',
          index > 0 && config.offset
        )}
        style={{ zIndex: items.length - index }}
      >
        <AvatarImage src={item.imageUrl} alt={item.name} />
        <AvatarFallback className={config.text}>
          {item.name.charAt(0).toUpperCase()}
        </AvatarFallback>
      </Avatar>
    );

    if (showTooltip) {
      return (
        <TooltipProvider key={item.id}>
          <Tooltip>
            <TooltipTrigger asChild>
              {avatar}
            </TooltipTrigger>
            <TooltipContent>
              <p>{item.name}</p>
            </TooltipContent>
          </Tooltip>
        </TooltipProvider>
      );
    }

    return <div key={item.id}>{avatar}</div>;
  };

  return (
    <div
      ref={ref}
      className={cn('flex items-center', className)}
      {...props}
    >
      {visibleItems.map((item, index) => renderAvatar(item, index))}
      
      {remainderCount > 0 && (
        <div
          className={cn(
            config.avatar,
            config.offset,
            'flex items-center justify-center rounded-full',
            'bg-muted border-2 border-background',
            'text-muted-foreground font-medium',
            config.text,
            onRemainderClick && 'cursor-pointer hover:bg-muted/80 transition-colors'
          )}
          style={{ zIndex: 0 }}
          onClick={onRemainderClick}
        >
          +{remainderCount}
        </div>
      )}
    </div>
  );
});

AvatarGroup.displayName = 'AvatarGroup';
