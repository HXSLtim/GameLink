import { forwardRef } from 'react';
import { Badge } from '@/components/ui/badge';
import { DropdownSelect, type DropdownOption } from '@/components/common/dropdown-select';
import { cn } from '@/lib/utils';
import type { LucideIcon } from 'lucide-react';
import { ArrowUpDown, ArrowUp, ArrowDown } from 'lucide-react';

// Badge 类型的筛选选项
export interface FilterOption {
  id: string | number;
  label: string;
  icon?: LucideIcon;
}

// Dropdown 类型的筛选选项
export interface DropdownFilterOption {
  id: string | number;
  label: string;
  icon?: LucideIcon;
  description?: string;
}

// 排序选项
export interface SortOption {
  id: string;
  label: string;
  icon?: LucideIcon;
}

// 筛选组
export interface FilterGroup {
  key: string;
  label?: string;
  type: 'badge' | 'dropdown' | 'sort';
  // Badge 类型配置
  options?: FilterOption[];
  selectedId?: string | number | null;
  // Dropdown 类型配置
  dropdownOptions?: DropdownFilterOption[];
  placeholder?: string;
  icon?: LucideIcon;
  loading?: boolean;
  // Sort 类型配置
  sortOptions?: SortOption[];
  sortOrder?: 'asc' | 'desc' | null;
  selectedSortId?: string | null;
}

interface FilterBarProps {
  /** 筛选组 */
  groups: FilterGroup[];
  /** 筛选变化回调 */
  onFilterChange: (groupKey: string, value: string | number | null) => void;
  /** 排序变化回调 */
  onSortChange?: (sortId: string | null, order: 'asc' | 'desc' | null) => void;
  /** 变体 */
  variant?: 'default' | 'pill' | 'card';
  /** 尺寸 */
  size?: 'sm' | 'md';
  /** 自定义类名 */
  className?: string;
}

// 转换为 DropdownSelect 的选项格式
function convertToDropdownOptions(options: DropdownFilterOption[]): DropdownOption[] {
  return options.map(opt => ({
    value: opt.id,
    label: opt.label,
    icon: opt.icon,
    description: opt.description,
  }));
}

// 获取排序图标
function getSortIcon(order: 'asc' | 'desc' | null) {
  if (order === 'asc') return ArrowUp;
  if (order === 'desc') return ArrowDown;
  return ArrowUpDown;
}

const variantClasses = {
  default: 'bg-muted/30 backdrop-blur-sm border border-border/30',
  pill: 'bg-muted/40 backdrop-blur-md border border-white/5',
  card: 'bg-card border border-border shadow-sm',
};

export const FilterBar = forwardRef<HTMLDivElement, FilterBarProps>(({
  groups,
  onFilterChange,
  onSortChange,
  variant = 'pill',
  size = 'sm',
  className,
}, ref) => {
  const handleSortClick = (group: FilterGroup, sortId: string) => {
    if (!onSortChange) return;
    
    // 如果点击的是当前已选中的排序项
    if (group.selectedSortId === sortId) {
      // 切换排序方向: null -> asc -> desc -> null
      const currentOrder = group.sortOrder;
      if (currentOrder === null) {
        onSortChange(sortId, 'asc');
      } else if (currentOrder === 'asc') {
        onSortChange(sortId, 'desc');
      } else {
        onSortChange(null, null);
      }
    } else {
      // 选择新的排序项，默认升序
      onSortChange(sortId, 'asc');
    }
  };

  return (
    <div
      ref={ref}
      className={cn(
        'flex items-center gap-1 sm:gap-2 p-1 sm:p-1.5 rounded-full',
        variantClasses[variant],
        className
      )}
    >
      {groups.map((group, groupIndex) => (
        <div key={group.key} className="contents">
          {/* 分隔线 */}
          {groupIndex > 0 && (
            <div className="w-px h-5 bg-border/50 mx-0.5 sm:mx-1" />
          )}
          
          {/* Badge 类型 */}
          {group.type === 'badge' && group.options && (
            <div className="flex items-center gap-1">
              {group.options.map((option) => {
                const isSelected = group.selectedId === option.id;
                const Icon = option.icon;
                return (
                  <Badge
                    key={option.id}
                    variant={isSelected ? 'default' : 'secondary'}
                    className={cn(
                      'cursor-pointer transition-all duration-200 select-none',
                      size === 'sm' ? 'px-2.5 py-1 text-xs' : 'px-3 py-1.5 text-sm',
                      isSelected
                        ? 'bg-primary text-primary-foreground shadow-sm'
                        : 'bg-transparent hover:bg-muted text-muted-foreground hover:text-foreground'
                    )}
                    onClick={() => onFilterChange(group.key, isSelected ? null : option.id)}
                  >
                    {Icon && <Icon className={cn('mr-1', size === 'sm' ? 'w-3 h-3' : 'w-3.5 h-3.5')} />}
                    {option.label}
                  </Badge>
                );
              })}
            </div>
          )}
          
          {/* Dropdown 类型 */}
          {group.type === 'dropdown' && (
            <DropdownSelect
              options={convertToDropdownOptions(group.dropdownOptions || [])}
              value={group.selectedId ?? null}
              onChange={(value) => onFilterChange(group.key, value)}
              placeholder={group.placeholder || '选择'}
              icon={group.icon}
              loading={group.loading}
              variant="pill"
              size={size}
              clearText="清除筛选"
            />
          )}
          
          {/* Sort 类型 */}
          {group.type === 'sort' && group.sortOptions && (
            <div className="flex items-center gap-1">
              {group.sortOptions.map((option) => {
                const isSelected = group.selectedSortId === option.id;
                const SortIcon = isSelected ? getSortIcon(group.sortOrder ?? null) : (option.icon || ArrowUpDown);
                return (
                  <Badge
                    key={option.id}
                    variant={isSelected ? 'default' : 'secondary'}
                    className={cn(
                      'cursor-pointer transition-all duration-200 select-none',
                      size === 'sm' ? 'px-2.5 py-1 text-xs' : 'px-3 py-1.5 text-sm',
                      isSelected
                        ? 'bg-primary text-primary-foreground shadow-sm'
                        : 'bg-transparent hover:bg-muted text-muted-foreground hover:text-foreground'
                    )}
                    onClick={() => handleSortClick(group, option.id)}
                  >
                    <SortIcon className={cn('mr-1', size === 'sm' ? 'w-3 h-3' : 'w-3.5 h-3.5')} />
                    {option.label}
                  </Badge>
                );
              })}
            </div>
          )}
        </div>
      ))}
    </div>
  );
});

FilterBar.displayName = 'FilterBar';
