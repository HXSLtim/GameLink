import { useTranslation } from 'react-i18next';
import { Crown, Check, User, MoreVertical } from 'lucide-react';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { cn } from '@/lib/utils';
import type { RoomMember } from '@/stores/modules/room-store';

interface RoomMemberListProps {
    members: RoomMember[];
    hostUserId: number;
    currentUserId?: number;
    isHost?: boolean;
    onKick?: (userId: number) => void;
}

export function RoomMemberList({
    members,
    hostUserId,
    currentUserId,
    isHost = false,
    onKick,
}: RoomMemberListProps) {
    const { t } = useTranslation();

    if (members.length === 0) {
        return (
            <div className="text-center py-8 text-muted-foreground">
                {t('room.noMembers')}
            </div>
        );
    }

    return (
        <div className="space-y-2">
            {members.map((member) => (
                <RoomMemberItem
                    key={member.id}
                    member={member}
                    isHostMember={member.userId === hostUserId}
                    isCurrentUser={member.userId === currentUserId}
                    canKick={isHost && member.userId !== hostUserId}
                    onKick={onKick}
                />
            ))}
        </div>
    );
}

interface RoomMemberItemProps {
    member: RoomMember;
    isHostMember: boolean;
    isCurrentUser: boolean;
    canKick: boolean;
    onKick?: (userId: number) => void;
}

function RoomMemberItem({
    member,
    isHostMember,
    isCurrentUser,
    canKick,
    onKick,
}: RoomMemberItemProps) {
    const { t } = useTranslation();

    return (
        <div
            className={cn(
                'flex items-center justify-between p-3 rounded-lg',
                'bg-muted/50 hover:bg-muted/80 transition-colors',
                isCurrentUser && 'ring-1 ring-primary/50'
            )}
        >
            <div className="flex items-center gap-3">
                <div className="relative">
                    <Avatar className="h-10 w-10">
                        <AvatarImage src={member.avatarUrl} />
                        <AvatarFallback>
                            {member.nickname?.charAt(0) || <User className="h-4 w-4" />}
                        </AvatarFallback>
                    </Avatar>
                    {isHostMember && (
                        <div className="absolute -top-1 -right-1 bg-yellow-500 rounded-full p-0.5">
                            <Crown className="h-3 w-3 text-white" />
                        </div>
                    )}
                </div>
                <div>
                    <div className="flex items-center gap-2">
                        <span className="font-medium text-foreground">
                            {member.nickname}
                        </span>
                        {isCurrentUser && (
                            <Badge variant="outline" className="text-xs">
                                {t('room.you')}
                            </Badge>
                        )}
                    </div>
                    <div className="flex items-center gap-2 text-sm text-muted-foreground">
                        {isHostMember ? (
                            <span className="text-yellow-600">{t('room.host')}</span>
                        ) : (
                            <span>{t('room.member')}</span>
                        )}
                    </div>
                </div>
            </div>

            <div className="flex items-center gap-2">
                {/* Ready status */}
                {member.isReady && (
                    <Badge variant="secondary" className="bg-green-500/20 text-green-600">
                        <Check className="h-3 w-3 mr-1" />
                        {t('room.actions.ready')}
                    </Badge>
                )}

                {/* Actions dropdown */}
                {canKick && onKick && (
                    <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                            <Button variant="ghost" size="icon" className="h-8 w-8">
                                <MoreVertical className="h-4 w-4" />
                            </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end">
                            <DropdownMenuItem
                                className="text-destructive"
                                onClick={() => onKick(member.userId)}
                            >
                                {t('room.actions.kick')}
                            </DropdownMenuItem>
                        </DropdownMenuContent>
                    </DropdownMenu>
                )}
            </div>
        </div>
    );
}

export function RoomMemberListSkeleton() {
    return (
        <div className="space-y-2">
            {[1, 2, 3].map((i) => (
                <div
                    key={i}
                    className="flex items-center justify-between p-3 rounded-lg bg-muted/50"
                >
                    <div className="flex items-center gap-3">
                        <div className="h-10 w-10 rounded-full bg-muted animate-pulse" />
                        <div>
                            <div className="h-4 w-24 bg-muted rounded animate-pulse mb-1" />
                            <div className="h-3 w-16 bg-muted rounded animate-pulse" />
                        </div>
                    </div>
                    <div className="h-6 w-16 bg-muted rounded animate-pulse" />
                </div>
            ))}
        </div>
    );
}
