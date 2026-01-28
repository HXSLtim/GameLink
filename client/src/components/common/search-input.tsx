import { forwardRef, useState, useCallback } from 'react';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Search, X, Loader2 } from 'lucide-react';
import { cn } from '@/lib/utils';

export interface SearchInputProps extends Omit<React.InputHTMLAttributes<HTMLInputElement>, 'onChange' | 'size'> {
  /** 值变化回调 */
  onChange?: (value: string) => void;
  /** 搜索回调 */
  onSearch?: (value: string) => void;
  /** 清除回调 */
  onClear?: () => void;
  /** 是否加载中 */
  loading?: boolean;
  /** 尺寸 */
  size?: 'sm' | 'md' | 'lg';
  /** 变体 */
  variant?: 'default' | 'filled' | 'ghost';
  /** 是否显示搜索按钮 */
  showSearchButton?: boolean;
  /** 是否显示清除按钮 */
  showClearButton?: boolean;
  /** 搜索按钮文本 */
  searchButtonText?: string;
  /** 防抖延迟（毫秒） */
  debounceMs?: number;
}

const sizeClasses = {
  sm: {
    wrapper: 'h-8',
    input: 'h-8 text-xs px-8',
    icon: 'w-3.5 h-3.5',
    iconLeft: 'left-2.5',
    iconRight: 'right-2.5',
    button: 'h-6 px-2 text-xs',
  },
  md: {
    wrapper: 'h-9',
    input: 'h-9 text-sm px-9',
    icon: 'w-4 h-4',
    iconLeft: 'left-3',
    iconRight: 'right-3',
    button: 'h-7 px-3 text-sm',
  },
  lg: {
    wrapper: 'h-11',
    input: 'h-11 text-base px-10',
    icon: 'w-5 h-5',
    iconLeft: 'left-3.5',
    iconRight: 'right-3.5',
    button: 'h-8 px-4 text-base',
  },
};

const variantClasses = {
  default: 'bg-background border border-border focus-within:border-primary/50',
  filled: 'bg-muted/50 border-0 focus-within:bg-muted',
  ghost: 'bg-transparent border-0 focus-within:bg-muted/50',
};

export const SearchInput = forwardRef<HTMLInputElement, SearchInputProps>(({
  value,
  onChange,
  onSearch,
  onClear,
  loading = false,
  size = 'md',
  variant = 'default',
  showSearchButton = false,
  showClearButton = true,
  searchButtonText = '搜索',
  placeholder = '搜索...',
  className,
  ...props
}, ref) => {
  const [internalValue, setInternalValue] = useState(value || '');
  const config = sizeClasses[size];
  
  const currentValue = value !== undefined ? String(value) : internalValue;
  const hasValue = currentValue.length > 0;

  const handleChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    const newValue = e.target.value;
    setInternalValue(newValue);
    onChange?.(newValue);
  }, [onChange]);

  const handleClear = useCallback(() => {
    setInternalValue('');
    onChange?.('');
    onClear?.();
  }, [onChange, onClear]);

  const handleKeyDown = useCallback((e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter' && !e.nativeEvent.isComposing) {
      e.preventDefault();
      onSearch?.(currentValue);
    }
    if (e.key === 'Escape') {
      handleClear();
    }
  }, [currentValue, onSearch, handleClear]);

  const handleSearch = useCallback(() => {
    onSearch?.(currentValue);
  }, [currentValue, onSearch]);

  return (
    <div className={cn(
      'relative flex items-center gap-2',
      className
    )}>
      <div className={cn(
        'relative flex-1 rounded-lg overflow-hidden transition-colors duration-200',
        variantClasses[variant],
        config.wrapper
      )}>
        {/* 搜索图标 */}
        <div className={cn(
          'absolute top-1/2 -translate-y-1/2 pointer-events-none',
          config.iconLeft
        )}>
          {loading ? (
            <Loader2 className={cn(config.icon, 'text-muted-foreground animate-spin')} />
          ) : (
            <Search className={cn(config.icon, 'text-muted-foreground')} />
          )}
        </div>
        
        {/* 输入框 */}
        <Input
          ref={ref}
          type="text"
          value={currentValue}
          onChange={handleChange}
          onKeyDown={handleKeyDown}
          placeholder={placeholder}
          className={cn(
            'border-0 bg-transparent focus-visible:ring-0 focus-visible:ring-offset-0',
            config.input,
            hasValue && showClearButton && 'pr-9'
          )}
          {...props}
        />
        
        {/* 清除按钮 */}
        {showClearButton && hasValue && !loading && (
          <Button
            type="button"
            variant="ghost"
            size="icon"
            onClick={handleClear}
            className={cn(
              'absolute top-1/2 -translate-y-1/2 h-6 w-6',
              'hover:bg-muted/80 rounded-md',
              config.iconRight
            )}
          >
            <X className={config.icon} />
            <span className="sr-only">清除</span>
          </Button>
        )}
      </div>
      
      {/* 搜索按钮 */}
      {showSearchButton && (
        <Button
          type="button"
          onClick={handleSearch}
          disabled={loading || !hasValue}
          className={config.button}
        >
          {loading ? (
            <Loader2 className={cn(config.icon, 'animate-spin')} />
          ) : (
            searchButtonText
          )}
        </Button>
      )}
    </div>
  );
});

SearchInput.displayName = 'SearchInput';
