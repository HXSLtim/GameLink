import { forwardRef } from 'react';
import { cn } from '@/lib/utils';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Camera, Crown } from 'lucide-react';

export interface UserAvatarProps extends React.HTMLAttributes<HTMLDivElement> {
  /** 用户头像 URL */
  src?: string;
  /** 用户名（用于 fallback） */
  name?: string;
  /** 尺寸 */
  size?: 'xs' | 'sm' | 'md' | 'lg' | 'xl' | '2xl';
  /** 在线状态 */
  status?: 'online' | 'offline' | 'busy' | 'away' | 'in_game';
  /** 是否显示状态指示器 */
  showStatus?: boolean;
  /** 是否显示编辑按钮 */
  editable?: boolean;
  /** 编辑按钮点击回调 */
  onEdit?: () => void;
  /** 是否是 VIP 用户 */
  isVip?: boolean;
  /** 是否显示边框 */
  bordered?: boolean;
}

const sizeConfig = {
  xs: { avatar: 'h-6 w-6', fallback: 'text-[10px]', status: 'w-2 h-2', edit: 'p-0.5', editIcon: 'h-2 w-2', vip: 'h-3 w-3' },
  sm: { avatar: 'h-8 w-8', fallback: 'text-xs', status: 'w-2.5 h-2.5', edit: 'p-1', editIcon: 'h-2.5 w-2.5', vip: 'h-3.5 w-3.5' },
  md: { avatar: 'h-10 w-10', fallback: 'text-sm', status: 'w-3 h-3', edit: 'p-1', editIcon: 'h-3 w-3', vip: 'h-4 w-4' },
  lg: { avatar: 'h-14 w-14', fallback: 'text-lg', status: 'w-3.5 h-3.5', edit: 'p-1.5', editIcon: 'h-3.5 w-3.5', vip: 'h-5 w-5' },
  xl: { avatar: 'h-20 w-20', fallback: 'text-2xl', status: 'w-4 h-4', edit: 'p-1.5', editIcon: 'h-4 w-4', vip: 'h-5 w-5' },
  '2xl': { avatar: 'h-28 w-28', fallback: 'text-4xl', status: 'w-5 h-5', edit: 'p-1.5', editIcon: 'h-4 w-4', vip: 'h-6 w-6' },
};

const statusColors = {
  online: 'bg-green-500',
  offline: 'bg-gray-400',
  busy: 'bg-red-500',
  away: 'bg-yellow-500',
  in_game: 'bg-blue-500',
};

export const UserAvatar = forwardRef<HTMLDivElement, UserAvatarProps>(({
  src,
  name,
  size = 'md',
  status,
  showStatus = true,
  editable = false,
  onEdit,
  isVip = false,
  bordered = false,
  className,
  ...props
}, ref) => {
  const config = sizeConfig[size];
  const initial = name?.charAt(0).toUpperCase() || 'U';

  return (
    <div
      ref={ref}
      className={cn('relative inline-block group', className)}
      {...props}
    >
      <Avatar className={cn(
        config.avatar,
        bordered && 'border-2 border-background ring-2 ring-primary/20',
        editable && 'cursor-pointer transition-transform group-hover:scale-105'
      )}>
        <AvatarImage src={src} alt={name} />
        <AvatarFallback className={cn(
          'font-bold bg-primary/10 text-primary',
          config.fallback
        )}>
          {initial}
        </AvatarFallback>
      </Avatar>

      {/* 在线状态指示器 */}
      {showStatus && status && (
        <span className={cn(
          'absolute bottom-0 right-0 rounded-full border-2 border-background',
          config.status,
          statusColors[status]
        )} />
      )}

      {/* VIP 徽章 */}
      {isVip && (
        <div className="absolute -top-1 -right-1">
          <Crown className={cn(
            'text-yellow-500 fill-yellow-500',
            config.vip
          )} />
        </div>
      )}

      {/* 编辑按钮 */}
      {editable && (
        <button
          onClick={onEdit}
          className={cn(
            'absolute bottom-0 right-0 bg-background rounded-full border shadow-sm',
            'cursor-pointer hover:bg-muted transition-colors',
            config.edit
          )}
        >
          <Camera className={cn('text-muted-foreground', config.editIcon)} />
        </button>
      )}
    </div>
  );
});

UserAvatar.displayName = 'UserAvatar';
