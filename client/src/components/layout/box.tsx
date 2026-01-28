import { forwardRef } from 'react';
import { cn } from '@/lib/utils';

type BoxProps = React.HTMLAttributes<HTMLDivElement> & {
    /** 内边距 */
    p?: 0 | 1 | 2 | 3 | 4 | 5 | 6 | 8 | 10 | 12;
    px?: 0 | 1 | 2 | 3 | 4 | 5 | 6 | 8;
    py?: 0 | 1 | 2 | 3 | 4 | 5 | 6 | 8;
    /** 外边距 */
    m?: 0 | 1 | 2 | 3 | 4 | 5 | 6 | 8 | 'auto';
    mx?: 0 | 1 | 2 | 3 | 4 | 5 | 6 | 8 | 'auto';
    my?: 0 | 1 | 2 | 3 | 4 | 5 | 6 | 8 | 'auto';
    /** 圆角 */
    rounded?: 'none' | 'sm' | 'md' | 'lg' | 'xl' | '2xl' | '3xl' | 'full';
    /** 阴影 */
    shadow?: 'none' | 'sm' | 'md' | 'lg' | 'xl' | '2xl';
    /** 边框 */
    border?: boolean;
    /** 背景 */
    bg?: 'transparent' | 'background' | 'muted' | 'card' | 'primary' | 'secondary' | 'accent';
};

const paddingMap = {
    0: 'p-0', 1: 'p-1', 2: 'p-2', 3: 'p-3', 4: 'p-4', 
    5: 'p-5', 6: 'p-6', 8: 'p-8', 10: 'p-10', 12: 'p-12'
};
const pxMap = { 0: 'px-0', 1: 'px-1', 2: 'px-2', 3: 'px-3', 4: 'px-4', 5: 'px-5', 6: 'px-6', 8: 'px-8' };
const pyMap = { 0: 'py-0', 1: 'py-1', 2: 'py-2', 3: 'py-3', 4: 'py-4', 5: 'py-5', 6: 'py-6', 8: 'py-8' };

const marginMap = { 0: 'm-0', 1: 'm-1', 2: 'm-2', 3: 'm-3', 4: 'm-4', 5: 'm-5', 6: 'm-6', 8: 'm-8', auto: 'm-auto' };
const mxMap = { 0: 'mx-0', 1: 'mx-1', 2: 'mx-2', 3: 'mx-3', 4: 'mx-4', 5: 'mx-5', 6: 'mx-6', 8: 'mx-8', auto: 'mx-auto' };
const myMap = { 0: 'my-0', 1: 'my-1', 2: 'my-2', 3: 'my-3', 4: 'my-4', 5: 'my-5', 6: 'my-6', 8: 'my-8', auto: 'my-auto' };

const roundedMap = {
    none: 'rounded-none', sm: 'rounded-sm', md: 'rounded-md', lg: 'rounded-lg',
    xl: 'rounded-xl', '2xl': 'rounded-2xl', '3xl': 'rounded-3xl', full: 'rounded-full'
};

const shadowMap = {
    none: 'shadow-none', sm: 'shadow-sm', md: 'shadow-md', 
    lg: 'shadow-lg', xl: 'shadow-xl', '2xl': 'shadow-2xl'
};

const bgMap = {
    transparent: 'bg-transparent',
    background: 'bg-background',
    muted: 'bg-muted',
    card: 'bg-card',
    primary: 'bg-primary',
    secondary: 'bg-secondary',
    accent: 'bg-accent',
};

/** 通用盒子容器 */
export const Box = forwardRef<HTMLDivElement, BoxProps>(({
    p, px, py,
    m, mx, my,
    rounded,
    shadow,
    border,
    bg,
    className,
    children,
    ...props
}, ref) => {
    return (
        <div
            ref={ref}
            className={cn(
                p !== undefined && paddingMap[p],
                px !== undefined && pxMap[px],
                py !== undefined && pyMap[py],
                m !== undefined && marginMap[m],
                mx !== undefined && mxMap[mx],
                my !== undefined && myMap[my],
                rounded && roundedMap[rounded],
                shadow && shadowMap[shadow],
                border && 'border border-border',
                bg && bgMap[bg],
                className
            )}
            {...props}
        >
            {children}
        </div>
    );
});
Box.displayName = 'Box';

/** 卡片容器 - 预设样式的 Box */
export const CardBox = forwardRef<HTMLDivElement, Omit<BoxProps, 'rounded' | 'shadow' | 'border' | 'bg'> & {
    variant?: 'default' | 'outlined' | 'elevated' | 'ghost';
    hoverable?: boolean;
}>(({ variant = 'default', hoverable = false, className, children, ...props }, ref) => {
    const variantStyles = {
        default: 'bg-card border border-border rounded-xl',
        outlined: 'bg-transparent border border-border rounded-xl',
        elevated: 'bg-card rounded-xl shadow-lg border-0',
        ghost: 'bg-muted/50 rounded-xl border-0',
    };

    return (
        <Box
            ref={ref}
            className={cn(
                variantStyles[variant],
                hoverable && 'transition-all duration-200 hover:shadow-md hover:border-primary/30 cursor-pointer',
                className
            )}
            {...props}
        >
            {children}
        </Box>
    );
});
CardBox.displayName = 'CardBox';

/** 玻璃拟态容器 */
export const GlassBox = forwardRef<HTMLDivElement, BoxProps & {
    blur?: 'sm' | 'md' | 'lg' | 'xl';
}>(({ blur = 'md', className, children, ...props }, ref) => {
    const blurMap = {
        sm: 'backdrop-blur-sm',
        md: 'backdrop-blur-md',
        lg: 'backdrop-blur-lg',
        xl: 'backdrop-blur-xl',
    };

    return (
        <Box
            ref={ref}
            className={cn(
                'bg-white/5 border border-white/10',
                blurMap[blur],
                className
            )}
            {...props}
        >
            {children}
        </Box>
    );
});
GlassBox.displayName = 'GlassBox';
