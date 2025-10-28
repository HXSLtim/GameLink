import { renderHook, waitFor } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { useTable } from './useTable';

describe('useTable', () => {
  it('should initialize with default values', () => {
    const mockFetchData = vi.fn().mockResolvedValue({ items: [], total: 0 });

    const { result } = renderHook(() =>
      useTable({
        fetchData: mockFetchData,
        autoFetch: false,
      }),
    );

    expect(result.current.data).toEqual([]);
    expect(result.current.loading).toBe(false);
    expect(result.current.pagination).toEqual({
      page: 1,
      pageSize: 10,
      total: 0,
    });
  });

  it('should fetch data on mount when autoFetch is true', async () => {
    const mockData = { items: [{ id: 1 }, { id: 2 }], total: 2 };
    const mockFetchData = vi.fn().mockResolvedValue(mockData);

    renderHook(() =>
      useTable({
        fetchData: mockFetchData,
        autoFetch: true,
      }),
    );

    await waitFor(() => {
      expect(mockFetchData).toHaveBeenCalledWith(1, 10);
    });
  });

  it('should not fetch data on mount when autoFetch is false', () => {
    const mockFetchData = vi.fn().mockResolvedValue({ items: [], total: 0 });

    renderHook(() =>
      useTable({
        fetchData: mockFetchData,
        autoFetch: false,
      }),
    );

    expect(mockFetchData).not.toHaveBeenCalled();
  });

  it('should update data after successful fetch', async () => {
    const mockData = { items: [{ id: 1 }, { id: 2 }], total: 2 };
    const mockFetchData = vi.fn().mockResolvedValue(mockData);

    const { result } = renderHook(() =>
      useTable({
        fetchData: mockFetchData,
        autoFetch: true,
      }),
    );

    await waitFor(() => {
      expect(result.current.data).toEqual(mockData.items);
      expect(result.current.pagination.total).toBe(2);
      expect(result.current.loading).toBe(false);
    });
  });

  it('should handle pagination change', async () => {
    const mockData = { items: [{ id: 3 }], total: 1 };
    const mockFetchData = vi.fn().mockResolvedValue(mockData);

    const { result } = renderHook(() =>
      useTable({
        fetchData: mockFetchData,
        autoFetch: false,
      }),
    );

    result.current.handlePageChange({ current: 2, pageSize: 20 });

    await waitFor(() => {
      expect(mockFetchData).toHaveBeenCalledWith(2, 20);
    });
  });

  it('should provide refetch function', async () => {
    const mockData = { items: [{ id: 1 }], total: 1 };
    const mockFetchData = vi.fn().mockResolvedValue(mockData);

    const { result } = renderHook(() =>
      useTable({
        fetchData: mockFetchData,
        autoFetch: true,
      }),
    );

    await waitFor(() => {
      expect(mockFetchData).toHaveBeenCalledTimes(1);
    });

    result.current.refetch();

    await waitFor(() => {
      expect(mockFetchData).toHaveBeenCalledTimes(2);
    });
  });

  it('should provide reset function', async () => {
    const mockData = { items: [], total: 0 };
    const mockFetchData = vi.fn().mockResolvedValue(mockData);

    const { result } = renderHook(() =>
      useTable({
        fetchData: mockFetchData,
        initialPage: 3,
        initialPageSize: 20,
        autoFetch: false,
      }),
    );

    // Change page
    result.current.handlePageChange({ current: 5, pageSize: 50 });

    await waitFor(() => {
      expect(mockFetchData).toHaveBeenCalledWith(5, 50);
    });

    // Reset to initial
    result.current.reset();

    await waitFor(() => {
      expect(mockFetchData).toHaveBeenCalledWith(3, 20);
    });
  });

  it('should handle fetch errors gracefully', async () => {
    const mockFetchData = vi.fn().mockRejectedValue(new Error('Network error'));

    const { result } = renderHook(() =>
      useTable({
        fetchData: mockFetchData,
        autoFetch: true,
      }),
    );

    await waitFor(() => {
      expect(result.current.data).toEqual([]);
      expect(result.current.loading).toBe(false);
      expect(result.current.pagination.total).toBe(0);
    });
  });

  it('should set loading state during fetch', async () => {
    const mockFetchData = vi.fn<
      (page: number, pageSize: number) => Promise<{ items: unknown[]; total: number }>
    >(() => new Promise((resolve) => setTimeout(() => resolve({ items: [], total: 0 }), 100)));

    const { result } = renderHook(() =>
      useTable({
        fetchData: mockFetchData,
        autoFetch: true,
      }),
    );

    // Should be loading immediately
    expect(result.current.loading).toBe(true);

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });
  });
});
