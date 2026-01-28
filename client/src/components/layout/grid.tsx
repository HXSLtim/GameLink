import { forwardRef } from 'react';
import { cn } from '@/lib/utils';

type GridProps = React.HTMLAttributes<HTMLDivElement> & {
    /** 列数 */
    cols?: 1 | 2 | 3 | 4 | 5 | 6 | 12;
    /** 响应式列数 */
    sm?: 1 | 2 | 3 | 4 | 5 | 6;
    md?: 1 | 2 | 3 | 4 | 5 | 6;
    lg?: 1 | 2 | 3 | 4 | 5 | 6;
    xl?: 1 | 2 | 3 | 4 | 5 | 6;
    /** 间距 */
    gap?: 0 | 1 | 2 | 3 | 4 | 5 | 6 | 8;
};

const colsMap = {
    1: 'grid-cols-1',
    2: 'grid-cols-2',
    3: 'grid-cols-3',
    4: 'grid-cols-4',
    5: 'grid-cols-5',
    6: 'grid-cols-6',
    12: 'grid-cols-12',
};

const smColsMap = {
    1: 'sm:grid-cols-1',
    2: 'sm:grid-cols-2',
    3: 'sm:grid-cols-3',
    4: 'sm:grid-cols-4',
    5: 'sm:grid-cols-5',
    6: 'sm:grid-cols-6',
};

const mdColsMap = {
    1: 'md:grid-cols-1',
    2: 'md:grid-cols-2',
    3: 'md:grid-cols-3',
    4: 'md:grid-cols-4',
    5: 'md:grid-cols-5',
    6: 'md:grid-cols-6',
};

const lgColsMap = {
    1: 'lg:grid-cols-1',
    2: 'lg:grid-cols-2',
    3: 'lg:grid-cols-3',
    4: 'lg:grid-cols-4',
    5: 'lg:grid-cols-5',
    6: 'lg:grid-cols-6',
};

const xlColsMap = {
    1: 'xl:grid-cols-1',
    2: 'xl:grid-cols-2',
    3: 'xl:grid-cols-3',
    4: 'xl:grid-cols-4',
    5: 'xl:grid-cols-5',
    6: 'xl:grid-cols-6',
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
};

/** 响应式网格 */
export const Grid = forwardRef<HTMLDivElement, GridProps>(({
    cols = 1,
    sm,
    md,
    lg,
    xl,
    gap = 4,
    className,
    children,
    ...props
}, ref) => {
    return (
        <div
            ref={ref}
            className={cn(
                'grid',
                colsMap[cols],
                sm && smColsMap[sm],
                md && mdColsMap[md],
                lg && lgColsMap[lg],
                xl && xlColsMap[xl],
                gapMap[gap],
                className
            )}
            {...props}
        >
            {children}
        </div>
    );
});
Grid.displayName = 'Grid';

type GridItemProps = React.HTMLAttributes<HTMLDivElement> & {
    /** 跨列数 */
    colSpan?: 1 | 2 | 3 | 4 | 5 | 6 | 12 | 'full';
    /** 跨行数 */
    rowSpan?: 1 | 2 | 3 | 4 | 5 | 6;
};

const colSpanMap = {
    1: 'col-span-1',
    2: 'col-span-2',
    3: 'col-span-3',
    4: 'col-span-4',
    5: 'col-span-5',
    6: 'col-span-6',
    12: 'col-span-12',
    full: 'col-span-full',
};

const rowSpanMap = {
    1: 'row-span-1',
    2: 'row-span-2',
    3: 'row-span-3',
    4: 'row-span-4',
    5: 'row-span-5',
    6: 'row-span-6',
};

/** 网格子项 */
export const GridItem = forwardRef<HTMLDivElement, GridItemProps>(({
    colSpan,
    rowSpan,
    className,
    children,
    ...props
}, ref) => {
    return (
        <div
            ref={ref}
            className={cn(
                colSpan && colSpanMap[colSpan],
                rowSpan && rowSpanMap[rowSpan],
                className
            )}
            {...props}
        >
            {children}
        </div>
    );
});
GridItem.displayName = 'GridItem';

/** 预设响应式网格 */
export const ResponsiveGrid = forwardRef<HTMLDivElement, Omit<GridProps, 'cols' | 'sm' | 'md' | 'lg' | 'xl'> & {
    /** 预设类型 */
    preset?: 'cards' | 'players' | 'games' | 'services' | 'gifts';
}>(({ preset = 'cards', gap = 4, className, children, ...props }, ref) => {
    const presetConfig = {
        cards: { cols: 1 as const, sm: 2 as const, md: 2 as const, lg: 3 as const, xl: 4 as const },
        players: { cols: 1 as const, sm: 2 as const, md: 3 as const, lg: 4 as const, xl: 5 as const },
        games: { cols: 2 as const, sm: 3 as const, md: 4 as const, lg: 6 as const, xl: 6 as const },
        services: { cols: 1 as const, sm: 2 as const, md: 2 as const, lg: 3 as const, xl: 3 as const },
        gifts: { cols: 2 as const, sm: 3 as const, md: 4 as const, lg: 5 as const, xl: 6 as const },
    };

    const config = presetConfig[preset];

    return (
        <Grid
            ref={ref}
            cols={config.cols}
            sm={config.sm}
            md={config.md}
            lg={config.lg}
            xl={config.xl}
            gap={gap}
            className={className}
            {...props}
        >
            {children}
        </Grid>
    );
});
ResponsiveGrid.displayName = 'ResponsiveGrid';
