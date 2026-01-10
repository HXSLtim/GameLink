/**
 * 动画常量定义
 */

/**
 * 预设动画变体
 */
export const animationVariants = {
  /** 淡入 */
  fadeIn: {
    initial: { opacity: 0 },
    animate: { opacity: 1 },
    exit: { opacity: 0 },
  },
  /** 从下方滑入 */
  slideUp: {
    initial: { opacity: 0, y: 20 },
    animate: { opacity: 1, y: 0 },
    exit: { opacity: 0, y: 20 },
  },
  /** 从上方滑入 */
  slideDown: {
    initial: { opacity: 0, y: -20 },
    animate: { opacity: 1, y: 0 },
    exit: { opacity: 0, y: -20 },
  },
  /** 从左侧滑入 */
  slideLeft: {
    initial: { opacity: 0, x: -20 },
    animate: { opacity: 1, x: 0 },
    exit: { opacity: 0, x: -20 },
  },
  /** 从右侧滑入 */
  slideRight: {
    initial: { opacity: 0, x: 20 },
    animate: { opacity: 1, x: 0 },
    exit: { opacity: 0, x: 20 },
  },
  /** 缩放 */
  scale: {
    initial: { opacity: 0, scale: 0.95 },
    animate: { opacity: 1, scale: 1 },
    exit: { opacity: 0, scale: 0.95 },
  },
} as const;

/**
 * 预设过渡配置
 */
export const transitions = {
  fast: { duration: 0.15, ease: 'easeOut' },
  normal: { duration: 0.3, ease: 'easeOut' },
  slow: { duration: 0.5, ease: 'easeOut' },
  spring: { type: 'spring', stiffness: 300, damping: 30 },
} as const;
