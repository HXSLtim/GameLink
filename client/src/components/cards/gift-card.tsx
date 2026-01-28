import { forwardRef } from 'react';
import { cn } from '@/lib/utils';
import { Card, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { ImageBox, VStack, HStack } from '@/components/layout';
import { PriceTag } from '@/components/common';
import { Gift, Heart, Sparkles } from 'lucide-react';

export interface GiftCardProps extends React.HTMLAttributes<HTMLDivElement> {
  /** 礼物 ID */
  giftId: number;
  /** 礼物名称 */
  name: string;
  /** 图片 URL */
  imageUrl?: string;
  /** 价格（分） */
  priceCents: number;
  /** 描述 */
  description?: string;
  /** 是否热门 */
  isHot?: boolean;
  /** 是否限时 */
  isLimited?: boolean;
  /** 变体 */
  variant?: 'default' | 'compact' | 'featured';
  /** 尺寸 */
  size?: 'sm' | 'md' | 'lg';
  /** 是否显示赠送按钮 */
  showSendButton?: boolean;
  /** 赠送按钮文本 */
  sendButtonText?: string;
  /** 赠送回调 */
  onSend?: (giftId: number) => void;
  /** 点击回调 */
  onSelect?: (giftId: number) => void;
}

const sizeConfig = {
  sm: {
    card: 'p-2',
    image: 'w-10 h-10',
    title: 'text-xs',
    price: 'xs' as const,
  },
  md: {
    card: 'p-3',
    image: 'w-14 h-14',
    title: 'text-sm',
    price: 'sm' as const,
  },
  lg: {
    card: 'p-4',
    image: 'w-20 h-20',
    title: 'text-base',
    price: 'md' as const,
  },
};

export const GiftCard = forwardRef<HTMLDivElement, GiftCardProps>(({
  giftId,
  name,
  imageUrl,
  priceCents,
  description,
  isHot = false,
  isLimited = false,
  variant = 'default',
  size = 'md',
  showSendButton = true,
  sendButtonText = '赠送',
  onSend,
  onSelect,
  className,
  ...props
}, ref) => {
  const config = sizeConfig[size];

  const handleClick = () => {
    onSelect?.(giftId);
  };

  const handleSend = (e: React.MouseEvent) => {
    e.stopPropagation();
    onSend?.(giftId);
  };

  // 特色变体 - 大卡片，展示更多信息
  if (variant === 'featured') {
    return (
      <Card
        ref={ref}
        className={cn(
          'cursor-pointer overflow-hidden group relative',
          'hover:shadow-xl hover:border-primary/30',
          'transition-all duration-300',
          'bg-gradient-to-br from-primary/5 via-background to-primary/10',
          className
        )}
        onClick={handleClick}
        {...props}
      >
        {/* 装饰 */}
        <div className="absolute top-2 right-2 flex gap-1">
          {isHot && (
            <span className="px-2 py-0.5 rounded-full bg-red-500/90 text-white text-xs font-medium flex items-center gap-1">
              <Heart className="w-3 h-3 fill-current" />
              热门
            </span>
          )}
          {isLimited && (
            <span className="px-2 py-0.5 rounded-full bg-amber-500/90 text-white text-xs font-medium flex items-center gap-1">
              <Sparkles className="w-3 h-3" />
              限时
            </span>
          )}
        </div>
        
        <CardContent className="p-6">
          <VStack spacing={4} align="center" className="text-center">
            <div className="relative">
              <ImageBox
                src={imageUrl}
                alt={name}
                className="w-24 h-24 rounded-2xl group-hover:scale-110 transition-transform duration-500"
                fallback={<Gift className="w-12 h-12 text-primary/50" />}
              />
              {/* 光晕效果 */}
              <div className="absolute inset-0 rounded-2xl bg-gradient-to-t from-primary/20 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300" />
            </div>
            
            <div>
              <h4 className="font-semibold text-base">{name}</h4>
              {description && (
                <p className="text-sm text-muted-foreground mt-1 line-clamp-2">
                  {description}
                </p>
              )}
            </div>
            
            <HStack justify="between" align="center" className="w-full">
              <PriceTag priceCents={priceCents} unit="" size="lg" />
              {showSendButton && (
                <Button onClick={handleSend}>
                  <Gift className="w-4 h-4 mr-2" />
                  {sendButtonText}
                </Button>
              )}
            </HStack>
          </VStack>
        </CardContent>
      </Card>
    );
  }

  // 紧凑变体 - 横向布局
  if (variant === 'compact') {
    return (
      <div
        ref={ref}
        className={cn(
          'flex items-center gap-3 p-2 rounded-lg cursor-pointer',
          'hover:bg-muted/50',
          'transition-all duration-200',
          className
        )}
        onClick={handleClick}
        {...props}
      >
        <ImageBox
          src={imageUrl}
          alt={name}
          className="w-10 h-10 rounded-lg shrink-0"
          fallback={<Gift className="w-5 h-5 text-muted-foreground" />}
        />
        <div className="flex-1 min-w-0">
          <span className="text-sm font-medium truncate block">{name}</span>
          <PriceTag priceCents={priceCents} unit="" size="xs" />
        </div>
        {showSendButton && (
          <Button size="sm" variant="ghost" onClick={handleSend}>
            <Gift className="w-4 h-4" />
          </Button>
        )}
      </div>
    );
  }

  // 默认变体 - 垂直布局小卡片
  return (
    <Card
      ref={ref}
      className={cn(
        'cursor-pointer overflow-hidden group',
        'hover:shadow-md hover:border-primary/30',
        'transition-all duration-200',
        className
      )}
      onClick={handleClick}
      {...props}
    >
      <CardContent className={cn('text-center', config.card)}>
        <VStack spacing={2} align="center">
          {/* 标签 */}
          {(isHot || isLimited) && (
            <div className="absolute top-1 right-1">
              {isHot && <Heart className="w-3 h-3 text-red-500 fill-current" />}
              {isLimited && !isHot && <Sparkles className="w-3 h-3 text-amber-500" />}
            </div>
          )}
          
          {/* 图片 */}
          <ImageBox
            src={imageUrl}
            alt={name}
            className={cn(
              config.image,
              'rounded-xl group-hover:scale-105 transition-transform duration-300'
            )}
            fallback={<Gift className="w-1/2 h-1/2 text-muted-foreground" />}
          />
          
          {/* 名称 */}
          <span className={cn('font-medium truncate w-full', config.title)}>
            {name}
          </span>
          
          {/* 价格 */}
          <PriceTag priceCents={priceCents} unit="" size={config.price} />
          
          {/* 赠送按钮 */}
          {showSendButton && (
            <Button
              size="sm"
              variant="outline"
              className="w-full mt-1"
              onClick={handleSend}
            >
              {sendButtonText}
            </Button>
          )}
        </VStack>
      </CardContent>
    </Card>
  );
});

GiftCard.displayName = 'GiftCard';
