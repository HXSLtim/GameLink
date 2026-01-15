import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { StarRating } from "@/components/ui/star-rating";
import { Badge } from "@/components/ui/badge";
import { format } from "date-fns";
import { useTranslation } from "react-i18next";

export interface Review {
    id: number;
    userId: number;
    username: string;
    avatar?: string;
    rating: number;
    content: string;
    tags: string[];
    createdAt: string;
}

interface ReviewListProps {
    reviews: Review[];
    className?: string;
}

export function ReviewList({ reviews, className }: ReviewListProps) {
    const { t } = useTranslation();

    if (!reviews || reviews.length === 0) {
        return (
            <div className="flex flex-col items-center justify-center py-12 text-center text-muted-foreground bg-muted/20 rounded-xl border border-dashed">
                <p>{t('review.no_reviews', { defaultValue: 'No reviews yet' })}</p>
            </div>
        );
    }

    return (
        <div className={className}>
            <div className="flex items-center gap-2 mb-4">
                <h3 className="text-lg font-semibold">{t('review.reviews', { defaultValue: 'Reviews' })}</h3>
                <Badge variant="secondary" className="rounded-full px-2 py-0.5 text-xs">
                    {reviews.length}
                </Badge>
            </div>

            <div className="space-y-6">
                {reviews.map((review) => (
                    <div key={review.id} className="flex gap-4">
                        <Avatar className="h-10 w-10 border border-border">
                            <AvatarImage src={review.avatar} />
                            <AvatarFallback>{review.username.slice(0, 2).toUpperCase()}</AvatarFallback>
                        </Avatar>

                        <div className="flex-1 space-y-1.5">
                            <div className="flex items-center justify-between">
                                <span className="font-semibold text-sm">{review.username}</span>
                                <span className="text-xs text-muted-foreground">{format(new Date(review.createdAt), 'MMM d, yyyy')}</span>
                            </div>

                            <StarRating value={review.rating} readOnly size="sm" />

                            {review.tags && review.tags.length > 0 && (
                                <div className="flex flex-wrap gap-1.5 mt-1">
                                    {review.tags.map(tag => (
                                        <span key={tag} className="inline-flex items-center rounded-sm bg-muted px-1.5 py-0.5 text-[10px] font-medium text-muted-foreground ring-1 ring-inset ring-gray-500/10">
                                            {tag}
                                        </span>
                                    ))}
                                </div>
                            )}

                            {review.content && (
                                <p className="text-sm text-foreground/80 leading-relaxed pt-1">
                                    {review.content}
                                </p>
                            )}
                        </div>
                    </div>
                ))}
            </div>
        </div>
    );
}
