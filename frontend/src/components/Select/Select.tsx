import React, { useState, useRef, useEffect } from 'react';
import styles from './Select.module.less';

export interface SelectOption {
  label: string;
  value: string | number;
  disabled?: boolean;
}

export interface SelectProps {
  value?: string | number;
  defaultValue?: string | number;
  options: SelectOption[];
  placeholder?: string;
  disabled?: boolean;
  onChange?: (value: string | number) => void;
  className?: string;
  style?: React.CSSProperties;
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
}) => {
  const [internalValue, setInternalValue] = useState(value ?? defaultValue);
  const [open, setOpen] = useState(false);
  const selectRef = useRef<HTMLDivElement>(null);

  const currentValue = value !== undefined ? value : internalValue;
  const selectedOption = options.find((opt) => opt.value === currentValue);

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

  const handleToggle = () => {
    if (!disabled) {
      setOpen(!open);
    }
  };

  const handleSelect = (option: SelectOption) => {
    if (option.disabled) return;

    setInternalValue(option.value);
    onChange?.(option.value);
    setOpen(false);
  };

  return (
    <div
      ref={selectRef}
      className={`${styles.select} ${disabled ? styles.disabled : ''} ${open ? styles.open : ''} ${className || ''}`}
      style={style}
    >
      <div className={styles.selector} onClick={handleToggle}>
        <span className={styles.value}>{selectedOption ? selectedOption.label : placeholder}</span>
        <ArrowIcon />
      </div>

      {open && (
        <div className={styles.dropdown}>
          {options.map((option) => (
            <div
              key={option.value}
              className={`${styles.option} ${option.disabled ? styles.optionDisabled : ''} ${
                option.value === currentValue ? styles.optionSelected : ''
              }`}
              onClick={() => handleSelect(option)}
            >
              {option.label}
              {option.value === currentValue && <CheckIcon />}
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
