import { useEffect } from 'react';
import { useChatStore } from '@/stores';
import { PageContainer, PageHeader } from '@/components/page-container';
import { ScrollArea } from '@/components/ui/scroll-area';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Badge } from '@/components/ui/badge';
import { useNavigate } from 'react-router-dom';
import { formatDistanceToNow } from 'date-fns';
import { zhCN } from 'date-fns/locale';

export default function ChatListPage() {
    const { conversations, fetchConversations, loading } = useChatStore();
    const navigate = useNavigate();

    useEffect(() => {
        fetchConversations();
    }, []); // eslint-disable-line react-hooks/exhaustive-deps

    const handleSelect = (id: string) => {
        navigate(`/chat/${id}`);
    };

    return (
        <PageContainer scrollable={false} className="h-full flex flex-col">
            <PageHeader title="Messages" description="Your conversations with players." className="shrink-0 px-6 py-4" />

            <ScrollArea className="flex-1">
                <div className="divide-y divide-border/50 border-t border-border/50">
                    {loading && conversations.length === 0 ? (
                        <div className="text-center py-10 text-muted-foreground">Loading chats...</div>
                    ) : conversations.length === 0 ? (
                        <div className="text-center py-20 text-muted-foreground">No conversations yet.</div>
                    ) : (
                        conversations.map(chat => (
                            <div
                                key={chat.id}
                                className="px-6 py-4 flex items-center gap-4 cursor-pointer hover:bg-muted/50 transition-colors w-full"
                                onClick={() => handleSelect(chat.id)}
                            >
                                <div className="relative shrink-0">
                                    <Avatar className="h-12 w-12">
                                        <AvatarImage src={chat.participantAvatar} />
                                        <AvatarFallback>{chat.participantName[0]}</AvatarFallback>
                                    </Avatar>
                                    {chat.online && (
                                        <span className="absolute bottom-0 right-0 w-3 h-3 bg-green-500 border-2 border-background rounded-full" />
                                    )}
                                </div>

                                <div className="flex-1 min-w-0">
                                    <div className="flex justify-between items-start mb-1">
                                        <h3 className="font-semibold truncate text-base">{chat.participantName}</h3>
                                        <span className="text-xs text-muted-foreground whitespace-nowrap ml-2">
                                            {chat.lastMessageTime && formatDistanceToNow(new Date(chat.lastMessageTime), { addSuffix: true, locale: zhCN })}
                                        </span>
                                    </div>
                                    <p className="text-sm text-muted-foreground truncate opacity-90">{chat.lastMessage}</p>
                                </div>

                                {chat.unreadCount > 0 && (
                                    <Badge variant="destructive" className="rounded-full h-5 w-5 flex items-center justify-center p-0 text-[10px] shrink-0">
                                        {chat.unreadCount}
                                    </Badge>
                                )}
                            </div>
                        ))
                    )}
                </div>
            </ScrollArea>
        </PageContainer>
    );
}
