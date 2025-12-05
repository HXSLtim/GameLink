/**
 * useCountUp Hook
 * 数字滚动动画Hook
 * 实现平滑的数字过渡效果
 */
import { useEffect, useRef, useState } from 'react';

interface UseCountUpOptions {
  /** 动画持续时间（毫秒） */
  duration?: number;
  /** 是否启用动画 */
  enabled?: boolean;
  /** 缓动函数 */
  easing?: (t: number) => number;
}

/**
 * 默认缓动函数 - easeOutExpo
 */
const defaultEasing = (t: number): number => {
  return t === 1 ? 1 : 1 - Math.pow(2, -10 * t);
};

/**
 * useCountUp - 数字滚动动画Hook
 * @param end 目标数值
 * @param options 配置选项
 * @returns 当前动画数值
 */
export const useCountUp = (
  end: number,
  options: UseCountUpOptions = {}
): number => {
  const {
    duration = 1000,
    enabled = true,
    easing = defaultEasing,
  } = options;

  const [count, setCount] = useState(enabled ? 0 : end);
  const startRef = useRef(0);
  const startTimeRef = useRef<number | null>(null);
  const rafRef = useRef<number | null>(null);

  useEffect(() => {
    // 如果动画未启用，直接设置为目标值
    if (!enabled) {
      setCount(end);
      return;
    }

    // 如果目标值没有变化，不需要动画
    if (count === end) {
      return;
    }

    // 记录起始值和起始时间
    startRef.current = count;
    startTimeRef.current = null;

    const animate = (timestamp: number) => {
      // 初始化起始时间
      if (startTimeRef.current === null) {
        startTimeRef.current = timestamp;
      }

      // 计算已经过的时间
      const elapsed = timestamp - startTimeRef.current;
      const progress = Math.min(elapsed / duration, 1);

      // 应用缓动函数
      const easedProgress = easing(progress);

      // 计算当前值
      const currentValue =
        startRef.current + (end - startRef.current) * easedProgress;

      setCount(currentValue);

      // 如果动画未完成，继续下一帧
      if (progress < 1) {
        rafRef.current = requestAnimationFrame(animate);
      }
    };

    // 开始动画
    rafRef.current = requestAnimationFrame(animate);

    // 清理函数
    return () => {
      if (rafRef.current !== null) {
        cancelAnimationFrame(rafRef.current);
      }
    };
  }, [end, duration, enabled, easing]);

  return count;
};

/**
 * 格式化数字显示
 * @param value 数值
 * @param decimals 小数位数
 * @returns 格式化后的字符串
 */
export const formatNumber = (value: number, decimals: number = 0): string => {
  return value.toFixed(decimals);
};

export default useCountUp;
