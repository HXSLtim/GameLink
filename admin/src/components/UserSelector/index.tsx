/**
 * 用户/陪玩选择器组件
 * 支持搜索、分页、单选/多选
 */
import React, { useState, useCallback, useEffect, useMemo } from 'react';
import {
  Select,
  Spin,
  Avatar,
  Tag,
  Space,
  Typography,
  Empty,
} from 'antd';
import { UserOutlined, SearchOutlined } from '@ant-design/icons';
import { adminApi } from '@/api/admin';
import styles from './index.module.css';

const { Text } = Typography;

interface User {
  id: number;
  name: string;
  phone: string;
  avatar?: string;
  status?: string;
  isPlayer?: boolean;
}

interface Player {
  id: number;
  userId: number;
  nickname: string;
  avatar?: string;
  status?: string;
  rating?: number;
  orderCount?: number;
}

type SelectorMode = 'user' | 'player' | 'both';

interface UserSelectorProps {
  /**
   * 选择模式: user - 仅用户, player - 仅陪玩, both - 两者皆可
   */
  mode?: SelectorMode;
  /**
   * 当前选中值
   */
  value?: number | number[];
  /**
   * 值变更回调
   */
  onChange?: (value: number | number[] | undefined, option: User | Player | (User | Player)[] | undefined) => void;
  /**
   * 是否多选
   */
  multiple?: boolean;
  /**
   * 占位文本
   */
  placeholder?: string;
  /**
   * 是否禁用
   */
  disabled?: boolean;
  /**
   * 是否允许清空
   */
  allowClear?: boolean;
  /**
   * 宽度
   */
  width?: number | string;
  /**
   * 最大显示标签数
   */
  maxTagCount?: number | 'responsive';
  /**
   * 过滤特定状态的用户/陪玩
   */
  statusFilter?: string | string[];
  /**
   * 自定义样式类名
   */
  className?: string;
}

