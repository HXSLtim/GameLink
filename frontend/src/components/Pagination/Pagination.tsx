import React from 'react';
import styles from './Pagination.module.less';

export interface PaginationProps {
  current: number;
  total: number;
  pageSize?: number;
  onChange?: (page: number) => void;
  showSizeChanger?: boolean;
  onSizeChange?: (size: number) => void;
  pageSizeOptions?: number[];
  className?: string;
}

export const Pagination: React.FC<PaginationProps> = ({
  current,
  total,
  pageSize = 10,
  onChange,
  showSizeChanger = false,
  onSizeChange,
  pageSizeOptions = [10, 20, 50, 100],
  className,
}) => {
  const totalPages = Math.ceil(total / pageSize);

  const handlePageChange = (page: number) => {
    if (page < 1 || page > totalPages || page === current) return;
    onChange?.(page);
  };

  const handleSizeChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const newSize = Number(e.target.value);
    onSizeChange?.(newSize);
  };

  const getPageNumbers = () => {
    const pages: (number | string)[] = [];
    const showPages = 5; // 显示的页码数量

    if (totalPages <= showPages + 2) {
      // 总页数较少，显示所有页码
      for (let i = 1; i <= totalPages; i++) {
        pages.push(i);
      }
    } else {
      // 总页数较多，显示部分页码
      pages.push(1);

      if (current <= 3) {
        for (let i = 2; i <= 4; i++) {
          pages.push(i);
        }
        pages.push('...');
        pages.push(totalPages);
      } else if (current >= totalPages - 2) {
        pages.push('...');
        for (let i = totalPages - 3; i < totalPages; i++) {
          pages.push(i);
        }
        pages.push(totalPages);
      } else {
        pages.push('...');
        for (let i = current - 1; i <= current + 1; i++) {
          pages.push(i);
        }
        pages.push('...');
        pages.push(totalPages);
      }
    }

    return pages;
  };

  if (totalPages <= 1) return null;

  return (
    <div className={`${styles.pagination} ${className || ''}`}>
      <button
        className={styles.button}
        onClick={() => handlePageChange(current - 1)}
        disabled={current === 1}
      >
        <ArrowLeftIcon />
      </button>

      {getPageNumbers().map((page, index) =>
        typeof page === 'number' ? (
          <button
            key={page}
            className={`${styles.button} ${page === current ? styles.active : ''}`}
            onClick={() => handlePageChange(page)}
          >
            {page}
          </button>
        ) : (
          <span key={`ellipsis-${index}`} className={styles.ellipsis}>
            {page}
          </span>
        ),
      )}

      <button
        className={styles.button}
        onClick={() => handlePageChange(current + 1)}
        disabled={current === totalPages}
      >
        <ArrowRightIcon />
      </button>

      {showSizeChanger && (
        <select className={styles.sizeChanger} value={pageSize} onChange={handleSizeChange}>
          {pageSizeOptions.map((size) => (
            <option key={size} value={size}>
              {size} / 页
            </option>
          ))}
        </select>
      )}

      <span className={styles.info}>
        共 {total} 条，第 {current} / {totalPages} 页
      </span>
    </div>
  );
};

const ArrowLeftIcon = () => (
  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor">
    <path strokeLinecap="square" strokeLinejoin="miter" strokeWidth="2" d="M15 18l-6-6 6-6" />
  </svg>
);

const ArrowRightIcon = () => (
  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor">
    <path strokeLinecap="square" strokeLinejoin="miter" strokeWidth="2" d="M9 18l6-6-6-6" />
  </svg>
);
