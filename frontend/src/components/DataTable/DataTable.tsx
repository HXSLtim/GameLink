import React, { ReactNode } from 'react';
import { Card, Table, Pagination } from '../index';
import type { TableColumn } from '../Table/Table';
import styles from './DataTable.module.less';

export interface FilterConfig {
  label: string;
  key: string;
  element: ReactNode;
}

export interface DataTableProps<T = any> {
  // 标题
  title: string;
  
  // 筛选配置
  filters?: FilterConfig[];
  filterActions?: ReactNode;
  
  // 表格数据
  columns: TableColumn<T>[];
  dataSource: T[];
  loading?: boolean;
  rowKey?: string | ((record: T) => string);
  
  // 分页
  pagination?: {
    current: number;
    pageSize: number;
    total: number;
    onChange: (page: number) => void;
  };
  
  // 自定义样式
  className?: string;
}

export const DataTable = <T extends Record<string, any>>({
  title,
  filters,
  filterActions,
  columns,
  dataSource,
  loading = false,
  rowKey = 'id',
  pagination,
  className,
}: DataTableProps<T>) => {
  return (
    <div className={`${styles.container} ${className || ''}`}>
      {/* 标题 */}
      <div className={styles.header}>
        <h1 className={styles.title}>{title}</h1>
      </div>

      {/* 筛选区域 */}
      {filters && filters.length > 0 && (
        <Card className={styles.filterCard}>
          <div className={styles.filters}>
            <div className={styles.filterRow}>
              {filters.map((filter) => (
                <div key={filter.key} className={styles.filterItem}>
                  <label className={styles.filterLabel}>{filter.label}</label>
                  {filter.element}
                </div>
              ))}
            </div>
            
            {filterActions && (
              <div className={styles.filterActions}>{filterActions}</div>
            )}
          </div>
        </Card>
      )}

      {/* 数据表格 */}
      <Card className={styles.tableCard}>
        <Table 
          columns={columns} 
          dataSource={dataSource} 
          loading={loading} 
          rowKey={rowKey} 
        />

        {/* 分页 */}
        {pagination && pagination.total > 0 && (
          <div className={styles.pagination}>
            <Pagination
              current={pagination.current}
              pageSize={pagination.pageSize}
              total={pagination.total}
              onChange={pagination.onChange}
            />
          </div>
        )}
      </Card>
    </div>
  );
};
