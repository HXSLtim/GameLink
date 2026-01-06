import { useState, useEffect, useRef, useCallback } from 'react';

/**
 * 检测容器是否准备好渲染（有有效尺寸）
 * 用于替代 setTimeout hack，更可靠地检测图表容器
 *
 * @param options 配置选项
 * @returns { ref, isReady, dimensions }
 */
export function useContainerReady<T extends HTMLElement = HTMLDivElement>(options?: {
  /** 最小宽度阈值，默认 1 */
  minWidth?: number;
  /** 最小高度阈值，默认 1 */
  minHeight?: number;
}) {
  const { minWidth = 1, minHeight = 1 } = options || {};
  const ref = useRef<T>(null);
  const [isReady, setIsReady] = useState(false);
  const [dimensions, setDimensions] = useState({ width: 0, height: 0 });

  const checkDimensions = useCallback(() => {
    if (ref.current) {
      const { width, height } = ref.current.getBoundingClientRect();
      if (width >= minWidth && height >= minHeight) {
        setDimensions({ width, height });
        setIsReady(true);
        return true;
      }
    }
    return false;
  }, [minWidth, minHeight]);

  useEffect(() => {
    // 首次检查
    if (checkDimensions()) return;

    // 使用 ResizeObserver 监听尺寸变化
    const observer = new ResizeObserver((entries) => {
      for (const entry of entries) {
        const { width, height } = entry.contentRect;
        if (width >= minWidth && height >= minHeight) {
          setDimensions({ width, height });
          setIsReady(true);
        }
      }
    });

    if (ref.current) {
      observer.observe(ref.current);
    }

    // 备用：使用 requestAnimationFrame 轮询（兼容性）
    let rafId: number;
    const pollDimensions = () => {
      if (!checkDimensions()) {
        rafId = requestAnimationFrame(pollDimensions);
      }
    };
    rafId = requestAnimationFrame(pollDimensions);

    return () => {
      observer.disconnect();
      cancelAnimationFrame(rafId);
    };
  }, [checkDimensions, minWidth, minHeight]);

  return { ref, isReady, dimensions };
}

/**
 * 响应式图表高度 hook
 * 根据屏幕宽度返回合适的图表高度
 */
export function useResponsiveChartHeight(options?: {
  /** 移动端高度，默认 200 */
  mobile?: number;
  /** 平板高度，默认 250 */
  tablet?: number;
  /** 桌面高度，默认 300 */
  desktop?: number;
  /** 移动端断点，默认 576 */
  mobileBreakpoint?: number;
  /** 平板断点，默认 992 */
  tabletBreakpoint?: number;
}) {
  const {
    mobile = 200,
    tablet = 250,
    desktop = 300,
    mobileBreakpoint = 576,
    tabletBreakpoint = 992,
  } = options || {};

  const [height, setHeight] = useState(desktop);

  useEffect(() => {
    const updateHeight = () => {
      const width = window.innerWidth;
      if (width < mobileBreakpoint) {
        setHeight(mobile);
      } else if (width < tabletBreakpoint) {
        setHeight(tablet);
      } else {
        setHeight(desktop);
      }
    };

    updateHeight();
    window.addEventListener('resize', updateHeight);
    return () => window.removeEventListener('resize', updateHeight);
  }, [mobile, tablet, desktop, mobileBreakpoint, tabletBreakpoint]);

  return height;
}

export default useContainerReady;
