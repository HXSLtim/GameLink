import { forwardRef } from 'react';
import { cn } from '@/lib/utils';

type ContainerProps = React.HTMLAttributes<HTMLDivElement> & {
    /** 最大宽度 */
    maxWidth?: 'sm' | 'md' | 'lg' | 'xl' | '2xl' | 'full';
    /** 是否居中 */
    centered?: boolean;
    /** 水平内边距 */
    px?: 0 | 2 | 4 | 6 | 8;
};

const maxWidthMap = {
    sm: 'max-w-screen-sm',
    md: 'max-w-screen-md',
    lg: 'max-w-screen-lg',
    xl: 'max-w-screen-xl',
    '2xl': 'max-w-screen-2xl',
    full: 'max-w-full',
};

const pxMap = {
    0: 'px-0',
    2: 'px-2',
    4: 'px-4',
    6: 'px-6',
    8: 'px-8',
};

/** 响应式容器 */
export const Container = forwardRef<HTMLDivElement, ContainerProps>(({
    maxWidth = '2xl',
    centered = true,
    px = 4,
    className,
    children,
    ...props
}, ref) => {
    return (
        <div
            ref={ref}
            className={cn(
                'w-full',
                maxWidthMap[maxWidth],
                centered && 'mx-auto',
                pxMap[px],
                className
            )}
            {...props}
        >
            {children}
        </div>
    );
});
Container.displayName = 'Container';

type SectionProps = React.HTMLAttributes<HTMLElement> & {
    /** 垂直间距 */
    spacing?: 'none' | 'sm' | 'md' | 'lg' | 'xl';
    /** 是否为 section 元素 */
    as?: 'section' | 'div' | 'article';
};

const spacingMap = {
    none: '',
    sm: 'space-y-2',
    md: 'space-y-4',
    lg: 'space-y-6',
    xl: 'space-y-8',
};

/** 内容区块 */
export const Section = forwardRef<HTMLElement, SectionProps>(({
    spacing = 'md',
    as: Component = 'section',
    className,
    children,
    ...props
}, ref) => {
    return (
        <Component
            ref={ref as React.Ref<HTMLElement>}
            className={cn(spacingMap[spacing], className)}
            {...props}
        >
            {children}
        </Component>
    );
});
Section.displayName = 'Section';

type DividerProps = React.HTMLAttributes<HTMLHRElement> & {
    /** 方向 */
    orientation?: 'horizontal' | 'vertical';
    /** 样式变体 */
    variant?: 'solid' | 'dashed' | 'dotted' | 'gradient';
    /** 标签文本 */
    label?: string;
};

/** 分隔线 */
export const Divider = forwardRef<HTMLHRElement, DividerProps>(({
    orientation = 'horizontal',
    variant = 'solid',
    label,
    className,
    ...props
}, ref) => {
    const variantStyles = {
        solid: 'border-border',
        dashed: 'border-dashed border-border',
        dotted: 'border-dotted border-border',
        gradient: 'border-0 bg-gradient-to-r from-transparent via-border to-transparent',
    };

    if (label) {
        return (
            <div className={cn('flex items-center gap-4', className)}>
                <div className={cn(
                    'flex-1 h-px',
                    variant === 'gradient' 
                        ? 'bg-gradient-to-r from-transparent via-border to-border/50'
                        : 'bg-border'
                )} />
                <span className="text-xs text-muted-foreground font-medium px-2">
                    {label}
                </span>
                <div className={cn(
                    'flex-1 h-px',
                    variant === 'gradient'
                        ? 'bg-gradient-to-r from-border/50 via-border to-transparent'
                        : 'bg-border'
                )} />
            </div>
        );
    }

    if (orientation === 'vertical') {
        return (
            <hr
                ref={ref}
                className={cn(
                    'w-px h-full border-0',
                    variant === 'gradient' 
                        ? 'bg-gradient-to-b from-transparent via-border to-transparent'
                        : 'bg-border',
                    className
                )}
                {...props}
            />
        );
    }

    return (
        <hr
            ref={ref}
            className={cn(
                'w-full border-t',
                variantStyles[variant],
                variant === 'gradient' && 'h-px',
                className
            )}
            {...props}
        />
    );
});
Divider.displayName = 'Divider';
