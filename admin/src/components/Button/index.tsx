/**
 * 统一按钮组件
 * 包装 Ant Design Button，提供标准化的尺寸和样式变体
 *
 * 使用方式:
 * import { Button } from '@/components';
 * <Button btnVariant="primary" btnSize="md">提交</Button>
 */
import React from 'react';
import { Button as AntButton } from 'antd';
import type { ButtonProps as AntButtonProps } from 'antd';

/**
 * 按钮尺寸映射
 */
const sizeMap = {
  xs: 'small',
  sm: 'small',
  md: 'middle',
  lg: 'large',
  xl: 'large',
} as const;

/**
 * 按钮变体配置
 * 映射到 Ant Design 6 的 color + variant 组合
 */
const btnVariantConfig = {
  primary: { color: 'primary' as const, variant: 'solid' as const },
  secondary: { color: 'default' as const, variant: 'outlined' as const },
  dashed: { color: 'default' as const, variant: 'dashed' as const },
  text: { color: 'default' as const, variant: 'text' as const },
  link: { color: 'default' as const, variant: 'link' as const },
  danger: { color: 'danger' as const, variant: 'solid' as const },
  ghost: { color: 'default' as const, variant: 'outlined' as const },
} as const;

export type ButtonSize = keyof typeof sizeMap;
export type BtnVariant = keyof typeof btnVariantConfig;

export interface ButtonProps extends Omit<AntButtonProps, 'size'> {
  /** 按钮尺寸 */
  btnSize?: ButtonSize;
  /** 按钮变体 */
  btnVariant?: BtnVariant;
}

/**
 * 统一按钮组件
 */
export const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ btnSize = 'md', btnVariant, color, variant, children, ...props }, ref) => {
    // 如果指定了 btnVariant，使用变体映射
    const variantConfig = btnVariant ? btnVariantConfig[btnVariant] : null;

    // 原生属性优先
    const finalColor = color ?? variantConfig?.color;
    const finalVariant = variant ?? variantConfig?.variant;

    return (
      <AntButton
        ref={ref}
        size={sizeMap[btnSize]}
        color={finalColor}
        variant={finalVariant}
        {...props}
      >
        {children}
      </AntButton>
    );
  }
);

Button.displayName = 'Button';

/**
 * 图标按钮 - 仅图标，无文字
 * 自动添加 aria-label 支持
 */
export interface IconButtonProps extends ButtonProps {
  /** 无障碍标签 (必填) */
  'aria-label': string;
}

export const IconButton: React.FC<IconButtonProps> = ({
  'aria-label': ariaLabel,
  children,
  ...props
}) => {
  return (
    <Button {...props} aria-label={ariaLabel}>
      {children}
    </Button>
  );
};

export default Button;
