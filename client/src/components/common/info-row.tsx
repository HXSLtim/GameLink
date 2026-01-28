import { forwardRef } from 'react';
import { cn } from '@/lib/utils';
import { HStack, VStack } from '@/components/layout';
import type { LucideIcon } from 'lucide-react';

export interface InfoRowProps extends React.HTMLAttributes<HTMLDivElement> {
  /** 标签 */
  label: string;
  /** 值 */
  value?: React.ReactNode;
  /** 图标 */
  icon?: LucideIcon;
  /** 尺寸 */
  size?: 'sm' | 'md' | 'lg';
  /** 布局方向 */
  direction?: 'horizontal' | 'vertical';
  /** 是否显示冒号 */
  colon?: boolean;
  /** 是否显示分隔线 */
  divider?: boolean;
  /** 值为空时显示的内容 */
  emptyText?: string;
}

const sizeConfig = {
  sm: { label: 'text-xs', value: 'text-xs', icon: 'w-3.5 h-3.5', gap: 1.5 },
  md: { label: 'text-sm', value: 'text-sm', icon: 'w-4 h-4', gap: 2 },
  lg: { label: 'text-base', value: 'text-base', icon: 'w-5 h-5', gap: 2.5 },
};

export const InfoRow = forwardRef<HTMLDivElement, InfoRowProps>(({
  label,
  value,
  icon: Icon,
  size = 'md',
  direction = 'horizontal',
  colon = false,
  divider = false,
  emptyText = '-',
  className,
  children,
  ...props
}, ref) => {
  const config = sizeConfig[size];
  const displayValue = value ?? children ?? emptyText;

  if (direction === 'vertical') {
    return (
      <VStack
        ref={ref}
        spacing={1}
        className={cn(
          divider && 'pb-3 border-b border-border/30',
          className
        )}
        {...props}
      >
        <HStack spacing={1} align="center">
          {Icon && <Icon className={cn(config.icon, 'text-muted-foreground')} />}
          <span className={cn(config.label, 'text-muted-foreground')}>
            {label}{colon && '：'}
          </span>
        </HStack>
        <div className={cn(config.value, 'text-foreground font-medium')}>
          {displayValue}
        </div>
      </VStack>
    );
  }

  return (
    <HStack
      ref={ref}
      spacing={config.gap}
      align="center"
      className={cn(
        'justify-between',
        divider && 'pb-3 border-b border-border/30',
        className
      )}
      {...props}
    >
      <HStack spacing={1.5} align="center" className="shrink-0">
        {Icon && <Icon className={cn(config.icon, 'text-muted-foreground')} />}
        <span className={cn(config.label, 'text-muted-foreground')}>
          {label}{colon && '：'}
        </span>
      </HStack>
      <div className={cn(config.value, 'text-foreground font-medium text-right truncate')}>
        {displayValue}
      </div>
    </HStack>
  );
});

InfoRow.displayName = 'InfoRow';

// InfoList 组件 - 用于显示多个 InfoRow
export interface InfoListProps extends React.HTMLAttributes<HTMLDivElement> {
  /** 信息项列表 */
  items: Array<Omit<InfoRowProps, 'ref'>>;
  /** 尺寸 */
  size?: 'sm' | 'md' | 'lg';
  /** 布局方向 */
  direction?: 'horizontal' | 'vertical';
  /** 是否显示分隔线 */
  divider?: boolean;
  /** 列数（仅 vertical 方向有效） */
  columns?: 1 | 2 | 3 | 4;
}

const columnClasses = {
  1: 'grid-cols-1',
  2: 'grid-cols-2',
  3: 'grid-cols-3',
  4: 'grid-cols-4',
};

export const InfoList = forwardRef<HTMLDivElement, InfoListProps>(({
  items,
  size = 'md',
  direction = 'horizontal',
  divider = false,
  columns = 1,
  className,
  ...props
}, ref) => {
  if (direction === 'vertical' && columns > 1) {
    return (
      <div
        ref={ref}
        className={cn('grid gap-4', columnClasses[columns], className)}
        {...props}
      >
        {items.map((item, index) => (
          <InfoRow
            key={index}
            {...item}
            size={item.size ?? size}
            direction={direction}
            divider={false}
          />
        ))}
      </div>
    );
  }

  return (
    <VStack
      ref={ref}
      spacing={divider ? 0 : 3}
      className={cn(divider && 'divide-y divide-border/30', className)}
      {...props}
    >
      {items.map((item, index) => (
        <InfoRow
          key={index}
          {...item}
          size={item.size ?? size}
          direction={direction}
          divider={false}
          className={divider ? 'py-3 first:pt-0 last:pb-0' : undefined}
        />
      ))}
    </VStack>
  );
});

InfoList.displayName = 'InfoList';
