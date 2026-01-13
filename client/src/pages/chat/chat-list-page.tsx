import { useEffect } from 'react';
import { useChatStore } from '@/stores';
import { PageContainer } from '@/components/page-container';
import { ScrollArea } from '@/components/ui/scroll-area';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { useNavigate } from 'react-router-dom';
import { formatDistanceToNow } from 'date-fns';
import { enUS, zhCN } from 'date-fns/locale';
import { useTranslation } from 'react-i18next';
import { Search, MessageSquarePlus } from 'lucide-react';

export default function ChatListPage() {
    const { conversations, fetchConversations, loading } = useChatStore();
    const navigate = useNavigate();
    const { t, i18n } = useTranslation();

    useEffect(() => {
        fetchConversations();
    }, []);

    const handleSelect = (id: string) => {
        navigate(`/chat/${id}`);
    };

    return (
        <PageContainer scrollable={false} className="h-full flex flex-col">
            <div className="flex-none px-6 py-6 space-y-4 animate-in fade-in slide-in-from-top-2">
                <div className="flex items-center justify-between">
                    <div>
                        <h1 className="text-3xl font-bold tracking-tight bg-clip-text text-transparent bg-gradient-to-r from-white to-white/60">
                            {t('nav.chat_title')}
                        </h1>
                    </div>
                    <Button size="icon" variant="ghost" className="rounded-full bg-white/5 hover:bg-white/10">
                        <MessageSquarePlus className="h-5 w-5" />
                    </Button>
                </div>

                <div className="relative">
                    <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                    <Input
                        placeholder={t('nav.search_placeholder')}
                        className="pl-9 bg-muted/40 border-white/5 rounded-full focus-visible:ring-primary/50"
                    />
                </div>
            </div>

            <ScrollArea className="flex-1 px-4">
                <div className="space-y-2 pb-4">
                    {loading && conversations.length === 0 ? (
                        <div className="text-center py-10 text-muted-foreground animate-pulse">Loading chats...</div>
                    ) : conversations.length === 0 ? (
                        <div className="flex flex-col items-center justify-center py-20 text-center space-y-4 text-muted-foreground">
                            <div className="p-4 rounded-full bg-muted/30">
                                <MessageSquarePlus className="h-8 w-8 opacity-50" />
                            </div>
                            <p>{t('nav.no_conversations')}</p>
                        </div>
                    ) : (
                        conversations.map((chat, i) => (
                            <div
                                key={chat.id}
                                className="group p-4 flex items-center gap-4 cursor-pointer hover:bg-muted/40 active:bg-muted/60 transition-all rounded-2xl border border-transparent hover:border-white/5 animate-in fade-in slide-in-from-bottom-2"
                                style={{ animationDelay: `${i * 50}ms` }}
                                onClick={() => handleSelect(chat.id)}
                            >
                                <div className="relative shrink-0">
                                    <Avatar className="h-12 w-12 border-2 border-background group-hover:scale-105 transition-transform">
                                        <AvatarImage src={chat.participantAvatar} />
                                        <AvatarFallback className="bg-primary/10 text-primary">{chat.participantName[0]}</AvatarFallback>
                                    </Avatar>
                                    {chat.online && (
                                        <span className="absolute bottom-0 right-0 w-3.5 h-3.5 bg-green-500 border-2 border-background rounded-full ring-1 ring-black/5" />
                                    )}
                                </div>

                                <div className="flex-1 min-w-0">
                                    <div className="flex justify-between items-center mb-1">
                                        <h3 className="font-semibold truncate text-base group-hover:text-primary transition-colors">{chat.participantName}</h3>
                                        <span className="text-[10px] text-muted-foreground font-medium">
                                            {chat.lastMessageTime && formatDistanceToNow(new Date(chat.lastMessageTime), {
                                                addSuffix: true,
                                                locale: i18n.language === 'zh-CN' ? zhCN : enUS
                                            })}
                                        </span>
                                    </div>
                                    <div className="flex items-center justify-between">
                                        <p className={`text-sm truncate max-w-[85%] ${chat.unreadCount > 0 ? 'text-foreground font-medium' : 'text-muted-foreground'}`}>
                                            {chat.lastMessage}
                                        </p>
                                        {chat.unreadCount > 0 && (
                                            <Badge variant="default" className="rounded-full h-5 min-w-[20px] px-1.5 flex items-center justify-center bg-primary text-[10px] shadow-sm animate-in zoom-in">
                                                {chat.unreadCount}
                                            </Badge>
                                        )}
                                    </div>
                                </div>
                            </div>
                        ))
                    )}
                </div>
            </ScrollArea>
        </PageContainer>
    );
}
