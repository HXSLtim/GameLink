import { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { http } from '@/lib/http';
import { Button } from '@/components/ui/button';
import { Textarea } from '@/components/ui/textarea';
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { StarRating } from '@/components/ui/star-rating';
import { Badge } from '@/components/ui/badge';
import { Loader2, ArrowLeft, Send } from 'lucide-react';
import { toast } from 'sonner';

interface ReviewTag {
    id: string;
    label: string;
}

const COMMON_TAGS: ReviewTag[] = [
    { id: 'pro', label: 'Pro Gamer' },
    { id: 'funny', label: 'Funny' },
    { id: 'sweet', label: 'Sweet Voice' },
    { id: 'carry', label: 'Hard Carry' },
    { id: 'patient', label: 'Patient' },
    { id: 'friendly', label: 'Friendly' },
];

export default function ReviewOrderPage() {
    const { id } = useParams();
    const navigate = useNavigate();
    const { t } = useTranslation();

    const [submitting, setSubmitting] = useState(false);
    const [rating, setRating] = useState(5);
    const [content, setContent] = useState('');
    const [selectedTags, setSelectedTags] = useState<string[]>([]);

    // In a real app, we might fetch order details here to show who we are reviewing
    // For now, we assume the user knows which order they clicked on

    const toggleTag = (tagId: string) => {
        setSelectedTags(prev =>
            prev.includes(tagId)
                ? prev.filter(t => t !== tagId)
                : [...prev, tagId]
        );
    };

    const handleSubmit = async () => {
        if (!rating) {
            toast.error(t('review.rating_required', { defaultValue: 'Please give a rating' }));
            return;
        }

        setSubmitting(true);
        try {
            await http.post(`/user/orders/${id}/review`, {
                rating,
                content,
                tags: selectedTags
            });

            toast.success(t('review.success', { defaultValue: 'Review submitted!' }));
            navigate('/orders');
        } catch (err) {
            console.error(err);
            toast.error(err instanceof Error ? err.message : t('review.failed', { defaultValue: 'Failed to submit review' }));
        } finally {
            setSubmitting(false);
        }
    };

    return (
        <div className="min-h-screen bg-muted/30 p-4 md:p-8 flex items-center justify-center">
            <div className="w-full max-w-lg space-y-6">
                <Button variant="ghost" className="pl-0 gap-2" onClick={() => navigate(-1)}>
                    <ArrowLeft className="h-4 w-4" />
                    {t('common.back', { defaultValue: 'Back' })}
                </Button>

                <div className="flex flex-col items-center space-y-2 text-center">
                    <h1 className="text-3xl font-bold tracking-tight">{t('review.title', { defaultValue: 'Rate Your Experience' })}</h1>
                    <p className="text-muted-foreground">{t('review.subtitle', { defaultValue: 'How was your gaming session?' })}</p>
                </div>

                <Card>
                    <CardHeader>
                        <CardTitle className="text-center">{t('review.rating_label', { defaultValue: 'Overall Rating' })}</CardTitle>
                    </CardHeader>
                    <CardContent className="space-y-6">
                        <div className="flex justify-center py-4">
                            <StarRating value={rating} onChange={setRating} size="lg" />
                        </div>

                        <div className="space-y-3">
                            <label className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
                                {t('review.tags_label', { defaultValue: 'What did you like?' })}
                            </label>
                            <div className="flex flex-wrap gap-2">
                                {COMMON_TAGS.map(tag => (
                                    <Badge
                                        key={tag.id}
                                        variant={selectedTags.includes(tag.id) ? "default" : "outline"}
                                        className="cursor-pointer px-3 py-1 hover:bg-primary/90 transition-colors"
                                        onClick={() => toggleTag(tag.id)}
                                    >
                                        {tag.label}
                                    </Badge>
                                ))}
                            </div>
                        </div>

                        <div className="space-y-2">
                            <label className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
                                {t('review.content_label', { defaultValue: 'Additional Comments' })}
                            </label>
                            <Textarea
                                placeholder={t('review.content_placeholder', { defaultValue: 'Share your thoughts about the player...' })}
                                className="resize-none min-h-[100px]"
                                value={content}
                                onChange={(e) => setContent(e.target.value)}
                            />
                        </div>
                    </CardContent>
                    <CardFooter>
                        <Button
                            className="w-full"
                            onClick={handleSubmit}
                            disabled={submitting}
                        >
                            {submitting ? (
                                <>
                                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                                    {t('common.submitting', { defaultValue: 'Submitting...' })}
                                </>
                            ) : (
                                <>
                                    <Send className="mr-2 h-4 w-4" />
                                    {t('review.submit', { defaultValue: 'Submit Review' })}
                                </>
                            )}
                        </Button>
                    </CardFooter>
                </Card>
            </div>
        </div>
    );
}
