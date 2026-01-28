import { forwardRef, Children, isValidElement } from 'react';
import { cn } from '@/lib/utils';
import { Divider } from './container';

type StackProps = React.HTMLAttributes<HTMLDivElement> & {
    /** 间距 */
    spacing?: 0 | 1 | 2 | 3 | 4 | 5 | 6 | 8;
    /** 方向 */
    direction?: 'horizontal' | 'vertical';
    /** 是否显示分隔线 */
    divider?: boolean;
    /** 对齐 */
    align?: 'start' | 'center' | 'end' | 'stretch';
};

const spacingMap = {
    0: 'gap-0', 1: 'gap-1', 2: 'gap-2', 3: 'gap-3',
    4: 'gap-4', 5: 'gap-5', 6: 'gap-6', 8: 'gap-8',
};

const alignMap = {
    start: 'items-start',
    center: 'items-center',
    end: 'items-end',
    stretch: 'items-stretch',
};

/** 堆叠容器 - 支持分隔线 */
export const Stack = forwardRef<HTMLDivElement, StackProps>(({
    spacing = 4,
    direction = 'vertical',
    divider = false,
    align = 'stretch',
    className,
    children,
    ...props
}, ref) => {
    const childArray = Children.toArray(children).filter(isValidElement);

    return (
        <div
            ref={ref}
            className={cn(
                'flex',
                direction === 'vertical' ? 'flex-col' : 'flex-row',
                !divider && spacingMap[spacing],
                alignMap[align],
                className
            )}
            {...props}
        >
            {divider
                ? childArray.map((child, index) => (
                    <div key={index} className={cn(direction === 'vertical' ? 'w-full' : '')}>
                        {child}
                        {index < childArray.length - 1 && (
                            <Divider 
                                orientation={direction === 'vertical' ? 'horizontal' : 'vertical'}
                                className={cn(
                                    direction === 'vertical' 
                                        ? `my-${spacing}` 
                                        : `mx-${spacing} h-auto self-stretch`
                                )}
                            />
                        )}
                    </div>
                ))
                : children
            }
        </div>
    );
});
Stack.displayName = 'Stack';

/** 重叠堆叠 - 用于头像组等 */
export const OverlapStack = forwardRef<HTMLDivElement, Omit<StackProps, 'divider'> & {
    /** 重叠偏移量 */
    offset?: number;
    /** 最大显示数量 */
    max?: number;
    /** 剩余数量显示 */
    showRemainder?: boolean;
}>(({ offset = -8, max, showRemainder = true, direction = 'horizontal', className, children, ...props }, ref) => {
    const childArray = Children.toArray(children).filter(isValidElement);
    const visibleChildren = max ? childArray.slice(0, max) : childArray;
    const remainderCount = max ? childArray.length - max : 0;

    return (
        <div
            ref={ref}
            className={cn(
                'flex',
                direction === 'vertical' ? 'flex-col' : 'flex-row',
                className
            )}
            {...props}
        >
            {visibleChildren.map((child, index) => (
                <div
                    key={index}
                    style={{
                        marginLeft: direction === 'horizontal' && index > 0 ? offset : undefined,
                        marginTop: direction === 'vertical' && index > 0 ? offset : undefined,
                        zIndex: visibleChildren.length - index,
                    }}
                    className="relative"
                >
                    {child}
                </div>
            ))}
            {showRemainder && remainderCount > 0 && (
                <div
                    style={{
                        marginLeft: direction === 'horizontal' ? offset : undefined,
                        marginTop: direction === 'vertical' ? offset : undefined,
                    }}
                    className="relative flex items-center justify-center w-8 h-8 rounded-full bg-muted text-xs font-medium text-muted-foreground border-2 border-background"
                >
                    +{remainderCount}
                </div>
            )}
        </div>
    );
});
OverlapStack.displayName = 'OverlapStack';
