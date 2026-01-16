import { Mic, MicOff, VolumeX } from 'lucide-react';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { cn } from '@/lib/utils';
import type { VoiceMember } from '@/stores/modules/voice-store';

interface VoiceMemberItemProps {
    member: VoiceMember;
}

export function VoiceMemberItem({ member }: VoiceMemberItemProps) {
    return (
        <div
            className={cn(
                'flex items-center gap-3 p-2 rounded-lg transition-colors',
                member.isSpeaking && 'bg-green-500/10 ring-2 ring-green-500/50'
            )}
        >
            <div className="relative">
                <Avatar className="h-8 w-8">
                    <AvatarImage src={member.avatarUrl} />
                    <AvatarFallback className="text-xs">
                        {member.nickname?.charAt(0) || 'U'}
                    </AvatarFallback>
                </Avatar>
                {/* Speaking indicator */}
                {member.isSpeaking && (
                    <div className="absolute -bottom-0.5 -right-0.5 h-3 w-3 bg-green-500 rounded-full border-2 border-background animate-pulse" />
                )}
            </div>

            <div className="flex-1 min-w-0">
                <span
                    className={cn(
                        'text-sm font-medium truncate block',
                        member.isSpeaking ? 'text-green-600' : 'text-foreground'
                    )}
                >
                    {member.nickname}
                </span>
            </div>

            {/* Status icons */}
            <div className="flex items-center gap-1">
                {member.isMuted && (
                    <MicOff className="h-4 w-4 text-red-500" />
                )}
                {member.isDeafened && (
                    <VolumeX className="h-4 w-4 text-red-500" />
                )}
                {!member.isMuted && !member.isDeafened && (
                    <Mic
                        className={cn(
                            'h-4 w-4',
                            member.isSpeaking ? 'text-green-500' : 'text-muted-foreground'
                        )}
                    />
                )}
            </div>
        </div>
    );
}
