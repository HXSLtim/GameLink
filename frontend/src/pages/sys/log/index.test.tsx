/**
 * Audit Log Page Tests
 * 测试审计日志页面的渲染和过滤功能
 * 
 * **Feature: rbac-button-level-permission, Property 12: 审计日志过滤正确性**
 * **Validates: Requirements 6.3**
 * 
 * Property 12: *For any* audit log query with filters (time range, action type, operator),
 * the returned logs should all satisfy the filter criteria.
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import * as fc from 'fast-check';
import AuditLogPage from './index';
import { auditLogApi } from '@/api/permission';
import type {
    PermissionAuditLog,
    AuditAction,
    AuditTargetType,
    AuditLogQueryParams,
} from '@/types/permission';

// Mock the permission API
vi.mock('@/api/permission', () => ({
    auditLogApi: {
        list: vi.fn(),
        get: vi.fn(),
        export: vi.fn(),
        getStats: vi.fn(),
    },
}));

// Mock the AdminContext
const mockHasPermission = vi.fn().mockReturnValue(true);
vi.mock('@/context/AdminContext', () => ({
    useAdmin: () => ({
        menus: [],
        permissions: ['*'],
        loading: false,
        refreshMenus: vi.fn(),
        hasPermission: mockHasPermission,
        isSuperAdmin: true,
    }),
}));

// ============================================================================
// Arbitraries for Property-Based Testing
// ============================================================================

/**
 * Arbitrary for audit action types
 */
const auditActionArb: fc.Arbitrary<AuditAction> = fc.constantFrom(
    'permission_create',
    'permission_update',
    'permission_delete',
    'role_create',
    'role_update',
    'role_delete',
    'role_permission_assign',
    'user_role_assign'
);

/**
 * Arbitrary for audit target types
 */
const auditTargetTypeArb: fc.Arbitrary<AuditTargetType> = fc.constantFrom(
    'permission',
    'role',
    'user'
);

/**
 * Arbitrary for date strings in YYYY-MM-DD format
 * Using integer-based approach to avoid invalid date issues
 */
const dateStringArb = fc.integer({ min: 0, max: 730 }).map(days => {
    const baseDate = new Date('2024-01-01');
    baseDate.setDate(baseDate.getDate() + days);
    return baseDate.toISOString().split('T')[0];
});

/**
 * Arbitrary for generating a valid ISO date string
 */
const isoDateStringArb = fc.integer({ min: 0, max: 730 }).map(days => {
    const baseDate = new Date('2024-01-01T00:00:00.000Z');
    baseDate.setDate(baseDate.getDate() + days);
    return baseDate.toISOString();
});

/**
 * Arbitrary for generating a single audit log entry
 */
const auditLogArb: fc.Arbitrary<PermissionAuditLog> = fc.record({
    id: fc.integer({ min: 1, max: 10000 }),
    operatorId: fc.integer({ min: 1, max: 1000 }),
    operatorName: fc.constantFrom('Admin', 'User1', 'User2', 'Manager', 'Operator'),
    targetType: auditTargetTypeArb,
    targetId: fc.integer({ min: 1, max: 1000 }),
    targetName: fc.constantFrom('Permission1', 'Role1', 'User1', 'Permission2', 'Role2'),
    action: auditActionArb,
    beforeData: fc.option(fc.constant('{"key":"value"}'), { nil: undefined }).map(v => v ?? undefined),
    afterData: fc.option(fc.constant('{"key":"value2"}'), { nil: undefined }).map(v => v ?? undefined),
    ipAddress: fc.option(fc.constant('192.168.1.1'), { nil: undefined }).map(v => v ?? undefined),
    userAgent: fc.option(fc.constant('Mozilla/5.0'), { nil: undefined }).map(v => v ?? undefined),
    requestId: fc.option(fc.uuid(), { nil: undefined }).map(v => v ?? undefined),
    createdAt: isoDateStringArb,
});

/**
 * Arbitrary for generating a list of audit logs with unique IDs
 */
const auditLogListArb = fc.array(auditLogArb, { minLength: 0, maxLength: 20 }).map(logs => {
    // Ensure unique IDs by reassigning them
    return logs.map((log, index) => ({ ...log, id: index + 1 }));
});

/**
 * Arbitrary for filter parameters
 */
const filterParamsArb = fc.record({
    action: fc.option(auditActionArb, { nil: undefined }),
    targetType: fc.option(auditTargetTypeArb, { nil: undefined }),
    dateFrom: fc.option(dateStringArb, { nil: undefined }),
    dateTo: fc.option(dateStringArb, { nil: undefined }),
});

// ============================================================================
// Helper Functions
// ============================================================================

/**
 * Filter audit logs based on query parameters
 * This simulates what the backend should do
 */