export const UserSelector: React.FC<UserSelectorProps> = ({
  mode = 'user',
  value,
  onChange,
  multiple = false,
  placeholder,
  disabled = false,
  allowClear = true,
  width = 300,
  maxTagCount = 'responsive',
  statusFilter,
  className,
}) => {
  const [options, setOptions] = useState<(User | Player)[]>([]);
  const [loading, setLoading] = useState(false);
  const [searchText, setSearchText] = useState('');

  // 默认占位文本
  const defaultPlaceholder = useMemo(() => {
    switch (mode) {
      case 'player':
        return '搜索陪玩师';
      case 'both':
        return '搜索用户或陪玩师';
      default:
        return '搜索用户';
    }
  }, [mode]);

  // 加载数据
  const fetchData = useCallback(async (keyword: string) => {
    setLoading(true);
    try {
      const results: (User | Player)[] = [];

      if (mode === 'user' || mode === 'both') {
        const userRes = await adminApi.getUsers({ keyword, page_size: 20 });
        if (userRes?.data?.data) {
          results.push(
            ...userRes.data.data.map((u: User) => ({
              ...u,
              _type: 'user' as const,
            }))
          );
        }
      }

      if (mode === 'player' || mode === 'both') {
        const playerRes = await adminApi.getPlayers({ keyword, page_size: 20 });
        if (playerRes?.data?.data) {
          results.push(
            ...playerRes.data.data.map((p: Player) => ({
              ...p,
              _type: 'player' as const,
            }))
          );
        }
      }

      // 过滤状态
      let filtered = results;
      if (statusFilter) {
        const statuses = Array.isArray(statusFilter) ? statusFilter : [statusFilter];
        filtered = results.filter((item) => 
          statuses.includes(item.status || '')
        );
      }

      setOptions(filtered);
    } catch (error) {
      console.error('Failed to fetch users/players:', error);
      setOptions([]);
    } finally {
      setLoading(false);
    }
  }, [mode, statusFilter]);

  // 防抖搜索 - 使用 setTimeout 实现
  const debouncedSearchRef = React.useRef<ReturnType<typeof setTimeout> | null>(null);
  const debouncedSearch = React.useCallback((keyword: string) => {
    if (debouncedSearchRef.current) {
      clearTimeout(debouncedSearchRef.current);
    }
    debouncedSearchRef.current = setTimeout(() => {
      fetchData(keyword);
    }, 300);
  }, [fetchData]);

  // 初始加载
  useEffect(() => {
    fetchData('');
  }, [fetchData]);

  // 搜索变更
  const handleSearch = (keyword: string) => {
    setSearchText(keyword);
    debouncedSearch(keyword);
  };

  // 选择变更
  const handleChange = (selectedValue: number | number[] | undefined) => {
    if (!onChange) return;

    if (selectedValue === undefined) {
      onChange(undefined, undefined);
      return;
    }

    if (multiple) {
      const selectedOptions = options.filter((opt) =>
        (selectedValue as number[]).includes(isPlayer(opt) ? opt.id : opt.id)
      );
      onChange(selectedValue, selectedOptions);
    } else {
      const selectedOption = options.find((opt) =>
        isPlayer(opt) ? opt.id === selectedValue : opt.id === selectedValue
      );
      onChange(selectedValue as number, selectedOption);
    }
  };

  // 判断是否为陪玩
  const isPlayer = (item: User | Player): item is Player => {
    return 'nickname' in item;
  };

  // 渲染选项
  const renderOption = (item: User | Player) => {
    const id = item.id;
    const name = isPlayer(item) ? item.nickname : item.name;
    const avatar = item.avatar;
    const status = item.status;

    return {
      value: id,
      label: (
        <div className={styles.option}>
          <Avatar
            size="small"
            src={avatar}
            icon={<UserOutlined />}
            className={styles.avatar}
          />
          <div className={styles.info}>
            <Text className={styles.name}>{name || `ID: ${id}`}</Text>
            {isPlayer(item) && (
              <Text type="secondary" className={styles.meta}>
                评分: {item.rating?.toFixed(1) || '-'} | 接单: {item.orderCount || 0}
              </Text>
            )}
            {!isPlayer(item) && item.phone && (
              <Text type="secondary" className={styles.meta}>
                {item.phone}
              </Text>
            )}
          </div>
          <Space>
            {isPlayer(item) && (
              <Tag color="purple" className={styles.typeTag}>
                陪玩
              </Tag>
            )}
            {status && (
              <Tag
                color={
                  status === 'active' || status === 'approved'
                    ? 'green'
                    : status === 'pending'
                    ? 'orange'
                    : 'default'
                }
              >
                {status === 'active' || status === 'approved'
                  ? '正常'
                  : status === 'pending'
                  ? '待审核'
                  : status === 'banned' || status === 'suspended'
                  ? '已禁用'
                  : status}
              </Tag>
            )}
          </Space>
        </div>
      ),
      item,
    };
  };

  return (
    <Select
      className={className}
      style={{ width }}
      mode={multiple ? 'multiple' : undefined}
      value={value}
      onChange={handleChange}
      onSearch={handleSearch}
      placeholder={placeholder || defaultPlaceholder}
      disabled={disabled}
      allowClear={allowClear}
      showSearch
      filterOption={false}
      loading={loading}
      maxTagCount={maxTagCount}
      notFoundContent={
        loading ? (
          <Spin size="small" />
        ) : searchText ? (
          <Empty description="未找到匹配结果" image={Empty.PRESENTED_IMAGE_SIMPLE} />
        ) : (
          <Empty description="请输入关键词搜索" image={Empty.PRESENTED_IMAGE_SIMPLE} />
        )
      }
      suffixIcon={loading ? <Spin size="small" /> : <SearchOutlined />}
      options={options.map(renderOption)}
    />
  );
};

export default UserSelector;
