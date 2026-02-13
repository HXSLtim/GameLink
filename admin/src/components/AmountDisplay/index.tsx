/**
 * 金额展示组件
 * 支持分转元、货币符号、格式化等功能
 */
import React from 'react';
import { Typography, Tooltip } from 'antd';
import { ArrowUpOutlined, ArrowDownOutlined, MinusOutlined } from '@ant-design/icons';
import styles from './index.module.css';

const { Text } = Typography;

type AmountType = 'default' | 'income' | 'expense' | 'frozen' | 'refund';
type AmountSize = 'small' | 'default' | 'large';

interface AmountDisplayProps {
  /**
   * 金额值（分）
   */
  value: number;
  /**
   * 是否自动将分转换为元
   * @default true
   */
  fromCents?: boolean;
  /**
   * 货币符号
   * @default '¥'
   */
  currency?: string;
  /**
   * 是否显示货币符号
   * @default true
   */
  showCurrency?: boolean;
  /**
   * 小数位数
   * @default 2
   */
  precision?: number;
  /**
   * 金额类型，影响颜色显示
   * @default 'default'
   */
  type?: AmountType;
  /**
   * 尺寸
   * @default 'default'
   */
  size?: AmountSize;
  /**
   * 是否显示趋势图标（正数显示上箭头，负数显示下箭头）
   * @default false
   */
  showTrend?: boolean;
  /**
   * 是否显示正负号
   * @default false
   */
  showSign?: boolean;
  /**
   * 原始金额（用于显示划线价）
   */
  originalValue?: number;
  /**
   * 提示文本
   */
  tooltip?: string;
  /**
   * 前缀文本
   */
  prefix?: React.ReactNode;
  /**
   * 后缀文本
   */
  suffix?: React.ReactNode;
  /**
   * 自定义类名
   */
  className?: string;
  /**
   * 自定义样式
   */
  style?: React.CSSProperties;
}

// 类型颜色配置
const TYPE_COLORS: Record<AmountType, string> = {
  default: 'inherit',
  income: '#52c41a',
  expense: '#ff4d4f',
  frozen: '#faad14',
  refund: '#1890ff',
};

// 尺寸配置
const SIZE_CONFIG: Record<AmountSize, { fontSize: number; fontWeight: number }> = {
  small: { fontSize: 12, fontWeight: 400 },
  default: { fontSize: 14, fontWeight: 500 },
  large: { fontSize: 20, fontWeight: 600 },
};

export const AmountDisplay: React.FC<AmountDisplayProps> = ({
  value,
  fromCents = true,
  currency = '¥',
  showCurrency = true,
  precision = 2,
  type = 'default',
  size = 'default',
  showTrend = false,
  showSign = false,
  originalValue,
  tooltip,
  prefix,
  suffix,
  className,
  style,
}) => {
  // 将分转换为元
  const displayValue = fromCents ? value / 100 : value;
  const displayOriginal = originalValue !== undefined 
    ? (fromCents ? originalValue / 100 : originalValue)
    : undefined;

  // 格式化金额
  const formatAmount = (amount: number): string => {
    const absValue = Math.abs(amount);
    return absValue.toLocaleString('zh-CN', {
      minimumFractionDigits: precision,
      maximumFractionDigits: precision,
    });
  };

  // 获取符号
  const getSign = (): string => {
    if (!showSign) return '';
    if (displayValue > 0) return '+';
    if (displayValue < 0) return '-';
    return '';
  };

  // 获取趋势图标
  const getTrendIcon = () => {
    if (!showTrend) return null;
    if (displayValue > 0) return <ArrowUpOutlined style={{ color: '#52c41a' }} />;
    if (displayValue < 0) return <ArrowDownOutlined style={{ color: '#ff4d4f' }} />;
    return <MinusOutlined style={{ color: '#999' }} />;
  };

  const sizeStyle = SIZE_CONFIG[size];
  const color = TYPE_COLORS[type];

  const content = (
    <span
      className={`${styles.amountDisplay} ${className || ''}`}
      style={{
        color,
        fontSize: sizeStyle.fontSize,
        fontWeight: sizeStyle.fontWeight,
        ...style,
      }}
    >
      {prefix && <span className={styles.prefix}>{prefix}</span>}
      
      {showTrend && <span className={styles.trend}>{getTrendIcon()}</span>}
      
      <span className={styles.amount}>
        {showSign && getSign()}
        {showCurrency && <span className={styles.currency}>{currency}</span>}
        {formatAmount(Math.abs(displayValue))}
      </span>

      {/* 划线价 */}
      {displayOriginal !== undefined && displayOriginal !== displayValue && (
        <span className={styles.original}>
          <Text delete type="secondary" style={{ fontSize: sizeStyle.fontSize * 0.8 }}>
            {showCurrency && currency}
            {formatAmount(displayOriginal)}
          </Text>
        </span>
      )}

      {suffix && <span className={styles.suffix}>{suffix}</span>}
    </span>
  );

  if (tooltip) {
    return (
      <Tooltip title={tooltip}>
        {content}
      </Tooltip>
    );
  }

  return content;
};

/**
 * 快捷方式：收入金额（绿色）
 */
export const IncomeAmount: React.FC<Omit<AmountDisplayProps, 'type'>> = (props) => (
  <AmountDisplay {...props} type="income" showSign />
);

/**
 * 快捷方式：支出金额（红色）
 */
export const ExpenseAmount: React.FC<Omit<AmountDisplayProps, 'type'>> = (props) => (
  <AmountDisplay {...props} type="expense" showSign />
);

/**
 * 快捷方式：冻结金额（黄色）
 */
export const FrozenAmount: React.FC<Omit<AmountDisplayProps, 'type'>> = (props) => (
  <AmountDisplay {...props} type="frozen" />
);

/**
 * 快捷方式：退款金额（蓝色）
 */
export const RefundAmount: React.FC<Omit<AmountDisplayProps, 'type'>> = (props) => (
  <AmountDisplay {...props} type="refund" />
);

export default AmountDisplay;
