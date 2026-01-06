/**
 * 动画容器组件
 * 提供常用的入场/退场动画效果
 */
import React from 'react';
import type { ReactNode } from 'react';
import { motion } from 'framer-motion';
import type { Variants, Transition } from 'framer-motion';

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

export type AnimationType = keyof typeof animationVariants;

/**
 * 预设过渡配置
 */
export const transitions = {
  fast: { duration: 0.15, ease: 'easeOut' },
  normal: { duration: 0.3, ease: 'easeOut' },
  slow: { duration: 0.5, ease: 'easeOut' },
  spring: { type: 'spring', stiffness: 300, damping: 30 },
} as const;

export type TransitionType = keyof typeof transitions;

export interface AnimatedContainerProps {
  children: ReactNode;
  /** 动画类型 */
  animation?: AnimationType;
  /** 过渡类型 */
  transition?: TransitionType;
  /** 延迟时间 (秒) */
  delay?: number;
  /** 自定义变体 */
  variants?: Variants;
  /** 自定义过渡 */
  customTransition?: Transition;
  /** 容器样式 */
  style?: React.CSSProperties;
  /** 容器类名 */
  className?: string;
}

/**
 * 动画容器
 */
export const AnimatedContainer: React.FC<AnimatedContainerProps> = ({
  children,
  animation = 'fadeIn',
  transition = 'normal',
  delay = 0,
  variants,
  customTransition,
  style,
  className,
}) => {
  const selectedVariants = variants || animationVariants[animation];
  const selectedTransition = customTransition || { ...transitions[transition], delay };

  return (
    <motion.div
      initial="initial"
      animate="animate"
      exit="exit"
      variants={selectedVariants}
      transition={selectedTransition}
      style={style}
      className={className}
    >
      {children}
    </motion.div>
  );
};

/**
 * 列表项动画容器
 * 用于列表中的每个项目，支持交错动画
 */
export interface AnimatedListItemProps {
  children: ReactNode;
  /** 项目索引，用于计算延迟 */
  index: number;
  /** 每项延迟间隔 (秒) */
  staggerDelay?: number;
  /** 动画类型 */
  animation?: AnimationType;
  /** 容器样式 */
  style?: React.CSSProperties;
}

export const AnimatedListItem: React.FC<AnimatedListItemProps> = ({
  children,
  index,
  staggerDelay = 0.05,
  animation = 'slideUp',
  style,
}) => {
  return (
    <AnimatedContainer animation={animation} delay={index * staggerDelay} style={style}>
      {children}
    </AnimatedContainer>
  );
};

/**
 * 页面过渡动画容器
 * 用于路由切换时的页面动画
 */
export const PageTransition: React.FC<{ children: ReactNode }> = ({ children }) => {
  return (
    <motion.div
      initial={{ opacity: 0, y: 10 }}
      animate={{ opacity: 1, y: 0 }}
      exit={{ opacity: 0, y: -10 }}
      transition={{ duration: 0.2, ease: 'easeOut' }}
    >
      {children}
    </motion.div>
  );
};

export default AnimatedContainer;
