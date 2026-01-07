/**
 * useNotification Hook
 * Manages browser notification permissions and SLA alerts
 */
import { useState, useEffect, useCallback, useRef } from 'react';
import {
  requestNotificationPermission,
  isNotificationEnabled,
  getNotificationPermission,
  notifySLABreach,
  playAlertSound,
  getNotificationPreferences,
  saveNotificationPreferences,
} from '@/utils/notification';

interface UseNotificationOptions {
  /** Enable SLA monitoring */
  enableSLAMonitoring?: boolean;
  /** SLA check interval in milliseconds (default: 60000 = 1 minute) */
  checkInterval?: number;
  /** Callback to get current SLA breach count */
  getSLABreachCount?: () => number;
  /** Callback when SLA alert is triggered */
  onSLAAlert?: () => void;
}

interface UseNotificationReturn {
  /** Whether browser notifications are enabled */
  isEnabled: boolean;
  /** Current permission status */
  permission: NotificationPermission;
  /** Request notification permission */
  requestPermission: () => Promise<boolean>;
  /** User preferences */
  preferences: {
    soundEnabled: boolean;
    browserNotificationEnabled: boolean;
    slaAlertEnabled: boolean;
  };
  /** Update preferences */
  updatePreferences: (prefs: Partial<UseNotificationReturn['preferences']>) => void;
  /** Manually trigger SLA alert */
  triggerSLAAlert: (count: number) => void;
  /** Test notification sound */
  testSound: (type?: 'warning' | 'urgent' | 'success') => void;
}

export function useNotification(options: UseNotificationOptions = {}): UseNotificationReturn {
  const {
    enableSLAMonitoring = false,
    checkInterval = 60000,
    getSLABreachCount,
    onSLAAlert,
  } = options;

  // Initialize state with functions (lazy initialization)
  const [isEnabled, setIsEnabled] = useState(() => isNotificationEnabled());
  const [permission, setPermission] = useState<NotificationPermission>(() => getNotificationPermission());
  const [preferences, setPreferences] = useState(() => getNotificationPreferences());
  
  // Track last alerted count to avoid duplicate alerts
  const lastAlertedCount = useRef(0);

  // Request permission
  const requestPermission = useCallback(async () => {
    const granted = await requestNotificationPermission();
    setIsEnabled(granted);
    setPermission(getNotificationPermission());
    return granted;
  }, []);

  // Update preferences
  const updatePreferences = useCallback((prefs: Partial<UseNotificationReturn['preferences']>) => {
    setPreferences(prev => {
      const updated = { ...prev, ...prefs };
      saveNotificationPreferences(updated);
      return updated;
    });
  }, []);

  // Trigger SLA alert
  const triggerSLAAlert = useCallback((count: number) => {
    if (count <= 0) return;
    
    // Only alert if count increased
    if (count <= lastAlertedCount.current) return;
    lastAlertedCount.current = count;

    if (preferences.soundEnabled) {
      playAlertSound('urgent');
    }

    if (preferences.browserNotificationEnabled && isEnabled) {
      notifySLABreach(count, onSLAAlert);
    }
  }, [preferences, isEnabled, onSLAAlert]);

  // Test sound
  const testSound = useCallback((type: 'warning' | 'urgent' | 'success' = 'warning') => {
    playAlertSound(type);
  }, []);

  // SLA monitoring interval
  useEffect(() => {
    if (!enableSLAMonitoring || !getSLABreachCount || !preferences.slaAlertEnabled) {
      return;
    }

    const checkSLA = () => {
      const count = getSLABreachCount();
      if (count > 0) {
        triggerSLAAlert(count);
      }
    };

    // Initial check
    checkSLA();

    // Set up interval
    const intervalId = setInterval(checkSLA, checkInterval);

    return () => clearInterval(intervalId);
  }, [enableSLAMonitoring, getSLABreachCount, checkInterval, preferences.slaAlertEnabled, triggerSLAAlert]);

  return {
    isEnabled,
    permission,
    requestPermission,
    preferences,
    updatePreferences,
    triggerSLAAlert,
    testSound,
  };
}

export default useNotification;
