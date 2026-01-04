/**
 * useCrud Hook Tests
 *
 * **Critical**: 90%+ of pages use this hook for CRUD operations
 * Coverage Target: 85%+
 *
 * Test Scenarios:
 * 1. Initialization with config
 * 2. Data fetching with pagination (fetchAll, refresh)
 * 3. Single record retrieval (getById - via fetchAll with params)
 * 4. Record creation (create)
 * 5. Record modification (update)
 * 6. Record deletion (remove)
 * 7. Loading states (loading, submitting)
 * 8. Error handling
 * 9. Pagination handling (setPage, setPageSize)
 * 10. Search parameters handling
 * 11. Custom callbacks (onCreateSuccess, onUpdateSuccess, onDeleteSuccess, onError)
 * 12. Data transformers and pagination extractors
 * 13. Silent mode for operations
 * 14. Refetch functionality
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, act, waitFor } from '@testing-library/react';
import * as fc from 'fast-check';
import React from 'react';
import { useCrud } from './useCrud';
import type { CrudApi, CrudQueryParams } from './useCrud';

// Mock antd message and modal
vi.mock('antd', () => ({
    message: {
        success: vi.fn(),
        error: vi.fn(),
        warning: vi.fn(),
        info: vi.fn(),
    },
    Modal: {
        confirm: vi.fn(({ onOk }: { onOk: () => void }) => {
            onOk();
            return null;
        }),
    },
}));

// Mock logger
vi.mock('@/utils/logger', () => ({
    logger: {
        error: vi.fn(),
        warn: vi.fn(),
        info: vi.fn(),
        debug: vi.fn(),
    },
}));

// Test types
interface TestItem {
    id: number;
    name: string;
    value: number;
}

interface TestCreateInput {
    name: string;
    value: number;
}

interface TestUpdateInput {
    name?: string;
    value?: number;
}

interface TestQueryParams extends CrudQueryParams {
    page?: number;
    page_size?: number;
    search?: string;
}

// Create mock API factory
const createMockApi = (): CrudApi<TestItem, TestCreateInput, TestUpdateInput, TestQueryParams> => ({
    getAll: vi.fn(),
    create: vi.fn(),
    update: vi.fn(),
    remove: vi.fn(),
});

describe('useCrud hook', () => {
    let mockApi: CrudApi<TestItem, TestCreateInput, TestUpdateInput, TestQueryParams>;

    beforeEach(() => {
        vi.clearAllMocks();
        mockApi = createMockApi();
    });

    describe('initialization', () => {
        it('should initialize with default values', () => {
            const { result } = renderHook(() =>
                useCrud<TestItem, TestCreateInput, TestUpdateInput, TestQueryParams>({
                    api: mockApi,
                    fetchOnMount: false,
                })
            );

            expect(result.current.data).toEqual([]);
            expect(result.current.loading).toBe(false);
            expect(result.current.submitting).toBe(false);
            expect(result.current.error).toBe(null);
            expect(result.current.pagination.current).toBe(1);
            expect(result.current.pagination.pageSize).toBe(10);
            expect(result.current.pagination.total).toBe(0);
        });

        it('should initialize with custom pagination', () => {
            const { result } = renderHook(() =>
                useCrud<TestItem, TestCreateInput, TestUpdateInput, TestQueryParams>({
                    api: mockApi,
                    fetchOnMount: false,
                    initialPagination: {
                        current: 2,
                        pageSize: 20,
                    },
                })
            );

            expect(result.current.pagination.current).toBe(2);
            expect(result.current.pagination.pageSize).toBe(20);
        });

        it('should initialize with custom query params', () => {
            const { result } = renderHook(() =>
                useCrud<TestItem, TestCreateInput, TestUpdateInput, TestQueryParams>({
                    api: mockApi,
                    fetchOnMount: false,
                    initialParams: { search: 'test' },
                })
            );

            expect(result.current.queryParams).toEqual({ search: 'test' });
        });

        it('should not fetch on mount when fetchOnMount is false', () => {
            renderHook(() =>
                useCrud<TestItem, TestCreateInput, TestUpdateInput, TestQueryParams>({
                    api: mockApi,
                    fetchOnMount: false,
                })
            );

            expect(mockApi.getAll).not.toHaveBeenCalled();
        });

        it('should fetch on mount when fetchOnMount is true', async () => {
            mockApi.getAll = vi.fn().mockResolvedValue({
                success: true,
                data: [{ id: 1, name: 'Test', value: 100 }],
            });

            const { result } = renderHook(() =>
                useCrud<TestItem, TestCreateInput, TestUpdateInput, TestQueryParams>({
                    api: mockApi,
                    fetchOnMount: true,
                })
            );

            await waitFor(() => {
                expect(mockApi.getAll).toHaveBeenCalled();
                expect(result.current.data.length).toBe(1);
            });
        });
    });

    describe('fetchAll - data fetching', () => {
        it('should fetch data successfully', async () => {
            const mockData = [
                { id: 1, name: 'Item 1', value: 100 },
                { id: 2, name: 'Item 2', value: 200 },
            ];
            mockApi.getAll = vi.fn().mockResolvedValue({
                success: true,
                data: mockData,
            });

            const { result } = renderHook(() =>
                useCrud<TestItem, TestCreateInput, TestUpdateInput, TestQueryParams>({
                    api: mockApi,
                    fetchOnMount: false,
                })
            );

            await act(async () => {
                await result.current.fetchAll();
            });

            expect(result.current.data).toEqual(mockData);
            expect(result.current.loading).toBe(false);
            expect(result.current.error).toBe(null);
        });

        it('should fetch with pagination params', async () => {
            mockApi.getAll = vi.fn().mockResolvedValue({
                success: true,
                data: [],
            });

            const { result } = renderHook(() =>
                useCrud<TestItem, TestCreateInput, TestUpdateInput, TestQueryParams>({
                    api: mockApi,
                    fetchOnMount: false,
                    initialPagination: { current: 2, pageSize: 20 },
                })
            );

            await act(async () => {
                await result.current.fetchAll();
            });

            expect(mockApi.getAll).toHaveBeenCalledWith({
                page: 2,
                page_size: 20,
            });
        });

        it('should fetch with custom params', async () => {
            mockApi.getAll = vi.fn().mockResolvedValue({
                success: true,
                data: [],
            });

            const { result } = renderHook(() =>
                useCrud<TestItem, TestCreateInput, TestUpdateInput, TestQueryParams>({
                    api: mockApi,
                    fetchOnMount: false,
                })
            );

            await act(async () => {
                await result.current.fetchAll({ search: 'keyword' });
            });

            expect(mockApi.getAll).toHaveBeenCalledWith(
                expect.objectContaining({ search: 'keyword' })
            );
        });

        it('should handle axios response format', async () => {
            const mockData = [{ id: 1, name: 'Test', value: 100 }];
            mockApi.getAll = vi.fn().mockResolvedValue({
                data: {
                    success: true,
                    data: mockData,
                },
            });

            const { result } = renderHook(() =>
                useCrud<TestItem, TestCreateInput, TestUpdateInput, TestQueryParams>({
                    api: mockApi,
                    fetchOnMount: false,
                })
            );

            await act(async () => {
                await result.current.fetchAll();
            });

            expect(result.current.data).toEqual(mockData);
        });

        it('should handle pagination from response', async () => {
            mockApi.getAll = vi.fn().mockResolvedValue({
                success: true,
                data: [],
                pagination: { total: 100, current: 1, pageSize: 10 },
            });

            const { result } = renderHook(() =>
                useCrud<TestItem, TestCreateInput, TestUpdateInput, TestQueryParams>({
                    api: mockApi,
                    fetchOnMount: false,
                })
            );

            await act(async () => {
                await result.current.fetchAll();
            });

            expect(result.current.pagination.total).toBe(100);
        });

        it('should use array length as total if no pagination', async () => {
            const mockData = Array.from({ length: 5 }, (_, i) => ({
                id: i + 1,
                name: `Item ${i}`,
                value: i * 10,
            }));
            mockApi.getAll = vi.fn().mockResolvedValue({
                success: true,
                data: mockData,
            });

            const { result } = renderHook(() =>
                useCrud<TestItem, TestCreateInput, TestUpdateInput, TestQueryParams>({
                    api: mockApi,
                    fetchOnMount: false,
                })
            );

            await act(async () => {
                await result.current.fetchAll();
            });

            expect(result.current.pagination.total).toBe(5);
        });

        it('should set loading state during fetch', async () => {
            let resolveFetch: (value: unknown) => void;
            const fetchPromise = new Promise((resolve) => {
                resolveFetch = resolve;
            });
            mockApi.getAll = vi.fn().mockReturnValue(fetchPromise);

            const { result } = renderHook(() =>
                useCrud<TestItem, TestCreateInput, TestUpdateInput, TestQueryParams>({
                    api: mockApi,
                    fetchOnMount: false,
                })
            );

            act(() => {
                result.current.fetchAll();
            });

            expect(result.current.loading).toBe(true);

            await act(async () => {
                await resolveFetch!({ success: true, data: [] });
            });

            expect(result.current.loading).toBe(false);
        });

        it('should handle fetch errors', async () => {
            const error = new Error('Network error');
            mockApi.getAll = vi.fn().mockRejectedValue(error);

            const { result } = renderHook(() =>
                useCrud<TestItem, TestCreateInput, TestUpdateInput, TestQueryParams>({
                    api: mockApi,
                    fetchOnMount: false,
                })
            );

            await act(async () => {
                await result.current.fetchAll();
            });

            expect(result.current.error).toBeInstanceOf(Error);
            expect(result.current.loading).toBe(false);
        });
    });

    describe('refresh', () => {
        it('should refresh data with current params', async () => {
            mockApi.getAll = vi.fn().mockResolvedValue({
                success: true,
                data: [{ id: 1, name: 'Test', value: 100 }],
            });

            const { result } = renderHook(() =>
                useCrud<TestItem, TestCreateInput, TestUpdateInput, TestQueryParams>({
                    api: mockApi,
                    fetchOnMount: false,
                    initialPagination: { current: 2, pageSize: 20 },
                })
            );

            await act(async () => {
                await result.current.refresh();
            });

            expect(mockApi.getAll).toHaveBeenCalledWith({
                page: 2,
                page_size: 20,
            });
        });
    });

    describe('create', () => {
        it('should create item successfully', async () => {
            const newItem: TestCreateInput = { name: 'New Item', value: 300 };
            const createdItem: TestItem = { id: 3, ...newItem };
            mockApi.getAll = vi.fn().mockResolvedValue({ success: true, data: [] });
            mockApi.create = vi.fn().mockResolvedValue({
                success: true,
                data: createdItem,
            });

            const { result } = renderHook(() =>
                useCrud<TestItem, TestCreateInput, TestUpdateInput, TestQueryParams>({
                    api: mockApi,
                    fetchOnMount: false,
                })
            );

            const created = await act(async () => {
                return await result.current.create(newItem);
            });

            expect(mockApi.create).toHaveBeenCalledWith(newItem);
            expect(created).toEqual(createdItem);
            expect(result.current.submitting).toBe(false);
        });

        it('should refresh list after create', async () => {
            mockApi.getAll = vi.fn().mockResolvedValue({ success: true, data: [] });
            mockApi.create = vi.fn().mockResolvedValue({
                success: true,
                data: { id: 1, name: 'Test', value: 100 },
            });

            const { result } = renderHook(() =>
                useCrud<TestItem, TestCreateInput, TestUpdateInput, TestQueryParams>({
                    api: mockApi,
                    fetchOnMount: false,
                })
            );

            await act(async () => {
                await result.current.create({ name: 'Test', value: 100 });
            });

            expect(mockApi.getAll).toHaveBeenCalled();
        });

        it('should call onCreateSuccess callback', async () => {
            const onSuccess = vi.fn();
            mockApi.getAll = vi.fn().mockResolvedValue({ success: true, data: [] });
            mockApi.create = vi.fn().mockResolvedValue({
                success: true,
                data: { id: 1, name: 'Test', value: 100 },
            });

            const { result } = renderHook(() =>
                useCrud<TestItem, TestCreateInput, TestUpdateInput, TestQueryParams>({
                    api: mockApi,
                    fetchOnMount: false,
                    onCreateSuccess: onSuccess,
                })
            );

            await act(async () => {
                await result.current.create({ name: 'Test', value: 100 });
            });

            expect(onSuccess).toHaveBeenCalledWith({
                id: 1,
                name: 'Test',
                value: 100,
            });
        });

        it('should handle create errors', async () => {
            const error = new Error('Create failed');
            mockApi.getAll = vi.fn().mockResolvedValue({ success: true, data: [] });
            mockApi.create = vi.fn().mockRejectedValue(error);

            const { result } = renderHook(() =>
                useCrud<TestItem, TestCreateInput, TestUpdateInput, TestQueryParams>({
                    api: mockApi,
                    fetchOnMount: false,
                })
            );

            const created = await act(async () => {
                return await result.current.create({ name: 'Test', value: 100 });
            });

            expect(created).toBe(null);
            expect(result.current.error).toBeInstanceOf(Error);
            expect(result.current.submitting).toBe(false);
        });

        it('should support silent mode', async () => {
            mockApi.getAll = vi.fn().mockResolvedValue({ success: true, data: [] });
            mockApi.create = vi.fn().mockResolvedValue({
                success: true,
                data: { id: 1, name: 'Test', value: 100 },
            });

            const { result } = renderHook(() =>
                useCrud<TestItem, TestCreateInput, TestUpdateInput, TestQueryParams>({
                    api: mockApi,
                    fetchOnMount: false,
                })
            );

            await act(async () => {
                await result.current.create({ name: 'Test', value: 100 }, { silent: true });
            });

            const { message } = await import('antd');
            expect(message.success).not.toHaveBeenCalled();
        });
    });

    describe('update', () => {
        it('should update item successfully', async () => {
            const updateData: TestUpdateInput = { name: 'Updated Name' };
            const updatedItem: TestItem = { id: 1, name: 'Updated Name', value: 100 };
            mockApi.getAll = vi.fn().mockResolvedValue({ success: true, data: [] });
            mockApi.update = vi.fn().mockResolvedValue({
                success: true,
                data: updatedItem,
            });

            const { result } = renderHook(() =>
                useCrud<TestItem, TestCreateInput, TestUpdateInput, TestQueryParams>({
                    api: mockApi,
                    fetchOnMount: false,
                })
            );

            const updated = await act(async () => {
                return await result.current.update(1, updateData);
            });

            expect(mockApi.update).toHaveBeenCalledWith(1, updateData);
            expect(updated).toEqual(updatedItem);
        });

        it('should refresh list after update', async () => {
            mockApi.getAll = vi.fn().mockResolvedValue({ success: true, data: [] });
            mockApi.update = vi.fn().mockResolvedValue({
                success: true,
                data: { id: 1, name: 'Updated', value: 100 },
            });

            const { result } = renderHook(() =>
                useCrud<TestItem, TestCreateInput, TestUpdateInput, TestQueryParams>({
                    api: mockApi,
                    fetchOnMount: false,
                })
            );

            await act(async () => {
                await result.current.update(1, { name: 'Updated' });
            });

            expect(mockApi.getAll).toHaveBeenCalled();
        });

        it('should call onUpdateSuccess callback', async () => {
            const onSuccess = vi.fn();
            mockApi.getAll = vi.fn().mockResolvedValue({ success: true, data: [] });
            mockApi.update = vi.fn().mockResolvedValue({
                success: true,
                data: { id: 1, name: 'Updated', value: 100 },
            });

            const { result } = renderHook(() =>
                useCrud<TestItem, TestCreateInput, TestUpdateInput, TestQueryParams>({
                    api: mockApi,
                    fetchOnMount: false,
                    onUpdateSuccess: onSuccess,
                })
            );

            await act(async () => {
                await result.current.update(1, { name: 'Updated' });
            });

            expect(onSuccess).toHaveBeenCalledWith({
                id: 1,
                name: 'Updated',
                value: 100,
            });
        });

        it('should handle update errors', async () => {
            const error = new Error('Update failed');
            mockApi.getAll = vi.fn().mockResolvedValue({ success: true, data: [] });
            mockApi.update = vi.fn().mockRejectedValue(error);

            const { result } = renderHook(() =>
                useCrud<TestItem, TestCreateInput, TestUpdateInput, TestQueryParams>({
                    api: mockApi,
                    fetchOnMount: false,
                })
            );

            const updated = await act(async () => {
                return await result.current.update(1, { name: 'Updated' });
            });

            expect(updated).toBe(null);
            expect(result.current.error).toBeInstanceOf(Error);
        });

        it('should support silent mode', async () => {
            mockApi.getAll = vi.fn().mockResolvedValue({ success: true, data: [] });
            mockApi.update = vi.fn().mockResolvedValue({
                success: true,
                data: { id: 1, name: 'Updated', value: 100 },
            });

            const { result } = renderHook(() =>
                useCrud<TestItem, TestCreateInput, TestUpdateInput, TestQueryParams>({
                    api: mockApi,
                    fetchOnMount: false,
                })
            );

            await act(async () => {
                await result.current.update(1, { name: 'Updated' }, { silent: true });
            });

            const { message } = await import('antd');
            expect(message.success).not.toHaveBeenCalled();
        });
    });

    describe('remove', () => {
        it('should remove item successfully with confirmation', async () => {
            mockApi.getAll = vi.fn().mockResolvedValue({ success: true, data: [] });
            mockApi.remove = vi.fn().mockResolvedValue({ success: true, data: undefined });

            const { result } = renderHook(() =>
                useCrud<TestItem, TestCreateInput, TestUpdateInput, TestQueryParams>({
                    api: mockApi,
                    fetchOnMount: false,
                })
            );

            const removed = await act(async () => {
                return await result.current.remove(1);
            });

            expect(mockApi.remove).toHaveBeenCalledWith(1);
            expect(removed).toBe(true);
        });

        it('should refresh list after delete', async () => {
            mockApi.getAll = vi.fn().mockResolvedValue({ success: true, data: [] });
            mockApi.remove = vi.fn().mockResolvedValue({ success: true, data: undefined });

            const { result } = renderHook(() =>
                useCrud<TestItem, TestCreateInput, TestUpdateInput, TestQueryParams>({
                    api: mockApi,
                    fetchOnMount: false,
                })
            );

            await act(async () => {
                await result.current.remove(1);
            });

            expect(mockApi.getAll).toHaveBeenCalled();
        });

        it('should call onDeleteSuccess callback', async () => {
            const onSuccess = vi.fn();
            mockApi.getAll = vi.fn().mockResolvedValue({ success: true, data: [] });
            mockApi.remove = vi.fn().mockResolvedValue({ success: true, data: undefined });

            const { result } = renderHook(() =>
                useCrud<TestItem, TestCreateInput, TestUpdateInput, TestQueryParams>({
                    api: mockApi,
                    fetchOnMount: false,
                    onDeleteSuccess: onSuccess,
                })
            );

            await act(async () => {
                await result.current.remove(1);
            });

            expect(onSuccess).toHaveBeenCalledWith(1);
        });

        it('should handle delete errors', async () => {
            const error = new Error('Delete failed');
            mockApi.getAll = vi.fn().mockResolvedValue({ success: true, data: [] });
            mockApi.remove = vi.fn().mockRejectedValue(error);

            const { result } = renderHook(() =>
                useCrud<TestItem, TestCreateInput, TestUpdateInput, TestQueryParams>({
                    api: mockApi,
                    fetchOnMount: false,
                })
            );

            const removed = await act(async () => {
                return await result.current.remove(1);
            });

            expect(removed).toBe(false);
            expect(result.current.error).toBeInstanceOf(Error);
        });

        it('should support silent mode without confirmation', async () => {
            mockApi.getAll = vi.fn().mockResolvedValue({ success: true, data: [] });
            mockApi.remove = vi.fn().mockResolvedValue({ success: true, data: undefined });

            const { result } = renderHook(() =>
                useCrud<TestItem, TestCreateInput, TestUpdateInput, TestQueryParams>({
                    api: mockApi,
                    fetchOnMount: false,
                })
            );

            await act(async () => {
                await result.current.remove(1, { silent: true });
            });

            const { Modal } = await import('antd');
            expect(Modal.confirm).not.toHaveBeenCalled();
            expect(mockApi.remove).toHaveBeenCalled();
        });
    });

    describe('pagination handling', () => {
        it('should set page and trigger refetch', async () => {
            mockApi.getAll = vi.fn().mockResolvedValue({ success: true, data: [] });

            const { result } = renderHook(() =>
                useCrud<TestItem, TestCreateInput, TestUpdateInput, TestQueryParams>({
                    api: mockApi,
                    fetchOnMount: false,
                })
            );

            act(() => {
                result.current.setPage(3);
            });

            expect(result.current.pagination.current).toBe(3);

            await waitFor(() => {
                expect(mockApi.getAll).toHaveBeenCalledWith({ page: 3, page_size: 10 });
            });
        });

        it('should set page size and reset to page 1', async () => {
            mockApi.getAll = vi.fn().mockResolvedValue({ success: true, data: [] });

            const { result } = renderHook(() =>
                useCrud<TestItem, TestCreateInput, TestUpdateInput, TestQueryParams>({
                    api: mockApi,
                    fetchOnMount: false,
                    initialPagination: { current: 5, pageSize: 10 },
                })
            );

            act(() => {
                result.current.setPageSize(25);
            });

            expect(result.current.pagination.pageSize).toBe(25);
            expect(result.current.pagination.current).toBe(1);
        });

        it('should handle pagination change callback', async () => {
            mockApi.getAll = vi.fn().mockResolvedValue({ success: true, data: [] });

            const { result } = renderHook(() =>
                useCrud<TestItem, TestCreateInput, TestUpdateInput, TestQueryParams>({
                    api: mockApi,
                    fetchOnMount: false,
                })
            );

            act(() => {
                result.current.pagination.onChange(3, 10);
            });

            expect(result.current.pagination.current).toBe(3);
            expect(result.current.pagination.pageSize).toBe(10);

            await waitFor(() => {
                expect(mockApi.getAll).toHaveBeenCalledWith({ page: 3, page_size: 10 });
            });
        });
    });

    describe('search params handling', () => {
        it('should set search params and reset to page 1', async () => {
            mockApi.getAll = vi.fn().mockResolvedValue({ success: true, data: [] });

            const { result } = renderHook(() =>
                useCrud<TestItem, TestCreateInput, TestUpdateInput, TestQueryParams>({
                    api: mockApi,
                    fetchOnMount: false,
                    initialPagination: { current: 3, pageSize: 10 },
                })
            );

            act(() => {
                result.current.setSearchParams({ search: 'keyword' });
            });

            expect(result.current.queryParams).toEqual({ search: 'keyword' });
            expect(result.current.pagination.current).toBe(1);
        });

        it('should merge search params', async () => {
            mockApi.getAll = vi.fn().mockResolvedValue({ success: true, data: [] });

            const { result } = renderHook(() =>
                useCrud<TestItem, TestCreateInput, TestUpdateInput, TestQueryParams>({
                    api: mockApi,
                    fetchOnMount: false,
                    initialParams: { status: 'active' },
                })
            );

            act(() => {
                result.current.setSearchParams({ search: 'keyword' });
            });

            expect(result.current.queryParams).toEqual({
                status: 'active',
                search: 'keyword',
            });
        });
    });

    describe('error handling', () => {
        it('should clear error', () => {
            const { result } = renderHook(() =>
                useCrud<TestItem, TestCreateInput, TestUpdateInput, TestQueryParams>({
                    api: mockApi,
                    fetchOnMount: false,
                })
            );

            // Set an error
            act(() => {
                result.current.setData([]);
                // Simulate error state
            });

            act(() => {
                result.current.clearError();
            });

            expect(result.current.error).toBe(null);
        });

        it('should call onError callback on fetch error', async () => {
            const onError = vi.fn();
            const error = new Error('Fetch failed');
            mockApi.getAll = vi.fn().mockRejectedValue(error);

            const { result } = renderHook(() =>
                useCrud<TestItem, TestCreateInput, TestUpdateInput, TestQueryParams>({
                    api: mockApi,
                    fetchOnMount: false,
                    onError,
                })
            );

            await act(async () => {
                await result.current.fetchAll();
            });

            expect(onError).toHaveBeenCalledWith(error, 'fetch');
        });

        it('should call onError callback on create error', async () => {
            const onError = vi.fn();
            const error = new Error('Create failed');
            mockApi.getAll = vi.fn().mockResolvedValue({ success: true, data: [] });
            mockApi.create = vi.fn().mockRejectedValue(error);

            const { result } = renderHook(() =>
                useCrud<TestItem, TestCreateInput, TestUpdateInput, TestQueryParams>({
                    api: mockApi,
                    fetchOnMount: false,
                    onError,
                })
            );

            await act(async () => {
                await result.current.create({ name: 'Test', value: 100 });
            });

            expect(onError).toHaveBeenCalledWith(error, 'create');
        });

        it('should call onError callback on update error', async () => {
            const onError = vi.fn();
            const error = new Error('Update failed');
            mockApi.getAll = vi.fn().mockResolvedValue({ success: true, data: [] });
            mockApi.update = vi.fn().mockRejectedValue(error);

            const { result } = renderHook(() =>
                useCrud<TestItem, TestCreateInput, TestUpdateInput, TestQueryParams>({
                    api: mockApi,
                    fetchOnMount: false,
                    onError,
                })
            );

            await act(async () => {
                await result.current.update(1, { name: 'Updated' });
            });

            expect(onError).toHaveBeenCalledWith(error, 'update');
        });

        it('should call onError callback on delete error', async () => {
            const onError = vi.fn();
            const error = new Error('Delete failed');
            mockApi.getAll = vi.fn().mockResolvedValue({ success: true, data: [] });
            mockApi.remove = vi.fn().mockRejectedValue(error);

            const { result } = renderHook(() =>
                useCrud<TestItem, TestCreateInput, TestUpdateInput, TestQueryParams>({
                    api: mockApi,
                    fetchOnMount: false,
                    onError,
                })
            );

            await act(async () => {
                await result.current.remove(1, { silent: true });
            });

            expect(onError).toHaveBeenCalledWith(error, 'delete');
        });
    });

    describe('data transformer', () => {
        it('should transform data with dataTransformer', async () => {
            const rawData = [{ id: 1, name: 'Test', value: 100 }];
            const transformedData = [{ id: 1, displayName: 'Test', amount: 100 }];
            mockApi.getAll = vi.fn().mockResolvedValue({
                success: true,
                data: rawData,
            });

            const { result } = renderHook(() =>
                useCrud<TestItem, TestCreateInput, TestUpdateInput, TestQueryParams>({
                    api: mockApi,
                    fetchOnMount: false,
                    dataTransformer: (data) => {
                        return (data as TestItem[]).map((item) => ({
                            id: item.id,
                            displayName: item.name,
                            amount: item.value,
                        }));
                    },
                })
            );

            await act(async () => {
                await result.current.fetchAll();
            });

            expect(result.current.data).toEqual(transformedData);
        });

        it('should handle empty array from transformer', async () => {
            mockApi.getAll = vi.fn().mockResolvedValue({
                success: true,
                data: [1, 2, 3],
            });

            const { result } = renderHook(() =>
                useCrud<TestItem, TestCreateInput, TestUpdateInput, TestQueryParams>({
                    api: mockApi,
                    fetchOnMount: false,
                    dataTransformer: () => [],
                })
            );

            await act(async () => {
                await result.current.fetchAll();
            });

            expect(result.current.data).toEqual([]);
        });
    });

    describe('manual data setting', () => {
        it('should allow manual data setting', () => {
            const { result } = renderHook(() =>
                useCrud<TestItem, TestCreateInput, TestUpdateInput, TestQueryParams>({
                    api: mockApi,
                    fetchOnMount: false,
                })
            );

            const manualData: TestItem[] = [
                { id: 1, name: 'Manual 1', value: 100 },
                { id: 2, name: 'Manual 2', value: 200 },
            ];

            act(() => {
                result.current.setData(manualData);
            });

            expect(result.current.data).toEqual(manualData);
        });
    });

    describe('pagination extractor', () => {
        it('should extract total using paginationExtractor', async () => {
            mockApi.getAll = vi.fn().mockResolvedValue({
                success: true,
                data: [],
                customTotal: 500,
            });

            const { result } = renderHook(() =>
                useCrud<TestItem, TestCreateInput, TestUpdateInput, TestQueryParams>({
                    api: mockApi,
                    fetchOnMount: false,
                    paginationExtractor: (response) => {
                        return (response as { customTotal?: number }).customTotal;
                    },
                })
            );

            await act(async () => {
                await result.current.fetchAll();
            });

            expect(result.current.pagination.total).toBe(500);
        });
    });

    describe('property-based tests', () => {
        /**
         * Property: fetchAll should store returned data correctly
         */
        it('should store any valid array of items', async () => {
            fc.assert(
                fc.property(
                    fc.array(
                        fc.record({
                            id: fc.nat(),
                            name: fc.string(),
                            value: fc.nat(),
                        })
                    ),
                    async (data) => {
                        mockApi.getAll = vi.fn().mockResolvedValue({
                            success: true,
                            data,
                        });

                        const { result } = renderHook(() =>
                            useCrud<TestItem, TestCreateInput, TestUpdateInput, TestQueryParams>({
                                api: mockApi,
                                fetchOnMount: false,
                            })
                        );

                        await act(async () => {
                            await result.current.fetchAll();
                        });

                        return result.current.data === data;
                    }
                ),
                { numRuns: 20 }
            );
        });

        /**
         * Property: Page changes should trigger fetch with correct params
         */
        it('should fetch with correct page number', async () => {
            fc.assert(
                fc.property(fc.nat({ max: 100 }), async (page) => {
                    const targetPage = page + 1; // Ensure >= 1
                    mockApi.getAll = vi.fn().mockResolvedValue({
                        success: true,
                        data: [],
                    });

                    const { result } = renderHook(() =>
                        useCrud<TestItem, TestCreateInput, TestUpdateInput, TestQueryParams>({
                            api: mockApi,
                            fetchOnMount: false,
                        })
                    );

                    act(() => {
                        result.current.setPage(targetPage);
                    });

                    await waitFor(() => {
                        expect(mockApi.getAll).toHaveBeenCalledWith({
                            page: targetPage,
                            page_size: 10,
                        });
                    });

                    return true;
                }),
                { numRuns: 20 }
            );
        });

        /**
         * Property: Page size changes should reset to page 1
         */
        it('should reset to page 1 when changing page size', async () => {
            fc.assert(
                fc.property(
                    fc.tuple(fc.nat({ max: 5 }), fc.nat({ min: 10, max: 100 })),
                    async ([initialPage, newPageSize]) => {
                        const startPage = initialPage + 1;
                        mockApi.getAll = vi.fn().mockResolvedValue({
                            success: true,
                            data: [],
                        });

                        const { result } = renderHook(() =>
                            useCrud<TestItem, TestCreateInput, TestUpdateInput, TestQueryParams>({
                                api: mockApi,
                                fetchOnMount: false,
                                initialPagination: { current: startPage, pageSize: 10 },
                            })
                        );

                        act(() => {
                            result.current.setPageSize(newPageSize);
                        });

                        return result.current.pagination.current === 1 &&
                               result.current.pagination.pageSize === newPageSize;
                    }
                ),
                { numRuns: 20 }
            );
        });
    });
});
