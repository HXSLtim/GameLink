import React, { ReactNode, useEffect } from 'react';
import { createPortal } from 'react-dom';
import { createRoot } from 'react-dom/client';
import styles from './Notification.module.less';

export type NoticeType = 'info' | 'success' | 'warning' | 'error';

export interface NotificationProps {
  type?: NoticeType;
  title?: React.ReactNode;
  description?: React.ReactNode;
  duration?: number; // ms
  onClose?: () => void;
}

const ensureContainer = () => {
  let container = document.getElementById('gl-notification-container');
  if (!container) {
    container = document.createElement('div');
    container.id = 'gl-notification-container';
    container.className = styles.container;
    document.body.appendChild(container);
  }
  return container;
};

export const Notification: React.FC<NotificationProps> = ({
  type = 'info',
  title,
  description,
  duration = 3000,
  onClose,
}) => {
  useEffect(() => {
    const timer = setTimeout(() => onClose?.(), duration);
    return () => clearTimeout(timer);
  }, [duration, onClose]);

  const node = (
    <div role="alert" className={`${styles.notice} ${styles[type]}`}>
      {title && <div className={styles.title}>{title}</div>}
      {description && <div className={styles.desc}>{description}</div>}
      <button type="button" className={styles.close} onClick={onClose} aria-label="关闭">
        ×
      </button>
    </div>
  );

  return createPortal(node, ensureContainer());
};

export interface NotifyOptions extends NotificationProps {}

export const notify = (options: NotifyOptions) => {
  const container = ensureContainer();
  const wrapper = document.createElement('div');
  container.appendChild(wrapper);

  const root = createRoot(wrapper);
  const close = () => {
    root.unmount();
    if (wrapper.parentNode) wrapper.parentNode.removeChild(wrapper);
    options.onClose?.();
  };

  const element = <Notification {...options} onClose={close} />;
  root.render(element);

  return { close };
};

// Message (toast)
export type MessageType = NoticeType | 'loading';

export interface MessageOptions {
  type?: MessageType;
  content: React.ReactNode;
  duration?: number;
}

export const Message: React.FC<MessageOptions & { onClose?: () => void }> = ({
  type = 'info',
  content,
  duration = 2000,
  onClose,
}) => {
  useEffect(() => {
    const timer = setTimeout(() => onClose?.(), duration);
    return () => clearTimeout(timer);
  }, [duration, onClose]);

  const containerId = 'gl-message-container';
  let container = document.getElementById(containerId);
  if (!container) {
    container = document.createElement('div');
    container.id = containerId;
    container.className = styles.messageContainer;
    document.body.appendChild(container);
  }

  const node = (
    <div role="status" aria-live="polite" className={`${styles.message} ${styles[type]}`}>
      {content}
    </div>
  );

  return createPortal(node, container);
};

const baseMessage = (options: MessageOptions) => {
  const container = document.getElementById('gl-message-container') || document.body;
  const wrapper = document.createElement('div');
  (container as HTMLElement).appendChild(wrapper);

  const root = createRoot(wrapper);
  const close = () => {
    root.unmount();
    if (wrapper.parentNode) wrapper.parentNode.removeChild(wrapper);
  };

  const element = <Message {...options} onClose={close} />;
  root.render(element);

  return { close };
};

type MessageShortcut = (content: ReactNode, duration?: number) => { close: () => void };

export const message = Object.assign(baseMessage, {
  success: ((content: ReactNode, duration?: number) => baseMessage({ type: 'success', content, duration })) as MessageShortcut,
  error: ((content: ReactNode, duration?: number) => baseMessage({ type: 'error', content, duration })) as MessageShortcut,
  info: ((content: ReactNode, duration?: number) => baseMessage({ type: 'info', content, duration })) as MessageShortcut,
  warning: ((content: ReactNode, duration?: number) => baseMessage({ type: 'warning', content, duration })) as MessageShortcut,
});