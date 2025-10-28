import React from 'react';
import styles from './Table.module.less';

export interface TableColumn<T = any> {
  key: string;
  title: string;
  dataIndex?: keyof T;
  render?: (value: any, record: T, index: number) => React.ReactNode;
  width?: string | number;
  align?: 'left' | 'center' | 'right';
  fixed?: 'left' | 'right';
  sorter?: boolean;
  sortOrder?: 'ascend' | 'descend' | null;
}

export interface TableProps<T = any> {
  columns: TableColumn<T>[];
  dataSource: T[];
  rowKey?: string | ((record: T) => string);
  loading?: boolean;
  emptyText?: string;
  onRow?: (record: T, index: number) => React.HTMLAttributes<HTMLTableRowElement>;
  scroll?: React.CSSProperties;
  className?: string;
  rowClassName?: string | ((record: T, index: number) => string);
}

export const Table = <T extends Record<string, any>>({
  columns = [],
  dataSource = [],
  rowKey = 'id',
  loading = false,
  emptyText = '暂无数据',
  onRow,
  scroll,
  className,
  rowClassName,
}: TableProps<T>) => {
  const getRowKey = (record: T, index: number): string => {
    if (typeof rowKey === 'function') {
      return rowKey(record);
    }
    return record[rowKey] || String(index);
  };

  const getRowClassName = (record: T, index: number): string => {
    const baseClass = styles.row;
    const customClass =
      typeof rowClassName === 'function' ? rowClassName(record, index) : rowClassName || '';
    return `${baseClass} ${customClass}`.trim();
  };

  const renderCell = (column: TableColumn<T>, record: T, index: number) => {
    if (column.render) {
      return column.render(column.dataIndex ? record[column.dataIndex] : undefined, record, index);
    }
    return column.dataIndex ? record[column.dataIndex] : null;
  };

  return (
    <div className={`${styles.tableWrapper} ${className || ''}`}>
      <div className={styles.tableScroll} style={scroll}>
        <table className={styles.table}>
          <thead className={styles.thead}>
            <tr>
              {(columns || []).map((column) => (
                <th
                  key={column.key}
                  className={styles.th}
                  style={{
                    width: column.width,
                    textAlign: column.align || 'left',
                    left: column.fixed === 'left' ? 0 : undefined,
                    right: column.fixed === 'right' ? 0 : undefined,
                    position: column.fixed ? 'sticky' : undefined,
                  }}
                >
                  {column.title}
                </th>
              ))}
            </tr>
          </thead>
          <tbody className={styles.tbody}>
            {loading ? (
              <tr>
                <td colSpan={(columns || []).length} className={styles.loading}>
                  <div className={styles.loadingContent}>加载中...</div>
                </td>
              </tr>
            ) : !dataSource || dataSource.length === 0 ? (
              <tr>
                <td colSpan={(columns || []).length} className={styles.empty}>
                  <div className={styles.emptyContent}>{emptyText}</div>
                </td>
              </tr>
            ) : (
              (dataSource || []).map((record, index) => (
                <tr
                  key={getRowKey(record, index)}
                  className={getRowClassName(record, index)}
                  {...(onRow ? onRow(record, index) : {})}
                >
                  {(columns || []).map((column) => (
                    <td
                      key={column.key}
                      className={styles.td}
                      style={{
                        width: column.width,
                        textAlign: column.align || 'left',
                        left: column.fixed === 'left' ? 0 : undefined,
                        right: column.fixed === 'right' ? 0 : undefined,
                        position: column.fixed ? 'sticky' : undefined,
                      }}
                    >
                      {renderCell(column, record, index)}
                    </td>
                  ))}
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
};
