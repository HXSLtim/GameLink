/**
 * Permission API Tests
 * 
 * Tests for RBAC permission management API
 * Requirements: 1.1, 2.1, 2.2, 2.3, 6.3, 6.5, 9.1, 9.2
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { permissionApi, roleApi, userRoleApi, auditLogApi } from './permission';
import apiClient from './client';

// Mock the API client
vi.mock('./client', () => ({
    default: {
        get: vi.fn(),
        post: vi.fn(),
        put: vi.fn(),
        patch: vi.fn(),
        delete: vi.fn(),
    },
}));

describe('Permission API', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    describe('permissionApi', () => {
        it('should list permissions with params', async () => {
            const mockResponse = { data: { items: [], totalCount: 0 } };
            (apiClient.get as ReturnType<typeof vi.fn>).mockResolvedValue(mockResponse);

            const params = { page: 1, pageSize: 10, keyword: 'admin' };
            await permissionApi.list(params);

            expect(apiClient.get).toHaveBeenCalledWith('/admin/permissions', { params });
        });

        it('should list permissions without params', async () => {
            const mockResponse = { data: { items: [], totalCount: 0 } };
            (apiClient.get as ReturnType<typeof vi.fn>).mockResolvedValue(mockResponse);

            await permissionApi.list();

            expect(apiClient.get).toHaveBeenCalledWith('/admin/permissions', { params: undefined });
        });

        it('should get permission by id', async () => {
            const mockResponse = { data: { id: 1, code: 'admin.users.list' } };
            (apiClient.get as ReturnType<typeof vi.fn>).mockResolvedValue(mockResponse);

            await permissionApi.get(1);

            expect(apiClient.get).toHaveBeenCalledWith('/admin/permissions/1');
        });

        it('should get permission tree', async () => {
            const mockResponse = { data: [] };
            (apiClient.get as ReturnType<typeof vi.fn>).mockResolvedValue(mockResponse);

            await permissionApi.getTree();

            expect(apiClient.get).toHaveBeenCalledWith('/admin/permissions/tree');
        });

        it('should get permission groups', async () => {
            const mockResponse = { data: ['用户管理', '订单管理'] };
            (apiClient.get as ReturnType<typeof vi.fn>).mockResolvedValue(mockResponse);

            await permissionApi.getGroups();

            expect(apiClient.get).toHaveBeenCalledWith('/admin/permissions/groups');
        });

        it('should create permission', async () => {
            const mockResponse = { data: { id: 1, code: 'admin.test.create' } };
            (apiClient.post as ReturnType<typeof vi.fn>).mockResolvedValue(mockResponse);

            const data = {
                code: 'admin.test.create',
                description: 'Test permission',
                group: '测试',
                method: 'POST' as const,
                path: '/api/admin/test',
            };
            await permissionApi.create(data);

            expect(apiClient.post).toHaveBeenCalledWith('/admin/permissions', data);
        });

        it('should update permission', async () => {
            const mockResponse = { data: { id: 1, code: 'admin.test.update' } };
            (apiClient.put as ReturnType<typeof vi.fn>).mockResolvedValue(mockResponse);

            const data = { description: 'Updated description' };
            await permissionApi.update(1, data);

            expect(apiClient.put).toHaveBeenCalledWith('/admin/permissions/1', data);
        });

        it('should patch permission', async () => {
            const mockResponse = { data: { id: 1 } };
            (apiClient.patch as ReturnType<typeof vi.fn>).mockResolvedValue(mockResponse);

            const data = { description: 'Patched' };
            await permissionApi.patch(1, data);

            expect(apiClient.patch).toHaveBeenCalledWith('/admin/permissions/1', data);
        });

        it('should delete permission', async () => {
            const mockResponse = { data: null };
            (apiClient.delete as ReturnType<typeof vi.fn>).mockResolvedValue(mockResponse);

            await permissionApi.delete(1);

            expect(apiClient.delete).toHaveBeenCalledWith('/admin/permissions/1');
        });

        it('should get my permissions', async () => {
            const mockResponse = { data: ['admin.users.list', 'admin.users.create'] };
            (apiClient.get as ReturnType<typeof vi.fn>).mockResolvedValue(mockResponse);

            await permissionApi.getMyPermissions();

            expect(apiClient.get).toHaveBeenCalledWith('/admin/me/permissions');
        });

        it('should get my menus', async () => {
            const mockResponse = { data: [] };
            (apiClient.get as ReturnType<typeof vi.fn>).mockResolvedValue(mockResponse);

            await permissionApi.getMyMenus();

            expect(apiClient.get).toHaveBeenCalledWith('/admin/me/menus');
        });
    });

    describe('roleApi', () => {
        it('should list roles', async () => {
            const mockResponse = { data: { items: [], totalCount: 0 } };
            (apiClient.get as ReturnType<typeof vi.fn>).mockResolvedValue(mockResponse);

            const params = { page: 1, pageSize: 10 };
            await roleApi.list(params);

            expect(apiClient.get).toHaveBeenCalledWith('/admin/roles', { params });
        });

        it('should get role by id', async () => {
            const mockResponse = { data: { id: 1, name: 'Admin' } };
            (apiClient.get as ReturnType<typeof vi.fn>).mockResolvedValue(mockResponse);

            await roleApi.get(1);

            expect(apiClient.get).toHaveBeenCalledWith('/admin/roles/1');
        });

        it('should get role with permissions', async () => {
            const mockResponse = { data: { id: 1, name: 'Admin', permissions: [] } };
            (apiClient.get as ReturnType<typeof vi.fn>).mockResolvedValue(mockResponse);

            await roleApi.getWithPermissions(1);

            expect(apiClient.get).toHaveBeenCalledWith('/admin/roles/1', {
                params: { with_permissions: true },
            });
        });

        it('should create role', async () => {
            const mockResponse = { data: { id: 1, name: 'NewRole' } };
            (apiClient.post as ReturnType<typeof vi.fn>).mockResolvedValue(mockResponse);

            const data = { name: 'NewRole', description: 'New role description' };
            await roleApi.create(data);

            expect(apiClient.post).toHaveBeenCalledWith('/admin/roles', data);
        });

        it('should update role', async () => {
            const mockResponse = { data: { id: 1, name: 'UpdatedRole' } };
            (apiClient.put as ReturnType<typeof vi.fn>).mockResolvedValue(mockResponse);

            const data = { name: 'UpdatedRole' };
            await roleApi.update(1, data);

            expect(apiClient.put).toHaveBeenCalledWith('/admin/roles/1', data);
        });

        it('should delete role', async () => {
            const mockResponse = { data: null };
            (apiClient.delete as ReturnType<typeof vi.fn>).mockResolvedValue(mockResponse);

            await roleApi.delete(1);

            expect(apiClient.delete).toHaveBeenCalledWith('/admin/roles/1');
        });

        it('should get role permissions', async () => {
            const mockResponse = { data: [1, 2, 3] };
            (apiClient.get as ReturnType<typeof vi.fn>).mockResolvedValue(mockResponse);

            await roleApi.getPermissions(1);

            expect(apiClient.get).toHaveBeenCalledWith('/admin/roles/1/permissions');
        });

        it('should batch assign permissions', async () => {
            const mockResponse = { data: null };
            (apiClient.put as ReturnType<typeof vi.fn>).mockResolvedValue(mockResponse);

            const data = { permissionIds: [1, 2, 3] };
            await roleApi.batchAssignPermissions(1, data);

            expect(apiClient.put).toHaveBeenCalledWith('/admin/roles/1/permissions/batch', data);
        });

        it('should add single permission', async () => {
            const mockResponse = { data: null };
            (apiClient.post as ReturnType<typeof vi.fn>).mockResolvedValue(mockResponse);

            await roleApi.addPermission(1, 5);

            expect(apiClient.post).toHaveBeenCalledWith('/admin/roles/1/permissions/5');
        });

        it('should remove single permission', async () => {
            const mockResponse = { data: null };
            (apiClient.delete as ReturnType<typeof vi.fn>).mockResolvedValue(mockResponse);

            await roleApi.removePermission(1, 5);

            expect(apiClient.delete).toHaveBeenCalledWith('/admin/roles/1/permissions/5');
        });

        it('should set parent role', async () => {
            const mockResponse = { data: { id: 1, parentId: 2 } };
            (apiClient.put as ReturnType<typeof vi.fn>).mockResolvedValue(mockResponse);

            await roleApi.setParent(1, 2);

            expect(apiClient.put).toHaveBeenCalledWith('/admin/roles/1/parent', { parentId: 2 });
        });

        it('should set parent role to null', async () => {
            const mockResponse = { data: { id: 1, parentId: null } };
            (apiClient.put as ReturnType<typeof vi.fn>).mockResolvedValue(mockResponse);

            await roleApi.setParent(1, null);

            expect(apiClient.put).toHaveBeenCalledWith('/admin/roles/1/parent', { parentId: null });
        });

        it('should get inheritance chain', async () => {
            const mockResponse = { data: [{ id: 1 }, { id: 2 }] };
            (apiClient.get as ReturnType<typeof vi.fn>).mockResolvedValue(mockResponse);

            await roleApi.getInheritanceChain(1);

            expect(apiClient.get).toHaveBeenCalledWith('/admin/roles/1/inheritance-chain');
        });
    });

    describe('userRoleApi', () => {
        it('should get user roles', async () => {
            const mockResponse = { data: [{ id: 1, name: 'Admin' }] };
            (apiClient.get as ReturnType<typeof vi.fn>).mockResolvedValue(mockResponse);

            await userRoleApi.getUserRoles(1);

            expect(apiClient.get).toHaveBeenCalledWith('/admin/users/1/roles');
        });

        it('should assign roles to user', async () => {
            const mockResponse = { data: null };
            (apiClient.put as ReturnType<typeof vi.fn>).mockResolvedValue(mockResponse);

            const data = { roleIds: [1, 2] };
            await userRoleApi.assignRoles(1, data);

            expect(apiClient.put).toHaveBeenCalledWith('/admin/users/1/roles', data);
        });

        it('should get user permissions', async () => {
            const mockResponse = { data: { permissions: [], roles: [] } };
            (apiClient.get as ReturnType<typeof vi.fn>).mockResolvedValue(mockResponse);

            await userRoleApi.getUserPermissions(1);

            expect(apiClient.get).toHaveBeenCalledWith('/admin/users/1/permissions');
        });

        it('should batch assign roles', async () => {
            const mockResponse = { data: { success: 2, failed: 0 } };
            (apiClient.post as ReturnType<typeof vi.fn>).mockResolvedValue(mockResponse);

            const data = { userIds: [1, 2], roleIds: [1] };
            await userRoleApi.batchAssignRoles(data);

            expect(apiClient.post).toHaveBeenCalledWith('/admin/users/batch/roles', data);
        });

        it('should check super admin', async () => {
            const mockResponse = { data: true };
            (apiClient.get as ReturnType<typeof vi.fn>).mockResolvedValue(mockResponse);

            await userRoleApi.checkSuperAdmin(1);

            expect(apiClient.get).toHaveBeenCalledWith('/admin/users/1/is-super-admin');
        });
    });

    describe('auditLogApi', () => {
        it('should list audit logs', async () => {
            const mockResponse = { data: { items: [], totalCount: 0 } };
            (apiClient.get as ReturnType<typeof vi.fn>).mockResolvedValue(mockResponse);

            const params = { page: 1, pageSize: 10 };
            await auditLogApi.list(params);

            expect(apiClient.get).toHaveBeenCalledWith('/admin/audit/permissions', { params });
        });

        it('should get audit log by id', async () => {
            const mockResponse = { data: { id: 1, action: 'create' } };
            (apiClient.get as ReturnType<typeof vi.fn>).mockResolvedValue(mockResponse);

            await auditLogApi.get(1);

            expect(apiClient.get).toHaveBeenCalledWith('/admin/audit/permissions/1');
        });

        it('should export audit logs', async () => {
            const mockBlob = new Blob(['test'], { type: 'text/csv' });
            (apiClient.get as ReturnType<typeof vi.fn>).mockResolvedValue(mockBlob);

            const params = { date_from: '2025-01-01', date_to: '2025-12-31' };
            await auditLogApi.export(params);

            expect(apiClient.get).toHaveBeenCalledWith('/admin/audit/permissions/export', {
                params,
                responseType: 'blob',
            });
        });

        it('should get audit log stats', async () => {
            const mockResponse = { data: { total: 100, byAction: {}, byTargetType: {} } };
            (apiClient.get as ReturnType<typeof vi.fn>).mockResolvedValue(mockResponse);

            const params = { date_from: '2025-01-01', date_to: '2025-12-31' };
            await auditLogApi.getStats(params);

            expect(apiClient.get).toHaveBeenCalledWith('/admin/audit/permissions/stats', { params });
        });
    });
});