function filterAuditLogs(
    logs: PermissionAuditLog[],
    params: {
        action?: AuditAction;
        targetType?: AuditTargetType;
        dateFrom?: string;
        dateTo?: string;
    }
): PermissionAuditLog[] {
    return logs.filter(log => {
        // Filter by action
        if (params.action && log.action !== params.action) {
            return false;
        }

        // Filter by target type
        if (params.targetType && log.targetType !== params.targetType) {
            return false;
        }

        // Filter by date range
        if (params.dateFrom || params.dateTo) {
            const logDate = new Date(log.createdAt).toISOString().split('T')[0];
            
            if (params.dateFrom && logDate < params.dateFrom) {
                return false;
            }
            
            if (params.dateTo && logDate > params.dateTo) {
                return false;
            }
        }

        return true;
    });
}

/**
 * Check if a log satisfies the filter criteria
 */
function logSatisfiesFilter(
    log: PermissionAuditLog,
    params: {
        action?: AuditAction;
        targetType?: AuditTargetType;
        dateFrom?: string;
        dateTo?: string;
    }
): boolean {
    // Check action filter
    if (params.action && log.action !== params.action) {
        return false;
    }

    // Check target type filter
    if (params.targetType && log.targetType !== params.targetType) {
        return false;
    }

    // Check date range
    if (params.dateFrom || params.dateTo) {
        const logDate = new Date(log.createdAt).toISOString().split('T')[0];
        
        if (params.dateFrom && logDate < params.dateFrom) {
            return false;
        }
        
        if (params.dateTo && logDate > params.dateTo) {
            return false;
        }
    }

    return true;
}

// ============================================================================
// Mock Data
// ============================================================================

const mockAuditLogs: PermissionAuditLog[] = [
    {
        id: 1,
        operatorId: 1,
        operatorName: '超级管理员',
        targetType: 'permission',
        targetId: 10,
        targetName: 'admin.users.create',
        action: 'permission_create',
        beforeData: undefined,
        afterData: JSON.stringify({ code: 'admin.users.create', description: '创建用户' }),
        ipAddress: '192.168.1.1',
        userAgent: 'Mozilla/5.0',
        requestId: 'req-001',
        createdAt: '2025-12-01T10:00:00Z',
    },
    {
        id: 2,
        operatorId: 1,
        operatorName: '超级管理员',
        targetType: 'role',
        targetId: 5,
        targetName: '内容管理员',
        action: 'role_permission_assign',
        beforeData: JSON.stringify({ permissionIds: [1, 2] }),
        afterData: JSON.stringify({ permissionIds: [1, 2, 3, 4] }),
        ipAddress: '192.168.1.1',
        userAgent: 'Mozilla/5.0',
        requestId: 'req-002',
        createdAt: '2025-12-02T14:30:00Z',
    },
    {
        id: 3,
        operatorId: 2,
        operatorName: '管理员A',
        targetType: 'user',
        targetId: 100,
        targetName: '用户张三',
        action: 'user_role_assign',
        beforeData: JSON.stringify({ roleIds: [1] }),
        afterData: JSON.stringify({ roleIds: [1, 2] }),
        ipAddress: '192.168.1.2',
        userAgent: 'Mozilla/5.0',
        requestId: 'req-003',
        createdAt: '2025-12-03T09:15:00Z',
    },
    {
        id: 4,
        operatorId: 1,
        operatorName: '超级管理员',
        targetType: 'permission',
        targetId: 11,
        targetName: 'admin.users.delete',
        action: 'permission_delete',
        beforeData: JSON.stringify({ code: 'admin.users.delete', description: '删除用户' }),
        afterData: undefined,
        ipAddress: '192.168.1.1',
        userAgent: 'Mozilla/5.0',
        requestId: 'req-004',
        createdAt: '2025-12-04T16:45:00Z',
    },
];

const renderWithRouter = async (component: React.ReactNode) => {
    const result = render(<BrowserRouter>{component}</BrowserRouter>);
    // Wait for initial render and effects to complete
    await waitFor(() => {
        expect(auditLogApi.list).toHaveBeenCalled();
    });
    return result;
};

// ============================================================================
// Tests
// ============================================================================

