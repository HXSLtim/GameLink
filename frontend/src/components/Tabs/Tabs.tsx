import React, { useState, useRef, useEffect } from 'react';
import styles from './Tabs.module.less';

export interface TabItem {
  key: string;
  label: React.ReactNode;
  disabled?: boolean;
  children: React.ReactNode;
}

export interface TabsProps {
  items: TabItem[];
  defaultActiveKey?: string;
  activeKey?: string;
  onChange?: (key: string) => void;
  className?: string;
  style?: React.CSSProperties;
}

export const Tabs: React.FC<TabsProps> = ({
  items,
  defaultActiveKey,
  activeKey,
  onChange,
  className = '',
  style,
}) => {
  const [internalKey, setInternalKey] = useState<string>(() => {
    return activeKey || defaultActiveKey || items.find((i) => !i.disabled)?.key || items[0]?.key;
  });

  useEffect(() => {
    if (activeKey) setInternalKey(activeKey);
  }, [activeKey]);

  const listRef = useRef<HTMLDivElement>(null);
  const focusTabByIndex = (index: number) => {
    const tabs = listRef.current?.querySelectorAll<HTMLButtonElement>('button[role="tab"]');
    const target = tabs?.[index];
    target?.focus();
  };

  const currentIndex = items.findIndex((i) => i.key === internalKey);

  const handleKeyDown = (e: React.KeyboardEvent) => {
    const enabledItems = items.filter((i) => !i.disabled);
    const enabledKeys = enabledItems.map((i) => i.key);
    const enabledIndex = enabledKeys.indexOf(internalKey);

    switch (e.key) {
      case 'ArrowRight':
      case 'ArrowDown': {
        const next = enabledKeys[(enabledIndex + 1) % enabledKeys.length];
        setInternalKey(next);
        onChange?.(next);
        focusTabByIndex(items.findIndex((i) => i.key === next));
        e.preventDefault();
        break;
      }
      case 'ArrowLeft':
      case 'ArrowUp': {
        const prev = enabledKeys[(enabledIndex - 1 + enabledKeys.length) % enabledKeys.length];
        setInternalKey(prev);
        onChange?.(prev);
        focusTabByIndex(items.findIndex((i) => i.key === prev));
        e.preventDefault();
        break;
      }
      case 'Home': {
        const first = enabledKeys[0];
        setInternalKey(first);
        onChange?.(first);
        focusTabByIndex(items.findIndex((i) => i.key === first));
        e.preventDefault();
        break;
      }
      case 'End': {
        const last = enabledKeys[enabledKeys.length - 1];
        setInternalKey(last);
        onChange?.(last);
        focusTabByIndex(items.findIndex((i) => i.key === last));
        e.preventDefault();
        break;
      }
      default:
        break;
    }
  };

  const classNames = [styles.tabs, className].filter(Boolean).join(' ');

  return (
    <div className={classNames} style={style}>
      <div
        className={styles.tabList}
        role="tablist"
        aria-orientation="horizontal"
        onKeyDown={handleKeyDown}
        ref={listRef}
      >
        {items.map((item, index) => {
          const selected = item.key === internalKey;
          return (
            <button
              key={item.key}
              role="tab"
              type="button"
              className={`${styles.tab} ${selected ? styles.active : ''} ${item.disabled ? styles.disabled : ''}`}
              aria-selected={selected}
              aria-controls={`panel-${item.key}`}
              id={`tab-${item.key}`}
              disabled={item.disabled}
              tabIndex={selected ? 0 : -1}
              onClick={() => {
                if (!item.disabled) {
                  setInternalKey(item.key);
                  onChange?.(item.key);
                }
              }}
            >
              {item.label}
            </button>
          );
        })}
      </div>

      {items.map((item) => {
        const selected = item.key === internalKey;
        return (
          <div
            key={item.key}
            role="tabpanel"
            id={`panel-${item.key}`}
            aria-labelledby={`tab-${item.key}`}
            className={`${styles.tabPanel} ${selected ? styles.visible : styles.hidden}`}
          >
            {selected && item.children}
          </div>
        );
      })}
    </div>
  );
};