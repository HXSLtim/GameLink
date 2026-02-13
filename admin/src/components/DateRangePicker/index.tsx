/**
 * 带预设的日期范围选择组件
 */
import React, { useMemo } from 'react';
import { DatePicker, Space, Button, Dropdown } from 'antd';
import type { MenuProps } from 'antd';
import { DownOutlined } from '@ant-design/icons';
import dayjs, { Dayjs } from 'dayjs';
import quarterOfYear from 'dayjs/plugin/quarterOfYear';
import styles from './index.module.css';

// 启用季度插件
dayjs.extend(quarterOfYear);

const { RangePicker } = DatePicker;

type RangeValue = [Dayjs | null, Dayjs | null] | null;

interface PresetRange {
  label: string;
  value: [Dayjs, Dayjs];
}

interface DateRangePickerProps {
  /**
   * 当前值
   */
  value?: RangeValue;
  /**
   * 值变更回调
   */
  onChange?: (dates: RangeValue, dateStrings: [string, string]) => void;
  /**
   * 日期格式
   * @default 'YYYY-MM-DD'
   */
  format?: string;
  /**
   * 是否显示预设按钮
   * @default true
   */
  showPresets?: boolean;
  /**
   * 自定义预设选项
   */
  presets?: PresetRange[];
  /**
   * 占位文本
   */
  placeholder?: [string, string];
  /**
   * 是否允许清空
   * @default true
   */
  allowClear?: boolean;
  /**
   * 是否禁用
   * @default false
   */
  disabled?: boolean;
  /**
   * 尺寸
   */
  size?: 'small' | 'middle' | 'large';
  /**
   * 是否显示快捷下拉菜单
   * @default false
   */
  showQuickMenu?: boolean;
  /**
   * 自定义样式
   */
  style?: React.CSSProperties;
  /**
   * 自定义类名
   */
  className?: string;
}

// 默认预设选项
const DEFAULT_PRESETS: PresetRange[] = [
  {
    label: '今天',
    value: [dayjs().startOf('day'), dayjs().endOf('day')],
  },
  {
    label: '昨天',
    value: [dayjs().subtract(1, 'day').startOf('day'), dayjs().subtract(1, 'day').endOf('day')],
  },
  {
    label: '本周',
    value: [dayjs().startOf('week'), dayjs().endOf('week')],
  },
  {
    label: '上周',
    value: [
      dayjs().subtract(1, 'week').startOf('week'),
      dayjs().subtract(1, 'week').endOf('week'),
    ],
  },
  {
    label: '本月',
    value: [dayjs().startOf('month'), dayjs().endOf('month')],
  },
  {
    label: '上月',
    value: [
      dayjs().subtract(1, 'month').startOf('month'),
      dayjs().subtract(1, 'month').endOf('month'),
    ],
  },
  {
    label: '最近7天',
    value: [dayjs().subtract(6, 'day').startOf('day'), dayjs().endOf('day')],
  },
  {
    label: '最近30天',
    value: [dayjs().subtract(29, 'day').startOf('day'), dayjs().endOf('day')],
  },
  {
    label: '最近90天',
    value: [dayjs().subtract(89, 'day').startOf('day'), dayjs().endOf('day')],
  },
  {
    label: '本季度',
    value: [dayjs().startOf('quarter'), dayjs().endOf('quarter')],
  },
  {
    label: '本年',
    value: [dayjs().startOf('year'), dayjs().endOf('year')],
  },
];

export const DateRangePicker: React.FC<DateRangePickerProps> = ({
  value,
  onChange,
  format = 'YYYY-MM-DD',
  showPresets = true,
  presets,
  placeholder = ['开始日期', '结束日期'],
  allowClear = true,
  disabled = false,
  size = 'middle',
  showQuickMenu = false,
  style,
  className,
}) => {
  const allPresets = presets || DEFAULT_PRESETS;

  // 转换为 Ant Design 的 presets 格式
  const antdPresets = useMemo(() => {
    if (!showPresets) return undefined;
    return allPresets.map((preset) => ({
      label: preset.label,
      value: preset.value,
    }));
  }, [showPresets, allPresets]);

  // 快捷菜单
  const quickMenuItems: MenuProps['items'] = useMemo(() => {
    return allPresets.map((preset, index) => ({
      key: index,
      label: preset.label,
      onClick: () => {
        if (onChange) {
          onChange(preset.value, [
            preset.value[0].format(format),
            preset.value[1].format(format),
          ]);
        }
      },
    }));
  }, [allPresets, onChange, format]);

  // 获取当前选中的预设标签
  const currentPresetLabel = useMemo(() => {
    if (!value || !value[0] || !value[1]) return null;
    
    const preset = allPresets.find((p) => {
      return (
        value[0]?.isSame(p.value[0], 'day') &&
        value[1]?.isSame(p.value[1], 'day')
      );
    });
    
    return preset?.label || null;
  }, [value, allPresets]);

  if (showQuickMenu) {
    return (
      <Space.Compact className={`${styles.dateRangePicker} ${className || ''}`} style={style}>
        <RangePicker
          value={value}
          onChange={onChange}
          format={format}
          placeholder={placeholder}
          allowClear={allowClear}
          disabled={disabled}
          size={size}
          presets={antdPresets}
        />
        <Dropdown menu={{ items: quickMenuItems }} placement="bottomRight">
          <Button size={size} disabled={disabled}>
            {currentPresetLabel || '快捷选择'}
            <DownOutlined />
          </Button>
        </Dropdown>
      </Space.Compact>
    );
  }

  return (
    <RangePicker
      className={`${styles.dateRangePicker} ${className || ''}`}
      style={style}
      value={value}
      onChange={onChange}
      format={format}
      placeholder={placeholder}
      allowClear={allowClear}
      disabled={disabled}
      size={size}
      presets={antdPresets}
    />
  );
};

export default DateRangePicker;