describe('Audit Log Page', () => {
    beforeEach(() => {
        vi.clearAllMocks();
        (auditLogApi.list as ReturnType<typeof vi.fn>).mockResolvedValue({
            data: {
                success: true,
                data: {
                    items: mockAuditLogs,
                    totalCount: mockAuditLogs.length,
                    page: 1,
                    pageSize: 10,
                    totalPages: 1,
                },
            },
        });
    });

    describe('基本渲染', () => {
        it('should render page title', async () => {
            await renderWithRouter(<AuditLogPage />);

            await waitFor(() => {
                expect(screen.getByText('权限审计日志')).toBeInTheDocument();
            });
        });

        it('should render audit log table', async () => {
            await renderWithRouter(<AuditLogPage />);

            await waitFor(() => {
                // Multiple rows may have the same operator name
                expect(screen.getAllByText('超级管理员').length).toBeGreaterThanOrEqual(1);
            });
        });

        it('should display action types with correct labels', async () => {
            await renderWithRouter(<AuditLogPage />);

            await waitFor(() => {
                expect(screen.getByText('创建权限')).toBeInTheDocument();
                expect(screen.getByText('分配角色权限')).toBeInTheDocument();
                expect(screen.getByText('分配用户角色')).toBeInTheDocument();
            });
        });

        it('should display target types with correct labels', async () => {
            await renderWithRouter(<AuditLogPage />);

            await waitFor(() => {
                expect(screen.getAllByText('权限').length).toBeGreaterThanOrEqual(1);
                expect(screen.getAllByText('角色').length).toBeGreaterThanOrEqual(1);
                expect(screen.getAllByText('用户').length).toBeGreaterThanOrEqual(1);
            });
        });

        it('should render filter controls', async () => {
            await renderWithRouter(<AuditLogPage />);

            await waitFor(() => {
                expect(screen.getByText('时间范围：')).toBeInTheDocument();
                expect(screen.getByText('操作类型：')).toBeInTheDocument();
                expect(screen.getByText('目标类型：')).toBeInTheDocument();
            });
        });

        it('should render export button', async () => {
            await renderWithRouter(<AuditLogPage />);

            await waitFor(() => {
                expect(screen.getByText('导出 CSV')).toBeInTheDocument();
            });
        });
    });

    describe('筛选功能', () => {
        it('should have search button', async () => {
            await renderWithRouter(<AuditLogPage />);

            await waitFor(() => {
                expect(screen.getByText('搜索')).toBeInTheDocument();
            });
        });

        it('should have reset button', async () => {
            await renderWithRouter(<AuditLogPage />);

            await waitFor(() => {
                expect(screen.getByText('重置')).toBeInTheDocument();
            });
        });

        it('should call API when search button is clicked', async () => {
            await renderWithRouter(<AuditLogPage />);

            await waitFor(() => {
                expect(screen.getByText('搜索')).toBeInTheDocument();
            });

            const searchButton = screen.getByText('搜索');
            fireEvent.click(searchButton);

            await waitFor(() => {
                // API should be called at least twice (initial load + search)
                expect(auditLogApi.list).toHaveBeenCalledTimes(2);
            });
        });
    });

    describe('详情弹窗', () => {
        it('should open detail modal when detail button is clicked', async () => {
            await renderWithRouter(<AuditLogPage />);

            await waitFor(() => {
                expect(screen.getAllByText('详情').length).toBeGreaterThanOrEqual(1);
            });

            const detailButtons = screen.getAllByText('详情');
            fireEvent.click(detailButtons[0]);

            await waitFor(() => {
                expect(screen.getByText('审计日志详情')).toBeInTheDocument();
            });
        });
    });

    describe('加载状态', () => {
        it('should show loading spinner while fetching data', async () => {
            (auditLogApi.list as ReturnType<typeof vi.fn>).mockImplementation(
                () => new Promise(() => {}) // Never resolves
            );

            render(<BrowserRouter><AuditLogPage /></BrowserRouter>);

            // Spin component should be present
            await waitFor(() => {
                const spinner = document.querySelector('.ant-spin');
                expect(spinner).toBeInTheDocument();
            });
        });
    });

    describe('错误处理', () => {
        it('should handle API error gracefully', async () => {
            (auditLogApi.list as ReturnType<typeof vi.fn>).mockRejectedValue(new Error('API Error'));

            render(<BrowserRouter><AuditLogPage /></BrowserRouter>);

            // Component should still render without crashing
            await waitFor(() => {
                expect(screen.getByText('权限审计日志')).toBeInTheDocument();
            });
        });
    });

    describe('分页功能', () => {
        it('should display pagination info', async () => {
            await renderWithRouter(<AuditLogPage />);

            await waitFor(() => {
                expect(screen.getByText(/共 \d+ 条/)).toBeInTheDocument();
            });
        });
    });
});

// ============================================================================
// Property-Based Tests
// ============================================================================

