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
import { ChevronDown, Check, X, Loader2 } from 'lucide-react';
import { cn } from '@/lib/utils';
import type { LucideIcon } from 'lucide-react';

export interface DropdownOption {
  value: string | number;
  label: string;
  icon?: LucideIcon;
  description?: string;
  disabled?: boolean;
}

export interface DropdownSelectProps {
  /** 选项列表 */
  options: DropdownOption[];
  /** 当前值 */
  value: string | number | null;
  /** 值变化回调 */
  onChange: (value: string | number | null) => void;
  /** 占位文本 */
  placeholder?: string;
  /** 标签 */
  label?: string;
  /** 前置图标 */
  icon?: LucideIcon;
  /** 尺寸 */
  size?: 'sm' | 'md' | 'lg';
  /** 变体 */
  variant?: 'default' | 'outline' | 'ghost' | 'pill' | 'glass';
  /** 是否禁用 */
  disabled?: boolean;
  /** 是否加载中 */
  loading?: boolean;
  /** 是否可清除 */
  clearable?: boolean;
  /** 清除文本 */
  clearText?: string;
  /** 对齐方式 */
  align?: 'start' | 'center' | 'end';
  /** 最小宽度 */
  minWidth?: number;
  /** 自定义类名 */
  className?: string;
}

const sizeClasses = {
  sm: 'h-8 px-3 text-xs gap-1.5',
  md: 'h-9 px-4 text-sm gap-2',
  lg: 'h-11 px-5 text-base gap-2.5',
};

const iconSizes = {
  sm: 'w-3.5 h-3.5',
  md: 'w-4 h-4',
  lg: 'w-5 h-5',
};

const variantClasses = {
  default: cn(
    'bg-background/80 backdrop-blur-sm border border-border/50',
    'hover:bg-accent/50 hover:border-primary/30',
    'shadow-sm hover:shadow-md',
    'transition-all duration-200'
  ),
  outline: cn(
    'bg-transparent border border-border',
    'hover:bg-accent hover:border-primary/50',
    'transition-all duration-200'
  ),
  ghost: cn(
    'bg-transparent border-transparent',
    'hover:bg-accent',
    'transition-all duration-200'
  ),
  pill: cn(
    'bg-muted/50 backdrop-blur-sm border-0 rounded-full',
    'hover:bg-muted',
    'transition-all duration-200'
  ),
  glass: cn(
    'bg-white/5 backdrop-blur-xl border border-white/10',
    'hover:bg-white/10 hover:border-white/20',
    'shadow-lg shadow-black/5',
    'transition-all duration-200'
  ),
};

