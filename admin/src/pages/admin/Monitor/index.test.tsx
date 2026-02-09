/**
 * Monitor Page Unit Tests
 */
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import MonitorPage from './index';
import { useMonitorStore } from '@/stores/modules/monitorStore';
import { wsManager } from '@/utils/websocket';
import { MessageType } from '@/utils/websocket';
import type { PageContainerProps } from '@/components/PageContainer';
import type { EventListeners, MessageHandler, WSMessage } from '@/utils/websocket';
import type { TMessageType } from '@/utils/websocket/types';

// Mock the dependencies
vi.mock('@/stores/modules/monitorStore');
vi.mock('@/utils/websocket');
vi.mock('@/components', () => ({
  PageContainer: ({ children, title, subTitle }: PageContainerProps) => (
    <div data-testid="page-container">
      <h1>{title}</h1>
      <h2>{subTitle}</h2>
      {children}
    </div>
  ),
}));

describe('MonitorPage', () => {
  const mockSetSystemStatus = vi.fn();
  const mockSetOnlineUsers = vi.fn();
  const mockSetOrderQueue = vi.fn();
  const mockAddAlert = vi.fn();
  const mockSetWsConnected = vi.fn();
  const mockMarkAlertAsRead = vi.fn();
  const mockClearAlert = vi.fn();
  const mockClearAllAlerts = vi.fn();

  const createMockStoreState = (
    overrides: Partial<ReturnType<typeof useMonitorStore>> = {}
  ): ReturnType<typeof useMonitorStore> => ({
    systemStatus: mockSystemStatus,
    systemStatusLoading: false,
    systemStatusLastUpdate: null,
    onlineUsers: mockOnlineUsers,
    onlineUsersLoading: false,
    orderQueue: mockOrderQueue,
    orderQueueLoading: false,
    alerts: mockAlerts,
    alertsLoading: false,
    wsConnected: true,
    setSystemStatus: mockSetSystemStatus,
    setOnlineUsers: mockSetOnlineUsers,
    setOrderQueue: mockSetOrderQueue,
    addAlert: mockAddAlert,
    markAlertAsRead: mockMarkAlertAsRead,
    clearAlert: mockClearAlert,
    clearAllAlerts: mockClearAllAlerts,
    setWsConnected: mockSetWsConnected,
    reset: vi.fn(),
    ...overrides,
  });

  const mockSystemStatus = {
    cpuUsage: 45,
    memoryUsed: 4294967296,
    memoryTotal: 8589934592,
    goroutines: 245,
    dbConnections: { active: 15, idle: 5, max: 20 },
    uptime: 259200,
    requestsPerSec: 123.45,
    status: 'healthy',
  };

  const mockOnlineUsers = {
    total: 1234,
    peak: 1500,
    updatedAt: Date.now(),
    byRole: {
      ADMIN: 10,
      PLAYER: 200,
      USER: 1024,
    },
  };

  const mockOrderQueue = {
    pending: 5,
    processing: 3,
    completed: 1500,
    hasBacklog: false,
    processingSpeed: 45.67,
    averageWaitTime: 120,
  };

  const mockAlerts = [
    {
      id: '1',
      title: '系统告警',
      message: 'CPU 使用率超过 80%',
      level: 'high' as const,
      type: 'system',
      source: 'monitor',
      isRead: false,
      createdAt: new Date().toISOString(),
    },
    {
      id: '2',
      title: '订单积压',
      message: '待处理订单超过 100 个',
      level: 'medium' as const,
      type: 'business',
      source: 'order',
      isRead: true,
      createdAt: new Date().toISOString(),
    },
  ];

  beforeEach(() => {
    vi.clearAllMocks();

    // Mock monitor store
    vi.mocked(useMonitorStore).mockReturnValue(createMockStoreState());

    // Mock wsManager
    vi.mocked(wsManager.on).mockReturnValue(vi.fn());
    vi.mocked(wsManager.setEventListeners).mockReturnValue(undefined);
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  describe('Page Rendering', () => {
    it('should render the monitor page with correct title', () => {
      render(<MonitorPage />);
      expect(screen.getByText('实时监控')).toBeInTheDocument();
      expect(screen.getByText('系统运行状态实时数据')).toBeInTheDocument();
    });

    it('should render system status card', () => {
      render(<MonitorPage />);
      expect(screen.getByText('系统状态')).toBeInTheDocument();
    });

    it('should render online users card', () => {
      render(<MonitorPage />);
      expect(screen.getByText('在线用户')).toBeInTheDocument();
    });

    it('should render order queue card', () => {
      render(<MonitorPage />);
      expect(screen.getByText('订单队列')).toBeInTheDocument();
    });

    it('should render alerts card', () => {
      render(<MonitorPage />);
      expect(screen.getByText('警告通知')).toBeInTheDocument();
    });
  });

  describe('System Status Card', () => {
    it('should display CPU usage', () => {
      render(<MonitorPage />);
      expect(screen.getByText('CPU 使用率')).toBeInTheDocument();
      expect(screen.getByText('45')).toBeInTheDocument();
    });

    it('should display memory usage', () => {
      render(<MonitorPage />);
      expect(screen.getByText('内存使用')).toBeInTheDocument();
      expect(screen.getByText(/4\.00/)).toBeInTheDocument();
      expect(screen.getByText(/8\.00 GB/)).toBeInTheDocument();
    });

    it('should display goroutines count', () => {
      render(<MonitorPage />);
      expect(screen.getByText('Goroutines')).toBeInTheDocument();
      expect(screen.getByText('245')).toBeInTheDocument();
    });

    it('should display database connections', () => {
      render(<MonitorPage />);
      expect(screen.getByText('数据库连接')).toBeInTheDocument();
      expect(screen.getByText('15')).toBeInTheDocument();
      expect(screen.getByText(/20/)).toBeInTheDocument();
    });

    it('should display uptime in hours', () => {
      render(<MonitorPage />);
      expect(screen.getByText('运行时间')).toBeInTheDocument();
      expect(screen.getByText('72')).toBeInTheDocument();
    });

    it('should display requests per second', () => {
      render(<MonitorPage />);
      expect(screen.getByText('请求/秒')).toBeInTheDocument();
      expect(screen.getByText('123.45')).toBeInTheDocument();
    });

    it('should show loading state when system status is null', () => {
      vi.mocked(useMonitorStore).mockReturnValue(
        createMockStoreState({ systemStatus: null })
      );

      render(<MonitorPage />);
      expect(screen.getByText('等待数据...')).toBeInTheDocument();
    });
  });

  describe('Online Users Card', () => {
    it('should display total online users', () => {
      render(<MonitorPage />);
      expect(screen.getByText('当前在线')).toBeInTheDocument();
      expect(screen.getByText('1,234')).toBeInTheDocument();
    });

    it('should display peak users', () => {
      render(<MonitorPage />);
      expect(screen.getByText('峰值')).toBeInTheDocument();
      expect(screen.getByText('1,500')).toBeInTheDocument();
    });

    it('should display users by role', () => {
      render(<MonitorPage />);
      expect(screen.getByText('按角色分布：')).toBeInTheDocument();
      expect(screen.getByText('ADMIN: 10')).toBeInTheDocument();
      expect(screen.getByText('PLAYER: 200')).toBeInTheDocument();
      expect(screen.getByText('USER: 1,024')).toBeInTheDocument();
    });

    it('should show loading state when online users is null', () => {
      vi.mocked(useMonitorStore).mockReturnValue(
        createMockStoreState({ onlineUsers: null })
      );

      render(<MonitorPage />);
      expect(screen.getByText('等待数据...')).toBeInTheDocument();
    });
  });

  describe('Order Queue Card', () => {
    it('should display pending orders', () => {
      render(<MonitorPage />);
      expect(screen.getByText('待处理')).toBeInTheDocument();
      expect(screen.getByText('5')).toBeInTheDocument();
    });

    it('should display processing orders', () => {
      render(<MonitorPage />);
      expect(screen.getByText('处理中')).toBeInTheDocument();
      expect(screen.getByText('3')).toBeInTheDocument();
    });

    it('should display completed orders', () => {
      render(<MonitorPage />);
      expect(screen.getByText('已完成')).toBeInTheDocument();
      expect(screen.getByText('1,500')).toBeInTheDocument();
    });

    it('should display queue status as normal', () => {
      render(<MonitorPage />);
      expect(screen.getByText('队列状态')).toBeInTheDocument();
      expect(screen.getByText('正常')).toBeInTheDocument();
    });

    it('should display queue status as backlog when hasBacklog is true', () => {
      vi.mocked(useMonitorStore).mockReturnValue(
        createMockStoreState({
          orderQueue: { ...mockOrderQueue, hasBacklog: true },
        })
      );

      render(<MonitorPage />);
      expect(screen.getByText('积压')).toBeInTheDocument();
    });

    it('should display processing speed', () => {
      render(<MonitorPage />);
      expect(screen.getByText('处理速度')).toBeInTheDocument();
      expect(screen.getByText('45.67')).toBeInTheDocument();
      expect(screen.getByText('单/分钟')).toBeInTheDocument();
    });

    it('should display average wait time', () => {
      render(<MonitorPage />);
      expect(screen.getByText('平均等待')).toBeInTheDocument();
      expect(screen.getByText('120')).toBeInTheDocument();
      expect(screen.getByText('秒')).toBeInTheDocument();
    });

    it('should show loading state when order queue is null', () => {
      vi.mocked(useMonitorStore).mockReturnValue(
        createMockStoreState({ orderQueue: null })
      );

      render(<MonitorPage />);
      expect(screen.getByText('等待数据...')).toBeInTheDocument();
    });
  });

  describe('Alerts Card', () => {
    it('should display alerts count badge', () => {
      render(<MonitorPage />);
      const badge = screen.getByText('1');
      expect(badge).toBeInTheDocument();
    });

    it('should display alert list', () => {
      render(<MonitorPage />);
      expect(screen.getByText('系统告警')).toBeInTheDocument();
      expect(screen.getByText('CPU 使用率超过 80%')).toBeInTheDocument();
      expect(screen.getByText('订单积压')).toBeInTheDocument();
      expect(screen.getByText('待处理订单超过 100 个')).toBeInTheDocument();
    });

    it('should show empty state when no alerts', () => {
      vi.mocked(useMonitorStore).mockReturnValue(
        createMockStoreState({ alerts: [] })
      );

      render(<MonitorPage />);
      expect(screen.getByText('暂无警告')).toBeInTheDocument();
      expect(screen.getByText('系统运行正常，没有新的警告通知')).toBeInTheDocument();
    });

    it('should call clearAllAlerts when clicking clear all button', async () => {
      const user = userEvent.setup();
      render(<MonitorPage />);

      const clearButton = screen.getByText('清空全部');
      await user.click(clearButton);

      expect(mockClearAllAlerts).toHaveBeenCalledTimes(1);
    });

    it('should call markAlertAsRead when clicking mark as read button', async () => {
      const user = userEvent.setup();
      render(<MonitorPage />);

      const markReadButton = screen.getAllByText('标记已读')[0];
      await user.click(markReadButton);

      expect(mockMarkAlertAsRead).toHaveBeenCalledWith('1');
    });

    it('should call clearAlert when clicking delete button', async () => {
      const user = userEvent.setup();
      render(<MonitorPage />);

      const deleteButtons = screen.getAllByText('删除');
      await user.click(deleteButtons[0]);

      expect(mockClearAlert).toHaveBeenCalledWith('1');
    });

    it('should not show mark as read button for read alerts', () => {
      render(<MonitorPage />);
      const markReadButtons = screen.getAllByText('标记已读');
      expect(markReadButtons).toHaveLength(1);
    });

    it('should show alert type tags', () => {
      render(<MonitorPage />);
      expect(screen.getByText('SYSTEM')).toBeInTheDocument();
      expect(screen.getByText('BUSINESS')).toBeInTheDocument();
    });

    it('should show alert level tags', () => {
      render(<MonitorPage />);
      expect(screen.getByText('HIGH')).toBeInTheDocument();
      expect(screen.getByText('MEDIUM')).toBeInTheDocument();
    });

    it('should show alert source and timestamp', () => {
      render(<MonitorPage />);
      expect(screen.getByText('monitor')).toBeInTheDocument();
      expect(screen.getByText('order')).toBeInTheDocument();
    });
  });

  describe('WebSocket Connection', () => {
    it('should register WebSocket handlers on mount', () => {
      render(<MonitorPage />);

      expect(wsManager.on).toHaveBeenCalledWith(
        MessageType.SystemStatus,
        expect.any(Function)
      );
      expect(wsManager.on).toHaveBeenCalledWith(
        MessageType.OnlineUsers,
        expect.any(Function)
      );
      expect(wsManager.on).toHaveBeenCalledWith(
        MessageType.OrderQueue,
        expect.any(Function)
      );
      expect(wsManager.on).toHaveBeenCalledWith(
        MessageType.Alert,
        expect.any(Function)
      );
    });

    it('should set event listeners on mount', () => {
      render(<MonitorPage />);

      expect(wsManager.setEventListeners).toHaveBeenCalledWith({
        onOpen: expect.any(Function),
        onClose: expect.any(Function),
        onError: expect.any(Function),
      });
    });

    it('should set wsConnected to true when WebSocket opens', () => {
      vi.mocked(wsManager.setEventListeners).mockImplementation((listeners: EventListeners) => {
        if (listeners?.onOpen) {
          listeners.onOpen(new Event('open'));
        }
      });

      render(<MonitorPage />);

      expect(mockSetWsConnected).toHaveBeenCalledWith(true);
    });

    it('should set wsConnected to false when WebSocket closes', () => {
      vi.mocked(wsManager.setEventListeners).mockImplementation((listeners: EventListeners) => {
        if (listeners?.onClose) {
          listeners.onClose(new CloseEvent('close'));
        }
      });

      render(<MonitorPage />);

      expect(mockSetWsConnected).toHaveBeenCalledWith(false);
    });

    it('should set wsConnected to false when WebSocket errors', () => {
      vi.mocked(wsManager.setEventListeners).mockImplementation((listeners: EventListeners) => {
        if (listeners?.onError) {
          listeners.onError(new Event('error'));
        }
      });

      render(<MonitorPage />);

      expect(mockSetWsConnected).toHaveBeenCalledWith(false);
    });

    it('should update system status when receiving SystemStatus message', () => {
      let systemStatusHandler: ((message: WSMessage) => void) | undefined;
      vi.mocked(wsManager.on).mockImplementation((messageType: TMessageType, handler: MessageHandler) => {
        if (messageType === MessageType.SystemStatus) {
          systemStatusHandler = handler;
        }
        return vi.fn();
      });

      render(<MonitorPage />);

      systemStatusHandler?.({
        type: MessageType.SystemStatus,
        timestamp: new Date().toISOString(),
        data: mockSystemStatus,
      });

      expect(mockSetSystemStatus).toHaveBeenCalledWith(mockSystemStatus);
    });

    it('should update online users when receiving OnlineUsers message', () => {
      let onlineUsersHandler: ((message: WSMessage) => void) | undefined;
      vi.mocked(wsManager.on).mockImplementation((messageType: TMessageType, handler: MessageHandler) => {
        if (messageType === MessageType.OnlineUsers) {
          onlineUsersHandler = handler;
        }
        return vi.fn();
      });

      render(<MonitorPage />);

      onlineUsersHandler?.({
        type: MessageType.OnlineUsers,
        timestamp: new Date().toISOString(),
        data: mockOnlineUsers,
      });

      expect(mockSetOnlineUsers).toHaveBeenCalledWith(mockOnlineUsers);
    });

    it('should update order queue when receiving OrderQueue message', () => {
      let orderQueueHandler: ((message: WSMessage) => void) | undefined;
      vi.mocked(wsManager.on).mockImplementation((messageType: TMessageType, handler: MessageHandler) => {
        if (messageType === MessageType.OrderQueue) {
          orderQueueHandler = handler;
        }
        return vi.fn();
      });

      render(<MonitorPage />);

      orderQueueHandler?.({
        type: MessageType.OrderQueue,
        timestamp: new Date().toISOString(),
        data: mockOrderQueue,
      });

      expect(mockSetOrderQueue).toHaveBeenCalledWith(mockOrderQueue);
    });

    it('should add alert when receiving Alert message', () => {
      let alertHandler: ((message: WSMessage) => void) | undefined;
      vi.mocked(wsManager.on).mockImplementation((messageType: TMessageType, handler: MessageHandler) => {
        if (messageType === MessageType.Alert) {
          alertHandler = handler;
        }
        return vi.fn();
      });

      render(<MonitorPage />);

      const newAlert = { ...mockAlerts[0], id: '3' };
      alertHandler?.({
        type: MessageType.Alert,
        timestamp: new Date().toISOString(),
        data: newAlert,
      });

      expect(mockAddAlert).toHaveBeenCalledWith(newAlert);
    });

    it('should cleanup handlers on unmount', () => {
      const mockUnsubscribe = vi.fn();
      vi.mocked(wsManager.on).mockReturnValue(mockUnsubscribe);

      const { unmount } = render(<MonitorPage />);
      unmount();

      expect(mockUnsubscribe).toHaveBeenCalledTimes(4);
    });
  });

  describe('WebSocket Connection Status', () => {
    it('should show connected tag when wsConnected is true', () => {
      render(<MonitorPage />);
      expect(screen.getByText('已连接')).toBeInTheDocument();
    });

    it('should show disconnected tag when wsConnected is false', () => {
      vi.mocked(useMonitorStore).mockReturnValue(
        createMockStoreState({ wsConnected: false })
      );

      render(<MonitorPage />);
      expect(screen.getByText('未连接')).toBeInTheDocument();
    });
  });

  describe('Progress Bars', () => {
    it('should show normal progress when CPU usage is below 80%', () => {
      render(<MonitorPage />);
      // CPU is 45%, should show as active (green)
      expect(screen.getByText('45%')).toBeInTheDocument();
    });

    it('should show exception progress when CPU usage is above 80%', () => {
      vi.mocked(useMonitorStore).mockReturnValue(
        createMockStoreState({
          systemStatus: { ...mockSystemStatus, cpuUsage: 85 },
        })
      );

      render(<MonitorPage />);
      expect(screen.getByText('85%')).toBeInTheDocument();
    });

    it('should calculate memory percentage correctly', () => {
      render(<MonitorPage />);
      // 4GB / 8GB = 50%
      expect(screen.getByText('4.00')).toBeInTheDocument();
    });
  });
});
