 /**
 * InfiniteList - 无限滚动列表组件
 *
 * 配合 useCrud hook 的 lazyLoad 模式使用
 *
 * @example
 * ```tsx
 * const { data, loading, loadingMore, hasMore, sentinelRef, scrollContainerRef } = useCrud({
 *   api: userApi,
 *   lazyLoad: true,
 * });
 *
 * <InfiniteList
 *   data={data}
 *   loading={loading}
 *   loadingMore={loadingMore}
 *   hasMore={hasMore}
 *   scrollContainerRef={scrollContainerRef}
 *   sentinelRef={sentinelRef}
 *   renderItem={(item) => <UserCard user={item} />}
 *   emptyText="暂无数据"
 * />
 * ```
 */

import React from 'react';
import { Spin, Empty } from 'antd';
import { LoadingOutlined } from '@ant-design/icons';

export interface InfiniteListProps<T> {
  /** 数据列表 */
  data: T[];
  /** 初始加载状态 */
  loading: boolean;
  /** 加载更多状态 */
  loadingMore: boolean;
  /** 是否还有更多数据 */
  hasMore: boolean;
  /** 滚动容器 ref */
  scrollContainerRef: (node: HTMLElement | null) => void;
  /** 哨兵元素 ref */
  sentinelRef: (node: HTMLElement | null) => void;
  /** 渲染单个项目 */
  renderItem: (item: T, index: number) => React.ReactNode;
  /** 获取项目唯一 key */
  rowKey?: keyof T | ((item: T, index: number) => string | number);
  /** 空状态文本 */
  emptyText?: string;
  /** 空状态描述 */
  emptyDescription?: React.ReactNode;
  /** 加载更多文本 */
  loadingMoreText?: string;
  /** 没有更多数据文本 */
  noMoreText?: string;
  /** 容器样式 */
  style?: React.CSSProperties;
  /** 容器类名 */
  className?: string;
  /** 列表项容器样式 */
  listStyle?: React.CSSProperties;
  /** 列表项容器类名 */
  listClassName?: string;
  /** 是否显示"没有更多"提示 */
  showNoMore?: boolean;
  /** 自定义加载指示器 */
  loadingIndicator?: React.ReactNode;
  /** 自定义空状态 */
  emptyElement?: React.ReactNode;
}

function InfiniteList<T>({
  data,
  loading,
  loadingMore,
  hasMore,
  scrollContainerRef,
  sentinelRef,
  renderItem,
  rowKey,
  emptyText = '暂无数据',
  emptyDescription,
  loadingMoreText = '加载中...',
  noMoreText = '没有更多了',
  style,
  className,
  listStyle,
  listClassName,
  showNoMore = true,
  loadingIndicator,
  emptyElement,
}: InfiniteListProps<T>) {
  // 获取项目 key
  const getRowKey = (item: T, index: number): string | number => {
    if (typeof rowKey === 'function') {
      return rowKey(item, index);
    }
    if (rowKey && typeof item === 'object' && item !== null) {
      const key = (item as Record<string, unknown>)[rowKey as string];
      if (key !== undefined && key !== null) {
        return String(key);
      }
    }
    return index;
  };

  // 默认加载指示器
  const defaultLoadingIndicator = (
    <div style={{ textAlign: 'center', padding: '16px 0' }}>
      <Spin indicator={<LoadingOutlined style={{ fontSize: 24 }} spin />} />
      <div style={{ marginTop: 8, color: '#999' }}>{loadingMoreText}</div>
    </div>
  );

  // 初始加载状态
  if (loading && data.length === 0) {
    return (
      <div
        ref={scrollContainerRef}
        style={{
          display: 'flex',
          justifyContent: 'center',
          alignItems: 'center',
          minHeight: 200,
          ...style,
        }}
        className={className}
      >
        <Spin size="large" />
      </div>
    );
  }

  // 空状态
  if (!loading && data.length === 0) {
    return (
      <div
        ref={scrollContainerRef}
        style={{
          display: 'flex',
          justifyContent: 'center',
          alignItems: 'center',
          minHeight: 200,
          ...style,
        }}
        className={className}
      >
        {emptyElement || (
          <Empty description={emptyDescription || emptyText} />
        )}
      </div>
    );
  }

  return (
    <div
      ref={scrollContainerRef}
      style={{
        overflow: 'auto',
        ...style,
      }}
      className={className}
    >
      {/* 列表内容 */}
      <div style={listStyle} className={listClassName}>
        {data.map((item, index) => (
          <React.Fragment key={getRowKey(item, index)}>
            {renderItem(item, index)}
          </React.Fragment>
        ))}
      </div>

      {/* 哨兵元素 - 用于触发加载更多 */}
      <div ref={sentinelRef} style={{ height: 1 }} />

      {/* 加载更多状态 */}
      {loadingMore && (loadingIndicator || defaultLoadingIndicator)}

      {/* 没有更多数据提示 */}
      {!hasMore && !loadingMore && showNoMore && data.length > 0 && (
        <div
          style={{
            textAlign: 'center',
            padding: '16px 0',
            color: '#999',
            fontSize: 14,
          }}
        >
          {noMoreText}
        </div>
      )}
    </div>
  );
}

export default InfiniteList;
