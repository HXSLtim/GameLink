import { useState, useEffect, useCallback } from 'react';
import { useSearchParams } from 'react-router-dom';
import { mergeUrlParams } from '../utils/urlParams';

export interface ListQueryParams {
  page?: number;
  page_size?: number;
  keyword?: string;
  [key: string]: any;
}

export interface UseListPageOptions<T, Q extends ListQueryParams> {
  initialParams: Q;
  fetchData: (params: Q) => Promise<{ list: T[]; total: number }>;
  /**
   * 是否在筛选条件变化时自动搜索
   * @default true
   */
  autoSearch?: boolean;
  /**
   * 需要从URL查询参数中读取的参数键
   * 例如: ['status', 'role']
   */
  urlParamKeys?: string[];
}

export interface UseListPageResult<T, Q extends ListQueryParams> {
  loading: boolean;
  data: T[];
  total: number;
  queryParams: Q;
  setQueryParams: React.Dispatch<React.SetStateAction<Q>>;
  handleSearch: () => Promise<void>;
  handleReset: (resetParams: Q) => Promise<void>;
  handlePageChange: (page: number) => void;
  reload: () => Promise<void>;
}

/**
 * 通用列表页Hook
 * 封装列表页的通用逻辑：加载数据、搜索、重置、分页
 */
export function useListPage<T, Q extends ListQueryParams>({
  initialParams,
  fetchData,
  _autoSearch = true,
  urlParamKeys = [],
}: UseListPageOptions<T, Q>): UseListPageResult<T, Q> {
  const [searchParams] = useSearchParams();
  const [loading, setLoading] = useState(false);
  const [data, setData] = useState<T[]>([]);
  const [total, setTotal] = useState(0);

  // 从URL参数合并初始参数
  const getInitialParams = useCallback((): Q => {
    if (urlParamKeys.length === 0) {
      return initialParams;
    }
    return mergeUrlParams(searchParams, initialParams, urlParamKeys);
  }, [searchParams, initialParams, urlParamKeys]);

  const [queryParams, setQueryParams] = useState<Q>(getInitialParams());
  const [isInitialized, setIsInitialized] = useState(false);

  // 加载数据
  const loadData = useCallback(async () => {
    setLoading(true);
    try {
      const result = await fetchData(queryParams);
      if (result && result.list) {
        setData(result.list);
        setTotal(result.total || 0);
      } else {
        setData([]);
        setTotal(0);
      }
    } catch (err) {
      console.error('加载数据失败:', err);
      setData([]);
      setTotal(0);
    } finally {
      setLoading(false);
    }
  }, [queryParams, fetchData]);

  // 搜索
  const handleSearch = useCallback(async () => {
    setQueryParams((prev) => ({ ...prev, page: 1 }));
    // 需要等待queryParams更新后再加载，所以直接在这里加载
    setTimeout(() => loadData(), 0);
  }, [loadData]);

  // 重置
  const handleReset = useCallback(
    async (resetParams: Q) => {
      setQueryParams(resetParams);
      setLoading(true);
      try {
        const result = await fetchData(resetParams);
        if (result && result.list) {
          setData(result.list);
          setTotal(result.total || 0);
        } else {
          setData([]);
          setTotal(0);
        }
      } catch (err) {
        console.error('加载数据失败:', err);
        setData([]);
        setTotal(0);
      } finally {
        setLoading(false);
      }
    },
    [fetchData],
  );

  // 分页变化
  const handlePageChange = useCallback((page: number) => {
    setQueryParams((prev) => ({ ...prev, page }));
  }, []);

  // 重新加载
  const reload = useCallback(async () => {
    await loadData();
  }, [loadData]);

  // 监听URL参数变化
  useEffect(() => {
    if (urlParamKeys.length === 0) return;

    const params = mergeUrlParams(searchParams, initialParams, urlParamKeys);
    setQueryParams((prev) => {
      // 只在有变化时更新
      const hasChanges = urlParamKeys.some((key) => prev[key] !== params[key]);
      if (hasChanges) {
        return { ...prev, ...params, page: 1 }; // 重置到第一页
      }
      return prev;
    });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [searchParams]);

  // 初始加载
  useEffect(() => {
    loadData();
    setIsInitialized(true);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  // 参数变化时自动加载（跳过初始化）
  useEffect(() => {
    if (!isInitialized) return;

    // 自动搜索模式：任何参数变化都触发加载
    // 手动搜索模式也要加载，因为分页、每页数量变化需要立即生效
    loadData();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [queryParams, isInitialized]);

  return {
    loading,
    data,
    total,
    queryParams,
    setQueryParams,
    handleSearch,
    handleReset,
    handlePageChange,
    reload,
  };
}
