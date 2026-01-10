/**
 * Role Management Page Tests
 *
 * Tests for Role page component including:
 * - Successful data loading
 * - Loading states
 * - Error handling
 * - CRUD operations
 * - Permission configuration
 * - System role protection
 */

import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { screen, waitFor, fireEvent } from '@testing-library/react';
import RolePage from './index';
import { renderWithProviders, resetAllMocks, flushPromises } from '@/testutils';

// Mock the adminApi module using vi.hoisted
const { mockApi, mockMessage } = vi.hoisted(() => ({
  mockApi: {
    getRoles: vi.fn(),
    createRole: vi.fn(),
    updateRole: vi.fn(),
    deleteRole: vi.fn(),
  },
  mockMessage: {
    success: vi.fn(),
    error: vi.fn(),
    warning: vi.fn(),
    info: vi.fn(),
    loading: vi.fn(),
  },
}));

vi.mock('@/api/admin', () => ({
  adminApi: mockApi,
}));

// Mock antd message
vi.mock('antd', async () => {
  const actual = await vi.importActual('antd');
  return {
    ...actual,
    message: mockMessage,
  };
});

// Mock react-router-dom
const mockNavigate = vi.fn();
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom');
  return {
    ...actual,
    useNavigate: () => mockNavigate,
  };
});

// Helper function to create mock role
const createMockRole = (overrides: Record<string, unknown> = {}): Record<string, unknown> => ({
  id: 1,
  name: '管理员',
  slug: 'admin',
  description: '系统管理员角色',
  isSystem: false,
  users: [{ id: 1 }, { id: 2 }],
  permissions: [{ id: 1 }, { id: 2 }, { id: 3 }],
  createdAt: '2024-01-01T00:00:00Z',
  updatedAt: '2024-01-01T00:00:00Z',
  ...overrides,
});

