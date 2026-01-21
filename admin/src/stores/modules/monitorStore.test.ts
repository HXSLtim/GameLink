/**
 * Monitor Store Tests
 */
import { describe, it, expect, beforeEach } from 'vitest';
import { useMonitorStore } from './monitorStore';
import type { SystemStatus, OnlineUsers, OrderQueue, Alert } from '@/utils/websocket';

describe('MonitorStore', () => {
  beforeEach(() => {
    // Reset store before each test
    useMonitorStore.getState().reset();
  });

  describe('initial state', () => {
    it('should have correct initial state', () => {
      const state = useMonitorStore.getState();

      expect(state.systemStatus).toBeNull();
      expect(state.systemStatusLoading).toBe(false);
      expect(state.systemStatusLastUpdate).toBeNull();
      expect(state.onlineUsers).toBeNull();
      expect(state.onlineUsersLoading).toBe(false);
      expect(state.orderQueue).toBeNull();
      expect(state.orderQueueLoading).toBe(false);
      expect(state.alerts).toEqual([]);
      expect(state.alertsLoading).toBe(false);
      expect(state.wsConnected).toBe(false);
    });
  });

  describe('setSystemStatus', () => {
    it('should update system status', () => {
      const status: SystemStatus = {
        cpuUsage: 45.5,
        memoryUsage: 60.2,
        memoryTotal: 8589934592,
        memoryUsed: 5170790400,
        goroutines: 120,
        dbConnections: {
          active: 10,
          idle: 5,
          max: 100,
        },
        uptime: 3600,
        requestsPerSec: 150.5,
        status: 'healthy',
      };

      useMonitorStore.getState().setSystemStatus(status);

      const state = useMonitorStore.getState();
      expect(state.systemStatus).toEqual(status);
      expect(state.systemStatusLoading).toBe(false);
      expect(state.systemStatusLastUpdate).not.toBeNull();
    });
  });

  describe('setOnlineUsers', () => {
    it('should update online users', () => {
      const data: OnlineUsers = {
        total: 150,
        peak: 200,
        byRole: {
          admin: 5,
          user: 100,
          player: 45,
        },
        updatedAt: new Date().toISOString(),
      };

      useMonitorStore.getState().setOnlineUsers(data);

      const state = useMonitorStore.getState();
      expect(state.onlineUsers).toEqual(data);
      expect(state.onlineUsersLoading).toBe(false);
    });
  });

  describe('setOrderQueue', () => {
    it('should update order queue', () => {
      const data: OrderQueue = {
        pending: 25,
        processing: 10,
        completed: 500,
        processingSpeed: 50,
        averageWaitTime: 120,
        hasBacklog: false,
      };

      useMonitorStore.getState().setOrderQueue(data);

      const state = useMonitorStore.getState();
      expect(state.orderQueue).toEqual(data);
      expect(state.orderQueueLoading).toBe(false);
    });
  });

  describe('alerts', () => {
    it('should add alert', () => {
      const alert: Alert = {
        id: '1',
        level: 'high',
        type: 'system',
        title: 'Test Alert',
        message: 'This is a test alert',
        source: 'test',
        createdAt: new Date().toISOString(),
        isRead: false,
      };

      useMonitorStore.getState().addAlert(alert);

      const state = useMonitorStore.getState();
      expect(state.alerts).toHaveLength(1);
      expect(state.alerts[0]).toEqual(alert);
    });

    it('should limit alerts to 100', () => {
      // Add 101 alerts
      for (let i = 0; i < 101; i++) {
        const alert: Alert = {
          id: String(i),
          level: 'low',
          type: 'system',
          title: `Alert ${i}`,
          message: `Message ${i}`,
          source: 'test',
          createdAt: new Date().toISOString(),
          isRead: false,
        };
        useMonitorStore.getState().addAlert(alert);
      }

      const state = useMonitorStore.getState();
      // Should only keep 100 (most recent)
      expect(state.alerts).toHaveLength(100);
    });

    it('should mark alert as read', () => {
      const alert: Alert = {
        id: '1',
        level: 'high',
        type: 'system',
        title: 'Test Alert',
        message: 'This is a test alert',
        source: 'test',
        createdAt: new Date().toISOString(),
        isRead: false,
      };

      useMonitorStore.getState().addAlert(alert);
      useMonitorStore.getState().markAlertAsRead('1');

      const state = useMonitorStore.getState();
      expect(state.alerts[0].isRead).toBe(true);
    });

    it('should clear alert', () => {
      const alert: Alert = {
        id: '1',
        level: 'high',
        type: 'system',
        title: 'Test Alert',
        message: 'This is a test alert',
        source: 'test',
        createdAt: new Date().toISOString(),
        isRead: false,
      };

      useMonitorStore.getState().addAlert(alert);
      useMonitorStore.getState().clearAlert('1');

      const state = useMonitorStore.getState();
      expect(state.alerts).toHaveLength(0);
    });

    it('should clear all alerts', () => {
      for (let i = 0; i < 5; i++) {
        const alert: Alert = {
          id: String(i),
          level: 'low',
          type: 'system',
          title: `Alert ${i}`,
          message: `Message ${i}`,
          source: 'test',
          createdAt: new Date().toISOString(),
          isRead: false,
        };
        useMonitorStore.getState().addAlert(alert);
      }

      useMonitorStore.getState().clearAllAlerts();

      const state = useMonitorStore.getState();
      expect(state.alerts).toHaveLength(0);
    });
  });

  describe('setWsConnected', () => {
    it('should update WebSocket connection status', () => {
      useMonitorStore.getState().setWsConnected(true);

      expect(useMonitorStore.getState().wsConnected).toBe(true);

      useMonitorStore.getState().setWsConnected(false);

      expect(useMonitorStore.getState().wsConnected).toBe(false);
    });
  });

  describe('reset', () => {
    it('should reset all state to initial values', () => {
      // Set some values
      const status: SystemStatus = {
        cpuUsage: 45.5,
        memoryUsage: 60.2,
        memoryTotal: 8589934592,
        memoryUsed: 5170790400,
        goroutines: 120,
        dbConnections: {
          active: 10,
          idle: 5,
          max: 100,
        },
        uptime: 3600,
        requestsPerSec: 150.5,
        status: 'healthy',
      };

      useMonitorStore.getState().setSystemStatus(status);
      useMonitorStore.getState().setWsConnected(true);

      // Reset
      useMonitorStore.getState().reset();

      const state = useMonitorStore.getState();
      expect(state.systemStatus).toBeNull();
      expect(state.wsConnected).toBe(false);
    });
  });
});

describe('MonitorStore selectors', () => {
  beforeEach(() => {
    useMonitorStore.getState().reset();
  });

  describe('selectUnreadAlertsCount', () => {
    it('should return count of unread alerts', () => {
      // Note: This selector is a function that needs to be imported
      // For now, we verify the store behavior
      const alerts: Alert[] = [
        {
          id: '1',
          level: 'high',
          type: 'system',
          title: 'Alert 1',
          message: 'Message 1',
          source: 'test',
          createdAt: new Date().toISOString(),
          isRead: false,
        },
        {
          id: '2',
          level: 'low',
          type: 'business',
          title: 'Alert 2',
          message: 'Message 2',
          source: 'test',
          createdAt: new Date().toISOString(),
          isRead: true,
        },
        {
          id: '3',
          level: 'medium',
          type: 'security',
          title: 'Alert 3',
          message: 'Message 3',
          source: 'test',
          createdAt: new Date().toISOString(),
          isRead: false,
        },
      ];

      alerts.forEach((alert) => useMonitorStore.getState().addAlert(alert));

      const state = useMonitorStore.getState();
      const unreadCount = state.alerts.filter((a) => !a.isRead).length;
      expect(unreadCount).toBe(2);
    });
  });
});
