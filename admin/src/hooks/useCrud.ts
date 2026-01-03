/**
 * Generic CRUD Hook
 *
 * A reusable React hook for managing CRUD operations with:
 * - Automatic data fetching and pagination
 * - Loading states and error handling
 * - Success/error messages via Ant Design message component
 * - Type-safe API integration
 *
 * @example
 * ```tsx
 * // Define your API interface
 * interface UserCrudApi {
 *   getAll: (params?: UserQueryParams) => Promise<ApiResponse<User[]>>;
 *   create: (data: CreateUserDto) => Promise<ApiResponse<User>>;
 *   update: (id: number, data: UpdateUserDto) => Promise<ApiResponse<User>>;
 *   remove: (id: number) => Promise<ApiResponse<void>>;
 * }
 *
 * // Use the hook in your component
 * const { data, loading, pagination, fetchAll, create, update, remove } = useCrud({
 *   api: userCrudApi,
 *   messages: {
 *     fetchError: '获取用户列表失败',
 *     createSuccess: '创建用户成功',
 *     updateSuccess: '更新用户成功',
 *     deleteSuccess: '删除用户成功',
 *   }
 * });
 * ```
 */

import { useState, useEffect, useCallback, useRef } from 'react';
import { message, Modal } from 'antd';
import type { ApiResponse, Pagination } from '@/types/api';
import { logger } from '@/utils/logger';

/**
 * Generic ID type - supports both number and string IDs
 */
export type CrudId = number | string;

/**
 * Query parameters for list operations
 */
export interface CrudQueryParams {
    page?: number;
    page_size?: number;
    [key: string]: unknown;
}

/**
 * Pagination state for table components
 */
export interface CrudPagination {
    current: number;
    pageSize: number;
    total: number;
    showSizeChanger?: boolean;
    showQuickJumper?: boolean;
    showTotal?: (total: number) => string;
    onChange: (page: number, pageSize: number) => void;
}

/**
 * Configuration options for useCrud hook
 */
export interface UseCrudOptions<T extends object, TCreate, TUpdate, TQuery extends CrudQueryParams> {
    /**
     * API functions for CRUD operations
     */
    api: CrudApi<T, TCreate, TUpdate, TQuery>;

    /**
     * Custom messages for user feedback
     */
    messages?: CrudMessages;

    /**
     * Initial query parameters
     */
    initialParams?: TQuery;

    /**
     * Initial pagination state
     */
    initialPagination?: {
        current?: number;
        pageSize?: number;
    };

    /**
     * Whether to fetch data on mount (default: true)
     */
    fetchOnMount?: boolean;

    /**
     * Custom success callback after create
     */
    onCreateSuccess?: (data: T) => void;

    /**
     * Custom success callback after update
     */
    onUpdateSuccess?: (data: T) => void;

    /**
     * Custom success callback after delete
     */
    onDeleteSuccess?: (id: CrudId) => void;

    /**
     * Custom error callback
     */
    onError?: (error: unknown, operation: CrudOperation) => void;

    /**
     * Transform function for API response data
     */
    dataTransformer?: (rawData: unknown) => T[];

    /**
     * Extract pagination from response
     */
    paginationExtractor?: (response: unknown) => number | undefined;
}

/**
 * CRUD operation types
 */
export type CrudOperation = 'fetch' | 'create' | 'update' | 'delete';

/**
 * Message configuration for CRUD operations
 */
export interface CrudMessages {
    fetchError?: string;
    createSuccess?: string;
    createError?: string;
    updateSuccess?: string;
    updateError?: string;
    deleteSuccess?: string;
    deleteError?: string;
    deleteConfirm?: string;
}

/**
 * API interface for CRUD operations
 * Compatible with axios response format
 */
export interface CrudApi<T extends object, TCreate, TUpdate, TQuery extends CrudQueryParams> {
    /**
     * Fetch all items with optional query parameters
     * Accepts: Promise<ApiResponse<T[]>> | Promise<AxiosResponse<ApiResponse<T[]>>>
     */
    getAll: (params?: TQuery) => Promise<ApiResponse<T[]> | { data: ApiResponse<T[]> }>;

    /**
     * Create a new item
     * Accepts: Promise<ApiResponse<T>> | Promise<AxiosResponse<ApiResponse<T>>>
     */
    create: (data: TCreate) => Promise<ApiResponse<T> | { data: ApiResponse<T> }>;