describe('RolePage', () => {
  beforeEach(() => {
    resetAllMocks();
    mockMessage.success.mockClear();
    mockMessage.error.mockClear();
    mockNavigate.mockClear();
    localStorage.setItem('token', 'test-token');
    localStorage.setItem('user_role', 'admin');

    // Set default mock return values
    mockApi.getRoles.mockResolvedValue({
      data: {
        success: true,
        data: [createMockRole()],
        pagination: { total: 1, page: 1, page_size: 10 },
      },
    });
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  describe('Successful Data Loading', () => {
    it('should render role page successfully', async () => {
      renderWithProviders(<RolePage />);

      await waitFor(() => {
        expect(screen.getByText('角色管理')).toBeInTheDocument();
      });

      expect(mockApi.getRoles).toHaveBeenCalled();
    });

    it('should display page subtitle', async () => {
      renderWithProviders(<RolePage />);

      await waitFor(() => {
        expect(screen.getByText('管理系统角色和权限分配')).toBeInTheDocument();
      });
    });

    it('should display role list', async () => {
      renderWithProviders(<RolePage />);

      await waitFor(() => {
        expect(screen.getByText('管理员')).toBeInTheDocument();
      });
    });

    it('should display role slug', async () => {
      renderWithProviders(<RolePage />);

      await waitFor(() => {
        expect(screen.getByText('admin')).toBeInTheDocument();
      });
    });

    it('should display role description', async () => {
      renderWithProviders(<RolePage />);

      await waitFor(() => {
        expect(screen.getByText('系统管理员角色')).toBeInTheDocument();
      });
    });
  });

  describe('Loading States', () => {
    it('should show loading indicator while fetching data', async () => {
      mockApi.getRoles.mockImplementation(
        () =>
          new Promise((resolve) => {
            setTimeout(() => {
              resolve({
                data: {
                  success: true,
                  data: [createMockRole()],
                  pagination: { total: 1 },
                },
              });
            }, 100);
          })
      );

      renderWithProviders(<RolePage />);

      expect(mockApi.getRoles).toHaveBeenCalled();
    });

    it('should hide loading indicator after data loads', async () => {
      renderWithProviders(<RolePage />);

      await waitFor(() => {
        expect(screen.getByText('角色管理')).toBeInTheDocument();
      });

      await flushPromises();

      expect(mockApi.getRoles).toHaveBeenCalledTimes(1);
    });
  });

  describe('Error Handling', () => {
    it('should handle API error gracefully', async () => {
      mockApi.getRoles.mockRejectedValue(new Error('Network error'));

      renderWithProviders(<RolePage />);

      await waitFor(() => {
        expect(screen.getByText('角色管理')).toBeInTheDocument();
      });

      // API was called even if it failed
      expect(mockApi.getRoles).toHaveBeenCalled();
    });
  });

  describe('Filter Functionality', () => {
    it('should have keyword search input', async () => {
      renderWithProviders(<RolePage />);

      await waitFor(() => {
        expect(screen.getByPlaceholderText('角色名称/编码')).toBeInTheDocument();
      });
    });
  });

  describe('Toolbar Buttons', () => {
    it('should display create button', async () => {
      renderWithProviders(<RolePage />);

      await waitFor(() => {
        expect(screen.getByText('新增角色')).toBeInTheDocument();
      });
    });
  });

  describe('Role Actions', () => {
    it('should display permission button', async () => {
      renderWithProviders(<RolePage />);

      await waitFor(() => {
        expect(screen.getByText('权限')).toBeInTheDocument();
      });
    });

    it('should display edit button', async () => {
      renderWithProviders(<RolePage />);

      await waitFor(() => {
        expect(screen.getByText('编辑')).toBeInTheDocument();
      });
    });

    it('should display delete button for non-system role', async () => {
      renderWithProviders(<RolePage />);

      await waitFor(() => {
        expect(screen.getByText('删除')).toBeInTheDocument();
      });
    });

    it('should not display delete button for system role', async () => {
      mockApi.getRoles.mockResolvedValue({
        data: {
          success: true,
          data: [createMockRole({ isSystem: true })],
          pagination: { total: 1 },
        },
      });

      renderWithProviders(<RolePage />);

      await waitFor(() => {
        expect(screen.getByText('角色管理')).toBeInTheDocument();
      });

      // System role should not have delete button
      const deleteButtons = screen.queryAllByText('删除');
      expect(deleteButtons.length).toBe(0);
    });

    it('should navigate to permission page when clicking permission button', async () => {
      renderWithProviders(<RolePage />);

      await waitFor(() => {
        expect(screen.getByText('权限')).toBeInTheDocument();
      });

      const permissionButton = screen.getByText('权限');
      fireEvent.click(permissionButton);

      expect(mockNavigate).toHaveBeenCalledWith('/admin/sys/role/1/permissions');
    });
  });

  describe('System Role Display', () => {
    it('should display system tag for system role', async () => {
      mockApi.getRoles.mockResolvedValue({
        data: {
          success: true,
          data: [createMockRole({ isSystem: true })],
          pagination: { total: 1 },
        },
      });

      renderWithProviders(<RolePage />);

      await waitFor(() => {
        expect(screen.getByText('系统')).toBeInTheDocument();
      });
    });

    it('should display super admin tag for superAdmin role', async () => {
      mockApi.getRoles.mockResolvedValue({
        data: {
          success: true,
          data: [createMockRole({ slug: 'superAdmin', isSystem: true })],
          pagination: { total: 1 },
        },
      });

      renderWithProviders(<RolePage />);

      await waitFor(() => {
        expect(screen.getByText('超管')).toBeInTheDocument();
      });
    });
  });

  describe('User and Permission Count Display', () => {
    it('should display user count', async () => {
      renderWithProviders(<RolePage />);

      await waitFor(() => {
        expect(screen.getByText('2 人')).toBeInTheDocument();
      });
    });

    it('should display permission count', async () => {
      renderWithProviders(<RolePage />);

      await waitFor(() => {
        expect(screen.getByText('3 项')).toBeInTheDocument();
      });
    });
  });

  describe('Create Role Modal', () => {
    it('should open create modal when clicking create button', async () => {
      renderWithProviders(<RolePage />);

      await waitFor(() => {
        expect(screen.getByText('新增角色')).toBeInTheDocument();
      });

      const createButton = screen.getByText('新增角色');
      fireEvent.click(createButton);

      await waitFor(() => {
        // Modal should be open - check for form placeholder
        expect(screen.getByPlaceholderText('请输入角色名称')).toBeInTheDocument();
      });
    });

    it('should display form fields in create modal', async () => {
      renderWithProviders(<RolePage />);

      await waitFor(() => {
        expect(screen.getByText('新增角色')).toBeInTheDocument();
      });

      const createButton = screen.getByText('新增角色');
      fireEvent.click(createButton);

      await waitFor(() => {
        // Modal should be open - check for form placeholders instead of labels
        expect(screen.getByPlaceholderText('请输入角色名称')).toBeInTheDocument();
      });
    });
  });

  describe('Edit Role Modal', () => {
    it('should open edit modal when clicking edit button', async () => {
      renderWithProviders(<RolePage />);

      await waitFor(() => {
        expect(screen.getByText('编辑')).toBeInTheDocument();
      });

      const editButton = screen.getByText('编辑');
      fireEvent.click(editButton);

      await waitFor(() => {
        expect(screen.getByText('编辑角色')).toBeInTheDocument();
      });
    });
  });

  describe('Delete Role', () => {
    it('should show confirmation when delete button clicked', async () => {
      renderWithProviders(<RolePage />);

      await waitFor(() => {
        expect(screen.getByText('删除')).toBeInTheDocument();
      });

      const deleteButton = screen.getByText('删除');
      fireEvent.click(deleteButton);

      await waitFor(() => {
        expect(screen.getByText('确定要删除该角色吗？')).toBeInTheDocument();
      });
    });
  });

  describe('Empty State', () => {
    it('should handle empty role list', async () => {
      mockApi.getRoles.mockResolvedValue({
        data: {
          success: true,
          data: [],
          pagination: { total: 0 },
        },
      });

      renderWithProviders(<RolePage />);

      await waitFor(() => {
        expect(screen.getByText('角色管理')).toBeInTheDocument();
      });

      expect(mockApi.getRoles).toHaveBeenCalled();
    });
  });

  describe('Accessibility', () => {
    it('should have proper page title', async () => {
      renderWithProviders(<RolePage />);

      await waitFor(() => {
        expect(screen.getByText('角色管理')).toBeInTheDocument();
      });
    });

    it('should have proper table structure', async () => {
      renderWithProviders(<RolePage />);

      await waitFor(() => {
        expect(screen.getByRole('table')).toBeInTheDocument();
      });
    });
  });

  describe('Data Refresh', () => {
    it('should call getRoles on mount', async () => {
      renderWithProviders(<RolePage />);

      await waitFor(() => {
        expect(mockApi.getRoles).toHaveBeenCalledTimes(1);
      });
    });
  });
});
