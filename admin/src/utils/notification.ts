/**
 * Browser Notification Utilities
 * Provides browser notifications and sound alerts for urgent events like SLA breaches
 */

// Audio context for sound alerts
let audioContext: AudioContext | null = null;

/**
 * Request browser notification permission
 * Should be called on user interaction (e.g., button click)
 */
export async function requestNotificationPermission(): Promise<boolean> {
  if (!('Notification' in window)) {
    console.warn('Browser does not support notifications');
    return false;
  }

  try {
    const permission = await Notification.requestPermission();
    return permission === 'granted';
  } catch (error) {
    console.error('Failed to request notification permission:', error);
    return false;
  }
}

/**
 * Check if notifications are enabled
 */
export function isNotificationEnabled(): boolean {
  return 'Notification' in window && Notification.permission === 'granted';
}

/**
 * Get current notification permission status
 */
export function getNotificationPermission(): NotificationPermission {
  if (!('Notification' in window)) {
    return 'denied';
  }
  return Notification.permission;
}

/**
 * Show a browser notification
 */
export function showNotification(
  title: string,
  options?: NotificationOptions & { onClick?: () => void }
): Notification | null {
  if (!isNotificationEnabled()) {
    console.warn('Notifications not enabled');
    return null;
  }

  try {
    const notification = new Notification(title, {
      icon: '/favicon.ico',
      badge: '/favicon.ico',
      ...options,
    });

    if (options?.onClick) {
      notification.onclick = () => {
        window.focus();
        options.onClick?.();
        notification.close();
      };
    }

    // Auto close after 10 seconds
    setTimeout(() => notification.close(), 10000);

    return notification;
  } catch (error) {
    console.error('Failed to show notification:', error);
    return null;
  }
}

/**
 * Play an alert sound
 * Uses Web Audio API to generate a simple beep sound
 */
export function playAlertSound(type: 'warning' | 'urgent' | 'success' = 'warning'): void {
  try {
    // Create audio context on first use (must be after user interaction)
    if (!audioContext) {
      audioContext = new (window.AudioContext || (window as unknown as { webkitAudioContext: typeof AudioContext }).webkitAudioContext)();
    }

    const oscillator = audioContext.createOscillator();
    const gainNode = audioContext.createGain();

    oscillator.connect(gainNode);
    gainNode.connect(audioContext.destination);

    // Different sounds for different alert types
    switch (type) {
      case 'urgent':
        // High-pitched urgent beep (repeated)
        oscillator.frequency.value = 880; // A5
        oscillator.type = 'square';
        gainNode.gain.value = 0.3;
        playBeepPattern(oscillator, gainNode, audioContext, [100, 50, 100, 50, 100]);
        break;
      case 'success':
        // Pleasant success chime
        oscillator.frequency.value = 523.25; // C5
        oscillator.type = 'sine';
        gainNode.gain.value = 0.2;
        playBeepPattern(oscillator, gainNode, audioContext, [150]);
        break;
      case 'warning':
      default:
        // Medium warning beep
        oscillator.frequency.value = 660; // E5
        oscillator.type = 'triangle';
        gainNode.gain.value = 0.25;
        playBeepPattern(oscillator, gainNode, audioContext, [200, 100, 200]);
        break;
    }
  } catch (error) {
    console.error('Failed to play alert sound:', error);
  }
}

/**
 * Play a beep pattern
 */
function playBeepPattern(
  oscillator: OscillatorNode,
  gainNode: GainNode,
  context: AudioContext,
  pattern: number[]
): void {
  let time = context.currentTime;
  
  oscillator.start(time);
  
  pattern.forEach((duration, index) => {
    if (index % 2 === 0) {
      // Sound on
      gainNode.gain.setValueAtTime(gainNode.gain.value, time);
    } else {
      // Sound off (gap)
      gainNode.gain.setValueAtTime(0, time);
    }
    time += duration / 1000;
  });
  
  // Final fade out
  gainNode.gain.setValueAtTime(gainNode.gain.value, time);
  gainNode.gain.exponentialRampToValueAtTime(0.001, time + 0.1);
  
  oscillator.stop(time + 0.1);
}

/**
 * Show SLA breach notification with sound
 */
export function notifySLABreach(count: number, onClick?: () => void): void {
  // Play urgent sound
  playAlertSound('urgent');

  // Show browser notification
  showNotification('⚠️ SLA 超时警告', {
    body: `有 ${count} 个纠纷已超过 SLA 时限，请立即处理！`,
    tag: 'sla-breach', // Prevents duplicate notifications
    requireInteraction: true, // Keep notification visible until user interacts
    onClick,
  });
}

/**
 * Show new dispute notification
 */
export function notifyNewDispute(orderNo: string, onClick?: () => void): void {
  playAlertSound('warning');

  showNotification('📋 新纠纷待处理', {
    body: `订单 ${orderNo} 有新的纠纷需要处理`,
    tag: `dispute-${orderNo}`,
    onClick,
  });
}

/**
 * Show dispute resolved notification
 */
export function notifyDisputeResolved(orderNo: string): void {
  playAlertSound('success');

  showNotification('✅ 纠纷已解决', {
    body: `订单 ${orderNo} 的纠纷已成功解决`,
    tag: `dispute-resolved-${orderNo}`,
  });
}

// Storage key for notification preferences
const NOTIFICATION_PREFS_KEY = 'gamelink_notification_prefs';

interface NotificationPreferences {
  soundEnabled: boolean;
  browserNotificationEnabled: boolean;
  slaAlertEnabled: boolean;
}

const defaultPrefs: NotificationPreferences = {
  soundEnabled: true,
  browserNotificationEnabled: true,
  slaAlertEnabled: true,
};

/**
 * Get notification preferences from localStorage
 */
export function getNotificationPreferences(): NotificationPreferences {
  try {
    const stored = localStorage.getItem(NOTIFICATION_PREFS_KEY);
    if (stored) {
      return { ...defaultPrefs, ...JSON.parse(stored) };
    }
  } catch {
    // Ignore parse errors
  }
  return defaultPrefs;
}

/**
 * Save notification preferences to localStorage
 */
export function saveNotificationPreferences(prefs: Partial<NotificationPreferences>): void {
  try {
    const current = getNotificationPreferences();
    const updated = { ...current, ...prefs };
    localStorage.setItem(NOTIFICATION_PREFS_KEY, JSON.stringify(updated));
  } catch {
    // Ignore storage errors
  }
}
