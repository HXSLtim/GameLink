import { forwardRef } from 'react';
import { cn } from '@/lib/utils';

type FlexProps = React.HTMLAttributes<HTMLDivElement> & {
    /** 主轴方向 */
    direction?: 'row' | 'column' | 'row-reverse' | 'column-reverse';
    /** 主轴对齐 */
    justify?: 'start' | 'end' | 'center' | 'between' | 'around' | 'evenly';
    /** 交叉轴对齐 */
    align?: 'start' | 'end' | 'center' | 'baseline' | 'stretch';
    /** 换行 */
    wrap?: 'nowrap' | 'wrap' | 'wrap-reverse';
    /** 间距 */
    gap?: 0 | 1 | 2 | 3 | 4 | 5 | 6 | 8 | 10 | 12;
    /** 是否内联 */
    inline?: boolean;
};

const justifyMap = {
    start: 'justify-start',
    end: 'justify-end',
    center: 'justify-center',
    between: 'justify-between',
    around: 'justify-around',
    evenly: 'justify-evenly',
};

const alignMap = {
    start: 'items-start',
    end: 'items-end',
    center: 'items-center',
    baseline: 'items-baseline',
    stretch: 'items-stretch',
};

const directionMap = {
    row: 'flex-row',
    column: 'flex-col',
    'row-reverse': 'flex-row-reverse',
    'column-reverse': 'flex-col-reverse',
};

const wrapMap = {
    nowrap: 'flex-nowrap',
    wrap: 'flex-wrap',
    'wrap-reverse': 'flex-wrap-reverse',
};

const gapMap = {
    0: 'gap-0',
    1: 'gap-1',
    2: 'gap-2',
    3: 'gap-3',
    4: 'gap-4',
    5: 'gap-5',
    6: 'gap-6',
    8: 'gap-8',
    10: 'gap-10',
    12: 'gap-12',
};

/** 通用 Flex 容器 */
export const Flex = forwardRef<HTMLDivElement, FlexProps>(({
    direction = 'row',
    justify = 'start',
    align = 'stretch',
    wrap = 'nowrap',
    gap = 0,
    inline = false,
    className,
    children,
    ...props
}, ref) => {
    return (
        <div
            ref={ref}
            className={cn(
                inline ? 'inline-flex' : 'flex',
                directionMap[direction],
                justifyMap[justify],
                alignMap[align],
                wrapMap[wrap],
                gapMap[gap],
                className
            )}
            {...props}
        >
            {children}
        </div>
    );
});
Flex.displayName = 'Flex';

type StackProps = Omit<FlexProps, 'direction'> & {
    /** 间距 */
    spacing?: 0 | 1 | 2 | 3 | 4 | 5 | 6 | 8 | 10 | 12;
};

/** 水平堆叠 */
export const HStack = forwardRef<HTMLDivElement, StackProps>(({
    spacing = 4,
    align = 'center',
    className,
    ...props
}, ref) => {
    return (
        <Flex
            ref={ref}
            direction="row"
            align={align}
            gap={spacing}
            className={className}
            {...props}
        />
    );
});
HStack.displayName = 'HStack';

/** 垂直堆叠 */
export const VStack = forwardRef<HTMLDivElement, StackProps>(({
    spacing = 4,
    align = 'stretch',
    className,
    ...props
}, ref) => {
    return (
        <Flex
            ref={ref}
            direction="column"
            align={align}
            gap={spacing}
            className={className}
            {...props}
        />
    );
});
VStack.displayName = 'VStack';

/** 居中容器 */
export const Center = forwardRef<HTMLDivElement, React.HTMLAttributes<HTMLDivElement>>(({
    className,
    children,
    ...props
}, ref) => {
    return (
        <div
            ref={ref}
            className={cn('flex items-center justify-center', className)}
            {...props}
        >
            {children}
        </div>
    );
});
Center.displayName = 'Center';

/** 弹性占位符 */
export const Spacer = forwardRef<HTMLDivElement, React.HTMLAttributes<HTMLDivElement>>(({
    className,
    ...props
}, ref) => {
    return (
        <div
            ref={ref}
            className={cn('flex-1', className)}
            {...props}
        />
    );
});
Spacer.displayName = 'Spacer';
