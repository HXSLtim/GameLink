import { forwardRef } from 'react';
import { cn } from '@/lib/utils';
import { Skeleton } from '@/components/ui/skeleton';
import { ResponsiveGrid } from '@/components/layout';

export interface LoadingGridProps extends React.HTMLAttributes<HTMLDivElement> {
  /** 加载项数量 */
  count?: number;
  /** 布局类型 */
  layout?: 'card' | 'player' | 'game' | 'service' | 'gift' | 'list';
  /** 卡片高度 */
  cardHeight?: number | string;
  /** 是否显示图片骨架 */
  showImage?: boolean;
  /** 图片宽高比 */
  imageAspect?: 'square' | 'video' | 'portrait';
}

const aspectClasses = {
  square: 'aspect-square',
  video: 'aspect-video',
  portrait: 'aspect-[3/4]',
};

// 不同类型的骨架卡片
function CardSkeleton({ showImage = true, imageAspect = 'video', cardHeight }: {
  showImage?: boolean;
  imageAspect?: 'square' | 'video' | 'portrait';
  cardHeight?: number | string;
}) {
  return (
    <div 
      className="rounded-xl border border-border/50 bg-card overflow-hidden"
      style={cardHeight ? { height: typeof cardHeight === 'number' ? `${cardHeight}px` : cardHeight } : undefined}
    >
      {showImage && (
        <Skeleton className={cn('w-full', aspectClasses[imageAspect])} />
      )}
      <div className="p-4 space-y-3">
        <Skeleton className="h-4 w-3/4" />
        <Skeleton className="h-3 w-1/2" />
        <div className="flex items-center justify-between pt-2">
          <Skeleton className="h-5 w-16" />
          <Skeleton className="h-8 w-20 rounded-lg" />
        </div>
      </div>
    </div>
  );
}

function PlayerCardSkeleton() {
  return (
    <div className="rounded-xl border border-border/50 bg-card p-4">
      <div className="flex items-start gap-3">
        <Skeleton className="w-16 h-16 rounded-full shrink-0" />
        <div className="flex-1 space-y-2">
          <Skeleton className="h-4 w-24" />
          <Skeleton className="h-3 w-32" />
          <div className="flex items-center gap-2">
            <Skeleton className="h-5 w-12" />
            <Skeleton className="h-5 w-16" />
          </div>
        </div>
      </div>
      <div className="mt-3 pt-3 border-t border-border/30">
        <div className="flex items-center justify-between">
          <Skeleton className="h-5 w-20" />
          <Skeleton className="h-8 w-16 rounded-lg" />
        </div>
      </div>
    </div>
  );
}

function GameCardSkeleton() {
  return (
    <div className="flex flex-col items-center gap-2">
      <Skeleton className="w-16 h-16 rounded-2xl" />
      <Skeleton className="h-3 w-12" />
    </div>
  );
}

function ServiceCardSkeleton() {
  return (
    <div className="rounded-xl border border-border/50 bg-card overflow-hidden">
      <Skeleton className="w-full aspect-video" />
      <div className="p-3 space-y-2">
        <Skeleton className="h-4 w-3/4" />
        <div className="flex items-center gap-2">
          <Skeleton className="w-6 h-6 rounded-full" />
          <Skeleton className="h-3 w-16" />
        </div>
        <div className="flex items-center justify-between pt-1">
          <Skeleton className="h-5 w-14" />
          <Skeleton className="h-6 w-12 rounded" />
        </div>
      </div>
    </div>
  );
}

function GiftCardSkeleton() {
  return (
    <div className="rounded-xl border border-border/50 bg-card p-3 text-center">
      <Skeleton className="w-12 h-12 rounded-xl mx-auto" />
      <Skeleton className="h-3 w-16 mx-auto mt-2" />
      <Skeleton className="h-4 w-12 mx-auto mt-1" />
    </div>
  );
}

function ListItemSkeleton() {
  return (
    <div className="flex items-center gap-4 p-4 border-b border-border/30 last:border-0">
      <Skeleton className="w-12 h-12 rounded-lg shrink-0" />
      <div className="flex-1 space-y-2">
        <Skeleton className="h-4 w-1/3" />
        <Skeleton className="h-3 w-2/3" />
      </div>
      <Skeleton className="h-8 w-16 rounded-lg" />
    </div>
  );
}

export const LoadingGrid = forwardRef<HTMLDivElement, LoadingGridProps>(({
  count = 6,
  layout = 'card',
  cardHeight,
  showImage = true,
  imageAspect = 'video',
  className,
  ...props
}, ref) => {
  // 列表布局不使用 Grid
  if (layout === 'list') {
    return (
      <div ref={ref} className={cn('divide-y divide-border/30', className)} {...props}>
        {Array.from({ length: count }).map((_, index) => (
          <ListItemSkeleton key={index} />
        ))}
      </div>
    );
  }

  // 选择骨架组件
  const SkeletonComponent = {
    card: () => <CardSkeleton showImage={showImage} imageAspect={imageAspect} cardHeight={cardHeight} />,
    player: PlayerCardSkeleton,
    game: GameCardSkeleton,
    service: ServiceCardSkeleton,
    gift: GiftCardSkeleton,
  }[layout];

  // 映射 layout 到 ResponsiveGrid preset
  const presetMap: Record<string, 'cards' | 'players' | 'games' | 'services' | 'gifts'> = {
    card: 'cards',
    player: 'players',
    game: 'games',
    service: 'services',
    gift: 'gifts',
  };

  return (
    <ResponsiveGrid
      ref={ref}
      preset={presetMap[layout] || 'cards'}
      className={className}
      {...props}
    >
      {Array.from({ length: count }).map((_, index) => (
        <SkeletonComponent key={index} />
      ))}
    </ResponsiveGrid>
  );
});

LoadingGrid.displayName = 'LoadingGrid';
