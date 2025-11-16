/**
 * 陪玩师列表数据管理Hook
 */
import { useState, useEffect, useCallback } from 'react';
import { getPlayers, type GetPlayersParams, type Player } from '@/api/player';

/**
 * 陪玩师列表Hook返回值
 */
interface UsePlayerListReturn {
  players: Player[];
  loading: boolean;
  error: Error | null;
  total: number;
  page: number;
  pageSize: number;
  refresh: () => Promise<void>;
  updateFilters: (filters: Partial<GetPlayersParams>) => void;
}

/**
 * 陪玩师列表Hook
 */
export const usePlayerList = (
  initialParams?: GetPlayersParams
): UsePlayerListReturn => {
  const [players, setPlayers] = useState<Player[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);
  const [total, setTotal] = useState(0);
  const [params, setParams] = useState<GetPlayersParams>(
    initialParams || {
      page: 1,
      pageSize: 20,
      sortBy: 'rating',
      sortOrder: 'desc',
    }
  );

  const fetchPlayers = useCallback(async () => {
    setLoading(true);
    setError(null);

    try {
      const response = await getPlayers(params);
      setPlayers(response.players);
      setTotal(response.total);
    } catch (err) {
      setError(err instanceof Error ? err : new Error('获取陪玩师列表失败'));
      console.error('Failed to fetch players:', err);
    } finally {
      setLoading(false);
    }
  }, [params]);

  useEffect(() => {
    fetchPlayers();
  }, [fetchPlayers]);

  const refresh = useCallback(async () => {
    await fetchPlayers();
  }, [fetchPlayers]);

  const updateFilters = useCallback((filters: Partial<GetPlayersParams>) => {
    setParams((prev) => ({
      ...prev,
      ...filters,
      page: filters.page ?? 1, // 重置页码除非明确指定
    }));
  }, []);

  return {
    players,
    loading,
    error,
    total,
    page: params.page || 1,
    pageSize: params.pageSize || 20,
    refresh,
    updateFilters,
  };
};
