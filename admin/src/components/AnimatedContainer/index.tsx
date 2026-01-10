/**
 * 动画容器组件
 * 提供常用的入场/退场动画效果
 */
import React from 'react';
import type { ReactNode } from 'react';
import { motion } from 'framer-motion';
import type { Variants, Transition } from 'framer-motion';
import { animationVariants, transitions } from './constants';

// 重新导出常量
export { animationVariants, transitions } from './constants';

export type AnimationType = keyof typeof animationVariants;
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
