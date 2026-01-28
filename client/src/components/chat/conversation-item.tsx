import { forwardRef } from 'react';
import { cn } from '@/lib/utils';
import { Badge } from '@/components/ui/badge';
import { UserAvatar } from '@/components/user';
import { HStack, VStack } from '@/components/layout';
import { formatDistanceToNow } from 'date-fns';
import { zhCN } from 'date-fns/locale';

export interface ConversationItemProps extends React.HTMLAttributes<HTMLDivElement> {
  /** 会话 ID */
  conversationId: number;
  /** 对方名称 */
  participantName: string;
  /** 对方头像 */
  participantAvatar?: string;
  /** 是否在线 */
  online?: boolean;
  /** 最后一条消息 */
  lastMessage?: string;
  /** 最后消息时间 */
  lastMessageTime?: string | Date;
  /** 未读数量 */
  unreadCount?: number;
  /** 是否选中 */
  selected?: boolean;
  /** 点击回调 */
  onSelect?: (id: number) => void;
}

export const ConversationItem = forwardRef<HTMLDivElement, ConversationItemProps>(({
  conversationId,
  participantName,
  participantAvatar,
  online = false,
  lastMessage,
  lastMessageTime,
  unreadCount = 0,
  selected = false,
  onSelect,
  className,
  ...props
}, ref) => {
  const handleClick = () => {
    onSelect?.(conversationId);
  };

  const formattedTime = lastMessageTime
    ? formatDistanceToNow(new Date(lastMessageTime), { addSuffix: true, locale: zhCN })
    : '';

  return (
    <div
      ref={ref}
      className={cn(
        'group p-4 flex items-center gap-4 cursor-pointer',
        'hover:bg-muted/40 active:bg-muted/60',
        'transition-all rounded-2xl border border-transparent',
        'hover:border-white/5',
        selected && 'bg-muted/50 border-primary/20',
        className
      )}
      onClick={handleClick}
      {...props}
    >
      <UserAvatar
        src={participantAvatar}
        name={participantName}
        size="lg"
        status={online ? 'online' : 'offline'}
        showStatus
        className="shrink-0 group-hover:scale-105 transition-transform"
      />

      <VStack spacing={1} className="flex-1 min-w-0">
        <HStack justify="between" align="center" className="w-full">
          <span className={cn(
            'font-semibold truncate text-base',
            'group-hover:text-primary transition-colors'
          )}>
            {participantName}
          </span>
          {formattedTime && (
            <span className="text-xs text-muted-foreground/60 shrink-0">
              {formattedTime}
            </span>
          )}
        </HStack>

        <HStack justify="between" align="center" className="w-full">
          <p className={cn(
            'text-sm text-muted-foreground truncate',
            unreadCount > 0 && 'text-foreground font-medium'
          )}>
            {lastMessage || '暂无消息'}
          </p>
          {unreadCount > 0 && (
            <Badge
              variant="default"
              className="ml-2 h-5 min-w-[20px] px-1.5 text-[10px] rounded-full shrink-0"
            >
              {unreadCount > 99 ? '99+' : unreadCount}
            </Badge>
          )}
        </HStack>
      </VStack>
    </div>
  );
});

ConversationItem.displayName = 'ConversationItem';