export function DropdownSelect({
  options,
  value,
  onChange,
  placeholder = '请选择',
  label,
  icon: Icon,
  size = 'md',
  variant = 'default',
  disabled = false,
  loading = false,
  clearable = true,
  clearText = '清除选择',
  align = 'start',
  minWidth = 180,
  className,
}: DropdownSelectProps) {
  const [open, setOpen] = useState(false);

  const selectedOption = options.find((opt) => opt.value === value);
  const hasValue = value !== null && value !== undefined;

  const getContentWidth = () => {
    return minWidth ? `${minWidth}px` : 'auto';
  };

  return (
    <DropdownMenu open={open} onOpenChange={setOpen}>
      <DropdownMenuTrigger asChild>
        <Button
          variant="ghost"
          size={size === 'lg' ? 'default' : 'sm'}
          disabled={disabled || loading}
          className={cn(
            'font-normal justify-between rounded-lg',
            sizeClasses[size],
            variantClasses[variant],
            hasValue && 'text-foreground',
            !hasValue && 'text-muted-foreground',
            disabled && 'opacity-50 cursor-not-allowed',
            className
          )}
        >
          <span className="flex items-center gap-2 truncate">
            {loading ? (
              <Loader2 className={cn(iconSizes[size], 'shrink-0 animate-spin')} />
            ) : Icon ? (
              <Icon className={cn(iconSizes[size], 'shrink-0 text-muted-foreground')} />
            ) : null}
            
            {selectedOption ? (
              <span className="flex items-center gap-2 truncate">
                {selectedOption.icon && (
                  <selectedOption.icon className={cn(iconSizes[size], 'shrink-0')} />
                )}
                <span className="truncate">{selectedOption.label}</span>
              </span>
            ) : (
              <span className="truncate">{placeholder}</span>
            )}
          </span>
          
          <ChevronDown className={cn(
            iconSizes[size],
            'shrink-0 text-muted-foreground/70 transition-transform duration-200',
            open && 'rotate-180'
          )} />
        </Button>
      </DropdownMenuTrigger>
      
      <DropdownMenuContent
        align={align}
        sideOffset={8}
        className={cn(
          'p-2 rounded-xl',
          'bg-gradient-to-b from-popover via-popover to-popover/95',
          'backdrop-blur-xl',
          'border border-white/10',
          'shadow-xl shadow-black/20',
          'max-h-[320px] overflow-y-auto overflow-x-hidden',
          'animate-in fade-in-0 zoom-in-95 slide-in-from-top-2 duration-200',
          '[&::-webkit-scrollbar]:w-1.5',
          '[&::-webkit-scrollbar-track]:bg-transparent',
          '[&::-webkit-scrollbar-thumb]:bg-border/50',
          '[&::-webkit-scrollbar-thumb]:rounded-full'
        )}
        style={{ minWidth: getContentWidth() }}
      >
        {label && (
          <DropdownMenuLabel className="px-3 py-2 text-xs text-muted-foreground font-medium">
            {label}
          </DropdownMenuLabel>
        )}
        
        {clearable && hasValue && (
          <>
            <DropdownMenuItem
              onClick={() => onChange(null)}
              className={cn(
                'group flex items-center gap-2 px-3 py-2 rounded-lg cursor-pointer',
                'text-muted-foreground hover:text-destructive',
                'hover:bg-destructive/10',
                'transition-all duration-200'
              )}
            >
              <X className="w-4 h-4" />
              <span className="text-sm">{clearText}</span>
            </DropdownMenuItem>
            <DropdownMenuSeparator className="my-1.5 mx-2 bg-border/50" />
          </>
        )}
        
        {loading ? (
          <div className="flex items-center justify-center py-8">
            <Loader2 className="w-5 h-5 animate-spin text-muted-foreground" />
          </div>
        ) : options.length > 0 ? (
          <div className="space-y-0.5">
            {options.map((option) => {
              const isSelected = value === option.value;
              return (
                <DropdownMenuItem
                  key={option.value}
                  onClick={() => onChange(option.value)}
                  disabled={option.disabled}
                  className={cn(
                    'group flex items-center gap-3 px-3 py-2.5 rounded-lg cursor-pointer',
                    'transition-all duration-200',
                    isSelected
                      ? 'bg-primary/10 text-primary'
                      : 'hover:bg-muted/80',
                    option.disabled && 'opacity-50 cursor-not-allowed'
                  )}
                >
                  {/* 图标或首字母 */}
                  {option.icon ? (
                    <div className={cn(
                      'flex items-center justify-center w-8 h-8 rounded-lg shrink-0',
                      isSelected ? 'bg-primary/15' : 'bg-muted/60'
                    )}>
                      <option.icon className={cn(
                        'w-4 h-4',
                        isSelected ? 'text-primary' : 'text-muted-foreground'
                      )} />
                    </div>
                  ) : (
                    <div className={cn(
                      'flex items-center justify-center w-8 h-8 rounded-lg shrink-0 text-sm font-medium',
                      isSelected
                        ? 'bg-primary text-primary-foreground'
                        : 'bg-muted/60 text-muted-foreground'
                    )}>
                      {option.label.charAt(0).toUpperCase()}
                    </div>
                  )}
                  
                  {/* 文本内容 */}
                  <div className="flex-1 min-w-0">
                    <div className={cn(
                      'text-sm font-medium truncate',
                      isSelected ? 'text-primary' : 'text-foreground'
                    )}>
                      {option.label}
                    </div>
                    {option.description && (
                      <div className="text-xs text-muted-foreground/80 truncate mt-0.5">
                        {option.description}
                      </div>
                    )}
                  </div>
                  
                  {/* 选中标记 */}
                  {isSelected && (
                    <div className="flex items-center justify-center w-5 h-5 rounded-full bg-primary shrink-0">
                      <Check className="w-3 h-3 text-primary-foreground" />
                    </div>
                  )}
                </DropdownMenuItem>
              );
            })}
          </div>
        ) : (
          <div className="py-8 text-center text-muted-foreground text-sm">
            暂无选项
          </div>
        )}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
