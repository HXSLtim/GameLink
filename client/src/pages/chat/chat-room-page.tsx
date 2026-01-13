import { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useChatStore, useAuthStore } from '@/stores';
import { PageContainer } from '@/components/page-container';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { ScrollArea } from '@/components/ui/scroll-area';
import { Send, ArrowLeft, Image as ImageIcon, MoreVertical, Phone } from 'lucide-react';
import { format } from 'date-fns';
import { useTranslation } from 'react-i18next';

export default function ChatRoomPage() {
    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate();
    const { t } = useTranslation();
    const { user } = useAuthStore();
    const {
        messages,
        conversations,
        selectConversation,
        sendMessage,
        loading
    } = useChatStore();

    const [inputText, setInputText] = useState('');

    const conversation = conversations.find(c => c.id === id);
    const currentMessages = id ? (messages[id] || []) : [];

    useEffect(() => {
        if (id) {
            selectConversation(id);
        }
    }, [id, selectConversation]);

    // Auto-scroll to bottom logic adapted for Radix ScrollArea
    useEffect(() => {
        const scrollToBottom = () => {
            const viewport = document.querySelector('[data-radix-scroll-area-viewport]');
            if (viewport) {
                viewport.scrollTop = viewport.scrollHeight;
            }
        };
        // Small timeout to ensure DOM is rendered
        setTimeout(scrollToBottom, 50);
    }, [currentMessages]);


    const handleSend = async () => {
        if (!inputText.trim() || !id) return;
        await sendMessage(inputText, 'text');
        setInputText('');
    };

    const handleKeyPress = (e: React.KeyboardEvent) => {
        if (e.key === 'Enter' && !e.shiftKey) {
            e.preventDefault();
            handleSend();
        }
    };

    if (!conversation && !loading) {
        return (
            <PageContainer>
                <div className="flex flex-col items-center justify-center h-full space-y-4 animate-in fade-in">
                    <p className="text-muted-foreground">Conversation not found.</p>
                    <Button variant="outline" onClick={() => navigate('/chat')}>Back to list</Button>
                </div>
            </PageContainer>
        );
    }

    return (
        <PageContainer scrollable={false} className="h-full flex flex-col">
            {/* Header */}
            <div className="flex items-center justify-between border-b border-white/5 px-4 py-3 bg-background/80 backdrop-blur-md shrink-0 z-20">
                <div className="flex items-center gap-3">
                    <Button variant="ghost" size="icon" onClick={() => navigate('/chat')} className="-ml-2 hover:bg-white/5 rounded-full">
                        <ArrowLeft className="h-5 w-5" />
                    </Button>

                    <div className="relative">
                        <Avatar className="h-10 w-10 border border-white/10">
                            <AvatarImage src={conversation?.participantAvatar} />
                            <AvatarFallback>{conversation?.participantName?.[0]}</AvatarFallback>
                        </Avatar>
                        {conversation?.online && (
                            <span className="absolute bottom-0 right-0 w-3 h-3 bg-green-500 border-2 border-background rounded-full" />
                        )}
                    </div>

                    <div>
                        <h2 className="font-semibold text-base leading-none">{conversation?.participantName}</h2>
                        <span className="text-xs text-muted-foreground flex items-center gap-1.5 mt-1">
                            {conversation?.online ? t('nav.status.online') : t('nav.status.online')}
                            {/* Ideally offline status, but keeping simple */}
                        </span>
                    </div>
                </div>

                <div className="flex items-center gap-1">
                    <Button variant="ghost" size="icon" className="hover:bg-white/5 rounded-full text-muted-foreground">
                        <Phone className="h-5 w-5" />
                    </Button>
                    <Button variant="ghost" size="icon" className="hover:bg-white/5 rounded-full text-muted-foreground">
                        <MoreVertical className="h-5 w-5" />
                    </Button>
                </div>
            </div>

            {/* Messages Area */}
            <div className="flex-1 overflow-hidden flex flex-col bg-background/30 relative">
                {/* Wallpaper effect */}
                <div className="absolute inset-0 bg-[url('https://images.unsplash.com/photo-1550745165-9bc0b252726f?q=80&w=2670&auto=format&fit=crop')] bg-cover bg-center opacity-5 pointer-events-none" />

                <ScrollArea className="flex-1 p-4">
                    <div className="space-y-6 px-2 pb-4 w-full pt-4">
                        {currentMessages.map((msg) => {
                            const isMe = msg.senderId === 0 || msg.senderId === Number(user?.id);
                            return (
                                <div key={msg.id} className={`flex ${isMe ? 'justify-end' : 'justify-start'} animate-in slide-in-from-bottom-2 duration-300`}>
                                    <div
                                        className={`max-w-[85%] sm:max-w-[70%] rounded-2xl px-5 py-3 text-sm shadow-sm backdrop-blur-sm ${isMe
                                            ? 'bg-primary text-primary-foreground rounded-br-sm'
                                            : 'bg-muted/80 text-foreground border border-white/5 rounded-bl-sm'
                                            }`}
                                    >
                                        <p className="leading-relaxed whitespace-pre-wrap break-words">{msg.content}</p>
                                        <span className={`text-[10px] block mt-1.5 opacity-70 ${isMe ? 'text-primary-foreground/90 text-right' : 'text-muted-foreground text-left'}`}>
                                            {format(new Date(msg.createdAt), 'HH:mm')}
                                            {isMe && <span className="ml-1">✓</span>}
                                        </span>
                                    </div>
                                </div>
                            );
                        })}
                        {loading && <div className="text-center text-xs text-muted-foreground py-4 animate-pulse">Updating...</div>}
                    </div>
                </ScrollArea>

                {/* Input Area */}
                <div className="p-4 bg-background/80 backdrop-blur-md border-t border-white/5 shrink-0 z-20">
                    <div className="flex items-end gap-2 bg-muted/50 p-1.5 rounded-[24px] border border-white/5 focus-within:ring-1 focus-within:ring-primary/50 focus-within:border-primary/50 transition-all w-full shadow-lg">
                        <Button variant="ghost" size="icon" className="shrink-0 rounded-full text-muted-foreground hover:text-foreground h-9 w-9 my-0.5 ml-0.5 hover:bg-background">
                            <ImageIcon className="h-5 w-5" />
                        </Button>

                        <Input
                            value={inputText}
                            onChange={(e) => setInputText(e.target.value)}
                            onKeyDown={handleKeyPress}
                            placeholder={t('nav.chat_placeholder')}
                            className="flex-1 border-none shadow-none focus-visible:ring-0 bg-transparent min-h-[40px] py-2.5 px-2 text-base md:text-sm"
                        />

                        <Button
                            onClick={handleSend}
                            size="icon"
                            disabled={!inputText.trim()}
                            className={`rounded-full shadow-sm transition-all duration-200 ${!inputText.trim() ? "opacity-50 scale-90" : "opacity-100 scale-100 shadow-primary/25"}`}
                            variant={inputText.trim() ? "default" : "ghost"}
                        >
                            <Send className="h-4 w-4 ml-0.5" />
                        </Button>
                    </div>
                </div>
            </div>
        </PageContainer>
    );
}