    /**
     * Update an existing item
     * Accepts: Promise<ApiResponse<T>> | Promise<AxiosResponse<ApiResponse<T>>>
     */
    update: (id: CrudId, data: TUpdate) => Promise<ApiResponse<T> | { data: ApiResponse<T> }>;

    /**
     * Delete an item
     * Accepts: Promise<ApiResponse<void>> | Promise<AxiosResponse<ApiResponse<void>>>
     */
    remove: (id: CrudId) => Promise<ApiResponse<void> | { data: ApiResponse<void> }>;
}

/**
 * Return type for useCrud hook
 */
export interface UseCrudReturn<T extends object, TCreate, TUpdate, TQuery extends CrudQueryParams> {
    /**
     * List of items
     */
    data: T[];

    /**
     * Loading state for list operations
     */
    loading: boolean;

    /**
     * Loading state for create/update operations
     */
    submitting: boolean;

    /**
     * Error state
     */
    error: Error | null;

    /**
     * Pagination state and handlers
     */
    pagination: CrudPagination;

    /**
     * Current query parameters
     */
    queryParams: Record<string, unknown>;

    /**
     * Fetch all items with optional parameters
     */
    fetchAll: (params?: TQuery | Record<string, unknown>) => Promise<void>;

    /**
     * Fetch all items and reset to page 1
     */
    refresh: () => Promise<void>;

    /**
     * Create a new item
     */
    create: (item: TCreate, options?: { silent?: boolean }) => Promise<T | null>;

    /**
     * Update an existing item
     */
    update: (id: CrudId, item: TUpdate, options?: { silent?: boolean }) => Promise<T | null>;

    /**
     * Delete an item
     */
    remove: (id: CrudId, options?: { silent?: boolean; confirmMessage?: string }) => Promise<boolean>;

    /**
     * Set current page
     */
    setPage: (page: number) => void;

    /**
     * Set page size
     */
    setPageSize: (pageSize: number) => void;

    /**
     * Update query parameters and refetch
     */
    setSearchParams: (params: Record<string, unknown>) => void;

    /**
     * Clear error state
     */
    clearError: () => void;

    /**
     * Manually set data (useful for optimistic updates)
     */
    setData: (data: T[]) => void;
}

/**
 * Default messages for CRUD operations
 */
const defaultMessages: Required<CrudMessages> = {
    fetchError: '加载数据失败',
    createSuccess: '创建成功',
    createError: '创建失败',
    updateSuccess: '更新成功',
    updateError: '更新失败',
    deleteSuccess: '删除成功',
    deleteError: '删除失败',
    deleteConfirm: '确定要删除吗？此操作不可恢复。',
};

/**
 * Generic CRUD hook
 *
 * @template T - Item type
 * @template TCreate - Create DTO type
 * @template TUpdate - Update DTO type
 * @template TQuery - Query parameters type
 *
 * @param options - Hook configuration options
 * @returns CRUD operations and state
 */
export function useCrud<
    T extends object,
    TCreate = Partial<T>,
    TUpdate = Partial<T>,
    TQuery extends CrudQueryParams = CrudQueryParams