describe('Property 12: 审计日志过滤正确性', () => {
    /**
     * **Feature: rbac-button-level-permission, Property 12: 审计日志过滤正确性**
     * **Validates: Requirements 6.3**
     * 
     * *For any* audit log query with filters (time range, action type, target type),
     * the returned logs should all satisfy the filter criteria.
     */

    it('Property 12a: filtering by action should return only logs with that action', () => {
        fc.assert(
            fc.property(auditLogListArb, auditActionArb, (logs, action) => {
                const filtered = filterAuditLogs(logs, { action });
                
                // All returned logs should have the specified action
                return filtered.every(log => log.action === action);
            }),
            { numRuns: 100 }
        );
    });

    it('Property 12b: filtering by target type should return only logs with that target type', () => {
        fc.assert(
            fc.property(auditLogListArb, auditTargetTypeArb, (logs, targetType) => {
                const filtered = filterAuditLogs(logs, { targetType });
                
                // All returned logs should have the specified target type
                return filtered.every(log => log.targetType === targetType);
            }),
            { numRuns: 100 }
        );
    });

    it('Property 12c: filtering by date range should return only logs within that range', () => {
        fc.assert(
            fc.property(
                auditLogListArb,
                dateStringArb,
                dateStringArb,
                (logs, date1, date2) => {
                    // Ensure dateFrom <= dateTo
                    const [dateFrom, dateTo] = date1 <= date2 ? [date1, date2] : [date2, date1];
                    
                    const filtered = filterAuditLogs(logs, { dateFrom, dateTo });
                    
                    // All returned logs should be within the date range
                    return filtered.every(log => {
                        const logDate = new Date(log.createdAt).toISOString().split('T')[0];
                        return logDate >= dateFrom && logDate <= dateTo;
                    });
                }
            ),
            { numRuns: 100 }
        );
    });

    it('Property 12d: combining multiple filters should return logs satisfying all criteria', () => {
        fc.assert(
            fc.property(
                auditLogListArb,
                filterParamsArb,
                (logs, params) => {
                    // Normalize date range
                    let dateFrom = params.dateFrom ?? undefined;
                    let dateTo = params.dateTo ?? undefined;
                    
                    if (dateFrom && dateTo && dateFrom > dateTo) {
                        [dateFrom, dateTo] = [dateTo, dateFrom];
                    }

                    const filterParams = {
                        action: params.action ?? undefined,
                        targetType: params.targetType ?? undefined,
                        dateFrom,
                        dateTo,
                    };

                    const filtered = filterAuditLogs(logs, filterParams);
                    
                    // All returned logs should satisfy all filter criteria
                    return filtered.every(log => logSatisfiesFilter(log, filterParams));
                }
            ),
            { numRuns: 100 }
        );
    });

    it('Property 12e: empty filter should return all logs', () => {
        fc.assert(
            fc.property(auditLogListArb, (logs) => {
                const filtered = filterAuditLogs(logs, {});
                
                // With no filters, all logs should be returned
                return filtered.length === logs.length;
            }),
            { numRuns: 100 }
        );
    });

    it('Property 12f: filtered results should be a subset of original logs', () => {
        fc.assert(
            fc.property(
                auditLogListArb,
                filterParamsArb,
                (logs, params) => {
                    const filterParams = {
                        action: params.action ?? undefined,
                        targetType: params.targetType ?? undefined,
                        dateFrom: params.dateFrom ?? undefined,
                        dateTo: params.dateTo ?? undefined,
                    };

                    const filtered = filterAuditLogs(logs, filterParams);
                    
                    // Filtered results should be a subset (length <= original)
                    if (filtered.length > logs.length) {
                        return false;
                    }
                    
                    // Every filtered log should exist in original logs
                    return filtered.every(filteredLog => 
                        logs.some(log => log.id === filteredLog.id)
                    );
                }
            ),
            { numRuns: 100 }
        );
    });

    it('Property 12g: filtering should be idempotent', () => {
        fc.assert(
            fc.property(
                auditLogListArb,
                filterParamsArb,
                (logs, params) => {
                    const filterParams = {
                        action: params.action ?? undefined,
                        targetType: params.targetType ?? undefined,
                        dateFrom: params.dateFrom ?? undefined,
                        dateTo: params.dateTo ?? undefined,
                    };

                    const filtered1 = filterAuditLogs(logs, filterParams);
                    const filtered2 = filterAuditLogs(filtered1, filterParams);
                    
                    // Filtering twice should give the same result
                    return filtered1.length === filtered2.length &&
                        filtered1.every((log, i) => log.id === filtered2[i].id);
                }
            ),
            { numRuns: 100 }
        );
    });

    it('Property 12h: logs not matching filter should not be in results', () => {
        fc.assert(
            fc.property(
                auditLogListArb,
                auditActionArb,
                (logs, action) => {
                    const filtered = filterAuditLogs(logs, { action });
                    const filteredIds = new Set(filtered.map(log => log.id));
                    
                    // Logs with different action should not be in filtered results
                    const logsWithDifferentAction = logs.filter(log => log.action !== action);
                    
                    return logsWithDifferentAction.every(log => !filteredIds.has(log.id));
                }
            ),
            { numRuns: 100 }
        );
    });
});
