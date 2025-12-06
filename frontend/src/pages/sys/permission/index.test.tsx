/**
 * Permission Management Page Tests
 * 测试权限管理页面的渲染和功能
 * Requirements: 1.1 - 权限定义与管理
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import PermissionPage from './index';
import { permissionApi } from '@/api/permission';

// Mock the permission API
vi.mock('@/api/permission', () => ({
    permissionApi: {
        list: vi.fn(),
        get: vi.fn(),
        getGroups: vi.fn(),
        create: vi.fn(),
        update: vi.fn(),
        delete: vi.fn(),
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

const mockPermissions = [
    {
        id: 1,
        code: 'admin.users.list',
        description: '用户列表',
        group: '用户管理',
        method: 'GET',
        path: '/api/admin/users',
        sortOrder: 0,
        isSystem: true,
        createdAt: '2025-12-01T00:00:00Z',
        updatedAt: '2025-12-01T00:00:00Z',
    },
    {
        id: 2,
        code: 'admin.users.create',
        description: '创建用户',
        group: '用户管理',
        method: 'POST',
        path: '/api/admin/users',
        sortOrder: 1,
        isSystem: false,
        createdAt: '2025-12-01T00:00:00Z',
        updatedAt: '2025-12-01T00:00:00Z',
    },
    {
        id: 3,
        code: 'admin.orders.list',
        description: '订单列表',
        group: '订单管理',
        method: 'GET',
        path: '/api/admin/orders',
        sortOrder: 0,
        isSystem: true,
        createdAt: '2025-12-01T00:00:00Z',
        updatedAt: '2025-12-01T00:00:00Z',
    },
];

const mockGroups = ['用户管理', '订单管理', '系统管理'];

const renderWithRouter = async (component: React.ReactNode) => {
    const result = render(<BrowserRouter>{component}</BrowserRouter>);
    // Wait for initial render and effects to complete
    await waitFor(() => {
        expect(permissionApi.list).toHaveBeenCalled();
    });
    return result;
};

describe('Permission Page', () => {
    beforeEach(() => {
        vi.clearAllMocks();
        (permissionApi.list as ReturnType<typeof vi.fn>).mockResolvedValue({
            data: {
                success: true,
                data: {
                    items: mockPermissions,
                    totalCount: mockPermissions.length,
                    page: 1,
                    pageSize: 10,
                    totalPages: 1,
                },
            },
        });
        (permissionApi.getGroups as ReturnType<typeof vi.fn>).mockResolvedValue({
            data: {
                success: true,
                data: mockGroups,
            },
        });
    });

    describe('基本渲染', () => {
        it('should render page title', async () => {
            await renderWithRouter(<PermissionPage />);

            await waitFor(() => {
                expect(screen.getByText('权限管理')).toBeInTheDocument();
            });
        });

        it('should render permission list table', async () => {
            await renderWithRouter(<PermissionPage />);

            await waitFor(() => {
                expect(screen.getByText('admin.users.list')).toBeInTheDocument();
                expect(screen.getByText('admin.users.create')).toBeInTheDocument();
                expect(screen.getByText('admin.orders.list')).toBeInTheDocument();
            });
        });

        it('should display permission descriptions', async () => {
            await renderWithRouter(<PermissionPage />);

            await waitFor(() => {
                expect(screen.getByText('用户列表')).toBeInTheDocument();
                expect(screen.getByText('创建用户')).toBeInTheDocument();
                expect(screen.getByText('订单列表')).toBeInTheDocument();
            });
        });

        it('should display permission groups', async () => {
            await renderWithRouter(<PermissionPage />);

            await waitFor(() => {
                // Groups appear multiple times in the table
                expect(screen.getAllByText('用户管理').length).toBeGreaterThanOrEqual(1);
                expect(screen.getAllByText('订单管理').length).toBeGreaterThanOrEqual(1);
            });
        });

        it('should display HTTP methods with correct colors', async () => {
            await renderWithRouter(<PermissionPage />);

            await waitFor(() => {
                // GET and POST methods should be displayed
                expect(screen.getAllByText('GET').length).toBeGreaterThanOrEqual(1);
                expect(screen.getAllByText('POST').length).toBeGreaterThanOrEqual(1);
            });
        });

        it('should show system badge for system permissions', async () => {
            await renderWithRouter(<PermissionPage />);

            await waitFor(() => {
                // System permissions should have a badge
                expect(screen.getAllByText('系统').length).toBeGreaterThanOrEqual(1);
            });
        });

        it('should render create button', async () => {
            await renderWithRouter(<PermissionPage />);

            await waitFor(() => {
                expect(screen.getByText('新增权限')).toBeInTheDocument();
            });
        });
    });

    describe('搜索功能', () => {
        it('should have search input field', async () => {
            await renderWithRouter(<PermissionPage />);

            await waitFor(() => {
                expect(screen.getByPlaceholderText('权限码/描述')).toBeInTheDocument();
            });
        });

        it('should have group filter label', async () => {
            await renderWithRouter(<PermissionPage />);

            await waitFor(() => {
                // Select component uses label - multiple elements exist (form label and table header)
                expect(screen.getAllByText('分组').length).toBeGreaterThanOrEqual(1);
            });
        });

        it('should call API with search params when search button is clicked', async () => {
            await renderWithRouter(<PermissionPage />);

            await waitFor(() => {
                expect(screen.getByPlaceholderText('权限码/描述')).toBeInTheDocument();
            });

            const searchInput = screen.getByPlaceholderText('权限码/描述');
            fireEvent.change(searchInput, { target: { value: 'users' } });

            const searchButton = screen.getByText('搜索');
            fireEvent.click(searchButton);

            await waitFor(() => {
                expect(permissionApi.list).toHaveBeenCalledWith(
                    expect.objectContaining({ keyword: 'users' })
                );
            });
        });
    });

    describe('创建权限', () => {
        it('should open create modal when create button is clicked', async () => {
            await renderWithRouter(<PermissionPage />);

            await waitFor(() => {
                expect(screen.getByText('新增权限')).toBeInTheDocument();
            });

            const createButton = screen.getByText('新增权限');
            fireEvent.click(createButton);

            await waitFor(() => {
                // Modal title and button both have "新增权限", use getAllByText
                expect(screen.getAllByText('新增权限').length).toBeGreaterThanOrEqual(2);
                expect(screen.getByLabelText('权限码')).toBeInTheDocument();
            });
        });

        it('should have permission code input in create modal', async () => {
            await renderWithRouter(<PermissionPage />);

            const createButton = screen.getByText('新增权限');
            fireEvent.click(createButton);

            await waitFor(() => {
                expect(screen.getByLabelText('权限码')).toBeInTheDocument();
            });

            // Verify the input is enabled (not disabled) in create mode
            const codeInput = screen.getByLabelText('权限码');
            expect(codeInput).not.toBeDisabled();
        });
    });

    describe('编辑权限', () => {
        it('should open edit modal when edit button is clicked', async () => {
            await renderWithRouter(<PermissionPage />);

            await waitFor(() => {
                expect(screen.getAllByText('编辑').length).toBeGreaterThanOrEqual(1);
            });

            const editButtons = screen.getAllByText('编辑');
            fireEvent.click(editButtons[0]);

            await waitFor(() => {
                expect(screen.getByText('编辑权限')).toBeInTheDocument();
            });
        });

        it('should disable permission code field in edit mode', async () => {
            await renderWithRouter(<PermissionPage />);

            await waitFor(() => {
                expect(screen.getAllByText('编辑').length).toBeGreaterThanOrEqual(1);
            });

            const editButtons = screen.getAllByText('编辑');
            fireEvent.click(editButtons[0]);

            await waitFor(() => {
                const codeInput = screen.getByLabelText('权限码');
                expect(codeInput).toBeDisabled();
            });
        });
    });

    describe('删除权限', () => {
        it('should not show delete button for system permissions', async () => {
            await renderWithRouter(<PermissionPage />);

            await waitFor(() => {
                // System permission (id=1) should not have delete button
                // Non-system permission (id=2) should have delete button
                const deleteButtons = screen.getAllByText('删除');
                // Only non-system permissions should have delete buttons
                expect(deleteButtons.length).toBe(1);
            });
        });

        it('should show delete confirmation modal when delete button is clicked', async () => {
            await renderWithRouter(<PermissionPage />);

            await waitFor(() => {
                expect(screen.getAllByText('删除').length).toBeGreaterThanOrEqual(1);
            });

            const deleteButtons = screen.getAllByText('删除');
            fireEvent.click(deleteButtons[0]);

            await waitFor(() => {
                expect(screen.getByText('确认删除权限')).toBeInTheDocument();
            });
        });

        it('should call delete API when confirmed', async () => {
            (permissionApi.delete as ReturnType<typeof vi.fn>).mockResolvedValue({
                data: { success: true },
            });

            await renderWithRouter(<PermissionPage />);

            await waitFor(() => {
                expect(screen.getAllByText('删除').length).toBeGreaterThanOrEqual(1);
            });

            const deleteButtons = screen.getAllByText('删除');
            fireEvent.click(deleteButtons[0]);

            await waitFor(() => {
                expect(screen.getByText('确认删除')).toBeInTheDocument();
            });

            const confirmButton = screen.getByText('确认删除');
            fireEvent.click(confirmButton);

            await waitFor(() => {
                expect(permissionApi.delete).toHaveBeenCalled();
            });
        });
    });

    describe('加载状态', () => {
        it('should show loading spinner while fetching data', async () => {
            (permissionApi.list as ReturnType<typeof vi.fn>).mockImplementation(
                () => new Promise(() => {}) // Never resolves
            );

            render(<BrowserRouter><PermissionPage /></BrowserRouter>);

            // Spin component should be present
            await waitFor(() => {
                const spinner = document.querySelector('.ant-spin');
                expect(spinner).toBeInTheDocument();
            });
        });
    });

    describe('错误处理', () => {
        it('should handle API error gracefully', async () => {
            (permissionApi.list as ReturnType<typeof vi.fn>).mockRejectedValue(new Error('API Error'));

            render(<BrowserRouter><PermissionPage /></BrowserRouter>);

            // Component should still render without crashing
            await waitFor(() => {
                expect(screen.getByText('权限管理')).toBeInTheDocument();
            });
        });
    });

    describe('分页功能', () => {
        it('should display pagination info', async () => {
            await renderWithRouter(<PermissionPage />);

            await waitFor(() => {
                expect(screen.getByText(/共 \d+ 条/)).toBeInTheDocument();
            });
        });
    });
});
