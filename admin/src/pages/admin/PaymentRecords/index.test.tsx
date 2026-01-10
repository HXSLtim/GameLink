/**
 * Payment Records Page Tests
 *
 * Tests for PaymentRecords page component including:
 * - Successful data loading
 * - Loading states
 * - Statistics display
 * - Filter functionality
 * - Status display
 */

import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { screen, waitFor } from '@testing-library/react';
import PaymentRecords from './index';
import { renderWithProviders, resetAllMocks, flushPromises } from '@/testutils';

// Mock antd App.useApp
const { mockMessage } = vi.hoisted(() => ({
  mockMessage: {
    success: vi.fn(),
    error: vi.fn(),
    warning: vi.fn(),
    info: vi.fn(),
    loading: vi.fn(),
  },
}));

vi.mock('antd', async () => {
  const actual = await vi.importActual('antd');
  return {
    ...actual,
    App: {
      ...((actual as Record<string, unknown>).App as Record<string, unknown>),
      useApp: () => ({
        message: mockMessage,
      }),
    },
  };
});

describe('PaymentRecords', () => {
  beforeEach(() => {
    resetAllMocks();
    mockMessage.success.mockClear();
    mockMessage.error.mockClear();
    localStorage.setItem('token', 'test-token');
    localStorage.setItem('user_role', 'admin');
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  describe('Successful Data Loading', () => {
    it('should render payment records page successfully', async () => {
      renderWithProviders(<PaymentRecords />);

      await waitFor(() => {
        expect(screen.getByText('支付记录')).toBeInTheDocument();
      });
    });

    it('should display statistics cards', async () => {
      renderWithProviders(<PaymentRecords />);

      await waitFor(() => {
        expect(screen.getByText('总支付记录')).toBeInTheDocument();
      });

      expect(screen.getByText('今日支付')).toBeInTheDocument();
      expect(screen.getByText('今日金额')).toBeInTheDocument();
      expect(screen.getByText('成功率')).toBeInTheDocument();
    });

    it('should display mock payment data', async () => {
      renderWithProviders(<PaymentRecords />);

      await waitFor(() => {
        expect(screen.getByText('ESC20251231113618283671')).toBeInTheDocument();
      });
    });

    it('should display user info', async () => {
      renderWithProviders(<PaymentRecords />);

      await waitFor(() => {
        expect(screen.getByText('用户 #4')).toBeInTheDocument();
      });
    });
  });

  describe('Table Structure', () => {
    it('should display ID column', async () => {
      renderWithProviders(<PaymentRecords />);

      await waitFor(() => {
        expect(screen.getByText('ID')).toBeInTheDocument();
      });
    });

    it('should display order number column', async () => {
      renderWithProviders(<PaymentRecords />);

      await waitFor(() => {
        expect(screen.getByText('订单号')).toBeInTheDocument();
      });
    });

    it('should display user column', async () => {
      renderWithProviders(<PaymentRecords />);

      await waitFor(() => {
        expect(screen.getByText('用户')).toBeInTheDocument();
      });
    });

    it('should display amount column', async () => {
      renderWithProviders(<PaymentRecords />);

      await waitFor(() => {
        expect(screen.getByText('金额')).toBeInTheDocument();
      });
    });

    it('should display transaction ID column', async () => {
      renderWithProviders(<PaymentRecords />);

      await waitFor(() => {
        expect(screen.getByText('交易流水号')).toBeInTheDocument();
      });
    });

    it('should display created time column', async () => {
      renderWithProviders(<PaymentRecords />);

      await waitFor(() => {
        expect(screen.getByText('创建时间')).toBeInTheDocument();
      });
    });

    it('should have proper table structure', async () => {
      renderWithProviders(<PaymentRecords />);

      await waitFor(() => {
        expect(screen.getByRole('table')).toBeInTheDocument();
      });
    });
  });

  describe('Status Display', () => {
    it('should display success status', async () => {
      renderWithProviders(<PaymentRecords />);

      await waitFor(() => {
        expect(screen.getByText('已支付')).toBeInTheDocument();
      });
    });

    it('should display refunded status', async () => {
      renderWithProviders(<PaymentRecords />);

      await waitFor(() => {
        expect(screen.getByText('已退款')).toBeInTheDocument();
      });
    });
  });

  describe('Payment Method Display', () => {
    it('should display wechat payment method', async () => {
      renderWithProviders(<PaymentRecords />);

      await waitFor(() => {
        expect(screen.getByText('微信支付')).toBeInTheDocument();
      });
    });

    it('should display alipay payment method', async () => {
      renderWithProviders(<PaymentRecords />);

      await waitFor(() => {
        expect(screen.getByText('支付宝')).toBeInTheDocument();
      });
    });
  });

  describe('Amount Display', () => {
    it('should display formatted amount', async () => {
      renderWithProviders(<PaymentRecords />);

      await waitFor(() => {
        expect(screen.getByText('¥219.00')).toBeInTheDocument();
      });
    });
  });

  describe('Toolbar Buttons', () => {
    it('should display refresh button', async () => {
      renderWithProviders(<PaymentRecords />);

      await waitFor(() => {
        expect(screen.getByText('刷新')).toBeInTheDocument();
      });
    });

    it('should display export button', async () => {
      renderWithProviders(<PaymentRecords />);

      await waitFor(() => {
        expect(screen.getByText('导出')).toBeInTheDocument();
      });
    });
  });

  describe('Filter Section', () => {
    it('should display search input', async () => {
      renderWithProviders(<PaymentRecords />);

      await waitFor(() => {
        expect(screen.getByPlaceholderText('搜索订单号')).toBeInTheDocument();
      });
    });

    it('should display payment method filter', async () => {
      renderWithProviders(<PaymentRecords />);

      await waitFor(() => {
        expect(screen.getByText('支付记录')).toBeInTheDocument();
      });

      // Payment method select should exist
      const selects = document.querySelectorAll('.ant-select');
      expect(selects.length).toBeGreaterThan(0);
    });

    it('should display search button', async () => {
      renderWithProviders(<PaymentRecords />);

      await waitFor(() => {
        expect(screen.getByText('搜索')).toBeInTheDocument();
      });
    });

    it('should display filter buttons', async () => {
      renderWithProviders(<PaymentRecords />);

      await waitFor(() => {
        expect(screen.getByText('搜索')).toBeInTheDocument();
      });

      // Check for buttons in the filter area
      const buttons = screen.getAllByRole('button');
      expect(buttons.length).toBeGreaterThan(0);
    });
  });

  describe('Pagination', () => {
    it('should display pagination', async () => {
      renderWithProviders(<PaymentRecords />);

      await waitFor(() => {
        expect(screen.getByText('支付记录')).toBeInTheDocument();
      });

      await flushPromises();

      // Pagination should exist
      const pagination = document.querySelector('.ant-pagination');
      expect(pagination).toBeInTheDocument();
    });
  });

  describe('Table Structure', () => {
    it('should have proper table structure', async () => {
      renderWithProviders(<PaymentRecords />);

      await waitFor(() => {
        expect(screen.getByRole('table')).toBeInTheDocument();
      });
    });
  });
});
