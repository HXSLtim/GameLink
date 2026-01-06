/**
 * 可折叠区块组件
 * 用于仪表盘等页面的渐进式披露
 */
import React, { useState } from 'react';
import type { ReactNode } from 'react';
import { Typography, Button, Space } from 'antd';
import { DownOutlined, UpOutlined } from '@ant-design/icons';
import { motion, AnimatePresence } from 'framer-motion';

const { Title } = Typography;

export interface CollapsibleSectionProps {
  /** 区块 ID，用于持久化 */
  id: string;
  /** 区块标题 */
  title: string;
  /** 子元素 */
  children: ReactNode;
  /** 是否默认折叠 */
  defaultCollapsed?: boolean;
  /** 折叠状态（受控） */
  collapsed?: boolean;
  /** 折叠状态变化回调 */
  onCollapsedChange?: (collapsed: boolean) => void;
  /** 额外操作 */
  extra?: ReactNode;
  /** 标题级别 */
  level?: 1 | 2 | 3 | 4 | 5;
}

/**
 * 可折叠区块
 */
export const CollapsibleSection: React.FC<CollapsibleSectionProps> = ({
  id,
  title,
  children,
  defaultCollapsed = false,
  collapsed: controlledCollapsed,
  onCollapsedChange,
  extra,
  level = 5,
}) => {
  const [internalCollapsed, setInternalCollapsed] = useState(defaultCollapsed);

  // 支持受控和非受控模式
  const isCollapsed = controlledCollapsed ?? internalCollapsed;

  const handleToggle = () => {
    const newValue = !isCollapsed;
    setInternalCollapsed(newValue);
    onCollapsedChange?.(newValue);
  };

  return (
    <section aria-labelledby={`section-${id}`} style={{ marginTop: 16 }}>
      <div
        style={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          marginBottom: isCollapsed ? 0 : 16,
        }}
      >
        <Space>
          <Title
            id={`section-${id}`}
            level={level}
            style={{ margin: 0, cursor: 'pointer' }}
            onClick={handleToggle}
          >
            {title}
          </Title>
          <Button
            type="text"
            size="small"
            icon={isCollapsed ? <DownOutlined /> : <UpOutlined />}
            onClick={handleToggle}
            aria-expanded={!isCollapsed}
            aria-controls={`content-${id}`}
            aria-label={isCollapsed ? '展开' : '折叠'}
          />
        </Space>
        {extra && <div>{extra}</div>}
      </div>

      <AnimatePresence initial={false}>
        {!isCollapsed && (
          <motion.div
            id={`content-${id}`}
            initial={{ height: 0, opacity: 0 }}
            animate={{ height: 'auto', opacity: 1 }}
            exit={{ height: 0, opacity: 0 }}
            transition={{ duration: 0.2, ease: 'easeInOut' }}
            style={{ overflow: 'hidden' }}
          >
            {children}
          </motion.div>
        )}
      </AnimatePresence>
    </section>
  );
};

export default CollapsibleSection;
