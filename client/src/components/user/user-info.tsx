import { forwardRef } from 'react';
import { cn } from '@/lib/utils';
import { Badge } from '@/components/ui/badge';
import { UserAvatar, type UserAvatarProps } from './user-avatar';
import { HStack, VStack } from '@/components/layout';
import { Crown, Shield, Star } from 'lucide-react';

export interface UserInfoProps extends React.HTMLAttributes<HTMLDivElement> {
  /** 用户 ID */
  userId?: number;
  /** 用户名 */
  username?: string;
  /** 昵称 */
  nickname?: string;
  /** 头像 URL */
  avatarUrl?: string;
  /** 角色 */
  role?: 'user' | 'player' | 'admin';
  /** 是否是 VIP */
  isVip?: boolean;
  /** VIP 等级 */
  _vipLevel?: number;
  /** 在线状态 */
  status?: 'online' | 'offline' | 'busy' | 'away' | 'in_game';
  /** 副标题/描述 */
  subtitle?: string;
  /** 头像尺寸 */
  avatarSize?: UserAvatarProps['size'];
  /** 布局方向 */
  direction?: 'horizontal' | 'vertical';
  /** 是否显示 ID */
  showId?: boolean;
  /** 是否显示角色徽章 */
  showRole?: boolean;
  /** 点击回调 */
  onClick?: () => void;
}

const roleConfig = {
  user: { label: '用户', variant: 'secondary' as const, icon: null },
  player: { label: '陪玩', variant: 'default' as const, icon: Star },
  admin: { label: '管理员', variant: 'destructive' as const, icon: Shield },
};

export const UserInfo = forwardRef<HTMLDivElement, UserInfoProps>(({
  userId,
  username,
  nickname,
  avatarUrl,
  role = 'user',
  isVip = false,
  vipLevel,
  status,
  subtitle,
  avatarSize = 'md',
  direction = 'horizontal',
  showId = false,
  showRole = true,
  onClick,
  className,
  ...props
}, ref) => {
  const displayName = nickname || username || '用户';
  const roleConf = roleConfig[role];
  const RoleIcon = roleConf.icon;

  const content = direction === 'vertical' ? (
    <VStack spacing={2} align="center" className="text-center">
      <UserAvatar
        src={avatarUrl}
        name={displayName}
        size={avatarSize}
        status={status}
        isVip={isVip}
        bordered
      />
      <VStack spacing={0.5} align="center">
        <HStack spacing={1} align="center">
          <span className="font-semibold truncate">{displayName}</span>
          {isVip && <Crown className="h-4 w-4 text-yellow-500 fill-yellow-500" />}
        </HStack>
        {subtitle && (
          <span className="text-xs text-muted-foreground">{subtitle}</span>
        )}
        {showId && userId && (
          <span className="text-xs text-muted-foreground/60 font-mono">
            ID: {userId}
          </span>
        )}
      </VStack>
      {showRole && (
        <Badge variant={roleConf.variant} className="text-xs">
          {RoleIcon && <RoleIcon className="h-3 w-3 mr-1" />}
          {roleConf.label}
        </Badge>
      )}
    </VStack>
  ) : (
    <HStack spacing={3} align="center">
      <UserAvatar
        src={avatarUrl}
        name={displayName}
        size={avatarSize}
        status={status}
        isVip={isVip}
      />
      <VStack spacing={0.5} className="min-w-0 flex-1">
        <HStack spacing={2} align="center">
          <span className="font-semibold truncate">{displayName}</span>
          {isVip && <Crown className="h-4 w-4 text-yellow-500 fill-yellow-500 shrink-0" />}
          {showRole && (
            <Badge variant={roleConf.variant} className="text-xs shrink-0">
              {RoleIcon && <RoleIcon className="h-3 w-3 mr-1" />}
              {roleConf.label}
            </Badge>
          )}
        </HStack>
        {subtitle && (
          <span className="text-xs text-muted-foreground truncate">{subtitle}</span>
        )}
        {showId && userId && (
          <span className="text-xs text-muted-foreground/60 font-mono">
            ID: {userId}
          </span>
        )}
      </VStack>
    </HStack>
  );

  return (
    <div
      ref={ref}
      className={cn(
        onClick && 'cursor-pointer hover:opacity-80 transition-opacity',
        className
      )}
      onClick={onClick}
      {...props}
    >
      {content}
    </div>
  );
});

UserInfo.displayName = 'UserInfo';
