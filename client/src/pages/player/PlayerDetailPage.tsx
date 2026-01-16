import { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { usePlayerStore, useOrderStore } from '@/stores';
import { PageContainer } from '@/components/page-container';
import { Button } from '@/components/ui/button';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Badge } from '@/components/ui/badge';
import { Card, CardContent } from '@/components/ui/card';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter } from '@/components/ui/dialog';
import { Separator } from '@/components/ui/separator';
import { useTranslation } from 'react-i18next';
import { Star, MessageSquare, ArrowLeft, Gamepad2, CheckCircle2, Shield, Loader2 } from 'lucide-react';
import { toast } from 'sonner';
import { FavoriteButton } from '@/components/player/favorite-button';
import { ReviewList, type Review } from '@/components/player/review-list';

export default function PlayerDetailPage() {
    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate();
    const { t } = useTranslation();
    const { currentPlayer, fetchPlayerById, loading: playerLoading } = usePlayerStore();
    const { createOrder, loading: orderLoading } = useOrderStore();

    const [isBookingOpen, setIsBookingOpen] = useState(false);
    const [quantity, setQuantity] = useState(1);
    const [bookingSuccess, setBookingSuccess] = useState(false);

    // Mock Reviews State (In real app, fetch from API)
    const [reviews] = useState<Review[]>([
        { id: 1, userId: 101, username: 'ProGamer_X', rating: 5, content: 'Amazing skill, carried hard!', tags: ['Pro', 'Carry'], createdAt: new Date().toISOString() },
        { id: 2, userId: 102, username: 'Alice_W', rating: 5, content: 'Very patient and friendly.', tags: ['Friendly', 'Patient'], createdAt: new Date(Date.now() - 86400000).toISOString() },
    ]);

    useEffect(() => {
        if (id) {
            fetchPlayerById(Number(id));
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [id]);

    const handleCreateOrder = async () => {
        if (!currentPlayer) return;

        try {
            await createOrder({
                playerId: currentPlayer.id,
                gameId: currentPlayer.gameId,
                quantity: quantity,
                amount: currentPlayer.price * quantity, // Total price (assuming price is per hour)
            });
            setBookingSuccess(true);
            toast.success(t('order.success', { defaultValue: 'Order created successfully!' }));
            setTimeout(() => {
                setIsBookingOpen(false);
                setBookingSuccess(false);
                navigate('/orders');
            }, 1500);
        } catch {
            toast.error(t('order.failed', { defaultValue: 'Failed to create order.' }));
        }
    };

    if (playerLoading || !currentPlayer) {
        return (
            <div className="flex items-center justify-center h-screen bg-background text-muted-foreground">
                <Loader2 className="h-8 w-8 animate-spin" />
            </div>
        );
    }

    return (
        <PageContainer>
            {/* Header Image Background */}
            <div className="relative h-64 md:h-80 w-full overflow-hidden">
                <div className="absolute inset-0 bg-gradient-to-t from-background via-background/60 to-transparent z-10" />
                <div className="absolute inset-0 bg-gradient-to-br from-primary/20 via-muted to-muted">
                    {currentPlayer.avatar && (
                        <img
                            src={currentPlayer.avatar}
                            alt="Cover"
                            className="w-full h-full object-cover"
                        />
                    )}
                </div>

                <div className="absolute top-4 left-4 z-20 flex w-full pr-8 justify-between items-center">
                    <Button variant="ghost" size="icon" onClick={() => navigate(-1)} className="rounded-full bg-background/20 backdrop-blur-md hover:bg-background/40 text-white">
                        <ArrowLeft className="h-5 w-5" />
                    </Button>

                    <FavoriteButton
                        playerId={Number(id)}
                        initialIsFavorite={false}
                        className="bg-black/20 backdrop-blur-md hover:bg-black/40 text-white hover:text-pink-500 rounded-full h-10 w-10"
                    />
                </div>
            </div>

            <div className="max-w-4xl mx-auto px-4 -mt-20 relative z-20 pb-24">
                {/* Profile Header */}
                <div className="flex flex-col md:flex-row gap-6 items-start">
                    <div className="relative">
                        <Avatar className="h-32 w-32 border-4 border-background shadow-2xl ring-2 ring-primary/20">
                            <AvatarImage src={currentPlayer.avatar} />
                            <AvatarFallback className="text-4xl">{currentPlayer.nickname[0]}</AvatarFallback>
                        </Avatar>
                        <Badge variant={currentPlayer.online ? "default" : "secondary"} className={`absolute bottom-2 right-2 border-2 border-background ${currentPlayer.online ? 'bg-green-500 hover:bg-green-600' : 'bg-gray-500'}`}>
                            {currentPlayer.online ? t('nav.status.online') : t('nav.status.offline')}
                        </Badge>
                    </div>

                    <div className="flex-1 space-y-2 pt-2 md:pt-12">
                        <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
                            <div>
                                <h1 className="text-3xl font-bold tracking-tight">{currentPlayer.nickname}</h1>
                                <div className="flex items-center gap-2 text-muted-foreground mt-1">
                                    <Gamepad2 className="h-4 w-4" />
                                    <span>{currentPlayer.gameName}</span>
                                    <span className="mx-1">•</span>
                                    <Star className="h-4 w-4 text-yellow-500 fill-yellow-500" />
                                    <span className="font-medium text-foreground">{currentPlayer.rating.toFixed(1)}</span>
                                </div>
                            </div>
                            <div className="flex flex-col items-end">
                                <div className="text-2xl font-bold text-primary">
                                    ¥{currentPlayer.price}<span className="text-sm font-normal text-muted-foreground">/hr</span>
                                </div>
                                <div className="text-xs text-muted-foreground">
                                    {currentPlayer.orderCount} orders completed
                                </div>
                            </div>
                        </div>

                        {/* Tags */}
                        <div className="flex flex-wrap gap-2 pt-2">
                            {currentPlayer.tags.map(tag => (
                                <Badge key={tag} variant="secondary" className="px-3 py-1 bg-secondary/50 backdrop-blur-sm border border-white/5">
                                    {tag}
                                </Badge>
                            ))}
                        </div>
                    </div>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-3 gap-8 mt-12">
                    {/* Main Content: Bio & Services */}
                    <div className="md:col-span-2 space-y-8">
                        {/* Service Card */}
                        <Card className="border-primary/20 bg-background/50 backdrop-blur-sm">
                            <CardContent className="p-6">
                                <h3 className="text-lg font-semibold mb-4">Service Details</h3>
                                <div className="space-y-4">
                                    <div className="flex gap-4">
                                        <div className="p-2 h-fit bg-primary/10 rounded-lg text-primary">
                                            <Gamepad2 className="h-6 w-6" />
                                        </div>
                                        <div>
                                            <div className="font-medium">Game Coaching / Playing</div>
                                            <p className="text-sm text-muted-foreground mt-1">
                                                Professional gaming accompaniment. Whether you need a carry, coaching, or just a fun duo partner, I'm here to help you win.
                                            </p>
                                        </div>
                                    </div>
                                    <Separator />
                                    <div className="flex gap-4">
                                        <div className="p-2 h-fit bg-blue-500/10 rounded-lg text-blue-500">
                                            <Shield className="h-6 w-6" />
                                        </div>
                                        <div>
                                            <div className="font-medium">Verified Pro</div>
                                            <p className="text-sm text-muted-foreground mt-1">
                                                Identity and skill level verified by GameLink platform.
                                            </p>
                                        </div>
                                    </div>
                                </div>
                            </CardContent>
                        </Card>

                        {/* Recent Reviews */}
                        <div>
                            <ReviewList reviews={reviews} />
                        </div>
                    </div>

                    {/* Sidebar: Sticky Actions (Desktop) */}
                    <div className="space-y-4">
                        {/* Notice */}
                        <Card className="bg-yellow-500/5 border-yellow-500/20">
                            <CardContent className="p-4 text-sm text-yellow-600/90 dark:text-yellow-500">
                                <span className="font-bold block mb-1">Safety Tip</span>
                                Always communicate through the GameLink chat to ensure your order is protected.
                            </CardContent>
                        </Card>
                    </div>
                </div>
            </div>

            {/* Sticky Bottom Bar (Mobile & Desktop) */}
            <div className="fixed bottom-0 left-0 right-0 p-4 border-t border-white/10 bg-background/80 backdrop-blur-xl z-50 md:hidden">
                <div className="flex gap-3">
                    <Button variant="outline" className="flex-1" onClick={() => navigate(`/chat/${currentPlayer.id}`)}>
                        <MessageSquare className="h-4 w-4 mr-2" />
                        Chat
                    </Button>
                    <Button className="flex-[2] font-bold" onClick={() => setIsBookingOpen(true)}>
                        Book Now
                    </Button>
                </div>
            </div>
            {/* Floating Action for Desktop (hidden on mobile) */}
            <div className="hidden md:flex fixed bottom-8 right-8 gap-3 z-50">
                <Button variant="secondary" size="lg" className="shadow-xl" onClick={() => navigate(`/chat/${currentPlayer.id}`)}>
                    <MessageSquare className="h-4 w-4 mr-2" />
                    Chat with Player
                </Button>
                <Button size="lg" className="shadow-xl font-bold px-8 scale-110 origin-right" onClick={() => setIsBookingOpen(true)}>
                    Book Session
                </Button>
            </div>

            {/* Booking Dialog */}
            <Dialog open={isBookingOpen} onOpenChange={setIsBookingOpen}>
                <DialogContent className="sm:max-w-md">
                    <DialogHeader>
                        <DialogTitle>Book Session</DialogTitle>
                        <DialogDescription>
                            Configure your order details below.
                        </DialogDescription>
                    </DialogHeader>

                    {bookingSuccess ? (
                        <div className="py-12 flex flex-col items-center justify-center text-center space-y-4 animate-in fade-in zoom-in">
                            <div className="h-16 w-16 bg-green-500/10 text-green-500 rounded-full flex items-center justify-center mb-2">
                                <CheckCircle2 className="h-8 w-8" />
                            </div>
                            <h3 className="text-xl font-bold">Booking Confirmed!</h3>
                            <p className="text-muted-foreground">Redirecting to your orders...</p>
                        </div>
                    ) : (
                        <div className="space-y-6 py-4">
                            <div className="flex items-center gap-4">
                                <Avatar className="h-12 w-12">
                                    <AvatarImage src={currentPlayer.avatar} />
                                    <AvatarFallback>{currentPlayer.nickname[0]}</AvatarFallback>
                                </Avatar>
                                <div>
                                    <div className="font-bold">{currentPlayer.nickname}</div>
                                    <div className="text-sm text-muted-foreground">{currentPlayer.gameName}</div>
                                </div>
                                <div className="ml-auto text-right">
                                    <div className="font-bold">¥{currentPlayer.price}</div>
                                    <div className="text-xs text-muted-foreground">per hour</div>
                                </div>
                            </div>

                            <Separator />

                            <div className="space-y-4">
                                <div className="flex items-center justify-between">
                                    <span className="text-sm font-medium">Duration</span>
                                    <div className="flex items-center gap-3">
                                        <Button
                                            variant="outline" size="icon" className="h-8 w-8 rounded-full"
                                            onClick={() => setQuantity(Math.max(1, quantity - 1))}
                                            disabled={quantity <= 1}
                                        >
                                            -
                                        </Button>
                                        <span className="w-4 text-center font-bold">{quantity}</span>
                                        <Button
                                            variant="outline" size="icon" className="h-8 w-8 rounded-full"
                                            onClick={() => setQuantity(quantity + 1)}
                                        >
                                            +
                                        </Button>
                                    </div>
                                </div>
                                <div className="flex items-center justify-between text-muted-foreground text-sm">
                                    <span>Total Hours</span>
                                    <span>{quantity}h</span>
                                </div>
                            </div>

                            <div className="bg-muted/50 p-4 rounded-lg flex justify-between items-center">
                                <span className="font-bold">Total Price</span>
                                <span className="text-2xl font-bold text-primary">¥{(currentPlayer.price * quantity).toFixed(2)}</span>
                            </div>
                        </div>
                    )}

                    {!bookingSuccess && (
                        <DialogFooter>
                            <Button variant="ghost" onClick={() => setIsBookingOpen(false)}>Cancel</Button>
                            <Button onClick={handleCreateOrder} disabled={orderLoading}>
                                {orderLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                                Confirm & Pay
                            </Button>
                        </DialogFooter>
                    )}
                </DialogContent>
            </Dialog>
        </PageContainer>
    );
}
