import { forwardRef } from 'react';
import { cn } from '@/lib/utils';
import { Card, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { ImageBox, HStack, VStack } from '@/components/layout';
import { PriceTag, RatingDisplay } from '@/components/common';
import { Users, Clock, ShoppingCart } from 'lucide-react';

export interface ServiceCardProps extends React.HTMLAttributes<HTMLDivElement> {
  /** 服务 ID */
  serviceId: number;
  /** 服务名称 */
  name: string;
  /** 描述 */
  description?: string;
  /** 封面图片 */
  coverUrl?: string;
  /** 价格（分） */
  priceCents: number;
  /** 单位 */
  unit?: string;
  /** 评分 */
  rating?: number;
  /** 评价数量 */
  reviewCount?: number;
  /** 销量 */
  salesCount?: number;
  /** 所需人数 */
  requiredPlayers?: number;
  /** 时长（分钟） */
  durationMinutes?: number;
  /** 分类 */
  category?: string;
  /** 关联陪玩师 */
  player?: {
    id: number;
    nickname: string;
    avatarUrl?: string;
  };
  /** 变体 */
  variant?: 'default' | 'compact' | 'horizontal';
  /** 是否显示购买按钮 */
  showBuyButton?: boolean;
  /** 购买按钮文本 */
  buyButtonText?: string;
  /** 购买回调 */
  onBuy?: (serviceId: number) => void;
  /** 点击回调 */
  onSelect?: (serviceId: number) => void;
}

export const ServiceCard = forwardRef<HTMLDivElement, ServiceCardProps>(({
  serviceId,
  name,
  description,
  coverUrl,
  priceCents,
  unit = '/次',
  rating,
  reviewCount,
  salesCount,
  requiredPlayers,
  durationMinutes,
  category,
  player,
  variant = 'default',
  showBuyButton = true,
  buyButtonText = '购买',
  onBuy,
  onSelect,
  className,
  ...props
}, ref) => {
  const handleClick = () => {
    onSelect?.(serviceId);
  };

  const handleBuy = (e: React.MouseEvent) => {
    e.stopPropagation();
    onBuy?.(serviceId);
  };

  // 横向布局变体
  if (variant === 'horizontal') {
    return (
      <Card
        ref={ref}
        className={cn(
          'cursor-pointer overflow-hidden',
          'hover:shadow-md hover:border-primary/30',
          'transition-all duration-200',
          className
        )}
        onClick={handleClick}
        {...props}
      >
        <CardContent className="p-0">
          <div className="flex gap-4">
            {/* 封面图 */}
            <ImageBox
              src={coverUrl}
              alt={name}
              className="w-32 h-24 shrink-0 rounded-l-lg"
              fallback={name.charAt(0)}
            />
            
            {/* 内容 */}
            <div className="flex-1 py-3 pr-4">
              <VStack spacing={2}>
                <div>
                  <h4 className="font-medium text-sm line-clamp-1">{name}</h4>
                  {description && (
                    <p className="text-xs text-muted-foreground line-clamp-1 mt-0.5">
                      {description}
                    </p>
                  )}
                </div>
                
                <HStack spacing={3} align="center" className="text-xs text-muted-foreground">
                  {rating !== undefined && (
                    <RatingDisplay rating={rating} count={reviewCount} size="xs" compact />
                  )}
                  {salesCount !== undefined && (
                    <span>已售 {salesCount}</span>
                  )}
                </HStack>
                
                <HStack justify="between" align="center">
                  <PriceTag priceCents={priceCents} unit={unit} size="sm" />
                  {showBuyButton && (
                    <Button size="sm" onClick={handleBuy}>
                      {buyButtonText}
                    </Button>
                  )}
                </HStack>
              </VStack>
            </div>
          </div>
        </CardContent>
      </Card>
    );
  }

  // 紧凑变体
  if (variant === 'compact') {
    return (
      <Card
        ref={ref}
        className={cn(
          'cursor-pointer overflow-hidden',
          'hover:shadow-md hover:border-primary/30',
          'transition-all duration-200',
          className
        )}
        onClick={handleClick}
        {...props}
      >
        <CardContent className="p-3">
          <HStack spacing={3} align="center">
            <ImageBox
              src={coverUrl}
              alt={name}
              className="w-12 h-12 shrink-0 rounded-lg"
              fallback={name.charAt(0)}
            />
            <div className="flex-1 min-w-0">
              <h4 className="font-medium text-sm truncate">{name}</h4>
              <PriceTag priceCents={priceCents} unit={unit} size="xs" className="mt-1" />
            </div>
            {showBuyButton && (
              <Button size="sm" variant="outline" onClick={handleBuy}>
                <ShoppingCart className="w-4 h-4" />
              </Button>
            )}
          </HStack>
        </CardContent>
      </Card>
    );
  }

  // 默认卡片变体
  return (
    <Card
      ref={ref}
      className={cn(
        'cursor-pointer overflow-hidden group',
        'hover:shadow-lg hover:border-primary/30',
        'transition-all duration-300',
        className
      )}
      onClick={handleClick}
      {...props}
    >
      {/* 封面图 */}
      <div className="relative overflow-hidden aspect-video">
        <ImageBox
          src={coverUrl}
          alt={name}
          className="w-full h-full group-hover:scale-105 transition-transform duration-500"
          fallback={name.charAt(0)}
        />
        {/* 分类标签 */}
        {category && (
          <Badge
            variant="secondary"
            className="absolute top-2 left-2 bg-black/50 text-white backdrop-blur-sm"
          >
            {category}
          </Badge>
        )}
        {/* 人数/时长标签 */}
        <div className="absolute bottom-2 right-2 flex gap-1.5">
          {requiredPlayers && requiredPlayers > 1 && (
            <Badge variant="secondary" className="bg-black/50 text-white backdrop-blur-sm">
              <Users className="w-3 h-3 mr-1" />
              {requiredPlayers}人
            </Badge>
          )}
          {durationMinutes && (
            <Badge variant="secondary" className="bg-black/50 text-white backdrop-blur-sm">
              <Clock className="w-3 h-3 mr-1" />
              {durationMinutes}分钟
            </Badge>
          )}
        </div>
      </div>
      
      <CardContent className="p-4">
        <VStack spacing={3}>
          {/* 标题和描述 */}
          <div>
            <h4 className="font-medium text-sm line-clamp-1">{name}</h4>
            {description && (
              <p className="text-xs text-muted-foreground line-clamp-2 mt-1">
                {description}
              </p>
            )}
          </div>
          
          {/* 关联陪玩师 */}
          {player && (
            <HStack spacing={2} align="center">
              <Avatar className="w-6 h-6">
                <AvatarImage src={player.avatarUrl} alt={player.nickname} />
                <AvatarFallback className="text-[10px]">
                  {player.nickname.charAt(0)}
                </AvatarFallback>
              </Avatar>
              <span className="text-xs text-muted-foreground truncate">
                {player.nickname}
              </span>
            </HStack>
          )}
          
          {/* 评分和销量 */}
          <HStack spacing={3} align="center" className="text-xs text-muted-foreground">
            {rating !== undefined && (
              <RatingDisplay rating={rating} count={reviewCount} size="xs" compact />
            )}
            {salesCount !== undefined && (
              <span>已售 {salesCount}</span>
            )}
          </HStack>
          
          {/* 价格和按钮 */}
          <HStack justify="between" align="center" className="pt-1">
            <PriceTag priceCents={priceCents} unit={unit} size="md" />
            {showBuyButton && (
              <Button size="sm" onClick={handleBuy}>
                {buyButtonText}
              </Button>
            )}
          </HStack>
        </VStack>
      </CardContent>
    </Card>
  );
});

ServiceCard.displayName = 'ServiceCard';
