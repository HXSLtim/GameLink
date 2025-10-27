import { useState, useCallback, useEffect } from 'react';
import { handleApiError } from '../utils/errorHandler';

export interface PaginationState {
  page: number;
  pageSize: number;
  total: number;
}

export interface TableState<T> {
  data: T[];
  loading: boolean;
  pagination: PaginationState;
}

export interface UseTableOptions<T> {
  /** Initial page number */
  initialPage?: number;
  /** Initial page size */
  initialPageSize?: number;
  /** Fetch data function */
  fetchData: (page: number, pageSize: number) => Promise<{ items: T[]; total: number }>;
  /** Error message prefix */
  errorMessage?: string;
  /** Enable auto fetch on mount */
  autoFetch?: boolean;
}

/**
 * Custom hook for table data management with pagination
 *
 * @template T - The type of data items
 * @param options - Table configuration options
 * @returns Table state and control functions
 *
 * @example
 * ```tsx
 * const { data, loading, pagination, handlePageChange, refetch } = useTable({
 *   fetchData: (page, pageSize) => userService.list({ page, page_size: pageSize }),
 *   errorMessage: '获取用户列表',
 * });
 * ```
 */
export const useTable = <T>({
  initialPage = 1,
  initialPageSize = 10,
  fetchData,
  errorMessage = '获取数据',
  autoFetch = true,
}: UseTableOptions<T>) => {
  const [state, setState] = useState<TableState<T>>({
    data: [],
    loading: false,
    pagination: {
      page: initialPage,
      pageSize: initialPageSize,
      total: 0,
    },
  });

  const fetchTableData = useCallback(
    async (page: number, pageSize: number) => {
      setState((prev) => ({ ...prev, loading: true }));

      try {
        const result = await fetchData(page, pageSize);
        setState({
          data: result.items,
          loading: false,
          pagination: {
            page,
            pageSize,
            total: result.total,
          },
        });
      } catch (error) {
        handleApiError(error, errorMessage);
        setState((prev) => ({
          ...prev,
          data: [],
          loading: false,
          pagination: {
            ...prev.pagination,
            total: 0,
          },
        }));
      }
    },
    [fetchData, errorMessage],
  );

  const handlePageChange = useCallback(
    (pagination: { current?: number; pageSize?: number }) => {
      const newPage = pagination.current || 1;
      const newPageSize = pagination.pageSize || state.pagination.pageSize;
      fetchTableData(newPage, newPageSize);
    },
    [fetchTableData, state.pagination.pageSize],
  );

  const refetch = useCallback(() => {
    fetchTableData(state.pagination.page, state.pagination.pageSize);
  }, [fetchTableData, state.pagination.page, state.pagination.pageSize]);

  const reset = useCallback(() => {
    fetchTableData(initialPage, initialPageSize);
  }, [fetchTableData, initialPage, initialPageSize]);

  useEffect(() => {
    if (autoFetch) {
      fetchTableData(initialPage, initialPageSize);
    }
  }, [autoFetch, fetchTableData, initialPage, initialPageSize]);

  return {
    data: state.data,
    loading: state.loading,
    pagination: state.pagination,
    handlePageChange,
    refetch,
    reset,
  };
};

