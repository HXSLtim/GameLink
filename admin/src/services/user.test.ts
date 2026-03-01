import { beforeEach, describe, expect, it, vi } from 'vitest';
import { adminApi } from '@/api/admin';
import {
  batchAddUserPoints,
  batchSendNotification,
  batchUpdateUserRole,
  batchUpdateUserStatus,
  userBatchService,
} from './user';
import type {
  BatchNotificationRequest,
  BatchPointsRequest,
  BatchRoleRequest,
  BatchStatusRequest,
} from '@/types/user';

vi.mock('@/api/admin', () => ({
  adminApi: {
    batchUpdateUserRole: vi.fn(),
    batchUpdateUserStatus: vi.fn(),
    batchAddUserPoints: vi.fn(),
    batchSendNotification: vi.fn(),
  },
}));

const asMock = (fn: unknown) => fn as ReturnType<typeof vi.fn>;

describe('user batch service', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('calls batch role api and returns payload', async () => {
    const payload: BatchRoleRequest = {
      userIds: [1, 2],
      role: 'player',
    };
    const responsePayload = { success: true, data: null, message: 'ok' };
    asMock(adminApi.batchUpdateUserRole).mockResolvedValue({ data: responsePayload });

    const result = await batchUpdateUserRole(payload);

    expect(adminApi.batchUpdateUserRole).toHaveBeenCalledWith(payload);
    expect(result).toEqual(responsePayload);
  });

  it('calls batch status api and returns payload', async () => {
    const payload: BatchStatusRequest = {
      userIds: [3, 4],
      status: 'banned',
    };
    const responsePayload = { success: true, data: null, message: 'ok' };
    asMock(adminApi.batchUpdateUserStatus).mockResolvedValue({ data: responsePayload });

    const result = await batchUpdateUserStatus(payload);

    expect(adminApi.batchUpdateUserStatus).toHaveBeenCalledWith(payload);
    expect(result).toEqual(responsePayload);
  });

  it('calls batch points api and returns payload', async () => {
    const payload: BatchPointsRequest = {
      target: 'users',
      userIds: [8],
      cents: 1000,
      reason: '活动补偿',
      type: 'activity',
    };
    const responsePayload = { success: true, data: null, message: 'ok' };
    asMock(adminApi.batchAddUserPoints).mockResolvedValue({ data: responsePayload });

    const result = await batchAddUserPoints(payload);

    expect(adminApi.batchAddUserPoints).toHaveBeenCalledWith(payload);
    expect(result).toEqual(responsePayload);
  });

  it('calls batch notification api and returns payload', async () => {
    const payload: BatchNotificationRequest = {
      target: 'role',
      roles: ['player'],
      title: '系统通知',
      content: '批量通知内容',
      type: 'system',
    };
    const responsePayload = { success: true, data: null, message: 'ok' };
    asMock(adminApi.batchSendNotification).mockResolvedValue({ data: responsePayload });

    const result = await batchSendNotification(payload);

    expect(adminApi.batchSendNotification).toHaveBeenCalledWith(payload);
    expect(result).toEqual(responsePayload);
  });

  it('exposes all batch apis from service object', () => {
    expect(userBatchService.batchUpdateRole).toBe(batchUpdateUserRole);
    expect(userBatchService.batchUpdateStatus).toBe(batchUpdateUserStatus);
    expect(userBatchService.batchAddPoints).toBe(batchAddUserPoints);
    expect(userBatchService.batchSendNotification).toBe(batchSendNotification);
  });
});
