import React, { useState, useRef, useEffect } from 'react';
import styles from './Select.module.less';

export interface SelectOption {
  label: string;
  value: string | number;
  disabled?: boolean;
}

export type SelectValue = string | number | Array<string | number>;

export interface SelectProps {
  value?: SelectValue;
  defaultValue?: SelectValue;
  options: SelectOption[];
  placeholder?: string;
  disabled?: boolean;
  onChange?: (value: SelectValue) => void;
  className?: string;
  style?: React.CSSProperties;
  // 多选
  multiple?: boolean;
  // 搜索过滤
  searchable?: boolean;
  filterOption?: (input: string, option: SelectOption) => boolean;
  onSearch?: (query: string) => void;
  // 异步搜索加载
  asyncSearch?: (query: string) => Promise<SelectOption[]>;
  loading?: boolean;
}

export const Select: React.FC<SelectProps> = ({
  value,
  defaultValue,
  options,
  placeholder = '请选择',
  disabled = false,
  onChange,
  className,
  style,
  multiple = false,
  searchable = false,
  filterOption,
  onSearch,
  asyncSearch,
  loading = false,
}) => {
  const [internalValue, setInternalValue] = useState<SelectValue | undefined>(value ?? defaultValue);
  const [open, setOpen] = useState(false);
  const [search, setSearch] = useState('');
  const [internalOptions, setInternalOptions] = useState<SelectOption[]>(options);
  const [highlightIndex, setHighlightIndex] = useState<number>(-1);
  const selectRef = useRef<HTMLDivElement>(null);

  // 当前选中值（单选）或选中列表（多选）
  const currentValue = value !== undefined ? value : internalValue;
  const currentValuesArray: Array<string | number> = Array.isArray(currentValue)
    ? currentValue
    : currentValue !== undefined && currentValue !== null
    ? [currentValue]
    : [];
  const selectedLabelText = multiple
    ? currentValuesArray
        .map((val) => internalOptions.find((opt) => opt.value === val)?.label)
        .filter(Boolean)
        .join(', ')
    : internalOptions.find((opt) => opt.value === currentValue)?.label;

  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (selectRef.current && !selectRef.current.contains(e.target as Node)) {
        setOpen(false);
      }
    };

    if (open) {
      document.addEventListener('mousedown', handleClickOutside);
    }

    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, [open]);

  // 当 props.options 变更时，同步内部 options（除非异步搜索正在接管）
  useEffect(() => {
    setInternalOptions(options);
  }, [options]);

  // 处理键盘导航（↑↓选择，Enter 确认，Esc 关闭）
  useEffect(() => {
    const el = selectRef.current;
    if (!el) return;

    const onKeyDown = (e: KeyboardEvent) => {
      if (!open) return;
      const displayed = getDisplayedOptions();
      if (e.key === 'ArrowDown') {
        e.preventDefault();
        setHighlightIndex((prev) => Math.min(displayed.length - 1, prev + 1));
      } else if (e.key === 'ArrowUp') {
        e.preventDefault();
        setHighlightIndex((prev) => Math.max(0, prev - 1));
      } else if (e.key === 'Enter') {
        e.preventDefault();
        const option = displayed[highlightIndex];
        if (option) handleSelect(option);
      } else if (e.key === 'Escape') {
        e.preventDefault();
        setOpen(false);
      }
    };

    el.addEventListener('keydown', onKeyDown);
    return () => el.removeEventListener('keydown', onKeyDown);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [open, highlightIndex, internalOptions, search, multiple, currentValuesArray]);

  const handleToggle = () => {
    if (!disabled) {
      setOpen(!open);
      // 打开时重置高亮
      if (!open) setHighlightIndex(0);
    }
  };

  const handleSelect = (option: SelectOption) => {
    if (option.disabled) return;
    if (multiple) {
      const exists = currentValuesArray.includes(option.value);
      const next = exists
        ? currentValuesArray.filter((v) => v !== option.value)
        : [...currentValuesArray, option.value];
      setInternalValue(next);
      onChange?.(next);
    } else {
      setInternalValue(option.value);
      onChange?.(option.value);
      setOpen(false);
    }
  };

  const getDisplayedOptions = () => {
    const base = internalOptions;
    if (!searchable) return base;
    if (filterOption) return base.filter((opt) => filterOption(search, opt));
    if (!search) return base;
    const lowered = search.toLowerCase();
    return base.filter((opt) => String(opt.label).toLowerCase().includes(lowered));
  };

  const handleSearchChange = async (val: string) => {
    setSearch(val);
    onSearch?.(val);
    if (asyncSearch) {
      try {
        const result = await asyncSearch(val);
        setInternalOptions(result);
        setHighlightIndex(0);
      } catch (e) {
        // swallow errors to keep UX smooth
      }
    }
  };

  return (
    <div
      ref={selectRef}
      className={`${styles.select} ${disabled ? styles.disabled : ''} ${open ? styles.open : ''} ${className || ''}`}
      style={style}
      role="combobox"
      aria-expanded={open}
      aria-controls="select-listbox"
      aria-haspopup="listbox"
      aria-disabled={disabled}
      aria-multiselectable={multiple || undefined}
      tabIndex={0}
    >
      <div className={styles.selector} onClick={handleToggle}>
        <span className={styles.value}>{selectedLabelText || placeholder}</span>
        <ArrowIcon />
      </div>

      {open && (
        <div className={styles.dropdown} role="listbox" id="select-listbox" aria-multiselectable={multiple || undefined}>
          {searchable && (
            <div className={styles.searchBar}>
              <input
                className={styles.searchInput}
                type="text"
                value={search}
                onChange={(e) => handleSearchChange(e.target.value)}
                placeholder="搜索..."
                aria-label="搜索选项"
                autoFocus
              />
            </div>
          )}
          {(loading && !asyncSearch) && (
            <div className={styles.loading}>加载中...</div>
          )}
          {getDisplayedOptions().map((option, idx) => (
            <div
              key={option.value}
              className={`${styles.option} ${option.disabled ? styles.optionDisabled : ''} ${
                (multiple
                  ? currentValuesArray.includes(option.value)
                  : option.value === currentValue)
                  ? styles.optionSelected
                  : ''
              }`}
              onClick={() => handleSelect(option)}
              role="option"
              aria-selected={multiple ? currentValuesArray.includes(option.value) : option.value === currentValue}
              aria-disabled={option.disabled}
              data-highlighted={highlightIndex === idx || undefined}
            >
              {option.label}
              {(multiple
                ? currentValuesArray.includes(option.value)
                : option.value === currentValue) && <CheckIcon />}
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

const ArrowIcon = () => (
  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor">
    <path strokeLinecap="square" strokeLinejoin="miter" strokeWidth="2" d="M6 9l6 6 6-6" />
  </svg>
);

const CheckIcon = () => (
  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor">
    <path strokeLinecap="square" strokeLinejoin="miter" strokeWidth="2" d="M5 13l4 4L19 7" />
  </svg>
);