>(options: UseCrudOptions<T, TCreate, TUpdate, TQuery>): UseCrudReturn<T, TCreate, TUpdate, TQuery> {
    const {
        api,
        messages = {},
        initialParams = {} as TQuery,
        initialPagination = {},
        fetchOnMount = true,
        onCreateSuccess,
        onUpdateSuccess,
        onDeleteSuccess,
        onError,
        dataTransformer,
        paginationExtractor,
    } = options;

    const mergedMessages = { ...defaultMessages, ...messages };

    // State
    const [data, setData] = useState<T[]>([]);
    const [loading, setLoading] = useState(false);
    const [submitting, setSubmitting] = useState(false);
    const [error, setError] = useState<Error | null>(null);
    const [current, setCurrent] = useState(initialPagination.current || 1);
    const [pageSize, setPageSize] = useState(initialPagination.pageSize || 10);
    const [total, setTotal] = useState(0);
    const [queryParams, setQueryParams] = useState<Record<string, unknown>>(initialParams);

    // Track if component is mounted
    const isMountedRef = useRef(true);

    /**
     * Normalize API response to extract data and pagination
     */
    const normalizeResponse = useCallback((response: ApiResponse<T[]> | { data: ApiResponse<T[]> }) => {
        // Handle axios response format { data: { success, data, pagination } }
        if ('data' in response && typeof response.data === 'object' && 'data' in response.data) {
            return {
                data: response.data.data as unknown,
                pagination: (response.data as unknown as { pagination?: Pagination }).pagination,
            };
        }
        // Handle direct API response format { success, data, pagination }
        return {
            data: (response as unknown as { data?: unknown }).data,
            pagination: (response as unknown as { pagination?: Pagination }).pagination,
        };
    }, []);

    /**
     * Normalize single item API response
     */
    const normalizeSingleResponse = useCallback((response: ApiResponse<T> | { data: ApiResponse<T> }) => {
        // Handle axios response format { data: { success, data } }
        if ('data' in response && typeof response.data === 'object' && 'data' in response.data) {
            return response.data.data as T;
        }
        // Handle direct API response format { success, data }
        return (response as unknown as { data?: T }).data as T;
    }, []);

    /**
     * Fetch all items
     */
    const fetchAll = useCallback(
        async (params?: TQuery | Record<string, unknown>) => {
            if (!isMountedRef.current) return;

            setLoading(true);
            setError(null);

            try {
                const mergedParams: TQuery = {
                    ...queryParams,
                    page: current,
                    page_size: pageSize,
                    ...(params as TQuery),
                };

                const response = await api.getAll(mergedParams);
                const { data: responseData, pagination: responsePagination } = normalizeResponse(response);

                if (isMountedRef.current) {
                    // Transform data if transformer provided
                    const items = dataTransformer
                        ? dataTransformer(responseData)
                        : (responseData as T[]);

                    setData(Array.isArray(items) ? items : []);

                    // Extract total count
                    if (paginationExtractor) {
                        const extractedTotal = paginationExtractor(response);
                        setTotal(extractedTotal || 0);
                    } else if (responsePagination?.total !== undefined) {
                        setTotal(responsePagination.total);
                    } else if (Array.isArray(items)) {
                        setTotal(items.length);
                    } else {
                        setTotal(0);
                    }
                }
            } catch (err) {
                if (isMountedRef.current) {
                    const errorObj = err instanceof Error ? err : new Error(String(err));
                    setError(errorObj);
                    logger.error('CRUD fetch error:', err);
                    if (!messages?.fetchError) {
                        message.error(mergedMessages.fetchError);
                    }
                    onError?.(err, 'fetch');
                }
            } finally {
                if (isMountedRef.current) {
                    setLoading(false);
                }
            }
        },
        [api, current, pageSize, queryParams, normalizeResponse, dataTransformer, paginationExtractor, mergedMessages, messages, onError]
    );

    /**
     * Refresh data (fetch with current params)
     */
    const refresh = useCallback(async () => {
        await fetchAll();
    }, [fetchAll]);

    /**
     * Create a new item
     */
    const create = useCallback(
        async (item: TCreate, options?: { silent?: boolean }): Promise<T | null> => {
            setSubmitting(true);
            setError(null);

            try {
                const response = await api.create(item);
                const createdItem = normalizeSingleResponse(response);

                if (!options?.silent) {
                    message.success(mergedMessages.createSuccess);
                }

                // Refresh the list
                await fetchAll();

                // Call custom success callback
                if (createdItem) {
                    onCreateSuccess?.(createdItem);
                }

                return createdItem || null;
            } catch (err) {
                const errorObj = err instanceof Error ? err : new Error(String(err));
                setError(errorObj);
                logger.error('CRUD create error:', err);
                if (!options?.silent) {
                    message.error(mergedMessages.createError);
                }
                onError?.(err, 'create');
                return null;
            } finally {
                setSubmitting(false);
            }
        },
        [api, fetchAll, mergedMessages, onCreateSuccess, onError, normalizeSingleResponse]
    );

    /**
     * Update an existing item
     */
    const update = useCallback(
        async (id: CrudId, item: TUpdate, options?: { silent?: boolean }): Promise<T | null> => {
            setSubmitting(true);
            setError(null);

            try {
                const response = await api.update(id, item);
                const updatedItem = normalizeSingleResponse(response);

                if (!options?.silent) {
                    message.success(mergedMessages.updateSuccess);
                }

                // Refresh the list
                await fetchAll();

                // Call custom success callback
                if (updatedItem) {
                    onUpdateSuccess?.(updatedItem);
                }

                return updatedItem || null;
            } catch (err) {
                const errorObj = err instanceof Error ? err : new Error(String(err));
                setError(errorObj);
                logger.error('CRUD update error:', err);
                if (!options?.silent) {
                    message.error(mergedMessages.updateError);
                }
                onError?.(err, 'update');
                return null;
            } finally {
                setSubmitting(false);
            }
        },
        [api, fetchAll, mergedMessages, onUpdateSuccess, onError, normalizeSingleResponse]
    );

    /**
     * Delete an item
     */
    const remove = useCallback(
        async (id: CrudId, options?: { silent?: boolean; confirmMessage?: string }): Promise<boolean> => {
            const confirmMsg = options?.confirmMessage || mergedMessages.deleteConfirm;

            // Show confirmation unless silent
            if (!options?.silent) {
                return new Promise<boolean>((resolve) => {
                    Modal.confirm({
                        title: '确认删除',
                        content: confirmMsg,
                        okText: '确认',
                        cancelText: '取消',
                        okButtonProps: { danger: true },
                        onOk: async () => {
                            setSubmitting(true);
                            setError(null);

                            try {
                                await api.remove(id);

                                if (!options?.silent) {
                                    message.success(mergedMessages.deleteSuccess);
                                }

                                // Refresh the list
                                await fetchAll();

                                // Call custom success callback
                                onDeleteSuccess?.(id);

                                resolve(true);
                            } catch (err) {
                                const errorObj = err instanceof Error ? err : new Error(String(err));
                                setError(errorObj);
                                logger.error('CRUD delete error:', err);
                                if (!options?.silent) {
                                    message.error(mergedMessages.deleteError);
                                }
                                onError?.(err, 'delete');
                                resolve(false);
                            } finally {
                                setSubmitting(false);
                            }
                        },
                        onCancel: () => resolve(false),
                    });
                });
            }

            // Silent delete without confirmation
            setSubmitting(true);
            setError(null);

            try {
                await api.remove(id);
                await fetchAll();
                onDeleteSuccess?.(id);
                return true;
            } catch (err) {
                const errorObj = err instanceof Error ? err : new Error(String(err));
                setError(errorObj);
                logger.error('CRUD delete error:', err);
                if (!options?.silent) {
                    message.error(mergedMessages.deleteError);
                }
                onError?.(err, 'delete');
                return false;
            } finally {
                setSubmitting(false);
            }
        },
        [api, fetchAll, mergedMessages, onDeleteSuccess, onError]
    );

    /**
     * Set page and refresh
     */
    const setPage = useCallback(
        (page: number) => {
            setCurrent(page);
        },
        []
    );

    /**
     * Set page size and refresh
     */
    const setPageSizeHandler = useCallback(
        (size: number) => {
            setPageSize(size);
            setCurrent(1); // Reset to first page when changing page size
        },
        []
    );

    /**
     * Update query parameters and refresh
     */
    const setSearchParams = useCallback(
        (params: Record<string, unknown>) => {
            setQueryParams(prev => ({ ...prev, ...params }));
            setCurrent(1); // Reset to first page when changing search params
        },
        []
    );

    /**
     * Clear error state
     */
    const clearError = useCallback(() => {
        setError(null);
    }, []);

    /**
     * Pagination handler
     */
    const handlePaginationChange = useCallback(
        (page: number, size: number) => {
            setCurrent(page);
            if (size !== pageSize) {
                setPageSize(size);
            }
        },
        [pageSize]
    );

    /**
     * Pagination object for Ant Design Table
     */
    const pagination: CrudPagination = {
        current,
        pageSize,
        total,
        showSizeChanger: true,
        showQuickJumper: true,
        showTotal: (total: number) => `共 ${total} 条`,
        onChange: handlePaginationChange,
    };

    /**
     * Fetch data on mount
     */
    useEffect(() => {
        if (fetchOnMount) {
            fetchAll();
        }

        return () => {
            isMountedRef.current = false;
        };
    }, [current, pageSize, queryParams]); // Only refetch when these dependencies change

    /**
     * Fetch data when pagination or query params change
     */
    useEffect(() => {
        if (fetchOnMount) {
            fetchAll();
        }
    }, [current, pageSize, queryParams]); // Note: fetchAll is not in dependencies to avoid infinite loop

    return {
        data,
        loading,
        submitting,
        error,
        pagination,
        queryParams,
        fetchAll,
        refresh,
        create,
        update,
        remove,
        setPage,
        setPageSize: setPageSizeHandler,
        setSearchParams,
        clearError,
        setData,
    };
}
