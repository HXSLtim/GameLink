import { useState } from 'react';
import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
  DropdownMenuSeparator,
  DropdownMenuLabel,
} from '@/components/ui/dropdown-menu';
import { MoreVertical, MoreHorizontal } from 'lucide-react';
import { cn } from '@/lib/utils';
import type { LucideIcon } from 'lucide-react';

export interface ActionMenuItem {
  key: string;
  label: string;
  icon?: LucideIcon;
  onClick: () => void;
  variant?: 'default' | 'destructive';
  disabled?: boolean;
  hidden?: boolean;
}

export interface ActionMenuGroup {
  label?: string;
  items: ActionMenuItem[];
}

interface ActionMenuProps {
  /** 扁平的菜单项 */
  items?: ActionMenuItem[];
  /** 分组的菜单项 */
  groups?: ActionMenuGroup[];
  /** 图标方向 */
  direction?: 'vertical' | 'horizontal';
  /** 尺寸 */
  size?: 'sm' | 'md' | 'lg';
  /** 对齐方式 */
  align?: 'start' | 'center' | 'end';
  /** 自定义触发器 */
  trigger?: React.ReactNode;
  /** 是否禁用 */
  disabled?: boolean;
  /** 自定义类名 */
  className?: string;
}

const sizeClasses = {
  sm: 'h-7 w-7',
  md: 'h-8 w-8',
  lg: 'h-9 w-9',
};

const iconSizes = {
  sm: 'w-3.5 h-3.5',
  md: 'w-4 h-4',
  lg: 'w-5 h-5',
};

function renderMenuItem(item: ActionMenuItem, iconSize: string) {
  if (item.hidden) return null;
  
  const Icon = item.icon;
  
  return (
    <DropdownMenuItem
      key={item.key}
      onClick={item.onClick}
      disabled={item.disabled}
      className={cn(
        'flex items-center gap-2 px-3 py-2 cursor-pointer',
        item.variant === 'destructive' && 'text-destructive focus:text-destructive',
        item.disabled && 'opacity-50 cursor-not-allowed'
      )}
    >
      {Icon && <Icon className={cn(iconSize, 'shrink-0')} />}
      <span>{item.label}</span>
    </DropdownMenuItem>
  );
}

function renderContent(
  items: ActionMenuItem[] | undefined,
  groups: ActionMenuGroup[] | undefined,
  iconSize: string
) {
  // 使用分组模式
  if (groups && groups.length > 0) {
    return (
      <>
        {groups.map((group, groupIndex) => {
          const visibleItems = group.items.filter((item) => !item.hidden);
          if (visibleItems.length === 0) return null;
          
          return (
            <div key={groupIndex}>
              {groupIndex > 0 && <DropdownMenuSeparator className="my-1" />}
              {group.label && (
                <DropdownMenuLabel className="px-3 py-1.5 text-xs text-muted-foreground">
                  {group.label}
                </DropdownMenuLabel>
              )}
              {visibleItems.map((item) => renderMenuItem(item, iconSize))}
            </div>
          );
        })}
      </>
    );
  }
  
  // 使用扁平模式
  if (items && items.length > 0) {
    const visibleItems = items.filter((item) => !item.hidden);
    return <>{visibleItems.map((item) => renderMenuItem(item, iconSize))}</>;
  }
  
  return (
    <div className="px-3 py-6 text-center text-muted-foreground text-sm">
      暂无操作
    </div>
  );
}

export function ActionMenu({
  items,
  groups,
  direction = 'vertical',
  size = 'md',
  align = 'end',
  trigger,
  disabled = false,
  className,
}: ActionMenuProps) {
  const [open, setOpen] = useState(false);
  
  const Icon = direction === 'vertical' ? MoreVertical : MoreHorizontal;

  return (
    <DropdownMenu open={open} onOpenChange={setOpen}>
      <DropdownMenuTrigger asChild>
        {trigger || (
          <Button
            variant="ghost"
            size="icon"
            disabled={disabled}
            className={cn(
              sizeClasses[size],
              'hover:bg-muted/80 rounded-lg',
              'transition-colors duration-200',
              className
            )}
            onClick={(e) => e.stopPropagation()}
          >
            <Icon className={iconSizes[size]} />
            <span className="sr-only">操作菜单</span>
          </Button>
        )}
      </DropdownMenuTrigger>
      
      <DropdownMenuContent
        align={align}
        sideOffset={4}
        className={cn(
          'min-w-[140px] p-1.5 rounded-xl',
          'bg-popover/95 backdrop-blur-xl',
          'border border-border/50',
          'shadow-lg shadow-black/10',
          'animate-in fade-in-0 zoom-in-95 duration-200'
        )}
        onClick={(e) => e.stopPropagation()}
      >
        {renderContent(items, groups, iconSizes[size])}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
