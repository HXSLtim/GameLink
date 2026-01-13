import { useEffect, useState, useRef } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useChatStore, useAuthStore } from '@/stores';
import { PageContainer } from '@/components/page-container';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { ScrollArea } from '@/components/ui/scroll-area';
import { Send, ArrowLeft, Image as ImageIcon } from 'lucide-react';
import { format } from 'date-fns';

export default function ChatRoomPage() {
    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate();
    const { user } = useAuthStore();
    const {
        messages,
        conversations,
        selectConversation,
        sendMessage,
        loading
    } = useChatStore();

    const [inputText, setInputText] = useState('');
    const scrollRef = useRef<HTMLDivElement>(null);

    const conversation = conversations.find(c => c.id === id);
    const currentMessages = id ? (messages[id] || []) : [];

    useEffect(() => {
        if (id) {
            selectConversation(id);
        }
    }, [id, selectConversation]);

    // Auto-scroll to bottom
    useEffect(() => {
        if (scrollRef.current) {
            scrollRef.current.scrollTop = scrollRef.current.scrollHeight;
        }
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
                <div className="flex flex-col items-center justify-center h-full space-y-4">
                    <p className="text-muted-foreground">Conversation not found.</p>
                    <Button variant="outline" onClick={() => navigate('/chat')}>Back to list</Button>
                </div>
            </PageContainer>
        );
    }

    return (
        <PageContainer scrollable={false} className="h-full flex flex-col">
            {/* Header */}
            <div className="flex items-center gap-4 border-b px-4 py-3 bg-card/50 backdrop-blur shrink-0">
                <Button variant="ghost" size="icon" onClick={() => navigate('/chat')} className="-ml-2">
                    <ArrowLeft className="h-5 w-5" />
                </Button>

                <Avatar className="h-9 w-9">
                    <AvatarImage src={conversation?.participantAvatar} />
                    <AvatarFallback>{conversation?.participantName?.[0]}</AvatarFallback>
                </Avatar>

                <div>
                    <h2 className="font-semibold text-base leading-none">{conversation?.participantName}</h2>
                    <span className="text-xs text-muted-foreground flex items-center gap-1.5 mt-0.5">
                        <span className={`w-1.5 h-1.5 rounded-full ${conversation?.online ? 'bg-green-500' : 'bg-gray-400'}`} />
                        {conversation?.online ? 'Online' : 'Offline'}
                    </span>
                </div>
            </div>

            {/* Messages Area */}
            <div className="flex-1 overflow-hidden flex flex-col bg-background/50">
                <ScrollArea className="flex-1 p-4" ref={scrollRef}>
                    <div className="space-y-4 px-4 w-full">
                        {currentMessages.map((msg) => {
                            const isMe = msg.senderId === 0 || msg.senderId === Number(user?.id);
                            return (
                                <div key={msg.id} className={`flex ${isMe ? 'justify-end' : 'justify-start'}`}>
                                    <div
                                        className={`max-w-[85%] sm:max-w-[70%] rounded-2xl px-4 py-2.5 text-sm shadow-sm ${isMe
                                            ? 'bg-primary text-primary-foreground rounded-br-none'
                                            : 'bg-card text-card-foreground border rounded-bl-none'
                                            }`}
                                    >
                                        <p className="leading-relaxed">{msg.content}</p>
                                        <span className={`text-[10px] block mt-1 opacity-70 ${isMe ? 'text-primary-foreground/80 text-right' : 'text-muted-foreground text-left'}`}>
                                            {format(new Date(msg.createdAt), 'HH:mm')}
                                        </span>
                                    </div>
                                </div>
                            );
                        })}
                        {loading && <div className="text-center text-xs text-muted-foreground py-4">Loading history...</div>}
                    </div>
                </ScrollArea>

                {/* Input Area */}
                <div className="p-4 bg-background border-t shrink-0">
                    <div className="flex items-end gap-2 bg-muted/50 p-1.5 rounded-xl border focus-within:ring-1 focus-within:ring-ring focus-within:border-primary/50 transition-all w-full">
                        <Button variant="ghost" size="icon" className="shrink-0 rounded-lg text-muted-foreground hover:text-foreground h-9 w-9 my-0.5 ml-0.5">
                            <ImageIcon className="h-5 w-5" />
                        </Button>

                        <Input
                            value={inputText}
                            onChange={(e) => setInputText(e.target.value)}
                            onKeyDown={handleKeyPress}
                            placeholder="Type a message..."
                            className="flex-1 border-none shadow-none focus-visible:ring-0 bg-transparent min-h-[40px] py-2.5 px-2"
                        />

                        <Button
                            onClick={handleSend}
                            size="icon"
                            disabled={!inputText.trim()}
                            className={!inputText.trim() ? "opacity-50" : "opacity-100"}
                            variant={inputText.trim() ? "default" : "ghost"}
                        >
                            <Send className="h-4 w-4" />
                        </Button>
                    </div>
                </div>
            </div>
        </PageContainer>
    );
}
