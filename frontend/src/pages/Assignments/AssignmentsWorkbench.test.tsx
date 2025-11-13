import { describe, it, expect, vi, beforeEach } from 'vitest';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, screen, waitFor } from '@testing-library/react';
import { renderWithProviders } from '../../test/utils/test-utils';
import { AssignmentsWorkbench } from './AssignmentsWorkbench';
import { assignmentApi } from '../../services/api/assignment';
import { orderApi } from '../../services/api/order';
import { message } from '../../components';
import { OrderStatus } from '../../types/order';

vi.mock('../../services/api/assignment', () => ({
  assignmentApi: {
    getPending: vi.fn(),
    getCandidates: vi.fn(),
    assign: vi.fn(),
    cancelAssignment: vi.fn(),
    getDisputes: vi.fn(),
    mediate: vi.fn(),
  },
}));

vi.mock('../../services/api/order', () => ({
  orderApi: {
    getDetail: vi.fn(),
    getLogs: vi.fn(),
  },
}));

const renderAssignments = () => {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
      },
    },
  });

  const result = renderWithProviders(
    <QueryClientProvider client={queryClient}>
      <AssignmentsWorkbench />
    </QueryClientProvider>,
  );

  return { queryClient, ...result };
};

describe('AssignmentsWorkbench', () => {
  beforeEach(() => {
    vi.resetAllMocks();
    (message as any).success = vi.fn();
    (message as any).error = vi.fn();
    vi.mocked(assignmentApi.getPending).mockResolvedValue({
      list: [
        {
          orderId: 1,
          userId: 10,
          status: OrderStatus.PENDING,
          assignmentSource: 'manual',
          createdAt: '2025-01-01T00:00:00Z',
          slaDeadline: '2025-01-01T00:30:00Z',
          slaRemainingSeconds: 1800,
          isOverdue: false,
        },
      ],
      total: 1,
      page: 1,
      page_size: 10,
    });
    vi.mocked(orderApi.getDetail).mockResolvedValue({
      id: 1,
      userId: 10,
      gameId: 9,
      title: '陪玩测试订单',
      status: OrderStatus.PENDING,
      priceCents: 10000,
      createdAt: '2025-01-01T00:00:00Z',
      updatedAt: '2025-01-01T00:05:00Z',
      assignmentSource: 'manual',
      disputeStatus: 'pending',
    } as any);
    vi.mocked(orderApi.getLogs).mockResolvedValue([]);
    vi.mocked(assignmentApi.getCandidates).mockResolvedValue([]);
    vi.mocked(assignmentApi.getDisputes).mockResolvedValue([]);
    vi.mocked(assignmentApi.assign).mockResolvedValue();
    vi.mocked(assignmentApi.cancelAssignment).mockResolvedValue();
    vi.mocked(assignmentApi.mediate).mockResolvedValue({} as any);
  });

  it('renders pending assignments and detail card', async () => {
    const { queryClient } = renderAssignments();

    await waitFor(() => {
      expect(assignmentApi.getPending).toHaveBeenCalled();
    });

    expect(await screen.findByText('订单 #1')).toBeInTheDocument();
    expect(screen.getByText('陪玩测试订单')).toBeInTheDocument();
    expect(orderApi.getDetail).toHaveBeenCalledWith(1);
    queryClient.clear();
  });

  it('allows assigning candidate players', async () => {
    vi.mocked(assignmentApi.getCandidates).mockResolvedValue([
      {
        playerId: 88,
        nickname: '推荐陪玩师',
        hourlyRateCents: 2000,
        score: 4.8,
        source: 'recommendation',
        reason: '响应迅速',
      },
    ]);

    const { queryClient } = renderAssignments();

    const assignButton = await screen.findByRole('button', { name: '指派' });
    fireEvent.click(assignButton);

    await waitFor(() => {
      expect(assignmentApi.assign).toHaveBeenCalledWith(1, {
        playerId: 88,
        source: 'recommendation',
      });
    });
    queryClient.clear();
  });
});
